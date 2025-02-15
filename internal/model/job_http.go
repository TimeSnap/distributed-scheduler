package model

import (
	error2 "github.com/TimeSnap/distributed-scheduler/internal/pkg/error"
	"gopkg.in/guregu/null.v4"
)

type AuthType string

const (
	AuthTypeNone   AuthType = "none"
	AuthTypeBasic  AuthType = "basic"
	AuthTypeBearer AuthType = "bearer"
)

func (at AuthType) Valid() bool {
	switch at {
	case AuthTypeNone, AuthTypeBasic, AuthTypeBearer:
		return true
	default:
		return false
	}
}

type HTTPJob struct {
	URL                string            `json:"url"`                       // e.g., "https://example.com"
	Method             string            `json:"method"`                    // e.g., "GET", "POST", "PUT", "PATCH", "DELETE"
	Headers            map[string]string `json:"headers"`                   // e.g., {"Content-Type": "application/json"}
	Body               null.String       `json:"body" swaggertype:"string"` // e.g., "{\"hello\": \"world\"}"
	ValidResponseCodes []int             `json:"valid_response_codes"`      // e.g., [200, 201, 202]
	Auth               Auth              `json:"auth"`                      // e.g., {"type": "basic", "username": "foo", "password": "bar"}
}

// Validate validates an HTTPJob struct.
func (httpJob *HTTPJob) Validate() error {
	if httpJob == nil {
		return error2.ErrHTTPJobNotDefined
	}

	if httpJob.URL == "" {
		return error2.ErrEmptyHTTPJobURL
	}

	if httpJob.Method == "" {
		return error2.ErrEmptyHTTPJobMethod
	}

	if err := httpJob.Auth.Validate(); err != nil {
		return err
	}

	return nil
}

func (httpJob *HTTPJob) RemoveCredentials() {
	httpJob.Auth.Username = null.String{}
	httpJob.Auth.Password = null.String{}
	httpJob.Auth.BearerToken = null.String{}
}

type Auth struct {
	Type        AuthType    `json:"type"`                                        // e.g., "none", "basic", "bearer"
	Username    null.String `json:"username,omitempty" swaggertype:"string"`     // for "basic"
	Password    null.String `json:"password,omitempty" swaggertype:"string"`     // for "basic"
	BearerToken null.String `json:"bearer_token,omitempty" swaggertype:"string"` // for "bearer"
}

func (auth *Auth) Validate() error {
	if auth == nil {
		return error2.ErrAuthMethodNotDefined
	}

	if !auth.Type.Valid() {
		return error2.ErrInvalidAuthType
	}

	if auth.Type == AuthTypeBasic {
		if !auth.Username.Valid || auth.Username.String == "" {
			return error2.ErrEmptyUsername
		}

		if !auth.Password.Valid || auth.Password.String == "" {
			return error2.ErrEmptyPassword
		}
	}

	if auth.Type == AuthTypeBearer && (!auth.BearerToken.Valid || auth.BearerToken.String == "") {
		return error2.ErrEmptyBearerToken
	}

	return nil
}
