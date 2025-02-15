package model

import (
	"testing"

	error2 "github.com/TimeSnap/distributed-scheduler/internal/pkg/error"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestAuthTypeValid(t *testing.T) {
	var authType AuthType = "INVALID"

	if authType.Valid() {
		t.Error("Expected false, got true")
	}

	authType = AuthTypeNone

	if !authType.Valid() {
		t.Error("Expected true, got false")
	}
}

func TestHTTPJobValidate(t *testing.T) {
	tests := []struct {
		name string
		job  HTTPJob
		want error
	}{
		{
			name: "valid job",
			job: HTTPJob{
				URL:    "https://example.com",
				Method: "GET",
				Auth: Auth{
					Type: AuthTypeNone,
				},
			},
			want: nil,
		},
		{
			name: "invalid job: empty URL",
			job: HTTPJob{
				URL:    "",
				Method: "GET",
				Auth: Auth{
					Type: AuthTypeNone,
				},
			},
			want: error2.ErrEmptyHTTPJobURL,
		},
		{
			name: "invalid job: empty Method",
			job: HTTPJob{
				URL:    "https://example.com",
				Method: "",
				Auth: Auth{
					Type: AuthTypeNone,
				},
			},
			want: error2.ErrEmptyHTTPJobMethod,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.job.Validate()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAuthValidate(t *testing.T) {
	tests := []struct {
		name string
		auth Auth
		want error
	}{
		{
			name: "valid auth: no auth",
			auth: Auth{
				Type: AuthTypeNone,
			},
			want: nil,
		},
		{
			name: "valid auth: basic auth",
			auth: Auth{
				Type:     AuthTypeBasic,
				Username: null.StringFrom("testuser"),
				Password: null.StringFrom("testpassword"),
			},
			want: nil,
		},
		{
			name: "invalid auth: missing username",
			auth: Auth{
				Type:     AuthTypeBasic,
				Password: null.StringFrom("testpassword"),
			},
			want: error2.ErrEmptyUsername,
		},
		{
			name: "invalid auth: missing password",
			auth: Auth{
				Type:     AuthTypeBasic,
				Username: null.StringFrom("testuser"),
			},
			want: error2.ErrEmptyPassword,
		},
		{
			name: "invalid auth: unsupported auth type",
			auth: Auth{
				Type: "unsupported_type",
			},
			want: error2.ErrInvalidAuthType,
		},
		{
			name: "valid auth: bearer token",
			auth: Auth{
				Type:        AuthTypeBearer,
				BearerToken: null.StringFrom("testtoken"),
			},
			want: nil,
		},
		{
			name: "invalid auth: missing bearer token",
			auth: Auth{
				Type: AuthTypeBearer,
			},
			want: error2.ErrEmptyBearerToken,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.auth.Validate()
			assert.Equal(t, tc.want, got)
		})
	}
}
