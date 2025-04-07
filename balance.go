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

// Balance ...
type Balance struct {
	session        Session
	authentication *Authentication
	httpClient     *LoggingHTTPClient
}

// NewBalance ...
func NewBalance(httpClient *http.Client, session Session) *Balance {
	return &Balance{
		session:        session,
		httpClient:     NewLoggingHTTPClient(httpClient),
		authentication: NewAuthentication(httpClient, session),
	}
}

// Balance ...
func (c *Balance) Balance(ctx context.Context, accountNumber string) (*BalanceResponse, error) {
	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"interface":  "Balance",
		"service":    "Balance",
	}

	u, err := url.Parse(c.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BalancePath)

	q := u.Query()

	q.Set("account", accountNumber)

	u.RawQuery = q.Encode()
	endpoint := u.String()

	logrus.WithFields(fields).WithField("celcoin_endpoint", endpoint).
		Info("requesting balance")

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
		var response *BalanceResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultBalance
		}

		logrus.WithFields(fields).WithField("celcoin_response", response).
			Info("received celcoin response")

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultBalance
	}

	if errResponse.Error != nil {
		err := FindBalanceError(*errResponse.Error.ErrorCode, *errResponse.Error.Message)
		logrus.WithField("celcoin_error", errResponse.Error).WithFields(fields).WithError(err).
			Error("celcoin get balance error")
		return nil, err
	}

	return nil, ErrDefaultBalance
}
