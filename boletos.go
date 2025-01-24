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

	"github.com/sirupsen/logrus"
)

type Boletos struct {
	session    Session
	httpClient *http.Client
}

// NewBoletos...
func NewBoletos(httpClient *http.Client, session Session) *Boletos {
	return &Boletos{
		session:    session,
		httpClient: httpClient,
	}
}

// Create creates a new Celcoin charge/boleto (POST /charge).
func (b *Boletos) Create(ctx context.Context, req CreateBoletoRequest) (CreateBoletoResponse, error) {
	fields := logrus.Fields{"method": "Create", "request": req}
	logrus.WithFields(fields).Info("Celcoin Boleto - Create")

	endpoint, err := b.buildEndpoint("charge", nil)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("failed to build endpoint for create")
		return CreateBoletoResponse{}, err
	}

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("marshal request error")
		return CreateBoletoResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("failed to create request")
		return CreateBoletoResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	resp, err := b.httpClient.Do(request)
	if err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Accept 200/201/202 as success
	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return CreateBoletoResponse{}, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Typically returns:
	// {
	//   "body": { "transactionId": "318f1f31-3553-4adf-b6be-cf6dc817e26c" },
	//   "version": "1.2.0",
	//   "status": "SUCCESS"
	// }
	var rawResp struct {
		Body struct {
			TransactionID string `json:"transactionId"`
		} `json:"body"`
		Version string `json:"version"`
		Status  string `json:"status"` // e.g. "SUCCESS"
	}

	if err := json.NewDecoder(resp.Body).Decode(&rawResp); err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return CreateBoletoResponse{
		TransactionID: rawResp.Body.TransactionID,
		Status:        rawResp.Status,
	}, nil
}

// Query fetches a boleto by transactionID (GET /charge?TransactionId=xxx).
func (b *Boletos) Query(ctx context.Context, transactionID string) (QueryBoletoResponse, error) {
	fields := logrus.Fields{"method": "Query", "transactionID": transactionID}
	logrus.WithFields(fields).Info("Celcoin Boleto - Query")

	queryParams := map[string]string{
		"TransactionId": transactionID,
	}
	endpoint, err := b.buildEndpoint("charge", queryParams)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("failed to build endpoint for query")
		return QueryBoletoResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return QueryBoletoResponse{}, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	// Typically returns:
	// {
	//   "body": {
	//     "transactionId":"bc47586f-6834-4581-88c8-40965b3e9307",
	//     "status":"PENDING",
	//     ...
	//   },
	//   "version":"1.2.0",
	//   "status":"SUCCESS"
	// }
	var rawResp struct {
		Body struct {
			TransactionID string `json:"transactionId"`
			Status        string `json:"status"`
		} `json:"body"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rawResp); err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return QueryBoletoResponse{
		TransactionID: rawResp.Body.TransactionID,
		Status:        rawResp.Body.Status,
	}, nil
}

// DownloadPDF...
func (b *Boletos) DownloadPDF(ctx context.Context, transactionID string) ([]byte, error) {
	fields := logrus.Fields{"method": "DownloadPDF", "transactionID": transactionID}
	logrus.WithFields(fields).Info("Celcoin Boleto - Download PDF")

	endpoint, err := b.buildEndpoint("charge/pdf/"+transactionID, nil)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("failed to build endpoint for pdf")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}


	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF bytes: %w", err)
	}

	return pdfBytes, nil
}

// Cancel a Celcoin charge...
func (b *Boletos) Cancel(ctx context.Context, transactionID, reason string) error {
	fields := logrus.Fields{"method": "Cancel", "transactionID": transactionID}
	logrus.WithFields(fields).Info("Celcoin Boleto - Cancel")

	endpoint, err := b.buildEndpoint("charge/"+transactionID, nil)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("failed to build endpoint for cancel")
		return err
	}

	payload := CancelInput{Reason: reason}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	resp, err := b.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Celcoin returns 200 for a successful cancel.
	// Typically returns:
	// {
	//   "body": {...},
	//   "version": "1.2.0",
	//   "status": "PROCESSING" (or "SUCCESS")
	// }
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (b *Boletos) buildEndpoint(pathStr string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(b.session.APIEndpoint)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, pathStr)

	if queryParams != nil {
		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
