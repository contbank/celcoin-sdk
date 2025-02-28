package celcoin_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"strconv"
	"testing"
	"time"

	"github.com/contbank/celcoin-sdk" // ajuste o caminho conforme sua estrutura
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockRoundTripper simula as chamadas HTTP.
type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if resp, ok := args.Get(0).(*http.Response); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

type PaymentTestSuite struct {
	suite.Suite
	assert        *assert.Assertions
	ctx           context.Context
	session       *celcoin.Session
	payment       *celcoin.Payment
	client        *http.Client
	mockTransport *MockRoundTripper
}

func TestPaymentTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentTestSuite))
}

func (s *PaymentTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	s.mockTransport = new(MockRoundTripper)
	s.client = &http.Client{Transport: s.mockTransport}

	clientID := "test-client-id"
	clientSecret := "test-client-secret"
	apiEndpoint := "https://sandbox.openfinance.celcoin.dev"
	loginEndpoint := "https://sandbox.openfinance.celcoin.dev"

	config := celcoin.Config{
		ClientID:      &clientID,
		ClientSecret:  &clientSecret,
		APIEndpoint:   &apiEndpoint,
		LoginEndpoint: &loginEndpoint,
		APIVersion:    "2.0",
	}
	session, err := celcoin.NewSession(config)
	s.assert.NoError(err)
	s.session = session

	s.payment = celcoin.NewPayment(s.client, *s.session)
}

func (s *PaymentTestSuite) TestValidatePayment() {
	correlationID := "test-correlation-id"
	reqModel := &celcoin.ValidatePaymentRequest{
		Code: "TESTCODE123",
	}
	expectedResp := celcoin.ValidatePaymentResponse{
		ID:             "PAYID123",
		Assignor:       "AssignorName",
		Code:           "TESTCODE123",
		Digitable:      "1234567890",
		Amount:         100.50,
		OriginalAmount: 100.50,
		MinAmount:      100.50,
		MaxAmount:      100.50,
		AllowChangeAmount: false,
		DueDate:        "2025-01-20",
		SettleDate:     "2025-01-21",
		NextSettle:     false,
	}
	respBytes, err := json.Marshal(expectedResp)
	s.assert.NoError(err)

	// Simula a resposta HTTP para ValidatePayment
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBytes)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResp, nil).Once()

	response, err := s.payment.ValidatePayment(s.ctx, correlationID, reqModel)
	s.assert.NoError(err)
	s.assert.Equal(expectedResp.ID, response.ID)
	s.mockTransport.AssertExpectations(s.T())
}

func (s *PaymentTestSuite) TestConfirmPayment() {
	correlationID := "test-correlation-id"
	description := "Test payment"
	reqModel := &celcoin.ConfirmPaymentRequest{
		ID:          "PAYID123",
		Amount:      100.50,
		Description: &description,
		BankBranch:  "0001",
		BankAccount: "123456",
	}
	expectedSettledDate := time.Now().UTC()
	expectedResp := celcoin.ConfirmPaymentResponse{
		AuthenticationCode: "AUTHCODE123",
		SettledDate:        expectedSettledDate,
	}
	respBytes, err := json.Marshal(expectedResp)
	s.assert.NoError(err)

	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBytes)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResp, nil).Once()

	response, err := s.payment.ConfirmPayment(s.ctx, correlationID, reqModel)
	s.assert.NoError(err)
	s.assert.Equal(expectedResp.AuthenticationCode, response.AuthenticationCode)
	s.assert.WithinDuration(expectedSettledDate, response.SettledDate, time.Second, "SettledDate differs")
	s.mockTransport.AssertExpectations(s.T())
}

func (s *PaymentTestSuite) TestFilterPayments() {
	correlationID := "test-correlation-id"
	pageSize := 10
	pageToken := "TOKEN123"
	reqModel := &celcoin.FilterPaymentsRequest{
		BankBranch:  "0001",
		BankAccount: "123456",
		PageSize:    pageSize,
		PageToken:   &pageToken,
	}
	expectedResp := celcoin.FilterPaymentsResponse{
		NextPageToken: "NEXTTOKEN",
		Data: []*celcoin.PaymentResponse{
			{
				AuthenticationCode: "AUTHCODE123",
				Status:             "SUCCESS",
				Amount:             100.50,
			},
		},
	}
	respBytes, err := json.Marshal(expectedResp)
	s.assert.NoError(err)

	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBytes)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResp, nil).Once()

	response, err := s.payment.FilterPayments(s.ctx, correlationID, reqModel)
	s.assert.NoError(err)
	s.assert.Equal(expectedResp.NextPageToken, response.NextPageToken)
	s.assert.Len(response.Data, 1)
	s.mockTransport.AssertExpectations(s.T())
}

func (s *PaymentTestSuite) TestDetailPayment() {
	correlationID := "test-correlation-id"
	reqModel := &celcoin.DetailPaymentRequest{
		BankBranch:         "0001",
		BankAccount:        "123456",
		AuthenticationCode: "AUTHCODE123",
	}
	expectedResp := celcoin.PaymentResponse{
		AuthenticationCode: "AUTHCODE123",
		Status:             "SUCCESS",
		Digitable:          "DGT123",
		BankBranch:         "0001",
		BankAccount:        "123456",
	}
	respBytes, err := json.Marshal(expectedResp)
	s.assert.NoError(err)

	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(respBytes)),
	}
	s.mockTransport.On("RoundTrip", mock.Anything).Return(mockResp, nil).Once()

	response, err := s.payment.DetailPayment(s.ctx, correlationID, reqModel)
	s.assert.NoError(err)
	s.assert.Equal(expectedResp.AuthenticationCode, response.AuthenticationCode)
	s.assert.Equal(expectedResp.Digitable, response.Digitable)
	s.mockTransport.AssertExpectations(s.T())
}
