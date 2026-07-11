package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

type Configuration struct {
	BaseURL    string
	HTTPClient *http.Client
	Debug      bool
	User       string
	Pass       string
}

type APIClient struct {
	cfg    *Configuration
	common service

	mu    sync.Mutex
	token string

	// API Services
	ZonesAPI   *ZonesAPIService
	RecordsAPI *RecordsAPIService
	UsersAPI   *UsersAPIService
}

type APIResponse[T any] struct {
	Data              T      `json:"response,omitempty"`
	Status            string `json:"status"`
	ErrorMessage      string `json:"errorMessage,omitempty"`
	StackTrace        string `json:"stackTrace,omitempty"`
	InnerErrorMessage string `json:"innerErrorMessage,omitempty"`
}

type service struct {
	client *APIClient
}

func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c

	c.ZonesAPI = (*ZonesAPIService)(&c.common)
	c.RecordsAPI = (*RecordsAPIService)(&c.common)
	c.UsersAPI = (*UsersAPIService)(&c.common)

	return c
}

// getToken returns a cached session token, logging in on first use.
func (c *APIClient) getToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" {
		return c.token, nil
	}

	loginRequest := &LoginRequest{
		User:        c.cfg.User,
		Pass:        c.cfg.Pass,
		IncludeInfo: false,
	}
	token, _, err := c.UsersAPI.Login(loginRequest)
	if err != nil {
		return "", err
	}

	c.token = token
	return token, nil
}

// invalidateToken discards the cached session token so the next call re-logs in.
func (c *APIClient) invalidateToken() {
	c.mu.Lock()
	c.token = ""
	c.mu.Unlock()
}

// callAPI authenticates the request with a cached session token and executes
// it. If the server reports an expired/invalid token, the cached token is
// discarded and the request is retried once with a fresh login.
func (c *APIClient) callAPI(req *http.Request) (*http.Response, error) {
	resp, body, err := c.doAuthenticated(req)
	if err != nil {
		return resp, err
	}

	if isInvalidTokenResponse(body) {
		c.invalidateToken()
		resp, body, err = c.doAuthenticated(req)
		if err != nil {
			return resp, err
		}
	}

	// Hand the buffered body back to the caller for decoding.
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return resp, nil
}

// doAuthenticated injects the session token, performs the request, and buffers
// the response body so callAPI can both inspect and return it.
func (c *APIClient) doAuthenticated(req *http.Request) (*http.Response, []byte, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, nil, fmt.Errorf("authentication failed: %w", err)
	}

	q := req.URL.Query()
	q.Set("token", token)
	req.URL.RawQuery = q.Encode()

	if c.cfg.Debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to dump request: %w", err)
		}
		slog.Debug(string(dump))
	}

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return resp, nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if c.cfg.Debug {
		slog.Debug(string(body))
	}

	return resp, body, nil
}

// isInvalidTokenResponse reports whether Technitium rejected the session token.
func isInvalidTokenResponse(body []byte) bool {
	var probe struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return false
	}
	return probe.Status == "invalid-token"
}

func structToQuery(s any) url.Values {
	values := url.Values{}
	val := reflect.ValueOf(s)

	// If it's a pointer, get the underlying element
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the json tag name, defaulting to field name if not present
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = strings.ToLower(fieldType.Name)
		}
		// Remove the omitempty suffix if present
		tag = strings.Split(tag, ",")[0]

		// Skip empty fields
		if field.IsZero() {
			continue
		}

		// Handle pointer fields
		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				// Get the underlying value
				value := field.Elem()
				values.Set(tag, fmt.Sprintf("%v", value.Interface()))
			}
			continue
		}

		// Handle non-pointer fields
		values.Set(tag, fmt.Sprintf("%v", field.Interface()))
	}

	return values
}
