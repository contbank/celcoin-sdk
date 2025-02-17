package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/contbank/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// StatementTestSuite ...
type StatementTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	ctx        context.Context
	session    *Session
	statement  *Statement
	mockClient *MockHTTPClient
}

// TestStatementTestSuite ...
func TestStatementTestSuite(t *testing.T) {
	suite.Run(t, new(StatementTestSuite))
}

// SetupTest ...
func (s *StatementTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	apiEndpoint := "https://sandbox.openfinance.celcoin.dev"
	session := &Session{APIEndpoint: apiEndpoint}
	s.session = session

	s.mockClient = new(MockHTTPClient)
	httpClient := &http.Client{Transport: s.mockClient}
	s.statement = NewStatement(httpClient, *s.session)
}

// TestFindAccounts ...
func (s *StatementTestSuite) TestGetStatements() {
	documentNumber := "12345678901"
	accountNumber := "123456"

	requestID := NewRequestID()

	expectedResponse := &StatementResponse{
		Status:       "SUCCESS",
		Version:      "1.0.0",
		TotalItems:   1,
		CurrentPage:  1,
		LimitPerPage: 10,
		TotalPages:   1,
		DateFrom:     "2025-01-13T00:00:00",
		DateTo:       "2025-01-19T23:59:59.9999999",
		Body: StatementBody{
			Account:        accountNumber,
			DocumentNumber: documentNumber,
			Movements: []StatementMovement{
				{
					ID:             NewRequestID(),
					ClientCode:     NewRequestID(),
					Description:    "teste de TEF entre contas p/ evento webhook.",
					CreateDate:     "2025-01-16T09:07:12",
					LastUpdateDate: "2025-01-16T09:11:29",
					Amount:         10.99,
					Status:         "Saldo Liberado",
					BalanceType:    "DEBIT",
					MovementType:   "TEFTRANSFEROUT",
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

	s.assert.Greater(len(mockResponseBody), 0)

	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	ctx := context.WithValue(s.ctx, requestIDKey("Request-Id"), requestID)

	dateFrom := "2025-01-13T00:00:00"
	dateTo := "2025-01-19T23:59:59"
	statementRequest := &StatementRequest{
		DocumentNumber: &documentNumber,
		Account:        &accountNumber,
		DateFrom:       &dateFrom,
		DateTo:         &dateTo,
		LimitPerPage:   aws.String("10"),
	}

	response, err := s.statement.GetStatements(ctx, statementRequest)
	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}

func (s *StatementTestSuite) TestGetStatementsInvalidInputError() {
	// Simula a resposta de erro da API
	errorResponse := map[string]interface{}{
		"status":  "ERROR",
		"version": "1.0.0",
		"error": map[string]interface{}{
			"errorCode": "CBE073",
			"message":   "É necessário informar pelo menos um dos campos: account, ou documentNumber.",
		},
	}
	mockResponseBody, err := json.Marshal(errorResponse)
	s.assert.NoError(err)

	mockResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	// Cria um mock do http.Client
	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Chama o método GetStatements
	_, err = s.statement.GetStatements(context.Background(), &StatementRequest{})

	// Verifica se o erro é do tipo *grok.Error
	grokErr, ok := err.(*grok.Error)
	s.assert.True(ok, "erro deve ser do tipo *grok.Error")

	// Verifica se o erro retornado é o esperado
	s.assert.Equal(400, grokErr.Code)
	s.assert.Equal("INVALID_INPUT", grokErr.Key)
	s.assert.Len(grokErr.Messages, 1)
	s.assert.Equal("É necessário informar pelo menos um dos campos: account, ou documentNumber.", grokErr.Messages[0])
}

func (s *StatementTestSuite) TestGetStatementsUnknownError() {
	// Simula a resposta de erro da API
	errorResponse := map[string]interface{}{
		"status":  "ERROR",
		"version": "1.0.0",
		"error": map[string]interface{}{
			"errorCode": "CBE999",
			"message":   "Unknown server error",
		},
	}
	mockResponseBody, err := json.Marshal(errorResponse)
	s.assert.NoError(err)

	mockResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	// Cria um mock do http.Client
	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Chama o método GetStatements
	_, err = s.statement.GetStatements(context.Background(), &StatementRequest{})

	// Verifica se o erro é do tipo *grok.Error
	grokErr, ok := err.(*grok.Error)
	s.assert.True(ok, "erro deve ser do tipo *grok.Error")

	// Verifica se o erro retornado é o esperado
	s.assert.Equal(http.StatusInternalServerError, grokErr.Code)
	s.assert.Equal("UNKNOWN_ERROR", grokErr.Key)
	s.assert.Len(grokErr.Messages, 1)
	s.assert.Equal("unknown error", grokErr.Messages[0])
}
