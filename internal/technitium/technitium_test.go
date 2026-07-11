package technitium

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"

	sdk "github.com/chrisatcho/external-dns-technitiumdns-webhook/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type mockDnsService struct {
	testErrorReturned bool
}

func TestNewProvider(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	domainFilter := endpoint.DomainFilter{}
	p := NewProvider(domainFilter, &Configuration{User: "", Pass: "", APIEndpointURL: ""})
	require.NotNilf(t, p.client, "client should not be nil")
}

func TestRecords(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	provider := &Provider{client: mockDnsService{testErrorReturned: false}}
	endpoints, err := provider.Records(nil)
	if err != nil {
		t.Errorf("should not fail, %s", err)
	}
	for _, e := range endpoints {
		log.Info(e)
	}
	require.Equal(t, 2, len(endpoints))

	var aRecord *endpoint.Endpoint
	for _, e := range endpoints {
		if e.DNSName == "a.au" && e.RecordType == "A" {
			aRecord = e
			break
		}
	}
	require.NotNil(t, aRecord)
	require.ElementsMatch(t, endpoint.Targets{"1.1.1.1", "1.1.1.2"}, aRecord.Targets)

	provider = &Provider{client: mockDnsService{testErrorReturned: true}}
	endpoints, err = provider.Records(nil)
	require.Equal(t, 0, len(endpoints))
}

func TestApplyChanges(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	provider := &Provider{client: mockDnsService{testErrorReturned: false}}
	err := provider.ApplyChanges(nil, changes())
	if err != nil {
		t.Errorf("should not fail, %s", err)
	}

	// 3 records must be deleted
	log.Infof("Deleted records: %v", deletedRecords)
	require.Equal(t, 3, len(deletedRecords))
	// deletes must identify the record by domain + data, not by an empty request
	if !isRecordDeleted("b.au", "A", "5.5.5.5") {
		t.Errorf("Record b.au A 5.5.5.5 not deleted")
	}
	if !isRecordDeleted("a.au", "A", "1.1.1.1") {
		t.Errorf("Record a.au A 1.1.1.1 not deleted")
	}
	if !isRecordDeleted("a.au", "A", "2.2.2.2") {
		t.Errorf("Record a.au A 2.2.2.2 not deleted")
	}
	// 3 records must be created
	if !isRecordCreated("a.au", "A", "3.3.3.3", 2000) {
		t.Errorf("Record a.au A 3.3.3.3 not created")
	}
	if !isRecordCreated("a.au", "A", "4.4.4.4", 2000) {
		t.Errorf("Record a.au A 4.4.4.4 not created")
	}
	if !isRecordCreated("new.a.au", "CNAME", "a.au", 0) {
		t.Errorf("Record new.a.au CNAME a.au not created")
	}

	provider = &Provider{client: mockDnsService{testErrorReturned: true}}
	err = provider.ApplyChanges(nil, nil)

	if err == nil {
		t.Errorf("expected to fail, %s", err)
	}
}

