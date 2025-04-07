package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Boletos provides methods to create, cancel, query, and download boleto (charge) documents.
type Boletos struct {
	session    Session
	httpClient *LoggingHTTPClient
}

// NewBoletos creates and returns a new instance of Boletos using the given httpClient and session.
func NewBoletos(httpClient *http.Client, session Session) *Boletos {
	return &Boletos{
		session:    session,
		httpClient: NewLoggingHTTPClient(httpClient),
	}
}

// CreateBoleto sends a request to create a new boleto (charge).
// It expects a CreateBoletoRequest and returns the unwrapped CreateBoletoResponse.
func (b *Boletos) CreateBoleto(ctx context.Context, req CreateBoletoRequest) (*CreateBoletoResponse, error) {
	// Validate the request payload.
	if err := grok.Validator.Struct(req); err != nil {
		logrus.WithError(err).Error("CreateBoleto: validation error")
		return nil, grok.FromValidationErros(err)
	}

	// Build the endpoint URL: {APIEndpoint}/baas/v2/charge
	endpoint := fmt.Sprintf("%s/baas/v2/charge", b.session.APIEndpoint)
	logrus.WithField("endpoint", endpoint).Info("CreateBoleto endpoint")

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("CreateBoleto: error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("CreateBoleto: error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	// The injected httpClient handles authentication.

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithError(err).Error("CreateBoleto: error performing HTTP request")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("CreateBoleto: error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CreateBoleto: request failed with status %d, body: %s", resp.StatusCode, string(respBody))
	}

	// Unwrap the response envelope.
	var envelope struct {
		Body    CreateBoletoResponse `json:"body"`
		Version string               `json:"version"`
		Status  string               `json:"status"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("CreateBoleto: error unmarshaling response: %v", err)
	}

	return &envelope.Body, nil
}

// CancelBoleto sends a request to cancel an existing boleto given its transaction ID and a cancellation reason.
func (b *Boletos) CancelBoleto(ctx context.Context, transactionID string, reason string) error {
	cancelPayload := CancelInput{Reason: reason}
	payload, err := json.Marshal(cancelPayload)
	if err != nil {
		return fmt.Errorf("CancelBoleto: error serializing request: %v", err)
	}

	// Build the endpoint URL: {APIEndpoint}/baas/v2/charge/{transactionID}
	endpoint := fmt.Sprintf("%s/baas/v2/charge/%s", b.session.APIEndpoint, transactionID)
	logrus.WithField("endpoint", endpoint).Info("CancelBoleto endpoint")

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("CancelBoleto: error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	// Authentication handled automatically.

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("CancelBoleto: HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("CancelBoleto: request failed with status %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// QueryBoleto queries a boleto by its transaction ID.
// It returns a QueryBoletoResponse containing the charge status.
func (b *Boletos) QueryBoleto(ctx context.Context, transactionID string) (*QueryBoletoResponse, error) {
	base, err := url.Parse(b.session.APIEndpoint)
	if err != nil {
		return nil, fmt.Errorf("QueryBoleto: error parsing API endpoint: %v", err)
	}
	base.Path = path.Join(base.Path, "baas/v2/charge")
	q := base.Query()
	q.Set("TransactionId", transactionID)
	base.RawQuery = q.Encode()
	endpoint := base.String()

	logrus.WithField("endpoint", endpoint).Info("QueryBoleto endpoint")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("QueryBoleto: error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	// Authentication handled by httpClient.

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("QueryBoleto: HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("QueryBoleto: error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("QueryBoleto: request failed with status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var envelope struct {
		Body    QueryBoletoResponse `json:"body"`
		Version string              `json:"version"`
		Status  string              `json:"status"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("QueryBoleto: error unmarshaling response: %v", err)
	}

	return &envelope.Body, nil
}

// DownloadBoletoPDF downloads the boleto PDF file and writes its content to the provided writer.
func (b *Boletos) DownloadBoletoPDF(ctx context.Context, transactionID string, writer io.Writer) error {
	// Build the endpoint URL: {APIEndpoint}/baas/v2/charge/pdf/{transactionID}
	endpoint := fmt.Sprintf("%s/baas/v2/charge/pdf/%s", b.session.APIEndpoint, transactionID)
	logrus.WithField("endpoint", endpoint).Info("DownloadBoletoPDF endpoint")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("DownloadBoletoPDF: error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	// Authentication is handled automatically.

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("DownloadBoletoPDF: HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DownloadBoletoPDF: request failed with status %d, body: %s", resp.StatusCode, string(respBody))
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("DownloadBoletoPDF: error writing PDF to writer: %v", err)
	}

	return nil
}

// GetCharge ...
func (r *Boletos) GetCharge(ctx context.Context,
	request *ChargeRequest) (*ChargeResponse, error) {

	requestID, _ := ctx.Value("Request-Id").(string)
	fields := logrus.Fields{
		"request_id": requestID,
		"interface":  "GetCharge",
		"service":    "charge",
	}

	if request != nil && request.TransactionID != nil && len(*request.TransactionID) > 0 {
		fields["transaction_id"] = request.TransactionID
	}

	if request != nil && request.ExternalID != nil && len(*request.ExternalID) > 0 {
		fields["external_id"] = request.ExternalID
	}

	u, err := url.Parse(r.session.APIEndpoint)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error parsing charge api endpoint")
		return nil, err
	}

	u.Path = path.Join(u.Path, BaasV2ChargePath)

	q := u.Query()

	if request != nil && request.TransactionID != nil && len(*request.TransactionID) > 0 {
		q.Set("transactionId", *request.TransactionID)
	}

	if request != nil && request.ExternalID != nil && len(*request.ExternalID) > 0 {
		q.Set("externalId", *request.ExternalID)
	}

	u.RawQuery = q.Encode()
	endpoint := u.String()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error creating charge request")
		return nil, err
	}

	req.Header.Set("accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error making charge request")
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error reading charge response body")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResponse *ErrorDefaultResponse
		if err := json.Unmarshal(bodyBytes, &errResponse); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error unmarshal charge response")
			return nil, err
		}

		if errResponse != nil && errResponse.Error != nil && len(*errResponse.Error.ErrorCode) > 0 {
			err := FindChargeError(*errResponse.Error.ErrorCode, &resp.StatusCode)
			logrus.WithFields(fields).WithError(err).
				Error("error getting charge response")
			return nil, err
		}
	}

	var chargeResponse *ChargeResponse
	if err := json.Unmarshal(bodyBytes, &chargeResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error unmarshal charge response")
		return nil, err
	}

	return chargeResponse, nil
}
