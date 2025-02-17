package celcoin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/contbank/celcoin-sdk"
)

type PixsTestSuite struct {
	suite.Suite
	assert        *assert.Assertions
	ctx           context.Context
	session       *celcoin.Session
	pixService    celcoin.Pixs
	mockTransport *MockRoundTripper
	client        *http.Client
}

func TestPixsTestSuite(t *testing.T) {
	suite.Run(t, new(PixsTestSuite))
}

func (s *PixsTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()
	s.mockTransport = new(MockRoundTripper)
	s.client = &http.Client{Transport: s.mockTransport}

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
	s.session = session

	s.pixService = celcoin.NewPixs(s.client, *s.session)
}

// TestCreatePixKey testa o método de criação de chave pix.
func (s *PixsTestSuite) TestCreatePixKey() {
	request := celcoin.PixKeyRequest{
		Key:     "test-key",
		KeyType: "CPF",
		Account: "123456",
	}

	fixedTime := time.Date(2025, 1, 17, 19, 43, 52, 0, time.UTC)

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

	mockResponseBody, err := json.Marshal(expectedResponse)
	s.assert.NoError(err)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockResponse, nil).
		Once()

	response, err := s.pixService.CreatePixKey(s.ctx, request)

	s.assert.NoError(err, "Erro inesperado na criação da chave PIX")
	s.assert.NotNil(response, "A resposta não deve ser nula")
	s.assert.Equal("CPF", response.Body.KeyType, "Tipo de chave incorreto")
	s.assert.Equal("test-key", response.Body.Key, "Chave incorreta")
	s.assert.WithinDuration(fixedTime, response.Body.Account.CreateDate, time.Second, "Data de criação incorreta")
	s.mockTransport.AssertExpectations(s.T())
}

