package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Boleto...
type Boleto struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewBoleto ...
func NewBoleto(baseURL, apiKey string) *Boleto {
	return &Boleto{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  &http.Client{},
	}
}

// CreateBoletoRequest is the payload for creating a new boleto.
type CreateBoletoRequest struct {
	ExternalID            string     `json:"externalId"`
	ExpirationAfterPayment int       `json:"expirationAfterPayment"`
	DueDate               string     `json:"dueDate"`
	Amount                float64    `json:"amount"`
	Key                   string     `json:"key,omitempty"` // optional 
	Debtor                Debtor     `json:"debtor"`
	Receiver              Receiver   `json:"receiver"`
	Instructions          Instructions `json:"instructions"`
}

type Debtor struct {
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	Name         string `json:"name"`
	Document     string `json:"document"`
	City         string `json:"city"`
	PublicArea   string `json:"publicArea"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
}

type Receiver struct {
	Account  string `json:"account"`
	Document string `json:"document"`
}

type Instructions struct {
	Fine     float64  `json:"fine"`
	Interest float64  `json:"interest"`
	Discount Discount `json:"discount"`
}

type Discount struct {
	Amount    float64 `json:"amount"`
	Modality  string  `json:"modality"`  // "fixed" or "percent"
	LimitDate string  `json:"limitDate"` // e.g. "2025-01-20T00:00:00.0000000"
}

// CreateBoletoResponse 
type CreateBoletoResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

// Create creates a new boleto in Celcoin.
func (b *Boleto) Create(ctx context.Context, req CreateBoletoRequest) (CreateBoletoResponse, error) {
	url := fmt.Sprintf("%s/charge", b.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+b.APIKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	resp, err := b.Client.Do(request)
	if err != nil {
		return CreateBoletoResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Accept 200 / 201 / 202 as "success"
	if resp.StatusCode < 200 || resp.StatusCode > 202 {
		errBody, _ := io.ReadAll(resp.Body)
		return CreateBoletoResponse{}, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	// Example success response:
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

// QueryBoletoResponse is the simplified struct for a GET response ...
type QueryBoletoResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

// Query fetches a boleto by transaction ID ...
func (b *Boleto) Query(ctx context.Context, transactionID string) (QueryBoletoResponse, error) {
	url := fmt.Sprintf("%s/charge?TransactionId=%s", b.BaseURL, transactionID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+b.APIKey)
	request.Header.Set("Accept", "application/json")

	resp, err := b.Client.Do(request)
	if err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return QueryBoletoResponse{}, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	// Example query response:
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
		Status  string `json:"status"` // "SUCCESS", "ERROR", etc.
	}

	if err := json.NewDecoder(resp.Body).Decode(&rawResp); err != nil {
		return QueryBoletoResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return QueryBoletoResponse{
		TransactionID: rawResp.Body.TransactionID,
		Status:        rawResp.Body.Status,
	}, nil
}

// DownloadPDF ...
func (b *Boleto) DownloadPDF(ctx context.Context, transactionID string) ([]byte, error) {
	url := fmt.Sprintf("%s/charge/pdf/%s", b.BaseURL, transactionID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+b.APIKey)

	resp, err := b.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	// The PDF or HTML bytes:
	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF bytes: %w", err)
	}

	return pdfBytes, nil
}

// CancelInput is the JSON body for cancel requests.
type CancelInput struct {
	Reason string `json:"reason"`
}

// Cancel cancels a charge with the given transactionID and reason.
func (b *Boleto) Cancel(ctx context.Context, transactionID, reason string) error {
	url := fmt.Sprintf("%s/charge/%s", b.BaseURL, transactionID)

	payload := CancelInput{Reason: reason}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+b.APIKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	resp, err := b.Client.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Celcoin returns 200 for a successful cancel.
	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	// Typically returns:
	// {
	//   "body": {...},
	//   "version": "1.2.0",
	//   "status": "PROCESSING"  (or "SUCCESS")
	// }

	return nil
}
