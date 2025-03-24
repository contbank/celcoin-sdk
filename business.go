package celcoin

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

// Business ...
type Business struct {
	session        Session
	httpClient     *LoggingHTTPClient
	authentication *Authentication
}

// NewBusiness ...
func NewBusiness(httpClient *http.Client, session Session) *Business {
	return &Business{
		session:        session,
		httpClient:     NewLoggingHTTPClient(httpClient),
		authentication: NewAuthentication(httpClient, session),
	}
}

// FindAccounts ...
func (c *Business) FindAccounts(ctx context.Context,
	documentNumber *string, accountNumber *string) (*BusinessResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"interface":  "FindAccounts",
		"service":    "business",
	}

	if documentNumber != nil {
		fields["document_number"] = documentNumber
	}

	if accountNumber != nil {
		fields["account_number"] = accountNumber
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BusinessPath)
	q := u.Query()

	if documentNumber != nil && len(*documentNumber) > 0 {
		q.Set("documentNumber", *documentNumber)
	}

	if accountNumber != nil && len(*accountNumber) > 0 {
		q.Set("account", *accountNumber)
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("requesting business account data")

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error new request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response BusinessResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		logrus.WithFields(fields).WithField("celcoin_response", response).
			Info("received celcoin response")

		return &response, nil
	} else if resp.StatusCode == http.StatusNotFound {
		logrus.WithFields(fields).WithError(ErrEntryNotFound).
			Error("error entry not found - FindBusiness")
		return nil, ErrEntryNotFound
	} else if resp.StatusCode == http.StatusInternalServerError {
		logrus.
			WithFields(fields).
			Error("internal server error - FindBusiness")
		return nil, ErrDefaultBusinessAccounts
	}

	var bodyErr *ErrorResponse

	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		logrus.WithFields(fields).
			Error("body error - FindBusiness")
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.
		WithFields(fields).
		Error("error default business accounts - FindBusiness")

	return nil, ErrDefaultBusinessAccounts
}

// CreateAccount ... cria uma nova conta de cliente
func (c *Business) CreateAccount(ctx context.Context, businessData *BusinessOnboardingRequest) (*BusinessOnboardingResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"interface":  "CreateAccount",
		"service":    "business",
	}

	endpoint := c.session.APIEndpoint
	reqBody, err := json.Marshal(businessData)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error marshalling request body")
		return nil, err
	}

	logrus.WithFields(fields).Info("request body marshalled successfully")

	url, err := url.Parse(endpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	url.Path = path.Join(url.Path, LegalPersonOnboardingPath)

	logrus.WithFields(fields).WithField("celcoin_endpoint", url.String()).
		WithField("celcoin_request", businessData).Info("requesting create business account")

	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), strings.NewReader(string(reqBody)))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	logrus.WithFields(fields).
		Info("request created successfully")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error executing request")
		return nil, err
	}
	defer resp.Body.Close()

	logrus.WithFields(fields).
		Info("request executed successfully")

	respBody, _ := ioutil.ReadAll(resp.Body)

	logrus.WithFields(fields).
		Infof("response received with status code %d", resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		var response BusinessOnboardingResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshalling response")
			return nil, err
		}

		logrus.WithFields(fields).WithField("celcoin_response", response).
			Info("response with success")
		return &response, nil
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshalling error response")
		return nil, err
	}

	if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
		err := FindOnboardingError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithFields(fields).WithError(err).
			Error("error creating business account")
		return nil, err
	}

	logrus.WithFields(fields).
		Error("error creating business account")
	return nil, ErrDefaultBusinessAccounts
}

// GetLegalPersonOnboardingProposal ... consulta o proposalId
func (c *Business) GetLegalPersonOnboardingProposal(ctx context.Context, proposalId string) (*OnboardingProposalResponseBody, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":  requestID,
		"interface":   "GetLegalPersonOnboardingProposal",
		"service":     "business",
		"proposal_id": proposalId,
	}

	endpoint := c.session.APIEndpoint
	u, err := url.Parse(endpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, ProposalsPath)

	q := u.Query()
	q.Set("proposalId", proposalId)
	u.RawQuery = q.Encode()

	logrus.WithFields(fields).WithField("celcoin_endpoint", u.String()).
		Info("legal person proposal request")

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error executing request")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response OnboardingProposalResponseBody
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshalling response")
			return nil, err
		}

		logrus.WithFields(fields).WithField("celcoin_response", response).
			Info("response with success")
		return &response, nil
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshalling error response")
		return nil, err
	}

	if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
		err := FindOnboardingError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithFields(fields).WithError(err).
			Error("error getting legal person onboarding proposal")
		return nil, err
	}

	logrus.WithFields(fields).
		Error("error getting legal person onboarding proposal")
	return nil, ErrDefaultBusinessAccounts
}

