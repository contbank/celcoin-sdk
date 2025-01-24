package celcoin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"


	"github.com/contbank/celcoin-sdk"
)

// MockRoundTripper implements http.RoundTripper
type MockRoundTripper struct {
	mock.Mock
}

// RoundTrip is the main mocked method: returns *http.Response and error.
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// BoletoTestSuite organizes our test for Celcoin Boletos.
type BoletoTestSuite struct {
	suite.Suite
	assert        *assert.Assertions
	ctx           context.Context
	mockTransport *MockRoundTripper
	client        *http.Client

	// The object under test (celcoin.Boletos):
	boleto *celcoin.Boletos
}

func TestBoletoTestSuite(t *testing.T) {
	suite.Run(t, new(BoletoTestSuite))
}

// SetupTest ...
func (s *BoletoTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	s.mockTransport = new(MockRoundTripper)
	s.client = &http.Client{Transport: s.mockTransport}

	// Create a Session...
	session := celcoin.Session{
		APIEndpoint: "https://sandbox.openfinance.celcoin.dev/baas/v2",
		// ... other fields if your Session struct has them
	}

	// Constructor NewBoletos...
	s.boleto = celcoin.NewBoletos(s.client, session)
}

// TestCreateQueryDownloadCancel mocks out all 4 operations.
func (s *BoletoTestSuite) TestCreateQueryDownloadCancel() {
	// 1) CREATE
	createReq := celcoin.CreateBoletoRequest{
		ExternalID:             "TesteSandbox_123",
		ExpirationAfterPayment: 1,
		DueDate:                "2025-01-30",
		Amount:                 10,
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

	// Mock the successful JSON body from Celcoin's response.
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
	createRespBytes, _ := json.Marshal(createResp)
	mockCreateHTTPResp := &http.Response{
		StatusCode: 201,
		Body:       ioutil.NopCloser(bytes.NewReader(createRespBytes)),
	}

	// Expect first RoundTrip call (Create) to return mockCreateHTTPResp
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockCreateHTTPResp, nil).
		Once()

	// Execute the Create call
	result, err := s.boleto.Create(s.ctx, createReq)
	s.assert.NoError(err, "Create should not return error")
	s.assert.Equal("9e4f8148-03cb-430a-aec9-558e83e17352", result.TransactionID)
	s.assert.Equal("SUCCESS", result.Status)

	// 2) QUERY
	queryResp := struct {
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
	queryRespBytes, _ := json.Marshal(queryResp)
	mockQueryHTTPResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(queryRespBytes)),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockQueryHTTPResp, nil).
		Once()

	qResult, err := s.boleto.Query(s.ctx, result.TransactionID)
	s.assert.NoError(err)
	s.assert.Equal("9e4f8148-03cb-430a-aec9-558e83e17352", qResult.TransactionID)
	s.assert.Equal("PENDING", qResult.Status)

	// 3) DOWNLOAD PDF
	mockPDFData := []byte("%PDF-1.7 \n ...Fake PDF data...")
	mockDownloadHTTPResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(mockPDFData)),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockDownloadHTTPResp, nil).
		Once()

	pdfBytes, err := s.boleto.DownloadPDF(s.ctx, qResult.TransactionID)
	s.assert.NoError(err)
	s.assert.True(len(pdfBytes) > 0, "Should download non-empty PDF")

	// 4) CANCEL
	cancelResp := struct {
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
	cancelRespBytes, _ := json.Marshal(cancelResp)
	mockCancelHTTPResp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(cancelRespBytes)),
	}
	s.mockTransport.
		On("RoundTrip", mock.Anything).
		Return(mockCancelHTTPResp, nil).
		Once()

	err = s.boleto.Cancel(s.ctx, qResult.TransactionID, "Cancelamento do contrato com o cliente.")
	s.assert.NoError(err)
}
