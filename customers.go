package celcoin

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

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