// GetLegalPersonOnboardingProposalFiles ... consulta os arquivos do proposalId
func (c *Business) GetLegalPersonOnboardingProposalFiles(ctx context.Context, proposalId string) (*OnboardingProposalFilesResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":  requestID,
		"interface":   "GetLegalPersonOnboardingProposalFiles",
		"service":     "business",
		"proposal_id": proposalId,
	}

	endpoint := c.session.APIEndpoint
	u, err := url.Parse(endpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, ProposalFilesPath)

	q := u.Query()
	q.Set("proposalId", proposalId)
	u.RawQuery = q.Encode()

	logrus.WithFields(fields).WithField("celcoin_endpoint", u.String()).
		Info("legal person proposal request")

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error executing request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error reading response body")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var response OnboardingProposalFilesResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshalling response")
			return nil, err
		}

		logrus.WithFields(fields).WithField("celcoin_response", response).
			Info("response with success")
		return &response, nil
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error unmarshalling error response")
		return nil, err
	}

	if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
		err := FindOnboardingError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithFields(fields).WithError(err).
			Error("error getting legal person onboarding proposal files")
		return nil, err
	}

	logrus.WithFields(fields).
		Error("error getting legal person onboarding proposal files")
	return nil, ErrDefaultBusinessAccounts
}

// CancelAccount ... consulta os arquivos do proposalId
func (c *Business) CancelAccount(ctx context.Context,
	accountNumber *string, documentNumber *string, reason *string) (*CancelAccountResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"interface":  "CancelAccount",
		"service":    "business",
	}

	if accountNumber != nil {
		fields["account_number"] = accountNumber
	}

	if reason != nil {
		fields["reason"] = reason
	}

	if documentNumber != nil {
		fields["document_number"] = documentNumber
	}

	endpoint := c.session.APIEndpoint
	u, err := url.Parse(endpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CancelAccountPath)

	q := u.Query()
	if accountNumber != nil {
		q.Set("account", *accountNumber)
	}

	if reason != nil {
		q.Set("reason", *reason)
	}

	if documentNumber != nil && accountNumber == nil {
		q.Set("documentNumber", *documentNumber)
	}

	u.RawQuery = q.Encode()

	logrus.WithFields(fields).WithField("celcoin_endpoint", u.String()).
		Info("celcoin cancel account request")

	req, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error executing request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error reading response body")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var response CancelAccountResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshalling response")
			return nil, err
		}

		logrus.WithFields(fields).WithField("celcoin_response", response).
			Info("response with success")
		return &response, nil
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error unmarshalling error response")
		return nil, err
	}

	if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
		err := FindCancelAccountError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithFields(fields).WithError(err).
			Error("error cancel account")
		return nil, err
	}

	logrus.WithFields(fields).
		Error("error cancel account")
	return nil, ErrDefaultBusinessAccounts
}

// UpdateAccountStatus ... faz atualização do status da conta
func (c *Business) UpdateAccountStatus(ctx context.Context,
	accountNumber *string, documentNumber *string, reason *string, status *string) (*UpdateAccountStatusResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"interface":  "CancelAccount",
		"service":    "business",
	}

	if accountNumber != nil {
		fields["account_number"] = accountNumber
	}

	if reason != nil {
		fields["reason"] = reason
	}

	if documentNumber != nil {
		fields["document_number"] = documentNumber
	}

	endpoint := c.session.APIEndpoint
	u, err := url.Parse(endpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CancelAccountPath)

	q := u.Query()
	if accountNumber != nil {
		q.Set("account", *accountNumber)
	}

	if reason != nil {
		q.Set("reason", *reason)
	}

	if status != nil {
		q.Set("status", *status)
	}

	if documentNumber != nil && accountNumber == nil {
		q.Set("documentNumber", *documentNumber)
	}

	u.RawQuery = q.Encode()

	logrus.WithFields(fields).WithField("celcoin_endpoint", u.String()).
		Info("celcoin update account status request")

	req, err := http.NewRequestWithContext(ctx, "PUT", u.String(), nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error executing request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error reading response body")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var response UpdateAccountStatusResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshalling response")
			return nil, err
		}

		logrus.WithFields(fields).WithField("celcoin_response", response).
			Info("response with success")
		return &response, nil
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error unmarshalling error response")
		return nil, err
	}

	if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
		err := FindUpdateAccountStatusAccountErrors(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithFields(fields).WithError(err).
			Error("error updating account status")
		return nil, err
	}

	logrus.WithFields(fields).
		Error("error updating account status")
	return nil, ErrDefaultBusinessAccounts
}
