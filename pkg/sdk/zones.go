package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ZonesAPIService service

type ListZonesResponse struct {
	PageNumber int    `json:"pageNumber"`
	TotalPages int    `json:"totalPages"`
	TotalZones int    `json:"totalZones"`
	Zones      []Zone `json:"zones"`
}

func (a *ZonesAPIService) ListZones() ([]Zone, *http.Response, error) {
	url := a.client.cfg.BaseURL + "/api/zones/list"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new ListZones request: %w", err)
	}

	res, err := a.client.callAPI(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do ListZones request: %w", err)
	}
	defer res.Body.Close()

	var body APIResponse[ListZonesResponse]
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return nil, nil, fmt.Errorf("decode ListZones response: %w", err)
	}

	if body.Status != "ok" {
		return nil, nil, fmt.Errorf("response status not 'ok': %v, %v", body.Status, body.ErrorMessage)
	}

	return body.Data.Zones, res, nil
}

type Zone struct {
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	DNSSECStatus    *string    `json:"dnssecStatus,omitempty"`
	SOASerial       *int       `json:"soaSerial,omitempty"`
	LastModified    *time.Time `json:"lastModified,omitempty"`
	Disabled        bool       `json:"disabled"`
	Internal        *bool      `json:"internal,omitempty"`
	Expiry          *time.Time `json:"expiry,omitempty"`
	IsExpired       *bool      `json:"isExpired,omitempty"`
	SyncFailed      *bool      `json:"syncFailed,omitempty"`
	NotifyFailed    *bool      `json:"notifyFailed,omitempty"`
	NotifyFailedFor []string   `json:"notifyFailedFor,omitempty"`
}
