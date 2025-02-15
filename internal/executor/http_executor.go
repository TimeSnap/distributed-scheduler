package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
	errors "github.com/TimeSnap/distributed-scheduler/internal/pkg/error"
)

// HTTPSPrefix and HTTPPrefix are prefixes for HTTP and HTTPS protocols
const (
	HTTPSPrefix = "https://"
	HTTPPrefix  = "http://"
)

type httpExecutor struct {
	Client HttpClient
}

// HttpClient interface
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (he *httpExecutor) Execute(ctx context.Context, j *model.Job) error {
	// Create the HTTP request
	req, err := he.createHTTPRequest(ctx, j)
	if err != nil {
		return err
	}

	// Send the request and get the response
	resp, err := he.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if status code is one of the valid response codes
	if !he.validResponseCode(resp.StatusCode, j.HTTPJob.ValidResponseCodes) {
		return errors.ErrInvalidResponseCode
	}

	return nil
}

func (he *httpExecutor) validResponseCode(code int, validCodes []int) bool {
	// If no valid response codes are defined, 200 is the default
	if len(validCodes) == 0 {
		return code == http.StatusOK
	}

	// Check if the response code is one of the valid response codes
	for _, c := range validCodes {
		if c == code {
			return true
		}
	}

	return false
}

func (he *httpExecutor) createHTTPRequest(ctx context.Context, j *model.Job) (*http.Request, error) {
	// Create the request body
	body := he.createHTTPRequestBody(j.HTTPJob.Body.String)

	// Create the request URL
	url := he.createHTTPRequestURL(j.HTTPJob.URL)

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, j.HTTPJob.Method, url, body)
	if err != nil {
		return nil, err
	}

	// Set the headers
	he.setHTTPRequestHeaders(req, j.HTTPJob.Headers)

	// Set the auth
	he.setHTTPRequestAuth(req, j.HTTPJob.Auth)

	return req, nil
}

func (he *httpExecutor) createHTTPRequestBody(body string) io.Reader {
	if body == "" {
		return nil
	}

	return strings.NewReader(body)
}

func (he *httpExecutor) createHTTPRequestURL(url string) string {
	if strings.HasPrefix(url, HTTPPrefix) || strings.HasPrefix(url, HTTPSPrefix) {
		return url
	}

	return HTTPSPrefix + url
}

func (he *httpExecutor) setHTTPRequestHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func (he *httpExecutor) setHTTPRequestAuth(req *http.Request, auth model.Auth) {
	switch auth.Type {
	case model.AuthTypeBasic:
		req.SetBasicAuth(auth.Username.String, auth.Password.String)
	case model.AuthTypeBearer:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.BearerToken.String))
	}
}
