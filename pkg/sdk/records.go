package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type RecordsAPIService service

type ListRecordsResponse struct {
	Zone    Zone     `json:"zone"`
	Records []Record `json:"records"`
}

func (a *RecordsAPIService) ListRecords(domain string) ([]Record, *http.Response, error) {
	reqURL := a.client.cfg.BaseURL + "/api/zones/records/get"
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new ListRecords request: %w", err)
	}

	q := url.Values{}
	q.Set("domain", domain)
	q.Set("listZone", "true")
	req.URL.RawQuery = q.Encode()

	res, err := a.client.callAPI(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do ListRecords request: %w", err)
	}
	defer res.Body.Close()

	var body APIResponse[ListRecordsResponse]
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, nil, fmt.Errorf("decode ListRecords response: %w", err)
	}

	if body.Status != "ok" {
		return nil, nil, fmt.Errorf("response ListRecords status not 'ok': %v, %v", body.Status, body.ErrorMessage)
	}

	return body.Data.Records, res, nil
}

type CreateRecordResponse struct {
	Zone        Zone   `json:"Zone"`
	AddedRecord Record `json:"addedRecord"`
}

type RecordRequest struct {
	Token                          string  `json:"token"`
	Domain                         string  `json:"domain"`
	Zone                           *string `json:"zone,omitempty"`
	Type                           string  `json:"type"`
	TTL                            *int    `json:"ttl,omitempty"`
	Overwrite                      *bool   `json:"overwrite,omitempty"`
	Comments                       *string `json:"comments,omitempty"`
	ExpiryTTL                      *int    `json:"expiryTtl,omitempty"`
	IPAddress                      *string `json:"ipAddress,omitempty"`
	PTR                            *bool   `json:"ptr,omitempty"`
	CreatePTRZone                  *bool   `json:"createPtrZone,omitempty"`
	UpdateSvcbHints                *bool   `json:"updateSvcbHints,omitempty"`
	NameServer                     *string `json:"nameServer,omitempty"`
	Glue                           *string `json:"glue,omitempty"`
	CNAME                          *string `json:"cname,omitempty"`
	PTRName                        *string `json:"ptrName,omitempty"`
	Exchange                       *string `json:"exchange,omitempty"`
	Preference                     *int    `json:"preference,omitempty"`
	Text                           *string `json:"text,omitempty"`
	SplitText                      *bool   `json:"splitText,omitempty"`
	Mailbox                        *string `json:"mailbox,omitempty"`
	TXTDName                       *string `json:"txtDomain,omitempty"`
	Priority                       *int    `json:"priority,omitempty"`
	Weight                         *int    `json:"weight,omitempty"`
	Port                           *int    `json:"port,omitempty"`
	Target                         *string `json:"target,omitempty"`
	NAPTROrder                     *int    `json:"naptrOrder,omitempty"`
	NAPTRPreference                *int    `json:"naptrPreference,omitempty"`
	NAPTRFlags                     *string `json:"naptrFlags,omitempty"`
	NAPTRServices                  *string `json:"naptrServices,omitempty"`
	NAPTRRegexp                    *string `json:"naptrRegexp,omitempty"`
	NAPTRReplacement               *string `json:"naptrReplacement,omitempty"`
	DNAME                          *string `json:"dname,omitempty"`
	KeyTag                         *int    `json:"keyTag,omitempty"`
	Algorithm                      *string `json:"algorithm,omitempty"`
	DigestType                     *string `json:"digestType,omitempty"`
	Digest                         *string `json:"digest,omitempty"`
	SSHFPAlgorithm                 *string `json:"sshfpAlgorithm,omitempty"`
	SSHFPFingerprintType           *string `json:"sshfpFingerprintType,omitempty"`
	SSHFPFingerprint               *string `json:"sshfpFingerprint,omitempty"`
	TLSACertificateUsage           *string `json:"tlsaCertificateUsage,omitempty"`
	TLSASelector                   *string `json:"tlsaSelector,omitempty"`
	TLSAMatchingType               *string `json:"tlsaMatchingType,omitempty"`
	TLSACertificateAssociationData *string `json:"tlsaCertificateAssociationData,omitempty"`
	SVCPriority                    *int    `json:"svcPriority,omitempty"`
	SVCTargetName                  *string `json:"svcTargetName,omitempty"`
	SVCParams                      *string `json:"svcParams,omitempty"`
	AutoIPv4Hint                   *bool   `json:"autoIpv4Hint,omitempty"`
	AutoIPv6Hint                   *bool   `json:"autoIpv6Hint,omitempty"`
	URIPriority                    *int    `json:"uriPriority,omitempty"`
	URIWeight                      *int    `json:"uriWeight,omitempty"`
	URI                            *string `json:"uri,omitempty"`
	Flags                          *int    `json:"flags,omitempty"`
	Tag                            *string `json:"tag,omitempty"`
	Value                          *string `json:"value,omitempty"`
	ANAME                          *string `json:"aname,omitempty"`
	Protocol                       *string `json:"protocol,omitempty"`
	Forwarder                      *string `json:"forwarder,omitempty"`
	ForwarderPriority              *int    `json:"forwarderPriority,omitempty"`
	DNSSECValidation               *bool   `json:"dnssecValidation,omitempty"`
	ProxyType                      *string `json:"proxyType,omitempty"`
	ProxyAddress                   *string `json:"proxyAddress,omitempty"`
	ProxyPort                      *int    `json:"proxyPort,omitempty"`
	ProxyUsername                  *string `json:"proxyUsername,omitempty"`
	ProxyPassword                  *string `json:"proxyPassword,omitempty"`
	AppName                        *string `json:"appName,omitempty"`
	ClassPath                      *string `json:"classPath,omitempty"`
	RecordData                     *string `json:"recordData,omitempty"`
	RData                          *string `json:"rdata,omitempty"`
}

