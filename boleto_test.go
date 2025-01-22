package celcoin_test

import (
	"bytes"
	"context"
	"encoding/json"
	//"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/contbank/celcoin-sdk" 
)

// MockRoundTripper implements http.RoundTripper (and thus can be used by http.Client).
type MockRoundTripper struct {
	mock.Mock
}

// RoundTrip mocks the actual HTTP round trip.
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	// Return the *http.Response and error from .On(...).Return(...)
	return args.Get(0).(*http.Response), args.Error(1)
}

// BoletoTestSuite ...
type BoletoTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	ctx        context.Context
	mockTransport *MockRoundTripper
	client     *http.Client
	boleto     *celcoin.Boleto
}

// TestBoletoTestSuite 
func TestBoletoTestSuite(t *testing.T) {
	suite.Run(t, new(BoletoTestSuite))
}

// SetupTest sets up each test with a new mock.
func (s *BoletoTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	// Create our mock ...
	s.mockTransport = new(MockRoundTripper)
	s.client = &http.Client{Transport: s.mockTransport}

	// Create the Boleto SDK instance ..
	baseURL := "https://sandbox.openfinance.celcoin.dev/baas/v2"
	apiKey := "fake-token"
	s.boleto = &celcoin.Boleto{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  s.client,
	}
}

// TestCreateQueryDownloadCancel ...
func (s *BoletoTestSuite) TestCreateQueryDownloadCancel() {
	// 1) CREATE mock
	createReq := celcoin.CreateBoletoRequest{
		ExternalID:            "TesteSandbox_123",
		ExpirationAfterPayment: 1,
		DueDate:               "2025-01-30",
		Amount:                10,
		Debtor: celcoin.Debtor{
			Number:       "123",
			Neighborhood: "Alphaville Residencial Um",
			Name:         "Erick Augusto Farias",
			Document:     "42318970858",
			City:         "Barueri",
			PublicArea:   "Alameda Holanda",
			State:        "SP",
			PostalCode:   "06474320",
		},
		Receiver: celcoin.Receiver{
			Account:  "30054999518",
			Document: "37786401865",
		},
		Instructions: celcoin.Instructions{
			Fine:     10,
			Interest: 5,
			Discount: celcoin.Discount{
				Amount:    1,
				Modality:  "fixed",
				LimitDate: "2025-01-20T00:00:00.0000000",
			},
		},
	}

	createResp := struct {
		Body struct {
			TransactionID string `json:"transactionId"`
		} `json:"body"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}{
		Body: struct {
			TransactionID string `json:"transactionId"`
		}{
			TransactionID: "9e4f8148-03cb-430a-aec9-558e83e17352",
		},
		Version: "1.2.0",
		Status:  "SUCCESS",
	}

	createRespBody, _ := json.Marshal(createResp)
	// Build the *http.Response to return from mock
	mockCreateResponse := &http.Response{
		StatusCode: 201, // typical success
		Body:       ioutil.NopCloser(bytes.NewReader(createRespBody)),
	}

	// Expect that we do exactly 1 POST call for creation
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockCreateResponse, nil).
		Once()

	// Execute the Create call
	result, err := s.boleto.Create(s.ctx, createReq)
	s.assert.NoError(err)
	s.assert.Equal("9e4f8148-03cb-430a-aec9-558e83e17352", result.TransactionID)
	s.assert.Equal("SUCCESS", result.Status)

	// 2) QUERY mock
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
			TransactionID: "9e4f8148-03cb-430a-aec9-558e83e17352",
			Status:        "PENDING",
		},
		Version: "1.2.0",
		Status:  "SUCCESS",
	}

	queryRespBody, _ := json.Marshal(queryRespJSON)
	mockQueryResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(queryRespBody)),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockQueryResponse, nil).
		Once()

	qResult, err := s.boleto.Query(s.ctx, result.TransactionID)
	s.assert.NoError(err)
	s.assert.Equal(result.TransactionID, qResult.TransactionID)
	s.assert.Equal("PENDING", qResult.Status)

	// 3) DOWNLOAD PDF mock
	mockPDFBody := []byte("%PDF-1.7 \n ...Fake PDF data...")
	mockDownloadResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(mockPDFBody)),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockDownloadResp, nil).
		Once()

	pdfData, err := s.boleto.DownloadPDF(s.ctx, result.TransactionID)
	s.assert.NoError(err)
	s.assert.True(len(pdfData) > 0)

	// 4) CANCEL mock
	cancelRespJSON := struct {
		Body struct {
			TransactionID string `json:"transactionId"`
		} `json:"body"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}{
		Body: struct {
			TransactionID string `json:"transactionId"`
		}{
			TransactionID: "9e4f8148-03cb-430a-aec9-558e83e17352",
		},
		Version: "1.2.0",
		Status:  "PROCESSING",
	}
	cancelRespBody, _ := json.Marshal(cancelRespJSON)

	mockCancelResponse := &http.Response{
		StatusCode: 200, // success on DELETE
		Body:       ioutil.NopCloser(bytes.NewReader(cancelRespBody)),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockCancelResponse, nil).
		Once()

	err = s.boleto.Cancel(s.ctx, result.TransactionID, "Cancelamento do contrato com o cliente.")
	s.assert.NoError(err)
}

