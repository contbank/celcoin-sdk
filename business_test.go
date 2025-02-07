package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/contbank/grok"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// BusinessTestSuite ...
type BusinessTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	ctx        context.Context
	session    *Session
	business   *Business
	mockClient *MockHTTPClient
}

// TestBusinessTestSuite ...
func TestBusinessTestSuite(t *testing.T) {
	suite.Run(t, new(BusinessTestSuite))
}

// SetupTest ...
func (s *BusinessTestSuite) SetupTest() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()

	apiEndpoint := "https://sandbox.openfinance.celcoin.dev"
	session := &Session{APIEndpoint: apiEndpoint}
	s.session = session

	s.mockClient = new(MockHTTPClient)
	httpClient := &http.Client{Transport: s.mockClient}
	s.business = NewBusiness(httpClient, *s.session)
}

// TestFindAccounts ...
func (s *BusinessTestSuite) TestFindAccounts() {
	documentNumber := grok.GeneratorCNPJ()
	accountNumber := "123456"
	defaultBranch := "0001"

	requestID := NewRequestID()
	nowString := time.Now().Format("2006-01-02T15:04:05")
	nowFormatted, err := ParseStringToCelcoinTime(nowString, "2006-01-02T15:04:05")
	s.assert.NoError(err)

	expectedResponse := &BusinessResponse{
		Body: BusinessResponseBody{
			StatusAccount:       "ATIVO",
			DocumentNumber:      documentNumber,
			ClientCode:          requestID,
			BusinessPhoneNumber: grok.GeneratorCellphone(),
			BusinessEmail:       "company@contbank.com",
			CreateDate:          CustomTime{Time: nowFormatted},
			BusinessName:        "Company Test S.A",
			TradingName:         "Company Test For Tests Purposes",
			Owners: []Owner{
				{
					OwnerType:                  "SOCIO",
					DocumentNumber:             grok.GeneratorCPF(),
					FullName:                   "Owner Test Name",
					BirthDate:                  "28-06-1990",
					PhoneNumber:                grok.GeneratorCellphone(),
					Email:                      "owner@contbank.com",
					MotherName:                 "Mother Name",
					IsPoliticallyExposedPerson: false,
					Address: Address{
						PostalCode:        "99999999",
						Street:            "Test Street",
						Number:            "9999",
						AddressComplement: "",
						Neighborhood:      "Test Neighborhood",
						City:              "City",
						State:             "ST",
					},
				},
			},
			BusinessAddress: Address{
				PostalCode:        "99999999",
				Street:            "Test Street",
				Number:            "9999",
				AddressComplement: "",
				Neighborhood:      "Test Neighborhood",
				City:              "City",
				State:             "ST",
			},
			BusinessAccount: BusinessAccount{
				Branch:  defaultBranch,
				Account: accountNumber,
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

	s.assert.Greater(len(mockResponseBody), 0)

	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	ctx := context.WithValue(s.ctx, requestIDKey("Request-Id"), requestID)
	response, err := s.business.FindAccounts(ctx, &documentNumber, &accountNumber)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}

// TestCreateAccount ...
func (s *BusinessTestSuite) TestCreateAccount() {
	requestID := NewRequestID()
	clientCode := NewRequestID()
	documentNumber := grok.GeneratorCNPJ()

	businessData := &BusinessOnboardingRequest{
		ClientCode:     clientCode,
		ContactNumber:  grok.GeneratorCellphone(),
		DocumentNumber: documentNumber,
		BusinessAddress: Address{
			PostalCode:        "99999999",
			Street:            "Test Street",
			Number:            "9999",
			AddressComplement: "",
			Neighborhood:      "Test Neighborhood",
			City:              "City",
			State:             "ST",
		},
		BusinessEmail:  "test@contbank.com",
		BusinessName:   "Company Test S.A",
		TradingName:    "Company Test For Tests Purposes",
		CompanyType:    "PJ",
		OnboardingType: "BAAS",
		Owner: []Owner{
			{
				OwnerType:                  LegalPersonOwnerTypeSocio,
				DocumentNumber:             grok.GeneratorCPF(),
				FullName:                   "Owner Name",
				BirthDate:                  "28-06-1990",
				PhoneNumber:                "+5511999999999",
				Email:                      "owner@contbank.com",
				MotherName:                 "Owner Mother Name",
				IsPoliticallyExposedPerson: false,
				Address: Address{
					PostalCode:        "99999999",
					Street:            "Test Street",
					Number:            "9999",
					AddressComplement: "",
					Neighborhood:      "Test Neighborhood",
					City:              "City",
					State:             "ST",
				},
			},
		},
	}

	expectedResponse := &BusinessOnboardingResponse{
		Body: BusinessOnboardingResponseBody{
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
	response, err := s.business.CreateAccount(ctx, businessData)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}

// TestGetOnboardingProposalPROCESSING ...
func (s *BusinessTestSuite) TestGetOnboardingProposalPROCESSING() {

	requestID := NewRequestID()
	clientCode := NewRequestID()
	documentNumber := grok.GeneratorCNPJ()
	proposalID := NewRequestID()
	documentsCopyID := NewRequestID()

	expectedResponse := &OnboardingProposalResponseBody{
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
				ProposalType:   ProposalTypePJ,
				CreatedAt:      "2025-01-01T00:00:00Z",
				UpdatedAt:      "2025-01-01T00:00:00Z",
				DocumentsCopys: []DocumentsCopy{
					{
						ProposalID:      proposalID,
						DocumentNumber:  documentNumber,
						DocumentsCopyID: documentsCopyID,
						Status:          OnboardingStatusProcessing,
						URL:             "https://example.com/document",
						CreatedAt:       "2025-01-01T00:00:00Z",
						UpdatedAt:       "2025-01-01T00:00:00Z",
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

	s.mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	ctx := context.WithValue(s.ctx, requestIDKey("Request-Id"), requestID)
	response, err := s.business.GetLegalPersonOnboardingProposal(ctx, proposalID)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}

// TestGetOnboardingProposalFiles ...
func (s *BusinessTestSuite) TestGetOnboardingProposalFiles() {
	requestID := NewRequestID()
	proposalID := NewRequestID()

	expectedResponse := &OnboardingProposalFilesResponse{
		Body: OnboardingProposalFilesResponseBody{
			Files: []OnboardingFile{
				{
					Type:           "CNH_BACK",
					URL:            "https://example.com/document1.jpg",
					ExpirationTime: time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
				},
				{
					Type:           "CNH_FRONT",
					URL:            "https://example.com/document2.jpg",
					ExpirationTime: time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
				},
				{
					Type:           "SELFIE",
					URL:            "https://example.com/document3.png",
					ExpirationTime: time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
				},
				{
					Type:           "CONTRATO_SOCIAL",
					URL:            "https://example.com/document4.png",
					ExpirationTime: time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
				},
			},
			ClientCode:     requestID,
			DocumentNumber: grok.GeneratorCNPJ(),
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
	response, err := s.business.GetLegalPersonOnboardingProposalFiles(ctx, proposalID)

	s.assert.NoError(err)
	s.assert.NotNil(response)
	s.assert.Equal(expectedResponse, response)
}
