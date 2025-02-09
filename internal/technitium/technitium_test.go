package technitium

import (
	"fmt"
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
	require.Equal(t, 3, len(endpoints))

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

func (m mockDnsService) DeleteRecord(record *sdk.Record) error {
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
	deletedRecords = []sdk.Record{}
)

func isRecordCreated(name string, recordType string, content string, ttl int) bool {
	for _, record := range createdRecords {
		if record.Domain == name && record.Type == recordType && *record.IPAddress == content && (ttl == 0 || *record.TTL == ttl) {
			return true
		}
	}

	return false
}
