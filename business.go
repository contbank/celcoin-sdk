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

// Business ...
type Business struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewBusiness ...
func NewBusiness(httpClient *http.Client, session Session) *Business {
	return &Business{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// FindAccounts ...
func (c *Business) FindAccounts(ctx context.Context,
	documentNumber *string, accountNumber *string) (*BusinessResponse, error) {

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

	endpoint, err := c.getApiEndpoint(requestID, documentNumber, accountNumber)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error getting business api endpoint")
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
		var response BusinessResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshal")
			return nil, err
		}

		fields["response"] = response
		logrus.WithFields(fields).
			Info("response with success - FindBusiness")

		return &response, nil
	} else if resp.StatusCode == http.StatusNotFound {
		logrus.
			WithFields(fields).
			WithError(ErrEntryNotFound).
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

// getApiEndpoint
func (c *Business) getApiEndpoint(requestID string, documentNumber *string,
	accountNumber *string) (*string, error) {

	fields := logrus.Fields{
		"request_id": requestID,
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

	fields["endpoint"] = endpoint
	logrus.WithFields(fields).Info("get endpoint success")

	return &endpoint, nil
}
