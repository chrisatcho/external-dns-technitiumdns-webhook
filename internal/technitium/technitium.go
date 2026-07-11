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
	DeleteRecord(record *sdk.RecordRequest) error
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

// GetRecords client get records method
func (c DnsClient) GetRecords() ([]sdk.Record, error) {
	zones, _, err := c.client.ZonesAPI.ListZones()
	if err != nil {
		return nil, fmt.Errorf("GetRecords: %w", err)
	}

	records := make([]sdk.Record, 0)
	for _, zone := range zones {
		// Skip Technitium's built-in internal zones (localhost, reverse zones)
		// and any nameless zone: external-dns never manages these, and querying
		// them wastes requests or errors outright (e.g. domain="").
		if zone.Name == "" || (zone.Internal != nil && *zone.Internal) {
			continue
		}

		rs, _, err := c.client.RecordsAPI.ListRecords(zone.Name)
		if err != nil {
			return nil, fmt.Errorf("GetRecords: %w", err)
		}
		records = append(records, rs...)
	}
	return records, nil
}

// CreateRecords client create records method
func (c DnsClient) CreateRecord(record *sdk.RecordRequest) error {
	_, _, err := c.client.RecordsAPI.CreateRecord(record)
	return err
}

// DeleteRecord client delete record method
func (c DnsClient) DeleteRecord(r *sdk.RecordRequest) error {
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
		BaseProvider: provider.BaseProvider{},
		client:       DnsClient{client: client},
		domainFilter: domainFilter,
	}

	return prov
}

// Records returns the list of resource records in all zones.
func (p *Provider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	endpoints := make([]*endpoint.Endpoint, 0)
	endpointByNameAndType := make(map[string]*endpoint.Endpoint)

	records, err := p.client.GetRecords()
	if err != nil {
		log.Errorf("failed to fetch records: %v", err)
		return nil, fmt.Errorf("failed to fetch records: %w", err)
	}

	for _, r := range records {
		ep := recordToEndpoint(r)
		if ep == nil {
			continue
		}

		if !p.domainFilter.Match(ep.DNSName) {
			continue
		}

		key := endpointKey(ep)
		existingEndpoint, ok := endpointByNameAndType[key]
		if ok {
			existingEndpoint.Targets = append(existingEndpoint.Targets, ep.Targets...)
			continue
		}

		endpointByNameAndType[key] = ep
		endpoints = append(endpoints, ep)
	}

	log.Debugf("Records() found %d endpoints", len(endpoints))
	return endpoints, nil
}

// ApplyChanges applies a given set of changes.
func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	if changes == nil {
		return fmt.Errorf("changes cannot be nil")
	}

	log.Infof("applying changes: create=%d, updateOld=%d, updateNew=%d, delete=%d",
		len(changes.Create), len(changes.UpdateOld), len(changes.UpdateNew), len(changes.Delete))

	toCreate := make([]*endpoint.Endpoint, len(changes.Create))
	copy(toCreate, changes.Create)

	toDelete := make([]*endpoint.Endpoint, len(changes.Delete))
	copy(toDelete, changes.Delete)

	if len(changes.UpdateOld) != len(changes.UpdateNew) {
		return fmt.Errorf("mismatched update endpoints: %d old, %d new",
			len(changes.UpdateOld), len(changes.UpdateNew))
	}
	for i, updateOldEndpoint := range changes.UpdateOld {
		if !sameEndpoints(*updateOldEndpoint, *changes.UpdateNew[i]) {
			toDelete = append(toDelete, updateOldEndpoint)
			toCreate = append(toCreate, changes.UpdateNew[i])
		}
	}

	for _, e := range toDelete {
		for _, target := range e.Targets {
			r := endpointToRecordRequest(e, target)
			if r == nil {
				log.Warnf("unsupported record type for deletion: %s", e.RecordType)
				continue
			}
			if err := p.client.DeleteRecord(r); err != nil {
				log.Errorf("failed to delete record %s: %v", e.DNSName, err)
				return fmt.Errorf("failed to delete record %s: %w", e.DNSName, err)
			}
			log.Infof("deleted record: %s %s -> %s", e.DNSName, e.RecordType, target)
		}
	}

	for _, e := range toCreate {
		for _, target := range e.Targets {
			r := endpointToRecordRequest(e, target)
			if r == nil {
				log.Warnf("unsupported record type for creation: %s", e.RecordType)
				continue
			}
			if err := p.client.CreateRecord(r); err != nil {
				log.Errorf("failed to create record %s: %v", e.DNSName, err)
				return fmt.Errorf("failed to create record %s: %w", e.DNSName, err)
			}
			log.Infof("created record: %s %s -> %s", e.DNSName, e.RecordType, target)
		}
	}

	log.Infof("successfully applied all changes")
	return nil
}

// endpointToRecordRequest builds a Technitium record request for a single
// target of an endpoint. It returns nil for unsupported record types.
func endpointToRecordRequest(e *endpoint.Endpoint, target string) *sdk.RecordRequest {
	r := &sdk.RecordRequest{
		Domain: e.DNSName,
		Type:   e.RecordType,
	}

	if e.RecordTTL != 0 {
		ttl := int(e.RecordTTL)
		r.TTL = &ttl
	}

	switch e.RecordType {
	case "A", "AAAA":
		r.IPAddress = &target
	case "CNAME":
		r.CNAME = &target
	case "TXT":
		r.Text = &target
	default:
		return nil
	}

	return r
}

// recordToEndpoint converts a record to an endpoint, or nil for record types
// this provider does not manage.
func recordToEndpoint(r sdk.Record) *endpoint.Endpoint {
	var target *string
	switch r.Type {
	case "A", "AAAA":
		target = r.RData.IPAddress
	case "CNAME":
		target = r.RData.CNAME
	case "TXT":
		target = r.RData.Text
	default:
		return nil
	}

	if target == nil {
		log.Warnf("record %s of type %s has no data", r.Name, r.Type)
		return nil
	}

	return endpoint.NewEndpointWithTTL(r.Name, r.Type, endpoint.TTL(r.TTL), *target)
}

func endpointKey(e *endpoint.Endpoint) string {
	return e.DNSName + "/" + e.RecordType
}

// sameEndpoints returns if the two endpoints have the same values.
func sameEndpoints(a endpoint.Endpoint, b endpoint.Endpoint) bool {
	return a.DNSName == b.DNSName && a.RecordType == b.RecordType && a.RecordTTL == b.RecordTTL && a.Targets.Same(b.Targets)
}
