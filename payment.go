package celcoin

import (
	"bytes"
	"context"
	"encoding/json"

	//"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Payment ...
type Payment struct {
	session    Session
	httpClient *LoggingHTTPClient
}

// NewPayment ... cria uma nova instância do serviço Payment.
func NewPayment(httpClient *http.Client, session Session) *Payment {
	return &Payment{
		session:    session,
		httpClient: NewLoggingHTTPClient(httpClient),
	}
}

// AuthorizePayment ... envia uma requisição para validar o pagamento via endpoint billpayment/authorize.
func (p *Payment) AuthorizePayment(ctx context.Context,
	request *ValidatePaymentRequest) (*PaymentResponse, error) {

	requestID := grok.GetRequestID(ctx)
	fields := logrus.Fields{
		"request_id": requestID,
		"request":    request,
	}

	u, err := url.Parse(p.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing API endpoint")
		return nil, err
	}

	// Concatena o endpoint billpayment e o subcaminho "validate"
	u.Path = path.Join(u.Path, BillPaymentAuthorizeBasePath, BillPaymentAuthorizePath)
	endpoint := u.String()

	reqBytes, err := json.Marshal(request)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error encoding model to JSON")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBytes))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error performing the request")
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
		var response *PaymentResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding JSON response")
			return nil, ErrDefaultPayment
		}
		return response, nil
	}

	var bodyErr ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding JSON error response")
		return nil, err
	}

	if bodyErr.Error != nil {
		err := FindPaymentError(*bodyErr.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", bodyErr.Error).WithFields(fields).WithError(err).
			Error("celcoin authorize payment error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// ExecutePayment ... envia uma requisição para confirmar o pagamento via endpoint billpayment/confirm.
func (p *Payment) ExecutePayment(ctx context.Context, request *ExecPaymentRequest) (*ExecPaymentResponse, error) {
	requestID := grok.GetRequestID(ctx)
	fields := logrus.Fields{
		"request_id": requestID,
		"request":    request,
	}

	u, err := url.Parse(p.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing API endpoint")
		return nil, err
	}

	// Concatena o endpoint billpayment e o subcaminho "confirm"
	u.Path = path.Join(u.Path, BillPaymentConfirmBasePath)
	endpoint := u.String()

	reqBytes, err := json.Marshal(request)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error encoding model to JSON")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBytes))
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error performing the request")
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
		var response ExecPaymentResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding JSON response")
			return nil, ErrDefaultPayment
		}
		return &response, nil
	}

	var bodyErr ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding JSON error response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Error != nil {
		err := FindPaymentError(*bodyErr.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", bodyErr.Error).WithFields(fields).WithError(err).
			Error("celcoin authorize payment error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// GetPayment ... envia uma requisição para confirmar o pagamento via endpoint billpayment/confirm.
func (p *Payment) Get(ctx context.Context, request *GetPaymentRequest) (*GetPaymentResponse, error) {
	requestID := grok.GetRequestID(ctx)
	fields := logrus.Fields{
		"request_id": requestID,
		"request":    request,
	}

	u, err := url.Parse(p.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing API endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BillPaymenStatusBasePath, BillPaymentStatusPath)

	q := u.Query()

	if request != nil && len(request.ClientRequestID) > 0 {
		q.Set("clientRequestId", request.ClientRequestID)
	}

	if request != nil && len(request.TransactionID) > 0 {
		q.Set("id", request.TransactionID)
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating request")
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error performing the request")
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
		var response GetPaymentResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding JSON response")
			return nil, ErrDefaultPayment
		}
		return &response, nil
	}

	var bodyErr ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding JSON error response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Error != nil {
		err := FindPaymentError(*bodyErr.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", bodyErr.Error).WithFields(fields).WithError(err).
			Error("celcoin authorize payment error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}
