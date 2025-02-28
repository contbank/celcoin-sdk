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
	"strconv"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// BillPaymentEndpoint define o caminho base para operações de payment na Celcoin.
const BillPaymentEndpoint = "/baas/v2/billpayment"

// ErrDefaultPayment (já consta no sdk)
// Payment Service
type Payment struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewPayment cria uma nova instância do serviço Payment.
func NewPayment(httpClient *http.Client, session Session) *Payment {
	return &Payment{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// ValidatePayment envia uma requisição para validar o pagamento via endpoint billpayment/validate.
func (p *Payment) ValidatePayment(ctx context.Context, correlationID string, model *ValidatePaymentRequest) (*ValidatePaymentResponse, error) {
	fields := logrus.Fields{"request_id": correlationID}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing API endpoint")
		return nil, err
	}

	// Concatena o endpoint billpayment e o subcaminho "validate"
	u.Path = path.Join(u.Path, BillPaymentEndpoint, "validate")
	endpoint := u.String()

	reqBytes, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error encoding model to JSON")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBytes))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error retrieving authentication token")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error performing the request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response ValidatePaymentResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error decoding JSON response")
			return nil, ErrDefaultPayment
		}
		return &response, nil
	}

	var bodyErr ErrorResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding JSON error response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.WithField("celcoin_error", bodyErr).WithFields(fields).WithError(err).Error("celcoin validate payment error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// ConfirmPayment envia uma requisição para confirmar o pagamento via endpoint billpayment/confirm.
func (p *Payment) ConfirmPayment(ctx context.Context, correlationID string, model *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error) {
	fields := logrus.Fields{"request_id": correlationID}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing API endpoint")
		return nil, err
	}

	// Concatena o endpoint billpayment e o subcaminho "confirm"
	u.Path = path.Join(u.Path, BillPaymentEndpoint, "confirm")
	endpoint := u.String()

	reqBytes, err := json.Marshal(model)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error encoding model to JSON")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBytes))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error retrieving authentication token")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error performing the request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response ConfirmPaymentResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error decoding JSON response")
			return nil, ErrDefaultPayment
		}
		return &response, nil
	}

	var bodyErr ErrorResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding JSON error response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.WithField("celcoin_error", bodyErr).WithFields(fields).WithError(err).Error("celcoin confirm payment error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// FilterPayments envia uma requisição GET para filtrar os pagamentos via endpoint billpayment com parâmetros de consulta.
func (p *Payment) FilterPayments(ctx context.Context, correlationID string, model *FilterPaymentsRequest) (*FilterPaymentsResponse, error) {
	fields := logrus.Fields{"request_id": correlationID}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing API endpoint")
		return nil, err
	}

	// Endpoint para filtro: billpayment (sem subcaminho específico)
	u.Path = path.Join(u.Path, BillPaymentEndpoint)

	q := u.Query()
	q.Set("bankAccount", model.BankAccount)
	q.Set("bankBranch", model.BankBranch)
	q.Set("pageSize", strconv.Itoa(model.PageSize))
	if model.PageToken != nil {
		q.Set("pageToken", *model.PageToken)
	}
	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error retrieving authentication token")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error performing the request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response FilterPaymentsResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error decoding JSON response")
			return nil, ErrDefaultPayment
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr ErrorResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding JSON error response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.WithField("celcoin_error", bodyErr).WithFields(fields).WithError(err).Error("celcoin filter payments error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}

// DetailPayment envia uma requisição GET para obter o detalhe de um pagamento via endpoint billpayment/detail.
func (p *Payment) DetailPayment(ctx context.Context, correlationID string, model *DetailPaymentRequest) (*PaymentResponse, error) {
	fields := logrus.Fields{"request_id": correlationID}

	if err := grok.Validator.Struct(model); err != nil {
		return nil, grok.FromValidationErros(err)
	}

	u, err := url.Parse(p.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error parsing API endpoint")
		return nil, err
	}

	// Monta o endpoint: billpayment/detail
	u.Path = path.Join(u.Path, BillPaymentEndpoint, "detail")

	q := u.Query()
	q.Set("bankAccount", model.BankAccount)
	q.Set("bankBranch", model.BankBranch)
	q.Set("authenticationCode", model.AuthenticationCode)
	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error creating request")
		return nil, err
	}

	token, err := p.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error retrieving authentication token")
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("api-version", p.session.APIVersion)
	req.Header.Add("x-correlation-id", correlationID)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("error performing the request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var response PaymentResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error decoding JSON response")
			return nil, ErrDefaultPayment
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var bodyErr ErrorResponse
	if err := json.Unmarshal(respBody, &bodyErr); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding JSON error response")
		return nil, ErrDefaultPayment
	}

	if bodyErr.Code != "" {
		err = FindError(bodyErr.Code, bodyErr.Message)
		logrus.WithField("celcoin_error", bodyErr).WithFields(fields).WithError(err).Error("celcoin get payment detail error")
		return nil, err
	}

	return nil, ErrDefaultPayment
}
