package sdk

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
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

func (c *APIClient) callAPI(req *http.Request) (*http.Response, error) {
	token, _, err := c.UsersAPI.Login(c.cfg.User, c.cfg.Pass)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	if c.cfg.Debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return nil, err
		}
		slog.Debug(string(dump))
	}

	resp, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return resp, err
	}

	if c.cfg.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}
		slog.Debug(string(dump))
	}

	return resp, err
}

func structToQuery(s interface{}) url.Values {
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
