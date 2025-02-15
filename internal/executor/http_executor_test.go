package executor

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestHTTPExecutor_Execute(t *testing.T) {
	ctx := context.Background()
	j := &model.Job{
		HTTPJob: &model.HTTPJob{
			Method: "GET",
			URL:    "www.example.com",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Auth: model.Auth{
				Type:     model.AuthTypeBasic,
				Username: null.StringFrom("username"),
				Password: null.StringFrom("password"),
			},
			ValidResponseCodes: []int{200, 201, 202},
		},
	}

	t.Run("successful request", func(t *testing.T) {
		mockHttpClient := &MockHttpClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       httptest.NewRecorder().Result().Body,
				}, nil
			},
		}

		httpExecutor := &httpExecutor{Client: mockHttpClient}
		err := httpExecutor.Execute(ctx, j)

		assert.Nil(t, err)
	})

	t.Run("client error", func(t *testing.T) {
		mockHttpClient := &MockHttpClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("client error")
			},
		}

		httpExecutor := &httpExecutor{Client: mockHttpClient}
		err := httpExecutor.Execute(ctx, j)

		assert.NotNil(t, err)
	})

	t.Run("invalid status code - 500", func(t *testing.T) {
		mockHttpClient := &MockHttpClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       httptest.NewRecorder().Result().Body,
				}, nil
			},
		}

		httpExecutor := &httpExecutor{Client: mockHttpClient}
		err := httpExecutor.Execute(ctx, j)

		assert.NotNil(t, err)
	})

	t.Run("valid 202 status code", func(t *testing.T) {
		mockHttpClient := &MockHttpClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusAccepted,
					Body:       httptest.NewRecorder().Result().Body,
				}, nil
			},
		}

		httpExecutor := &httpExecutor{Client: mockHttpClient}
		err := httpExecutor.Execute(ctx, j)

		assert.Nil(t, err)
	})
}

func TestHTTPExecutor_createHTTPRequest(t *testing.T) {
	ctx := context.Background()
	j := &model.Job{
		HTTPJob: &model.HTTPJob{
			Method: "GET",
			URL:    "www.example.com",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Auth: model.Auth{
				Type:     model.AuthTypeBasic,
				Username: null.StringFrom("username"),
				Password: null.StringFrom("password"),
			},
		},
	}

	httpExecutor := &httpExecutor{}
	req, err := httpExecutor.createHTTPRequest(ctx, j)

	assert.Nil(t, err)
	assert.Equal(t, j.HTTPJob.Method, req.Method)
	assert.Equal(t, j.HTTPJob.Headers["Content-Type"], req.Header.Get("Content-Type"))
	assert.Equal(t, "Basic "+base64.StdEncoding.EncodeToString([]byte(j.HTTPJob.Auth.Username.String+":"+j.HTTPJob.Auth.Password.String)), req.Header.Get("Authorization"))
}

func TestHTTPExecutor_validResponseCode(t *testing.T) {
	httpExecutor := &httpExecutor{}

	validResponseCodes := []int{200, 201}

	assert.True(t, httpExecutor.validResponseCode(200, validResponseCodes))
	assert.False(t, httpExecutor.validResponseCode(404, validResponseCodes))

	validResponseCodes = []int{}
	assert.True(t, httpExecutor.validResponseCode(http.StatusOK, validResponseCodes))
	assert.False(t, httpExecutor.validResponseCode(http.StatusInternalServerError, validResponseCodes))
}
