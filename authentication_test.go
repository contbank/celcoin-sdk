package celcoin_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/celcoin-sdk"

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

	clientID := celcoin.GetEnvCelcoinClientID()
	clientSecret := celcoin.GetEnvCelcoinClientSecret()

	celcoinConfig := celcoin.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Mtls:         celcoin.Bool(false),
	}

	session, err := celcoin.NewSession(celcoinConfig)

	s.assert.NoError(err)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	s.session = session
	s.authentication = celcoin.NewAuthentication(httpClient, *s.session)
}

func (s *AuthenticationTestSuite) TestToken() {
	token, err := s.authentication.Token(s.ctx)

	s.assert.NoError(err)
	s.assert.Contains(token, "Bearer")
}
