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
)

// Boletos defines the interface:.
type Boletos interface {
	Create(ctx context.Context, req CreateBoletoRequest) (CreateBoletoResponse, error)
	Query(ctx context.Context, transactionID string) (QueryBoletoResponse, error)
	DownloadPDF(ctx context.Context, transactionID string) ([]byte, error)
	Cancel(ctx context.Context, transactionID, reason string) error
}

// boletosService
type boletosService struct {
	session    Session
	httpClient *http.Client
}

// NewBoletos
func NewBoletos(httpClient *http.Client, session Session) Boletos {
	return &boletosService{
		session:    session,
		httpClient: httpClient,
	}
}

// Create calls POST /charge with a CreateBoletoRequest.
func (b *boletosService) Create(ctx context.Context, req CreateBoletoRequest) (CreateBoletoResponse, error) {
	endpoint, err := b.buildEndpoint("charge", nil)
	if err != nil {
		return CreateBoletoResponse{}, err
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Accept 200..202 as success
	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		errBody, _ := io.ReadAll(resp.Body)
		return CreateBoletoResponse{}, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	var raw struct {
		Body struct {
			TransactionID string `json:"transactionId"`
		} `json:"body"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return CreateBoletoResponse{
		TransactionID: raw.Body.TransactionID,
		Status:        raw.Status,
	}, nil
}

// Query calls GET /charge?TransactionId=xxx to retrieve Boleto data.
func (b *boletosService) Query(ctx context.Context, transactionID string) (QueryBoletoResponse, error) {
	queryParams := map[string]string{"TransactionId": transactionID}
	endpoint, err := b.buildEndpoint("charge", queryParams)
	if err != nil {
		return QueryBoletoResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Accept", "application/json")

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return QueryBoletoResponse{}, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	var raw struct {
		Body struct {
			TransactionID string `json:"transactionId"`
			Status        string `json:"status"`
		} `json:"body"`
		Version string `json:"version"`
		Status  string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return QueryBoletoResponse{
		TransactionID: raw.Body.TransactionID,
		Status:        raw.Body.Status,
	}, nil
}

// DownloadPDF calls GET /charge/pdf/:transactionID to download the PDF/HTML content.
func (b *boletosService) DownloadPDF(ctx context.Context, transactionID string) ([]byte, error) {
	// build URL with subpath "charge/pdf/{transactionID}"
	endpoint, err := b.buildEndpoint(path.Join("charge", "pdf", transactionID), nil)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF bytes: %w", err)
	}

	return pdfBytes, nil
}

// Cancel calls DELETE /charge/{transactionID} with a JSON body containing { "reason": "..."}.
func (b *boletosService) Cancel(ctx context.Context, transactionID, reason string) error {
	endpoint, err := b.buildEndpoint(path.Join("charge", transactionID), nil)
	if err != nil {
		return err
	}

	payload := CancelInput{Reason: reason}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := b.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Usually 200 on success
	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}
	return nil
}

// buildEndpoint is a helper method
func (b *boletosService) buildEndpoint(basePath string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(b.session.APIEndpoint)
	if err != nil {
		return "", err
	}

	// join: basePath might be "charge" or "charge/pdf/xxx"
	u.Path = path.Join(u.Path, basePath)

	if queryParams != nil {
		q := u.Query()
		for k, v := range queryParams {
			if v != "" {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}