func (a *RecordsAPIService) CreateRecord(r *RecordRequest) (*Record, *http.Response, error) {
	reqURL := a.client.cfg.BaseURL + `/api/zones/records/add`
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new CreateRecord request: %w", err)
	}

	q := structToQuery(r)
	req.URL.RawQuery = q.Encode()

	res, err := a.client.callAPI(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do CreateRecord request: %w", err)
	}
	defer res.Body.Close()

	var body APIResponse[CreateRecordResponse]
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, nil, fmt.Errorf("decode CreateRecord response: %w", err)
	}

	if body.Status != "ok" {
		return nil, nil, fmt.Errorf("response CreateRecords status not 'ok': %v, %v", body.Status, body.ErrorMessage)
	}

	return &body.Data.AddedRecord, res, nil
}

func (a *RecordsAPIService) DeleteRecord(r *Record) (*http.Response, error) {
	q := url.Values{}
	q.Set("domain", r.Name)
	q.Set("type", r.Type)

	switch r.Type {
	case "A":
		q.Set("ipAddress", *r.RData.IPAddress)
	case "AAAA":
		q.Set("ipAddress", *r.RData.IPAddress)
	case "CNAME":
		q.Set("cname", *r.RData.CNAME)
	case "TXT":
		q.Set("text", *r.RData.Text)
	}

	url := a.client.cfg.BaseURL + `/api/zones/records/delete`
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new DeleteRecord request: %w", err)
	}

	req.URL.RawQuery = q.Encode()

	res, err := a.client.callAPI(req)
	if err != nil {
		return nil, fmt.Errorf("do DeleteRecord request: %w", err)
	}
	defer res.Body.Close()

	var body APIResponse[interface{}]
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("decode DeleteRecord response: %w", err)
	}

	if body.Status != "ok" {
		return nil, fmt.Errorf("response DeleteRecords status not 'ok': %v, %v", body.Status, body.ErrorMessage)
	}

	return res, nil
}

type Record struct {
	Disabled     bool    `json:"disabled"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	TTL          int     `json:"ttl"`
	RData        RData   `json:"rData"`
	DNSSecStatus string  `json:"dnssecStatus"`
	LastUsedOn   *string `json:"lastUsedOn,omitempty"`
}

type RData struct {
	IPAddress  *string `json:"ipAddress,omitempty"`
	CNAME      *string `json:"cname,omitempty"`
	NameServer *string `json:"nameServer,omitempty"`
	Text       *string `json:"text,omitempty"`
}
