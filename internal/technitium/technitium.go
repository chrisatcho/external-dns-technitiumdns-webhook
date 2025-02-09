package technitium

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	sdk "github.com/chrisatcho/external-dns-technitiumdns-webhook/pkg/sdk"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

// Provider implements the DNS provider for Technitium DNS.
type Provider struct {
	provider.BaseProvider

	client       DnsService
	domainFilter endpoint.DomainFilter
}

// Configuration holds configuration from environmental variables
type Configuration struct {
	User           string `env:"TECHNITIUM_USER,notEmpty"`
	Pass           string `env:"TECHNITIUM_PASS,notEmpty"`
	APIEndpointURL string `env:"TECHNITIUM_API_URL,notEmpty"`
	Debug          bool   `env:"TECHNITIUM_DEBUG" envDefault:"false"`
}

// DnsService interface to the dns backend, also needed for creating mocks in tests
type DnsService interface {
	GetZones() ([]sdk.Zone, error)
	GetRecords() ([]sdk.Record, error)
	CreateRecord(records *sdk.RecordRequest) error
	DeleteRecord(record *sdk.Record) error
}

// DnsClient client of the dns api
type DnsClient struct {
	client *sdk.APIClient
}

// GetZones client get zones method
func (c DnsClient) GetZones() ([]sdk.Zone, error) {
	zones, _, err := c.client.ZonesAPI.ListZones()
	return zones, err
}

// GetZone client get zone method
func (c DnsClient) GetZone(zoneName string) (*sdk.Zone, error) {
	zones, _, err := c.client.ZonesAPI.ListZones()
	if err != nil {
		return nil, err
	}

	for _, zone := range zones {
		if zone.Name == zoneName {
			return &zone, err
		}
	}

	return nil, fmt.Errorf("Zone %v not found", zoneName)
}

// GetRecords client get records method
func (c DnsClient) GetRecords() ([]sdk.Record, error) {
	zones, _, err := c.client.ZonesAPI.ListZones()
	records := make([]sdk.Record, 0)
	for _, zone := range zones {
		rs, _, err := c.client.RecordsAPI.ListRecords(zone.Name)
		if err != nil {
			return nil, fmt.Errorf("GetRecords: %w", err)
		}
		records = append(records, rs...)
	}
	return records, err
}

// CreateRecords client create records method
func (c DnsClient) CreateRecord(record *sdk.RecordRequest) error {
	_, _, err := c.client.RecordsAPI.CreateRecord(record)
	return err
}

// DeleteRecord client delete record method
func (c DnsClient) DeleteRecord(r *sdk.Record) error {
	_, err := c.client.RecordsAPI.DeleteRecord(r)
	return err
}

// NewProvider creates a new Technitium DNS provider.
func NewProvider(domainFilter endpoint.DomainFilter, configuration *Configuration) *Provider {
	cfg := &sdk.Configuration{
		BaseURL: configuration.APIEndpointURL,
		User:    configuration.User,
		Pass:    configuration.Pass,
		Debug:   configuration.Debug,
	}
	client := sdk.NewAPIClient(cfg)

	prov := &Provider{
		BaseProvider: *&provider.BaseProvider{},
		client:       DnsClient{client: client},
		domainFilter: domainFilter,
	}

	return prov
}

// Records returns the list of resource records in all zones.
func (p *Provider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	endpoints := make([]*endpoint.Endpoint, 0)

	records, err := p.client.GetRecords()
	if err != nil {
		log.Warnf("Failed to fetch records: %v", err)
	}

	for _, r := range records {
		endpoint := recordToEndpoint(r)
		if endpoint == nil {
			continue
		}

		if !p.domainFilter.Match(endpoint.DNSName) {
			continue
		}

		endpoints = append(endpoints, endpoint)
	}

	log.Debugf("Records() found %d endpoints: %v", len(endpoints), endpoints)
	return endpoints, nil
}

// ApplyChanges applies a given set of changes.
func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	if changes == nil {
		return fmt.Errorf("changes cannot be nil")
	}

	log.Warnf("Request to ApplyChanges: %v", changes.Create)
	toCreate := make([]*endpoint.Endpoint, len(changes.Create))
	copy(toCreate, changes.Create)

	toDelete := make([]*endpoint.Endpoint, len(changes.Delete))
	copy(toDelete, changes.Delete)

	for i, updateOldEndpoint := range changes.UpdateOld {
		if !sameEndpoints(*updateOldEndpoint, *changes.UpdateNew[i]) {
			toDelete = append(toDelete, updateOldEndpoint)
			toCreate = append(toCreate, changes.UpdateNew[i])
		}
	}

	for _, e := range toDelete {
		rs := endpointToRecords(e)
		for _, r := range rs {
			p.client.DeleteRecord(&r)
		}
	}

	for _, e := range toCreate {
		ttl := int(e.RecordTTL)
		for _, t := range e.Targets {
			ipAddress := t
			r := &sdk.RecordRequest{
				Domain:    e.DNSName,
				Type:      e.RecordType,
				TTL:       &ttl,
				IPAddress: &ipAddress,
			}
			p.client.CreateRecord(r)
		}
	}

	return nil
}

// endpointToRecords converts an endpoint to a slice of records.
func endpointToRecords(endpoint *endpoint.Endpoint) []sdk.Record {
	records := make([]sdk.Record, 0)

	for _, target := range endpoint.Targets {
		record := &sdk.Record{}

		record.Name = endpoint.DNSName
		record.Type = endpoint.RecordType

		switch record.Type {
		case "A":
			record.RData.IPAddress = &target
		case "AAAA":
			record.RData.IPAddress = &target
		case "CNAME":
			record.RData.CNAME = &target
		case "TXT":
			record.RData.Text = &target
		}

		ttl := int(endpoint.RecordTTL)
		if ttl != 0 {
			record.TTL = ttl
		}

		records = append(records, *record)
	}

	return records
}

// recordToEndpoint converts a record to an endpoint.
func recordToEndpoint(r sdk.Record) *endpoint.Endpoint {
	switch r.Type {
	case "A":
		return endpoint.NewEndpointWithTTL(r.Name, r.Type, endpoint.TTL(r.TTL), *r.RData.IPAddress)
	case "AAAA":
		return endpoint.NewEndpointWithTTL(r.Name, r.Type, endpoint.TTL(r.TTL), *r.RData.IPAddress)
	case "CNAME":
		return endpoint.NewEndpointWithTTL(r.Name, r.Type, endpoint.TTL(r.TTL), *r.RData.CNAME)
	case "TXT":
		return endpoint.NewEndpointWithTTL(r.Name, r.Type, endpoint.TTL(r.TTL), *r.RData.Text)
	}
	return nil
}

// sameEndpoints returns if the two endpoints have the same values.
func sameEndpoints(a endpoint.Endpoint, b endpoint.Endpoint) bool {
	result := (a.DNSName == b.DNSName && a.RecordType == b.RecordType && a.RecordTTL == b.RecordTTL && a.Targets.Same(b.Targets))
	if !result {
		log.Warnf("Endpoints do not match: %v, %v", a, b)
	}
	return result
}
