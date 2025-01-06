package celcoin_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/celcoin-sdk"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticationTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	ctx            context.Context
	session        *celcoin.Session
	authentication *celcoin.Authentication
}

func TestAuthenticationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}

func (s *AuthenticationTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	clientID := "test-client-id"
	clientSecret := "test-client-secret"
	apiEndpoint := "https://sandbox.openfinance.celcoin.dev"
	loginEndpoint := "https://sandbox.openfinance.celcoin.dev"

	celcoinConfig := celcoin.Config{
		ClientID:      &clientID,
		ClientSecret:  &clientSecret,
		Mtls:          celcoin.Bool(false),
		APIEndpoint:   &apiEndpoint,
		LoginEndpoint: &loginEndpoint,
	}

	session, err := celcoin.NewSession(celcoinConfig)
	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.authentication = celcoin.NewAuthentication(httpClient, *s.session)

	// Initialize httpmock
	httpmock.ActivateNonDefault(httpClient)
}

func (s *AuthenticationTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *AuthenticationTestSuite) TestToken() {
	mockURL := fmt.Sprintf("%s/%s", s.session.LoginEndpoint, celcoin.LoginPath)

	httpmock.RegisterResponder("POST", mockURL,
		httpmock.NewStringResponder(200, `{
            "access_token": "mock-access-token",
            "expires_in": 3600
        }`),
	)

	token, err := s.authentication.Token(s.ctx)

	s.assert.NoError(err)
	s.assert.Equal("Bearer mock-access-token", token)
}

func (s *AuthenticationTestSuite) TestTokenError() {
	mockURL := fmt.Sprintf("%s/%s", s.session.LoginEndpoint, celcoin.LoginPath)

	httpmock.RegisterResponder("POST", mockURL,
		httpmock.NewStringResponder(400, `{
			"error": "invalid_request",
			"message": "Invalid client credentials."
		}`),
	)

	token, err := s.authentication.Token(s.ctx)

	s.assert.Error(err)
	s.assert.Empty(token)
	s.assert.Contains(err.Error(), "invalid_request")
}