func TestGetRecordsSkipsInternalAndUnnamedZones(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/user/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"token":"t","status":"ok"}`)
	})
	mux.HandleFunc("GET /api/zones/list", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"response":{"zones":[
			{"name":"","type":"Secondary","disabled":false},
			{"name":"localhost","type":"Primary","internal":true,"disabled":false},
			{"name":"0.in-addr.arpa","type":"Primary","internal":true,"disabled":false},
			{"name":"example.com","type":"Primary","internal":false,"disabled":false}
		]},"status":"ok"}`)
	})

	var queried []string
	mux.HandleFunc("GET /api/zones/records/get", func(w http.ResponseWriter, r *http.Request) {
		queried = append(queried, r.URL.Query().Get("domain"))
		fmt.Fprint(w, `{"response":{"zone":{"name":"example.com"},"records":[
			{"name":"example.com","type":"A","ttl":3600,"rData":{"ipAddress":"1.1.1.1"}}
		]},"status":"ok"}`)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := DnsClient{client: sdk.NewAPIClient(&sdk.Configuration{BaseURL: server.URL, User: "u", Pass: "p"})}
	records, err := client.GetRecords()
	require.NoError(t, err)

	// Only the real, externally-managed zone should be queried for records.
	require.Equal(t, []string{"example.com"}, queried)
	require.Len(t, records, 1)
}

func (m mockDnsService) GetZones() ([]sdk.Zone, error) {
	if m.testErrorReturned {
		return nil, fmt.Errorf("GetZones failed")
	}

	a := &sdk.Zone{}
	a.Name = "a.au"
	a.Type = "Secondary"
	a.Disabled = false

	b := &sdk.Zone{}
	b.Name = "b.au"
	b.Type = "Primary"
	b.Disabled = false

	return []sdk.Zone{*a, *b}, nil
}

func (m mockDnsService) GetRecords() ([]sdk.Record, error) {
	if m.testErrorReturned {
		return nil, fmt.Errorf("GetZone failed")
	}

	records := make([]sdk.Record, 0)

	ipAddress1 := "1.1.1.1"
	a := sdk.Record{
		Disabled:     false,
		Name:         "a.au",
		Type:         "A",
		TTL:          3000,
		RData:        sdk.RData{IPAddress: &ipAddress1},
		DNSSecStatus: "Unknown",
	}

	ipAddress2 := "1.1.1.2"
	a2 := sdk.Record{
		Disabled:     false,
		Name:         "a.au",
		Type:         "A",
		TTL:          3000,
		RData:        sdk.RData{IPAddress: &ipAddress2},
		DNSSecStatus: "Unknown",
	}

	ipAddress3 := "2.2.2.2"
	b := sdk.Record{
		Disabled:     false,
		Name:         "b.au",
		Type:         "A",
		TTL:          3000,
		RData:        sdk.RData{IPAddress: &ipAddress3},
		DNSSecStatus: "Unknown",
	}

	records = append(records, a, a2, b)

	return records, nil
}

func (m mockDnsService) CreateRecord(record *sdk.RecordRequest) error {
	createdRecords = append(createdRecords, *record)
	return nil
}

func (m mockDnsService) DeleteRecord(record *sdk.RecordRequest) error {
	log.Infof("Deleting: %v", record)
	deletedRecords = append(deletedRecords, *record)
	return nil
}

func changes() *plan.Changes {
	changes := &plan.Changes{}

	changes.Create = []*endpoint.Endpoint{
		{DNSName: "new.a.au", Targets: endpoint.Targets{"a.au"}, RecordType: "CNAME"},
	}
	changes.Delete = []*endpoint.Endpoint{{DNSName: "b.au", RecordType: "A", Targets: endpoint.Targets{"5.5.5.5"}}}
	changes.UpdateOld = []*endpoint.Endpoint{{DNSName: "a.au", RecordType: "A", Targets: endpoint.Targets{"1.1.1.1", "2.2.2.2"}, RecordTTL: 1000}}
	changes.UpdateNew = []*endpoint.Endpoint{{DNSName: "a.au", RecordType: "A", Targets: endpoint.Targets{"3.3.3.3", "4.4.4.4"}, RecordTTL: 2000}}

	return changes
}

var (
	createdRecords = []sdk.RecordRequest{}
	deletedRecords = []sdk.RecordRequest{}
)

func isRecordCreated(name string, recordType string, content string, ttl int) bool {
	for _, record := range createdRecords {
		if record.Domain != name || record.Type != recordType {
			continue
		}

		// Check TTL if specified
		if ttl != 0 && (record.TTL == nil || *record.TTL != ttl) {
			continue
		}

		// Check content based on record type
		switch recordType {
		case "A", "AAAA":
			if record.IPAddress != nil && *record.IPAddress == content {
				return true
			}
		case "CNAME":
			if record.CNAME != nil && *record.CNAME == content {
				return true
			}
		case "TXT":
			if record.Text != nil && *record.Text == content {
				return true
			}
		}
	}

	return false
}

func isRecordDeleted(name string, recordType string, content string) bool {
	for _, record := range deletedRecords {
		if record.Domain != name || record.Type != recordType {
			continue
		}

		switch recordType {
		case "A", "AAAA":
			if record.IPAddress != nil && *record.IPAddress == content {
				return true
			}
		case "CNAME":
			if record.CNAME != nil && *record.CNAME == content {
				return true
			}
		case "TXT":
			if record.Text != nil && *record.Text == content {
				return true
			}
		}
	}

	return false
}
