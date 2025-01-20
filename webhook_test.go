package celcoin_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/contbank/celcoin-sdk"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WebhookTestSuite é a suite de testes para os webhooks.
type WebhookTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	ctx            context.Context
	session        *celcoin.Session
	webhooks       celcoin.Webhooks
	mockServer     *httptest.Server
	authentication *celcoin.Authentication
	httpClient     *http.Client
}

// TestWebhookTestSuite executa a suite de testes.
func TestWebhookTestSuite(t *testing.T) {
	suite.Run(t, new(WebhookTestSuite))
}

// SetupTest inicializa os recursos para os testes.
func (s *WebhookTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	clientID := "test-client-id"
	clientSecret := "test-client-secret"
	apiEndpoint := "https://sandbox.openfinance.celcoin.dev"
	loginEndpoint := "https://sandbox.openfinance.celcoin.dev"

	// Configuração da sessão mockada
	celcoinConfig := celcoin.Config{
		ClientID:      &clientID,
		ClientSecret:  &clientSecret,
		Mtls:          celcoin.Bool(false),
		APIEndpoint:   &apiEndpoint,
		LoginEndpoint: &loginEndpoint,
	}

	session, err := celcoin.NewSession(celcoinConfig)
	s.assert.NoError(err)

	// Cliente HTTP
	s.httpClient = &http.Client{}
	s.session = session
	s.authentication = celcoin.NewAuthentication(s.httpClient, *s.session)

	// Inicialização do mock HTTP
	httpmock.ActivateNonDefault(s.httpClient)

	// Registrar mock do endpoint de autenticação
	mockAuthURL := fmt.Sprintf("%s/%s", s.session.LoginEndpoint, celcoin.LoginPath)
	httpmock.RegisterResponder("POST", mockAuthURL,
		httpmock.NewStringResponder(200, `{
			"access_token": "mock-access-token",
			"expires_in": 3600
		}`),
	)

	// Mock do servidor para simular a API Celcoin
	s.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/baas-webhookmanager/v1/webhook/subscription") {
			// Validar cabeçalhos
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			// Validar payload
			body, _ := io.ReadAll(r.Body)
			var req celcoin.WebhookSubscriptionRequest
			err := json.Unmarshal(body, &req)
			if err != nil || req.Entity == "" || req.WebhookURL == "" {
				http.Error(w, "invalid request payload", http.StatusBadRequest)
				return
			}

			// Retornar resposta simulada
			w.WriteHeader(http.StatusOK)
			response := celcoin.WebhookSubscriptionResponse{
				Version: "1.0.0",
				Status:  "SUCCESS",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))

	// Registrar mock para CreateSubscription
	mockSubscriptionURL := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/subscription", s.mockServer.URL)
	httpmock.RegisterResponder("POST", mockSubscriptionURL,
		httpmock.NewStringResponder(200, `{
			"version": "1.0.0",
			"status": "SUCCESS"
		}`),
	)

	// Atualizar o endpoint da sessão para apontar para o mock server
	s.session.APIEndpoint = s.mockServer.URL
	s.webhooks = celcoin.NewWebhooks(s.httpClient, *s.session)
}

// TearDownTest finaliza os recursos criados para os testes.
func (s *WebhookTestSuite) TearDownTest() {
	if s.mockServer != nil {
		s.mockServer.Close()
	}
	httpmock.DeactivateAndReset()
}

// TestCreateSubscription verifica o comportamento do método CreateSubscription.
func (s *WebhookTestSuite) TestCreateSubscription() {
	// Mock do token de autenticação
	token, err := s.authentication.Token(s.ctx)
	s.assert.NoError(err, "erro ao obter token de autenticação")
	s.assert.Equal("Bearer mock-access-token", token, "esperava o token mockado")

	// Dados de entrada para o teste
	request := celcoin.WebhookSubscriptionRequest{
		Entity:     "pix-payment-out",
		WebhookURL: "http://example.com/webhook",
		Auth: celcoin.WebhookAuth{
			Login:    "test-user",
			Password: "test-pass",
			Type:     "basic",
		},
	}

	// Chamada ao método CreateSubscription
	response, err := s.webhooks.CreateSubscription(s.ctx, request)

	// Validações do teste
	s.assert.NoError(err, "esperava nenhum erro em CreateSubscription")
	s.assert.NotNil(response, "esperava uma resposta não nula")
	s.assert.Equal("1.0.0", response.Version, "esperava versão '1.0.0'")
	s.assert.Equal("SUCCESS", response.Status, "esperava status 'SUCCESS'")
}

// TestCreateSubscriptionInvalidRequest verifica erro em requisições inválidas.
func (s *WebhookTestSuite) TestCreateSubscriptionInvalidRequest() {
	// Mock do token de autenticação
	token, err := s.authentication.Token(s.ctx)
	s.assert.NoError(err, "erro ao obter token de autenticação")
	s.assert.Equal("Bearer mock-access-token", token, "esperava o token mockado")

	// Requisição com valores inválidos
	request := celcoin.WebhookSubscriptionRequest{
		Auth: celcoin.WebhookAuth{
			Login:    "test-user",
			Password: "test-pass",
			Type:     "basic",
		},
	}

	// Chamada ao método CreateSubscription
	response, err := s.webhooks.CreateSubscription(s.ctx, request)

	// Validações do teste
	s.assert.Error(err, "validation failed for entity")
	s.assert.Nil(response, "esperava uma resposta nula em caso de erro")
}
