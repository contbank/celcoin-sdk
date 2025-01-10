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

// Customers ...
type Customers struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewCustomers ...
func NewCustomers(httpClient *http.Client, session Session) *Customers {
	return &Customers{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// FindAccounts ...
func (c *Customers) FindAccounts(ctx context.Context,
	documentNumber *string, accountNumber *string) (*CustomerResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	if documentNumber != nil {
		fields["document_number"] = documentNumber
	}

	if accountNumber != nil {
		fields["account_number"] = accountNumber
	}

	endpoint, err := c.getCustomerAPIEndpoint(requestID, documentNumber, accountNumber)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error getting customer api endpoint")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
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
		var response CustomerResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		fields["response"] = response
		logrus.WithFields(fields).
			Info("response with success")

		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr *ErrorResponse

	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindError(errModel.Code, errModel.Messages...)
	}

	logrus.WithFields(fields).
		Error("error default customers accounts - FindAccounts")

	return nil, ErrDefaultCustomersAccounts
}

// CreateAccount ... cria uma nova conta de cliente
func (c *Customers) CreateAccount(ctx context.Context, customerData *Customer) (*CustomerOnboardingResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
	}

	endpoint := c.session.APIEndpoint
	reqBody, err := json.Marshal(customerData)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error marshalling request body")
		return nil, err
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing endpoint")
		return nil, err
	}

	url.Path = path.Join(url.Path, NaturalPersonOnboardingPath)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(string(reqBody)))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error creating request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error executing request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response CustomerOnboardingResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error unmarshalling response")
			return nil, err
		}

		fields["response"] = response
		logrus.WithFields(fields).Info("response with success")
		return &response, nil
	}

	var bodyErr *ErrorResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error unmarshalling error response")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindOnboardingError(errModel.Code, &resp.StatusCode)
	}

	logrus.WithFields(fields).Error("error default create account")
	return nil, ErrDefaultCustomersAccounts
}

// GetOnboardingProposal ... consulta o proposalId
func (c *Customers) GetOnboardingProposal(ctx context.Context, proposalId string) (*OnboardingProposalResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id":  requestID,
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

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error creating request")
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
		var response OnboardingProposalResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error unmarshalling response")
			return nil, err
		}

		fields["response"] = response
		logrus.WithFields(fields).Info("response with success")
		return &response, nil
	}

	var bodyErr *ErrorResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error unmarshalling error response")
		return nil, err
	}

	if len(bodyErr.Errors) > 0 {
		errModel := bodyErr.Errors[0]
		return nil, FindOnboardingError(errModel.Code, &resp.StatusCode)
	}

	logrus.WithFields(fields).Error("error default onboarding proposal")
	return nil, ErrDefaultCustomersAccounts
}

// getCustomerAPIEndpoint
func (c *Customers) getCustomerAPIEndpoint(requestID string, documentNumber *string,
	accountNumber *string) (*string, error) {

	fields := logrus.Fields{
		"request_id":      requestID,
		"document_number": documentNumber,
		"account_number":  accountNumber,
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).
			WithError(err).Error("error api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, CustomersPath)
	q := u.Query()

	if documentNumber != nil && len(*documentNumber) > 0 {
		q.Set("documentNumber", *documentNumber)
	}

	if accountNumber != nil && len(*accountNumber) > 0 {
		q.Set("account", *accountNumber)
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	fields["endpoint"] = endpoint
	logrus.WithFields(fields).Info("get endpoint success")

	return &endpoint, nil
}
