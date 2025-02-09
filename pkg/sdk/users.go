package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type UsersAPIService service

type LoginResponse struct {
	DisplayName       *string `json:"displayName,omitempty"`
	Username          *string `json:"username,omitempty"`
	Token             *string `json:"token,omitempty"`
	Status            string  `json:"status"`
	ErrorMessage      string  `json:"errorMessage,omitempty"`
	StackTrace        string  `json:"stackTrace,omitempty"`
	InnerErrorMessage string  `json:"innerErrorMessage,omitempty"`
}

func (a *UsersAPIService) Login(user, pass string) (string, *http.Response, error) {
	reqURL := a.client.cfg.BaseURL + "/api/user/login"

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("new Login request: %w", err)
	}

	q := url.Values{}
	q.Set("user", user)
	q.Set("pass", pass)
	q.Set("includeInfo", "false")
	req.URL.RawQuery = q.Encode()

	res, err := a.client.cfg.HTTPClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("do Login request: %w", err)
	}
	defer res.Body.Close()

	var body LoginResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return "", nil, fmt.Errorf("decode Login response: %w", err)
	}

	if body.Status != "ok" {
		return "", nil, fmt.Errorf("response status not 'ok': %v, %v", body.Status, body.ErrorMessage)
	}

	return *body.Token, nil, nil
}
