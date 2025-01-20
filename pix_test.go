package celcoin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/celcoin-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockHTTPClient é um mock da interface http.RoundTripper
type MockHTTPClient struct {
	mock.Mock
}

// RoundTrip implementa a interface http.RoundTripper
func (m *MockHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	// Registrar a requisição recebida no mock
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// PixsTestSuite é a suite de testes para PixsService
type PixsTestSuite struct {
	suite.Suite
	assert         *assert.Assertions
	ctx            context.Context
	session        *celcoin.Session
	pixService     celcoin.Pixs
	mockClient     *MockHTTPClient
	authentication *celcoin.Authentication
}

// TestPixsTestSuite inicializa a suite de testes
func TestPixsTestSuite(t *testing.T) {
	suite.Run(t, new(PixsTestSuite))
}

// SetupTest configura o ambiente para os testes
func (s *PixsTestSuite) SetupTest() {
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

	// Criar o cliente HTTP mockado
	s.mockClient = new(MockHTTPClient)
	httpClient := &http.Client{Transport: s.mockClient}

	// Configuração de sessão e autenticação
	s.session = session
	s.authentication = celcoin.NewAuthentication(httpClient, *s.session)

	// Inicializar o serviço Pix
	s.pixService = celcoin.NewPixs(httpClient, *s.session)
}

// TestCreatePixKey testa o método CreatePixKey
func (s *PixsTestSuite) TestCreatePixKey() {
	// Request que será enviado
	request := celcoin.PixKeyRequest{
		Key:     "test-key",
		KeyType: "CPF",
		Account: "123456",
	}

	// Data fixa para evitar diferenças de tempo
	fixedTime := time.Date(2025, 1, 17, 19, 43, 52, 0, time.UTC)

	// Response esperado com base no novo layout
	expectedResponse := &celcoin.PixKeyResponse{
		Body: celcoin.PixKeyResponseBody{
			KeyType: request.KeyType,
			Key:     request.Key,
			Account: celcoin.PixKeyAccount{
				Participant: "12345678",
				Branch:      "0001",
				Account:     request.Account,
				AccountType: "TRAN",
				CreateDate:  fixedTime,
			},
			Owner: celcoin.PixKeyOwner{
				Type:           "NATURAL_PERSON",
				DocumentNumber: "123456",
				Name:           "Test Owner",
			},
		},
		Version: "1.0",
		Status:  "ACTIVE",
	}

	// Serializar o expectedResponse para simular o corpo da resposta do mock
	mockResponseBody, err := json.Marshal(expectedResponse)
	s.assert.NoError(err)

	// Simular a resposta HTTP do cliente mockado
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	// Configurar o mock do cliente HTTP
	s.mockClient.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "POST" && req.URL.Path == "/v5/token"
	})).Return(mockResponse, nil)

	// Chamar o serviço sendo testado
	response, err := s.pixService.CreatePixKey(s.ctx, request)

	// Validar o comportamento e o resultado
	s.assert.NoError(err, "Erro inesperado na criação da chave PIX")
	s.assert.NotNil(response, "A resposta não deve ser nula")
	s.assert.Equal("CPF", response.Body.KeyType, "Tipo de chave incorreto")
	s.assert.Equal("test-key", response.Body.Key, "Chave incorreta")
	s.assert.WithinDuration(fixedTime, response.Body.Account.CreateDate, time.Second, "Data de criação incorreta")
	s.mockClient.AssertExpectations(s.T())
}
