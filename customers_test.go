package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// CustomersTestSuite ...
type CustomersTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	ctx        context.Context
	session    *Session
	customers  *Customers
	mockClient *MockHTTPClient
}

// TestCustomersTestSuite ...
func TestCustomersTestSuite(t *testing.T) {
	suite.Run(t, new(CustomersTestSuite))
}

// SetupTest ...
func (s *CustomersTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	apiEndpoint := "https://sandbox.openfinance.celcoin.dev"
	session := &Session{APIEndpoint: apiEndpoint}
	s.session = session

	s.mockClient = new(MockHTTPClient)
	httpClient := &http.Client{Transport: s.mockClient}
	s.customers = NewCustomers(httpClient, *s.session)
}

// TestFindAccounts ...
func (s *CustomersTestSuite) TestFindAccounts() {
	documentNumber := "12345678901"
	accountNumber := "123456"
	defaultBranch := "0001"

	requestID := NewRequestID()
	nowString := time.Now().Format("2006-01-02T15:04:05")
	nowFormatted, err := ParseStringToCelcoinTime(nowString, "2006-01-02T15:04:05")
	s.assert.NoError(err)

	expectedResponse := &CustomerResponse{
		Body: CustomerResponseBody{
			StatusAccount:  "ACTIVE",
			DocumentNumber: documentNumber,
			PhoneNumber:    "+5511999999999",
			Email:          "test@contbank.com",
			ClientCode:     requestID,
			MotherName:     "Mother Name",
			FullName:       "Test User Name",
			SocialName:     "",
			BirthDate:      "28-06-1990",
			Address: Address{
				PostalCode:        "99999999",
				Street:            "Test Street",
				Number:            "9999",
				AddressComplement: "",
				Neighborhood:      "Test Neighborhood",
				City:              "City",
				State:             "ST",
			},
			IsPoliticallyExposedPerson: false,
			Account: Account{
				Branch:  defaultBranch,
				Account: accountNumber,
			},
			CreateDate: CustomTime{Time: nowFormatted},
		},
		Version: "1.0.0",
		Status:  "SUCCESS",
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
	response, err := s.customers.FindAccounts(ctx, &documentNumber, &accountNumber)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}

// TestCreateAccount ...
func (s *CustomersTestSuite) TestCreateAccount() {
	requestID := NewRequestID()
	clientCode := NewRequestID()
	documentNumber := "99999999999"

	customerData := &Customer{
		ClientCode:     clientCode,
		DocumentNumber: documentNumber,
		PhoneNumber:    "+5511999999999",
		Email:          "test@contbank.com",
		MotherName:     "Mother Name",
		FullName:       "Test User Name",
		SocialName:     "",
		BirthDate:      "28-06-1990",
		Address: CustomerAddress{
			PostalCode:        "99999999",
			Street:            "Test Street",
			Number:            "9999",
			AddressComplement: "",
			Neighborhood:      "Test Neighborhood",
			City:              "City",
			State:             "ST",
		},
		IsPoliticallyExposedPerson: false,
		OnboardingType:             "BAAS",
	}

	expectedResponse := &CustomerOnboardingResponse{
		Body: CustomerOnboardingResponseBody{
			ProposalID:     NewRequestID(),
			ClientCode:     clientCode,
			DocumentNumber: documentNumber,
		},
		Version: "1.0.0",
		Status:  "PROCESSING",
	}

	mockResponseBody, _ := json.Marshal(expectedResponse)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	ctx := context.WithValue(s.ctx, requestIDKey("Request-Id"), requestID)
	response, err := s.customers.CreateAccount(ctx, customerData)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}

// TestGetOnboardingProposalPROCESSING ...
func (s *CustomersTestSuite) TestGetOnboardingProposalPROCESSING() {

	requestID := NewRequestID()
	clientCode := NewRequestID()
	documentNumber := "99999999999"
	proposalID := NewRequestID()

	expectedResponse := &OnboardingProposalResponse{
		Body: OnboardingProposalResponseBody{
			Limit:        10,
			CurrentPage:  1,
			LimitPerPage: 10,
			TotalPages:   1,
			TotalItems:   1,
			Proposals: []Proposal{
				{
					ProposalID:     proposalID,
					ClientCode:     clientCode,
					DocumentNumber: documentNumber,
					Status:         OnboardingStatusProcessing,
					ProposalType:   ProposalTypePF,
					CreatedAt:      "2023-01-01T00:00:00Z",
					UpdatedAt:      "2023-01-01T00:00:00Z",
					DocumentsCopys: []DocumentsCopy{
						{
							ProposalID:      proposalID,
							DocumentNumber:  documentNumber,
							DocumentsCopyID: NewRequestID(),
							Status:          OnboardingStatusProcessing,
							URL:             "https://example.com/document",
							CreatedAt:       "2023-01-01T00:00:00Z",
							UpdatedAt:       "2023-01-01T00:00:00Z",
						},
					},
				},
			},
		},
		Version: "1.0.0",
		Status:  "SUCCESS",
	}

	mockResponseBody, err := json.Marshal(expectedResponse)
	s.assert.NoError(err)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	ctx := context.WithValue(s.ctx, requestIDKey("Request-Id"), requestID)
	response, err := s.customers.GetOnboardingProposal(ctx, proposalID)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}

// TestGetOnboardingProposalFiles ...
func (s *CustomersTestSuite) TestGetOnboardingProposalFiles() {
	proposalID := "d92fffcf-6d39-4b89-bc63-bdd2e3077a23"
	requestID := NewRequestID()

	expectedResponse := &OnboardingProposalFilesResponse{
		Body: OnboardingProposalFilesResponseBody{
			Files: []OnboardingFile{
				{
					Type:           "CNH_BACK",
					URL:            "https://example.com/document1.jpg",
					ExpirationTime: time.Date(2025, 1, 11, 18, 8, 10, 0, time.UTC),
				},
				{
					Type:           "CNH_FRONT",
					URL:            "https://example.com/document2.jpg",
					ExpirationTime: time.Date(2025, 1, 11, 18, 8, 10, 0, time.UTC),
				},
				{
					Type:           "SELFIE",
					URL:            "https://example.com/document3.png",
					ExpirationTime: time.Date(2025, 1, 11, 18, 8, 11, 0, time.UTC),
				},
			},
			ClientCode:     "c150d3ea-8ad6-4d2f-89da-d2ad5c316527",
			DocumentNumber: "99999999999",
			ProposalID:     proposalID,
		},
		Version: "1.0.0",
		Status:  "SUCCESS",
	}

	mockResponseBody, err := json.Marshal(expectedResponse)
	s.assert.NoError(err)

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(mockResponseBody)),
	}

	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	ctx := context.WithValue(s.ctx, requestIDKey("Request-Id"), requestID)
	response, err := s.customers.GetOnboardingProposalFiles(ctx, proposalID)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}