// TestGetPixKeys testa o método de pesquisa de chave pix
func (s *PixsTestSuite) TestGetPixKeys() {
	// Cenário: Resposta bem-sucedida
	account := "123456"
	expectedResponse := &celcoin.PixKeyListResponse{
		Version: "1.0",
		Status:  "SUCCESS",
		Body: celcoin.PixKeyListResponseBody{
			ListKeys: []celcoin.PixKeyListItem{
				{
					KeyType: "EMAIL",
					Key:     "test@example.com",
					Account: celcoin.PixKeyAccount{
						Participant: "12345678",
						Branch:      "0001",
						Account:     "123456",
						AccountType: "CONTA_CORRENTE",
					},
					Owner: celcoin.PixKeyOwner{
						Type:           "NATURAL_PERSON",
						DocumentNumber: "12345678901",
						Name:           "Test Owner",
					},
				},
			},
		},
	}

	mockResponseBody, err := json.Marshal(expectedResponse)
	s.assert.NoError(err)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockResponse, nil).
		Once()

	response, err := s.pixService.GetPixKeys(s.ctx, account)

	s.assert.NoError(err, "Erro inesperado ao obter as chaves Pix")
	s.assert.NotNil(response, "A resposta não deve ser nula")
	s.assert.Equal("1.0", response.Version, "Versão incorreta")
	s.assert.Equal("SUCCESS", response.Status, "Status incorreto")
	s.assert.Len(response.Body.ListKeys, 1, "Número incorreto de chaves")
	s.assert.Equal("EMAIL", response.Body.ListKeys[0].KeyType, "Tipo de chave incorreto")
	s.assert.Equal("test@example.com", response.Body.ListKeys[0].Key, "Chave incorreta")

	s.mockTransport.ExpectedCalls = nil

	s.pixService = celcoin.NewPixs(s.client, *s.session) // Resetando o serviço
	mockHTTPError := fmt.Errorf("error making HTTP request")
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(nil, mockHTTPError).
		Once()

	response, err = s.pixService.GetPixKeys(s.ctx, account)
	s.assert.Error(err, "Erro esperado na requisição HTTP")
	s.assert.Nil(response, "A resposta deve ser nula")

	// Cenário: Status HTTP 404
	s.mockTransport.ExpectedCalls = nil
	mockNotFoundResponse := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{}`))),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockNotFoundResponse, nil).
		Once()

	response, err = s.pixService.GetPixKeys(s.ctx, account)
	s.assert.ErrorIs(err, celcoin.ErrEntryNotFound, "Erro esperado: entrada não encontrada")
	s.assert.Nil(response, "A resposta deve ser nula")

	// Cenário: Erro genérico de resposta
	s.mockTransport.ExpectedCalls = nil
	mockErrorResponseBody := `{"error": {"errorCode": "GENERIC_ERROR"}}`
	mockGenericErrorResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(mockErrorResponseBody))),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockGenericErrorResponse, nil).
		Once()

	response, err = s.pixService.GetPixKeys(s.ctx, account)
	s.assert.Error(err, "Erro esperado para resposta genérica")
	s.assert.Nil(response, "A resposta deve ser nula")
}

func (s *PixsTestSuite) TestPaymentPixCashOut() {
	// Cenário: Resposta bem-sucedida para QR Code Dinâmico
	transactionIdentification := "dc8cf02b81b54bd59323453b207e704a"
	dynamicQRCodeRequest := &celcoin.PixCashOutRequest{
		Amount:                    25.55,
		ClientCode:                "1458854",
		TransactionIdentification: transactionIdentification,
		EndToEndId:                "E3030629420200808185300887639654",
		InitiationType:            "DYNAMIC_QRCODE",
		PaymentType:               "IMMEDIATE",
		Urgency:                   "HIGH",
		TransactionType:           "TRANSFER",
		DebitParty: celcoin.DebitParty{
			Account: "444444",
		},
		CreditParty: celcoin.CreditParty{
			Bank:        "30306294",
			Key:         "5244f4e-15ff-413d-808d-7837652ebdc2",
			Account:     "10545584",
			Branch:      "10545584",
			TaxId:       "11122233344",
			Name:        "Celcoin",
			AccountType: "CACC",
		},
		RemittanceInformation: "Texto de mensagem",
	}

	expectedDynamicQRCodeResponse := &celcoin.PixCashOutResponse{
		Status:  "SUCCESS",
		Version: "1.0",
		Body: celcoin.PixCashOutResponseBody{
			ID:                        "123",
			ClientCode:                "1458854",
			EndToEndID:                "E3030629420200808185300887639654",
			PaymentType:               "IMMEDIATE",
			Urgency:                   "HIGH",
			TransactionType:           "TRANSFER",
			Amount:                    25.55,
			TransactionIdentification: &transactionIdentification,
			InitiationType:            "DYNAMIC_QRCODE",
			DebitParty: celcoin.DebitParty{
				Account: "444444",
			},
			CreditParty: celcoin.CreditParty{
				Bank:        "30306294",
				Key:         "5244f4e-15ff-413d-808d-7837652ebdc2",
				Account:     "10545584",
				Branch:      "10545584",
				TaxId:       "11122233344",
				Name:        "Celcoin",
				AccountType: "CACC",
			},
			RemittanceInformation: "Texto de mensagem",
		},
	}

	mockResponseBody, err := json.Marshal(expectedDynamicQRCodeResponse)
	s.Require().NoError(err, "Erro ao serializar resposta simulada")

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockResponse, nil).
		Once()

	response, err := s.pixService.PaymentPixCashOut(s.ctx, *dynamicQRCodeRequest)

	s.Assert().NoError(err, "Erro inesperado ao processar PixCashOut por QR Code Dinâmico")
	s.Assert().NotNil(response, "A resposta não deve ser nula")
	s.Assert().Equal("SUCCESS", response.Status, "O status esperado era SUCCESS")
	s.Assert().Equal(dynamicQRCodeRequest.Amount, response.Body.Amount, "Valor da transação incorreto")

	// Cenário: Resposta bem-sucedida para QR Code Estático
	staticQRCodeRequest := &celcoin.PixCashOutRequest{
		Amount:                    25.55,
		ClientCode:                "1458854",
		TransactionIdentification: "dc8cf02b81b54bd59323453b",
		EndToEndId:                "E3030629420200808185300887639654",
		InitiationType:            "STATIC_QRCODE",
		PaymentType:               "IMMEDIATE",
		Urgency:                   "HIGH",
		TransactionType:           "TRANSFER",
		DebitParty: celcoin.DebitParty{
			Account: "444444",
		},
		CreditParty: celcoin.CreditParty{
			Bank:        "30306294",
			Key:         "5244f4e-15ff-413d-808d-7837652ebdc2",
			Account:     "10545584",
			Branch:      "10545584",
			TaxId:       "11122233344",
			Name:        "Celcoin",
			AccountType: "CACC",
		},
		RemittanceInformation: "Texto de mensagem",
	}

	expectedStaticQRCodeResponse := &celcoin.PixCashOutResponse{
		Status:  "SUCCESS",
		Version: "1.0",
		Body: celcoin.PixCashOutResponseBody{
			ID:                        "124",
			ClientCode:                staticQRCodeRequest.ClientCode,
			TransactionIdentification: &staticQRCodeRequest.TransactionIdentification,
			EndToEndID:                staticQRCodeRequest.EndToEndId,
			PaymentType:               staticQRCodeRequest.PaymentType,
			Urgency:                   staticQRCodeRequest.Urgency,
			TransactionType:           staticQRCodeRequest.TransactionType,
			Amount:                    staticQRCodeRequest.Amount,
			InitiationType:            staticQRCodeRequest.InitiationType,
			DebitParty:                staticQRCodeRequest.DebitParty,
			CreditParty:               staticQRCodeRequest.CreditParty,
			RemittanceInformation:     staticQRCodeRequest.RemittanceInformation,
		},
	}

	mockResponseBody, err = json.Marshal(expectedStaticQRCodeResponse)
	s.Require().NoError(err, "Erro ao serializar resposta simulada")

	mockResponse = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResponse, nil).Once()

	response, err = s.pixService.PaymentPixCashOut(s.ctx, *staticQRCodeRequest)

	s.Assert().NoError(err, "Erro inesperado ao processar PixCashOut por QR Code Estático")
	s.Assert().NotNil(response, "A resposta não deve ser nula")
	s.Assert().Equal("SUCCESS", response.Status, "O status esperado era SUCCESS")
	s.Assert().Equal(staticQRCodeRequest.Amount, response.Body.Amount, "Valor da transação incorreto")

	// Cenário: Resposta bem-sucedida para Chaves Pix
	pixKeysRequest := &celcoin.PixCashOutRequest{
		Amount:          25.55,
		ClientCode:      "1458854",
		EndToEndId:      "E3030629420200808185300887639654",
		InitiationType:  "DICT",
		PaymentType:     "IMMEDIATE",
		Urgency:         "HIGH",
		TransactionType: "TRANSFER",
		DebitParty: celcoin.DebitParty{
			Account: "444444",
		},
		CreditParty: celcoin.CreditParty{
			Bank:        "30306294",
			Key:         "5244f4e-15ff-413d-808d-7837652ebdc2",
			Name:        "Celcoin",
			AccountType: "CACC",
		},
		RemittanceInformation: "Texto de mensagem",
	}

	expectedPixKeysResponse := &celcoin.PixCashOutResponse{
		Status:  "SUCCESS",
		Version: "1.0",
		Body: celcoin.PixCashOutResponseBody{
			ID:                        "txn-12345",
			ClientCode:                pixKeysRequest.ClientCode,
			TransactionIdentification: &pixKeysRequest.TransactionIdentification,
			EndToEndID:                pixKeysRequest.EndToEndId,
			InitiationType:            "DICT",
			PaymentType:               "IMMEDIATE",
			Urgency:                   "HIGH",
			TransactionType:           "TRANSFER",
			DebitParty:                pixKeysRequest.DebitParty,
			CreditParty:               pixKeysRequest.CreditParty,
			RemittanceInformation:     pixKeysRequest.RemittanceInformation,
			Amount:                    pixKeysRequest.Amount,
		},
	}

	mockResponseBody, err = json.Marshal(expectedPixKeysResponse)
	s.Require().NoError(err, "Erro ao serializar resposta simulada")

	mockResponse = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResponse, nil).Once()

	response, err = s.pixService.PaymentPixCashOut(s.ctx, *pixKeysRequest)

	s.Assert().NoError(err, "Erro inesperado ao processar PixCashOut por Chaves Pix")
	s.Assert().NotNil(response, "A resposta não deve ser nula")
	s.Assert().Equal("SUCCESS", response.Status, "O status esperado era SUCCESS")
	s.Assert().Equal(pixKeysRequest.Amount, response.Body.Amount, "Valor da transação incorreto")

	// Cenário: Resposta bem-sucedida para dados bancários
	bankDetailsRequest := &celcoin.PixCashOutRequest{
		Amount:          25.55,
		ClientCode:      "1458854",
		InitiationType:  "MANUAL",
		PaymentType:     "IMMEDIATE",
		Urgency:         "HIGH",
		TransactionType: "TRANSFER",
		DebitParty: celcoin.DebitParty{
			Account: "444444",
		},
		CreditParty: celcoin.CreditParty{
			Bank:        "30306294",
			Account:     "10545584",
			Branch:      "10545584",
			TaxId:       "11122233344",
			Name:        "Celcoin",
			AccountType: "CACC",
		},
		RemittanceInformation: "Texto de mensagem",
	}

	expectedBankDetailsResponse := &celcoin.PixCashOutResponse{
		Status:  "SUCCESS",
		Version: "1.0",
		Body: celcoin.PixCashOutResponseBody{
			InitiationType:        bankDetailsRequest.InitiationType,
			Amount:                bankDetailsRequest.Amount,
			DebitParty:            bankDetailsRequest.DebitParty,
			CreditParty:           bankDetailsRequest.CreditParty,
			RemittanceInformation: bankDetailsRequest.RemittanceInformation,
		},
	}

	mockResponseBody, err = json.Marshal(expectedBankDetailsResponse)
	s.Require().NoError(err, "Erro ao serializar resposta simulada")

	mockResponse = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResponse, nil).Once()

	response, err = s.pixService.PaymentPixCashOut(s.ctx, *bankDetailsRequest)

	s.Assert().NoError(err, "Erro inesperado ao processar PixCashOut por dados bancários")
	s.Assert().NotNil(response, "A resposta não deve ser nula")
	s.Assert().Equal("SUCCESS", response.Status, "O status esperado era SUCCESS")
	s.Assert().Equal(bankDetailsRequest.Amount, response.Body.Amount, "Valor da transação incorreto")
}

func (s *PixsTestSuite) TestGetPixCashoutStatus() {
	// Cenário: Resposta bem-sucedida
	id := "60ec4471-71dd-43a3-a848-efe7a314d76f"
	endtoendId := "E1393589320221110144001306556986"
	clientCode := "1458856889"
	expectedResponse := celcoin.PixCashoutStatusTransactionResponse{
		Status:  "CONFIRMED",
		Version: "1.0.0",
		Body: celcoin.PixCashoutStatusTransactionBody{
			ID:                        id,
			Amount:                    50,
			ClientCode:                clientCode,
			TransactionIdentification: nil,
			EndToEndID:                endtoendId,
			InitiationType:            "MANUAL",
			PaymentType:               "IMMEDIATE",
			Urgency:                   "HIGH",
			TransactionType:           "TRANSFER",
			DebitParty: celcoin.DebitParty{
				Account:     "30053913714179",
				Branch:      "0001",
				TaxId:       "77859635097",
				Name:        "Hernani  Conrado",
				AccountType: "TRAN",
			},
			CreditParty: celcoin.CreditParty{
				Bank:        "30306294",
				Account:     "42161",
				Branch:      "20",
				TaxId:       "12312312300",
				Name:        "Fulano de Tal",
				AccountType: "CACC",
			},
			RemittanceInformation: "Texto de mensagem",
			Error:                 nil,
		},
	}

	mockResponseBody, err := json.Marshal(expectedResponse)
	s.Require().NoError(err, "Erro ao serializar resposta simulada")

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResponse, nil).Once()

	response, err := s.pixService.GetPixCashoutStatus(s.ctx, id, endtoendId, clientCode)

	s.Assert().NoError(err, "Erro inesperado ao obter status do Pix Cashout")
	s.Assert().NotNil(response, "A resposta não deve ser nula")
	s.Assert().Equal("CONFIRMED", response.Status, "O status esperado era CONFIRMED")
	s.Assert().Equal("MANUAL", response.Body.InitiationType, "Tipo de iniciação incorreto")
	s.Assert().Equal(id, response.Body.ID, "ID da transação incorreto")
	s.Assert().Equal(50.0, response.Body.Amount, "Valor da transação incorreto")

	// Cenário: Erro na requisição HTTP
	s.mockTransport.ExpectedCalls = nil
	mockHTTPError := fmt.Errorf("error making HTTP request")
	s.mockTransport.On("RoundTrip", mock.Anything).Return(nil, mockHTTPError).Once()

	response, err = s.pixService.GetPixCashoutStatus(s.ctx, id, endtoendId, clientCode)
	s.Assert().Error(err, "Erro esperado para falha na requisição HTTP")
	s.Assert().Nil(response, "A resposta deve ser nula em caso de erro HTTP")

	// Cenário: Status HTTP 404
	s.mockTransport.ExpectedCalls = nil
	mockNotFoundResponse := &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{}`))),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockNotFoundResponse, nil).Once()

	response, err = s.pixService.GetPixCashoutStatus(s.ctx, id, endtoendId, clientCode)
	s.Assert().ErrorIs(err, celcoin.ErrEntryNotFound, "Erro esperado: entrada não encontrada")
	s.Assert().Nil(response, "A resposta deve ser nula")

	// Cenário: Erro genérico de resposta
	s.mockTransport.ExpectedCalls = nil
	mockErrorResponseBody := `{"error": {"errorCode": "GENERIC_ERROR"}}`
	mockGenericErrorResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(mockErrorResponseBody))),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockGenericErrorResponse, nil).Once()

	response, err = s.pixService.GetPixCashoutStatus(s.ctx, id, endtoendId, clientCode)
	s.Assert().Error(err, "Erro esperado para resposta genérica")
	s.Assert().Nil(response, "A resposta deve ser nula")

	// Cenário: Parâmetros inválidos
	response, err = s.pixService.GetPixCashoutStatus(s.ctx, "", "", "")
	s.Assert().Error(err, "Erro esperado para parâmetros inválidos")
	s.Assert().Nil(response, "A resposta deve ser nula para parâmetros inválidos")
}
