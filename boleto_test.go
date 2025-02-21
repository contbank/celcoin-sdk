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

// MockRoundTripper implements http.RoundTripper for testing.
type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*http.Response), args.Error(1)
}

type BoletoTestSuite struct {
	suite.Suite
	assert        *assert.Assertions
	ctx           context.Context
	mockTransport *MockRoundTripper
	client        *http.Client

	session *celcoin.Session
	boletos *celcoin.Boletos
}

func TestBoletoTestSuite(t *testing.T) {
	suite.Run(t, new(BoletoTestSuite))
}

func (s *BoletoTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	s.mockTransport = new(MockRoundTripper)
	s.client = &http.Client{
		Transport: s.mockTransport,
		Timeout:   120 * time.Second,
	}

	// Create a mock session configuration.
	clientID := "test-client-id"
	clientSecret := "test-client-secret"
	apiEndpoint := "https://sandbox.openfinance.celcoin.dev"
	loginEndpoint := "https://sandbox.openfinance.celcoin.dev"
	config := celcoin.Config{
		ClientID:      &clientID,
		ClientSecret:  &clientSecret,
		Mtls:          celcoin.Bool(false),
		APIEndpoint:   &apiEndpoint,
		LoginEndpoint: &loginEndpoint,
	}
	sess, err := celcoin.NewSession(config)
	s.assert.NoError(err)
	s.session = sess

	// Instantiate the Boletos service.
	s.boletos = celcoin.NewBoletos(s.client, *s.session)
}

func (s *BoletoTestSuite) TestCreateQueryDownloadCancel() {
	// --- 1) CREATE BOLETO ---
	createReq := celcoin.CreateBoletoRequest{
		ExternalID:             "TestBoleto123",
		ExpirationAfterPayment: 1,
		DueDate:                "2025-01-30",
		Amount:                 100.0,
		Debtor: celcoin.Debtor{
			Number:       "123",
			Neighborhood: "Alphaville",
			Name:         "João da Silva",
			Document:     "42318970858",
			City:         "Barueri",
			PublicArea:   "Alameda Teste",
			State:        "SP",
			PostalCode:   "06474320",
		},
		Receiver: celcoin.Receiver{
			Account:  "123456",
			Document: "98765432100",
		},
		Instructions: celcoin.Instructions{
			Fine:     10,
			Interest: 5,
			Discount: celcoin.Discount{
				Amount:    2,
				Modality:  "fixed",
				LimitDate: "2025-01-20T00:00:00.0000000",
			},
		},
	}

	// Simulated response for a successful boleto creation:
	createRespJSON := struct {
		Body    celcoin.CreateBoletoResponse `json:"body"`
		Version string                       `json:"version"`
		Status  string                       `json:"status"`
	}{
		Body: celcoin.CreateBoletoResponse{
			TransactionID: "5a0f8148-03cb-430a-aec9-558e83e17352",
			Status:        "SUCCESS",
		},
		Version: "1.2.0",
		Status:  "SUCCESS",
	}
	createBytes, _ := json.Marshal(createRespJSON)
	mockCreateResp := &http.Response{
		StatusCode: http.StatusCreated,
		Body:       ioutil.NopCloser(bytes.NewReader(createBytes)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockCreateResp, nil).Once()

	// Call CreateBoleto
	createResult, err := s.boletos.CreateBoleto(s.ctx, createReq)
	s.assert.NoError(err)
	s.assert.Equal("5a0f8148-03cb-430a-aec9-558e83e17352", createResult.TransactionID)
	s.assert.Equal("SUCCESS", createResult.Status)

	// --- 2) QUERY BOLETO ---
	queryRespJSON := struct {
		Body struct {
			TransactionID string `json:"transactionId"`
			Status        string `json:"status"`
		} `json:"body"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}{
		Body: struct {
			TransactionID string `json:"transactionId"`
			Status        string `json:"status"`
		}{
			TransactionID: createResult.TransactionID,
			Status:        "PENDING",
		},
		Version: "1.2.0",
		Status:  "SUCCESS",
	}
	queryBytes, _ := json.Marshal(queryRespJSON)
	mockQueryResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(queryBytes)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockQueryResp, nil).Once()

	queryResult, err := s.boletos.QueryBoleto(s.ctx, createResult.TransactionID)
	s.assert.NoError(err)
	s.assert.Equal(createResult.TransactionID, queryResult.TransactionID)
	s.assert.Equal("PENDING", queryResult.Status)

	// --- 3) DOWNLOAD PDF ---
	fakePDFData := []byte("%PDF-1.7 \n ...Fake PDF Data...")
	mockPDFResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(fakePDFData)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockPDFResp, nil).Once()

	// Create a buffer to capture the PDF output.
	var buf bytes.Buffer
	err = s.boletos.DownloadBoletoPDF(s.ctx, queryResult.TransactionID, &buf)
	s.assert.NoError(err)
	pdfData := buf.Bytes()
	s.assert.Equal(fakePDFData, pdfData)

	// --- 4) CANCEL BOLETO ---
	cancelRespJSON := struct {
		Body    celcoin.CreateBoletoResponse `json:"body"`
		Version string                       `json:"version"`
		Status  string                       `json:"status"`
	}{
		Body: celcoin.CreateBoletoResponse{
			TransactionID: createResult.TransactionID,
			Status:        "PROCESSING",
		},
		Version: "1.2.0",
		Status:  "PROCESSING",
	}
	cancelBytes, _ := json.Marshal(cancelRespJSON)
	mockCancelResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(cancelBytes)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockCancelResp, nil).Once()

	err = s.boletos.CancelBoleto(s.ctx, queryResult.TransactionID, "Cliente desistiu do contrato.")
	s.assert.NoError(err)
}
