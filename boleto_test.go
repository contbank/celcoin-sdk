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

// MockRoundTripper
type MockRoundTripper struct {
    mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    args := m.Called(req)
    if args.Get(0) == nil {
        
        return nil, args.Error(1)
    }
    return args.Get(0).(*http.Response), args.Error(1)
}

type BoletoTestSuite struct {
    suite.Suite
    assert        *assert.Assertions
    ctx           context.Context
    mockTransport *MockRoundTripper
    client        *http.Client

    session  *celcoin.Session
    boletos  celcoin.Boletos // Interface in boletos.go
}

func TestBoletoTestSuite(t *testing.T) {
    suite.Run(t, new(BoletoTestSuite))
}

func (s *BoletoTestSuite) SetupTest() {
    s.assert = assert.New(s.T())
    s.ctx = context.Background()

    s.mockTransport = new(MockRoundTripper)
    s.client = &http.Client{Transport: s.mockTransport}

    // Create a mock session config.
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

    // Create the Boletos service with mocked http.Client + session
    s.boletos = celcoin.NewBoletos(s.client, *s.session)
}

// TestCreateQueryDownloadCancel
func (s *BoletoTestSuite) TestCreateQueryDownloadCancel() {
    // 1) CREATE
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

    // Example Celcoin-ish JSON response for a successful creation:
    createRespJSON := struct {
        Body struct {
            TransactionID string `json:"transactionId"`
        } `json:"body"`
        Version string `json:"version"`
        Status  string `json:"status"`
    }{
        Body: struct {
            TransactionID string `json:"transactionId"`
        }{
            TransactionID: "5a0f8148-03cb-430a-aec9-558e83e17352",
        },
        Version: "1.2.0",
        Status:  "SUCCESS",
    }
    createBytes, _ := json.Marshal(createRespJSON)
    mockCreateResp := &http.Response{
        StatusCode: 201,
        Body:       ioutil.NopCloser(bytes.NewReader(createBytes)),
    }

    // Mock the first RoundTrip call
    s.mockTransport.On("RoundTrip", mock.Anything).Return(mockCreateResp, nil).Once()

    // 1) Actually call s.boletos.Create
    createResult, err := s.boletos.Create(s.ctx, createReq)
    s.assert.NoError(err)
    s.assert.Equal("5a0f8148-03cb-430a-aec9-558e83e17352", createResult.TransactionID)
    s.assert.Equal("SUCCESS", createResult.Status)

    // 2) QUERY
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
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader(queryBytes)),
    }
    s.mockTransport.On("RoundTrip", mock.Anything).Return(mockQueryResp, nil).Once()

    queryResult, err := s.boletos.Query(s.ctx, createResult.TransactionID)
    s.assert.NoError(err)
    s.assert.Equal(createResult.TransactionID, queryResult.TransactionID)
    s.assert.Equal("PENDING", queryResult.Status)

    // 3) DOWNLOAD PDF
    fakePDFData := []byte("%PDF-1.7 \n ...Fake PDF Data...")
    mockPDFResp := &http.Response{
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader(fakePDFData)),
    }
    s.mockTransport.On("RoundTrip", mock.Anything).Return(mockPDFResp, nil).Once()

    pdfData, err := s.boletos.DownloadPDF(s.ctx, queryResult.TransactionID)
    s.assert.NoError(err)
    s.assert.NotEmpty(pdfData)

    // 4) CANCEL
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
            TransactionID: createResult.TransactionID,
        },
        Version: "1.2.0",
        Status:  "PROCESSING",
    }
    cancelBytes, _ := json.Marshal(cancelRespJSON)
    mockCancelResp := &http.Response{
        StatusCode: 200,
        Body:       ioutil.NopCloser(bytes.NewReader(cancelBytes)),
    }
    s.mockTransport.On("RoundTrip", mock.Anything).Return(mockCancelResp, nil).Once()

    err = s.boletos.Cancel(s.ctx, queryResult.TransactionID, "Cliente desistiu do contrato.")
    s.assert.NoError(err)

}
