package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// IncomeReportTestSuite ...
type IncomeReportTestSuite struct {
	suite.Suite
	assert       *assert.Assertions
	ctx          context.Context
	session      *Session
	incomeReport *IncomeReport
	mockClient   *MockHTTPClient
}

// TestCustomersTestSuite ...
func TestIncomeReportTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeReportTestSuite))
}

// SetupTest ...
func (s *IncomeReportTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	session := &Session{APIEndpoint: ApiEndpoint}
	s.session = session

	s.mockClient = new(MockHTTPClient)
	httpClient := &http.Client{Transport: s.mockClient}
	s.incomeReport = NewIncomeReport(httpClient, *s.session)
}

// TestGetIncomeReport ...
func (s *IncomeReportTestSuite) TestGetIncomeReport() {
	documentNumber := "12345678901"
	accountNumber := "123456"
	userName := "TEST USER"
	calendarYear := strconv.Itoa(time.Now().Year() - 1)

	requestID := NewRequestID()

	expectedResponse := &IncomeReportResponse{
		Version: "1.0.0",
		Status:  "SUCCESS",
		Body: IncomeReportBody{
			PayerSource: IncomeReportPayerSource{
				Name:           CelcoinBankName,
				DocumentNumber: CelcoinBankISPB,
			},
			Owner: IncomeReportOwner{
				DocumentNumber: documentNumber,
				Name:           userName,
				Type:           "NATURAL_PERSON",
				CreateDate:     time.Now().Format(time.RFC3339),
			},
			Account: IncomeReportAccount{
				Branch:  "0001",
				Account: accountNumber,
			},
			Balances: []IncomeReportBalance{
				{
					CalendarYear: calendarYear,
					Amount:       150,
					Currency:     "BRL",
					Type:         "SALDO",
				},
			},
			IncomeFile: "JVBERi0xLjcKJeLjz9MKNyAwIG9iago8PC9GaWx0ZXIvRmxhdGVE",
			FileType:   "pdf",
		},
	}

	mockResponseBody, err := json.Marshal(expectedResponse)
	s.assert.NoError(err)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.assert.Greater(len(mockResponseBody), 0)

	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	ctx := context.WithValue(s.ctx, requestIDKey("Request-Id"), requestID)
	response, err := s.incomeReport.GetIncomeReport(ctx, &calendarYear, &accountNumber)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}
