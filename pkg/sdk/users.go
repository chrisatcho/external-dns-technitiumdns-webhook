package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UsersAPIService service

type LoginRequest struct {
	User        string `json:"user"`
	Pass        string `json:"pass"`
	IncludeInfo bool   `json:"includeInfo"`
}

type LoginResponse struct {
	DisplayName       *string `json:"displayName,omitempty"`
	Username          *string `json:"username,omitempty"`
	Token             *string `json:"token,omitempty"`
	Status            string  `json:"status"`
	ErrorMessage      string  `json:"errorMessage,omitempty"`
	StackTrace        string  `json:"stackTrace,omitempty"`
	InnerErrorMessage string  `json:"innerErrorMessage,omitempty"`
}

func (a *UsersAPIService) Login(r *LoginRequest) (string, *http.Response, error) {
	reqURL := a.client.cfg.BaseURL + "/api/user/login"

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("new Login request: %w", err)
	}

	q := structToQuery(r)
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

	if body.Token == nil {
		return "", nil, fmt.Errorf("login response missing token")
	}

	return *body.Token, nil, nil
}
