package celcoin_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/contbank/celcoin-sdk"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MockAuthentication simula um Token válido.
type MockAuthentication struct{}

func (m *MockAuthentication) Token(ctx context.Context) (string, error) {
	// Simulando a resposta do token.
	return "Bearer fake-api-key", nil
}

// WebhookTestSuite é a suite de testes para os webhooks.
type WebhookTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	ctx            context.Context
	session        *celcoin.Session
	authentication *MockAuthentication
	webhooks       *celcoin.Webhooks
	logger         *logrus.Logger
	server         *httptest.Server
}

// TestWebhookTestSuite executa a suite de testes.
func TestWebhookTestSuite(t *testing.T) {
	suite.Run(t, new(WebhookTestSuite))
}

// SetupTest inicializa os mocks e configurações.
func (s *WebhookTestSuite) SetupTest() {
	// Inicializar o assert
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	// Mock das variáveis de ambiente
	celcoin.GetClientID = func() string {
		return "mock_client_id"
	}
	celcoin.GetClientSecret = func() string {
		return "mock_client_secret"
	}

	clientID := celcoin.GetEnvCelcoinClientID()
	clientSecret := celcoin.GetEnvCelcoinClientSecret()

	celcoinConfig := celcoin.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Mtls:         celcoin.Bool(false),
	}

	// Criar sessão Celcoin
	session, err := celcoin.NewSession(celcoinConfig)
	s.assert.NoError(err)

	// Configurar o logger para os testes
	s.logger = logrus.New()
	s.logger.SetFormatter(&logrus.JSONFormatter{})
	s.logger.SetOutput(ioutil.Discard)

	// Iniciar servidor mock
	s.server = celcoin.NewMockServer()

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Configurar a sessão para usar o mock server
	session.APIEndpoint = s.server.URL
	session.LoginEndpoint = s.server.URL

	// Inicializar dependências
	s.session = session
	s.authentication = &MockAuthentication{}
	s.webhooks = celcoin.NewWebhooks(httpClient, *s.session)
}

// TearDownTest finaliza os recursos do teste.
func (s *WebhookTestSuite) TearDownTest() {
	if s.server != nil {
		s.server.Close()
	}
}

// TestCreateSubscriptionSuccess verifica o sucesso do método CreateSubscription.
func (s *WebhookTestSuite) TestCreateSubscriptionSuccess() {
	// Definindo a requisição de teste
	req := celcoin.WebhookSubscriptionRequest{
		Entity:     "pix-payment-out",
		WebhookURL: "http://example.com/webhook",
		Auth: celcoin.WebhookAuth{
			Login:    "test-user",
			Password: "test-pass",
			Type:     "basic",
		},
	}

	// Chamando o método que queremos testar
	resp, err := s.webhooks.CreateSubscription(s.ctx, req)
	s.assert.NoError(err)

	// Verificando o resultado
	s.assert.Equal("1.0.0", resp.Version)
	s.assert.Equal("SUCCESS", resp.Status)
}

// TestCreateSubscriptionError verifica o comportamento do método em caso de erro
func (s *WebhookTestSuite) TestCreateSubscriptionError() {
	// Mockando a resposta de erro da API
	handler := http.NewServeMux()
	handler.HandleFunc("/baas-webhookmanager/v1/webhook/subscription", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(celcoin.WebhookSubscriptionResponse{
			Version: "1.0.0",
			Status:  "ERROR",
			Error: &celcoin.WebhookError{
				ErrorCode: "500",
				Message:   "Internal Server Error",
			},
		})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	// Definindo a requisição de teste
	req := celcoin.WebhookSubscriptionRequest{
		Entity:     "pix-payment-out",
		WebhookURL: "http://example.com/webhook",
	}

	// Chamando o método que queremos testar
	resp, err := s.webhooks.CreateSubscription(context.Background(), req)

	// Verificando o erro
	s.assert.Error(err)
	s.assert.Nil(resp)
	s.assert.EqualError(err, "erro ao cadastrar webhook: status inesperado da API")
}

// TestCreateSubscriptionInvalidPayload verifica o comportamento com um payload inválido
func (s *WebhookTestSuite) TestCreateSubscriptionInvalidPayload() {
	// Chamando o método com um payload inválido
	req := celcoin.WebhookSubscriptionRequest{}
	resp, err := s.webhooks.CreateSubscription(context.Background(), req)

	// Verificando o erro
	s.assert.Nil(resp)
	s.assert.Error(err)
	s.assert.Contains(err.Error(), "erro ao serializar a requisição")
}
