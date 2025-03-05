package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Pixs define a interface para operações relacionadas ao serviço de Pix.
type Pixs interface {
	BuildEndpoint(basePath string, queryParams map[string]string, pathParams ...string) (*string, error)

	// CHAVE PIX
	CreatePixKey(ctx context.Context, req PixKeyRequest) (*PixKeyResponse, error)
	GetPixKeys(ctx context.Context, account string) (*PixKeyListResponse, error)
	// REIVIDICAÇÃO CHAVE PIX
	GetPixClaim(ctx context.Context, claimID string) (*PixClaimResponse, error)
	GetPixClaimList(ctx context.Context, dateFrom, dateTo string, limit, page int, status, claimType string) (*PixClaimListResponse, error)
	CancelPixClaim(ctx context.Context, req PixClaimActionRequest) (*PixClaimResponse, error)
	ConfirmPixClaim(ctx context.Context, req PixClaimActionRequest) (*PixClaimResponse, error)
	CreatePixClaim(ctx context.Context, req PixClaimRequest) (*PixClaimResponse, error)
	// IMMEDIATE, STATIC
	GetExternalPixKey(ctx context.Context, account string, key string, ownerTaxId string) (*PixExternalKeyResponse, error)
	// DUEDATE
	GetExternalPixKeyDueDate(ctx context.Context, documentNumberReceiver string, key string) (*PixExternalKeyDueDateResponse, error)
	DeletePixKey(ctx context.Context, account, key string) error

	// CASH OUT - EFETUANDO PAGAMENTO
	PaymentPixCashOut(ctx context.Context, req PixCashOutRequest) (*PixCashOutResponse, error)
	GetPixCashoutStatus(ctx context.Context, id, endtoendId, clientCode string) (*PixCashoutStatusTransactionResponse, error)
	GetPixCashinStatus(ctx context.Context, returnIdentification, transactionId, clientCode string) (*PixCashinStatusTransactionResponse, error)
	PixCashInStatic(ctx context.Context, req PixCashInStaticRequest) (*PixCashInStaticResponse, error)

	// CASH IN - DUE DATE - EMITINDO COBRANÇA COM VENCIMENTO
	CreatePixCashInDueDate(ctx context.Context, req PixCashInDueDateRequest) (*PixCashInDueDateResponse, error)
	GetPixCashInDueDate(ctx context.Context, transactionId *string) (*PixCashInDueDateResponse, error)
	PutPixCashInDueDate(ctx context.Context, transactionId string, req PixCashInDueDateRequest) (*PixCashInDueDateResponse, error)
	DeletePixCashInDueDate(ctx context.Context, transactionId *string) (*PixDeleteResponse, error)

	// CASH IN - DUE DATE - EMITINDO COBRANÇA IMEDIATA
	CreatePixCashInImmediate(ctx context.Context, req PixCashInImmediateRequest) (*PixCashInImmediateResponse, error)
	GetPixCashInImmediate(ctx context.Context, transactionId *string) (*PixCashInImmediateResponse, error)
	PutPixCashInImmediate(ctx context.Context, transactionId string, req PixCashInImmediateRequest) (*PixCashInImmediateResponse, error)
	DeletePixCashInImmediate(ctx context.Context, transactionId *string) (*PixDeleteResponse, error)

	// QRCODE
	GetAddressKey(ctx context.Context, key, currentIdentity, account string, searchDict *bool) (*PixAddressKeyResponse, error)
	DecodeEmvQRCode(ctx context.Context, emv string) (*QRCodeResponse, error)
	GetEmvQRCodeImmediate(ctx context.Context, merchanturl *string) (*QRCodeImmediateResponse, error)
	GetEmvQRCodeDueDate(ctx context.Context, merchanturl *string) (*QRCodeDueDateResponse, error)

	CreateQrCodeLocation(ctx context.Context, req PixQrCodeLocationRequest) (*PixQrCodeLocationResponse, error)
}

// PixsService implementa a interface Pixs.
type PixsService struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewPixs cria uma nova instância de PixsService.
func NewPixs(httpClient *http.Client, session Session) Pixs {
	return &PixsService{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// Método genérico para construir URLs
func (s *PixsService) BuildEndpoint(basePath string, queryParams map[string]string, pathParams ...string) (*string, error) {
	u, err := url.Parse(s.session.APIEndpoint)
	if err != nil {
		logrus.WithError(err).Error("Error parsing API endpoint")
		return nil, err
	}

	// Construção do caminho completo
	fullPath := path.Join(basePath, path.Join(pathParams...))
	u.Path = path.Join(u.Path, fullPath)

	// Adição de query parameters
	if queryParams != nil {
		q := u.Query()
		for key, value := range queryParams {
			if value != "" {
				q.Set(key, value)
			}
		}
		u.RawQuery = q.Encode()
	}

	endpoint := u.String()
	logrus.WithField("endpoint", endpoint).Info("Endpoint built successfully")
	return &endpoint, nil
}

// CreatePixKey cadastra uma nova chave Pix.
func (s *PixsService) CreatePixKey(ctx context.Context, req PixKeyRequest) (*PixKeyResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Create Pix Key")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixDictPath, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for CreatePixKey")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling CreatePixKey")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixKeyResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetPixKeys consulta as chaves Pix de uma conta.
func (s *PixsService) GetPixKeys(ctx context.Context, account string) (*PixKeyListResponse, error) {
	fields := logrus.Fields{"account": account}
	logrus.WithFields(fields).Info("Get Pix Keys")

	endpoint, err := s.BuildEndpoint(PixDictPath, nil, account)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for GetPixKeys")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling GetPixKeys")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixKeyListResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix

}

// DeletePixKey exclui uma chave Pix.
func (s *PixsService) DeletePixKey(ctx context.Context, account, key string) error {
	fields := logrus.Fields{"account": account, "key": key}
	logrus.WithFields(fields).Info("Delete Pix Key")

	if account == "" || key == "" {
		err := fmt.Errorf("account and key are required")
		logrus.WithFields(fields).WithError(err).Error("Invalid input parameters")
		return err
	}

	endpoint, err := s.BuildEndpoint(PixDictPath, nil, key)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for DeletePixKey")
		return err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling DeletePixKey")

	payload := map[string]string{
		"account": account,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing payload")
		return fmt.Errorf("error serializing payload: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", *endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		logrus.WithFields(fields).Info("Pix key deleted successfully")
		return nil
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return err
	}

	return ErrDefaultPix
}

// GetExternalPixKey consulta uma chave Pix externa (DICT).
func (s *PixsService) GetExternalPixKey(ctx context.Context, account string, key string, ownerTaxId string) (*PixExternalKeyResponse, error) {
	fields := logrus.Fields{
		"key":        key,
		"ownerTaxId": ownerTaxId,
		"account":    account,
	}
	logrus.WithFields(fields).Info("Get External Pix Key")

	// Parâmetros para a query string
	params := map[string]string{
		"key":        key,
		"ownerTaxId": ownerTaxId,
	}

	endpoint, err := s.BuildEndpoint(PixDictPath, params, "external", account)

	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixExternalKeyResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetExternalPixKeyDueDate realiza uma consulta POST para o endpoint Celcoin para obter informações sobre uma chave Pix com vencimento(duedate).
func (s *PixsService) GetExternalPixKeyDueDate(ctx context.Context, documentNumberReceiver string, key string) (*PixExternalKeyDueDateResponse, error) {
	fields := logrus.Fields{
		"payerId": documentNumberReceiver,
		"key":     key,
	}
	logrus.WithFields(fields).Info("Get External Pix Key Due Date")

	requestBody := map[string]string{
		"payerId": documentNumberReceiver,
		"key":     key,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error marshaling request body")
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	endpoint, err := s.BuildEndpoint(PixDictDueDatePath, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json-patch+json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixExternalKeyDueDateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// PerformPixCashOut ...
// Realizar um Pix Cash-Out por Chaves Pix
// Realizar um Pix Cash-Out por Agência e Conta
// Realizar um Pix Cash-out por QR Code Estático
// Realizar um Pix Cash-out por QR Code Dinâmico
func (s *PixsService) PaymentPixCashOut(ctx context.Context, req PixCashOutRequest) (*PixCashOutResponse, error) {

	// deixando como upper pois contbank usa parametros minusculos mas é obrigatório na celcoin maiusculo
	req.InitiationType = strings.ToUpper(req.InitiationType)
	req.TransactionType = strings.ToUpper(req.TransactionType)

	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Create Pix PerformPixCashOut")

	if err := validatePixCashOut(req); err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating fields")
		return nil, err
	}

	endpoint, err := s.BuildEndpoint(PixCashOutPath, nil)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashOutResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// DecodeEmvQRCode... Decofificando o qrcode do pix copia e cola
func (s *PixsService) DecodeEmvQRCode(ctx context.Context, emv string) (*QRCodeResponse, error) {
	fields := logrus.Fields{"emv": emv}
	logrus.WithFields(fields).Info("Decoding QR Code")

	endpoint, err := s.BuildEndpoint(PixEmvPath, nil)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload := map[string]string{
		"emv": emv,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(data))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *QRCodeResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetPixCashoutStatus consulta o status de uma transferência Pix-Out.
func (s *PixsService) GetPixCashoutStatus(ctx context.Context, id, endtoendId, clientCode string) (*PixCashoutStatusTransactionResponse, error) {
	fields := logrus.Fields{"id": id, "endtoendId": endtoendId, "clientCode": clientCode}
	logrus.WithFields(fields).Info("Consultando status do Pix Cashout")

	if id == "" && endtoendId == "" && clientCode == "" {
		logrus.WithFields(fields).Error("é necessário informar pelo menos um dos campos: id, endtoendId, ou clientCode")
		return nil, fmt.Errorf("é necessário informar pelo menos um dos campos: id, endtoendId, ou clientCode")
	}

	params := map[string]string{
		"id":         id,
		"endtoendId": endtoendId,
		"clientCode": clientCode,
	}

	endpoint, err := s.BuildEndpoint(PixCashOutPath, params, "status")
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro criando a requisição HTTP")
		return nil, fmt.Errorf("erro criando requisição HTTP: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro na requisição HTTP")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashoutStatusTransactionResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetPixCashinStatus consulta o status de uma devolução Pix (Pix Cash-In).
func (s *PixsService) GetPixCashinStatus(ctx context.Context, returnIdentification, transactionId, clientCode string) (*PixCashinStatusTransactionResponse, error) {
	fields := logrus.Fields{
		"returnIdentification": returnIdentification,
		"transactionId":        transactionId,
		"clientCode":           clientCode,
	}
	logrus.WithFields(fields).Info("Consultando status do Pix Cash-In")

	if returnIdentification == "" || transactionId == "" || clientCode == "" {
		logrus.WithFields(fields).Error("é necessário informar pelo menos um dos campos: returnIdentification, transactionId, ou clientCode")
		return nil, fmt.Errorf("é necessário informar pelo menos um dos campos: returnIdentification, transactionId, ou clientCode")
	}

	params := map[string]string{
		"returnIdentification": returnIdentification,
		"transactionId":        transactionId,
		"clientCode":           clientCode,
	}

	endpoint, err := s.BuildEndpoint(PixCashInStatusPath, params)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro criando a requisição HTTP")
		return nil, fmt.Errorf("erro criando requisição HTTP: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro na requisição HTTP")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashinStatusTransactionResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// PixCashInStatic realiza um Pix Cash-in por Cobrança Estática.
func (s *PixsService) PixCashInStatic(ctx context.Context, req PixCashInStaticRequest) (*PixCashInStaticResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Create PixCashInStatic")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixStaticPath, nil)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashInStaticResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// PixCashInDueDate realiza um Pix Cash-in por Cobrança com Vencimento.
func (s *PixsService) CreatePixCashInDueDate(ctx context.Context, req PixCashInDueDateRequest) (*PixCashInDueDateResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Create PixCashInDueDate")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "duedate")
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashInDueDateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetPixCashInDueDate realiza um Pix Cash-in por Cobrança com Vencimento.
func (s *PixsService) GetPixCashInDueDate(ctx context.Context, transactionId *string) (*PixCashInDueDateResponse, error) {
	fields := logrus.Fields{"transactionId": transactionId}
	logrus.WithFields(fields).Info("Create GetPixCashInDueDate")

	if transactionId == nil {
		logrus.WithFields(fields).Error("Error transactionId is required in request")
		return nil, fmt.Errorf("Error transactionId is required in request")
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "duedate", *transactionId)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashInDueDateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// PutPixCashInDueDate ...
func (s *PixsService) PutPixCashInDueDate(ctx context.Context, transactionId string, req PixCashInDueDateRequest) (*PixCashInDueDateResponse, error) {
	fields := logrus.Fields{"transactionId": transactionId, "request": req}
	logrus.WithFields(fields).Info("Create PutPixCashInDueDate")

	if transactionId == "" {
		logrus.WithFields(fields).Error("Error: transactionId is required in request")
		return nil, fmt.Errorf("transactionId is required in request")
	}

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "duedate", transactionId)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashInDueDateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("Erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("Erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("Erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// DeletePixCashInDueDate remove um Pix Cash-in por Cobrança com Vencimento.
func (s *PixsService) DeletePixCashInDueDate(ctx context.Context, transactionId *string) (*PixDeleteResponse, error) {
	fields := logrus.Fields{"transactionId": transactionId}
	logrus.WithFields(fields).Info("DeletePixCashInDueDate called")

	if transactionId == nil {
		logrus.WithFields(fields).Error("Error: transactionId is required")
		return nil, fmt.Errorf("Error: transactionId is required")
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "duedate", *transactionId)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error reading response body")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusNoContent {
		var response *PixDeleteResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("Error decoding JSON response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error decoding JSON response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).WithFields(fields).WithError(err).Error("Pix service error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// PixCashInImmediate realiza um Pix Cash-in por Cobrança Imediata.
func (s *PixsService) CreatePixCashInImmediate(ctx context.Context, req PixCashInImmediateRequest) (*PixCashInImmediateResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Create PixCashInImmediate")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "immediate")
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashInImmediateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetPixCashInImmediate realiza um Pix Cash-in por Cobrança imediata.
func (s *PixsService) GetPixCashInImmediate(ctx context.Context, transactionId *string) (*PixCashInImmediateResponse, error) {
	fields := logrus.Fields{"transactionId": transactionId}
	logrus.WithFields(fields).Info("Create GetPixCashInImmediate")

	if transactionId == nil {
		logrus.WithFields(fields).Error("Error transactionId is required in request")
		return nil, fmt.Errorf("Error transactionId is required in request")
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "immediate", *transactionId)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashInImmediateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// PutPixCashInImmediate ...
func (s *PixsService) PutPixCashInImmediate(ctx context.Context, transactionId string, req PixCashInImmediateRequest) (*PixCashInImmediateResponse, error) {
	fields := logrus.Fields{"transactionId": transactionId, "request": req}
	logrus.WithFields(fields).Info("Create PutPixCashInImmediate")

	if transactionId == "" {
		logrus.WithFields(fields).Error("Error: transactionId is required in request")
		return nil, fmt.Errorf("transactionId is required in request")
	}

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "immediate", transactionId)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixCashInImmediateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("Erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("Erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("Erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// DeletePixCashInImmediate remove um Pix Cash-in por Cobrança imediata
func (s *PixsService) DeletePixCashInImmediate(ctx context.Context, transactionId *string) (*PixDeleteResponse, error) {
	fields := logrus.Fields{"transactionId": transactionId}
	logrus.WithFields(fields).Info("DeletePixCashInDueDate called")

	if transactionId == nil {
		logrus.WithFields(fields).Error("Error: transactionId is required")
		return nil, fmt.Errorf("Error: transactionId is required")
	}

	endpoint, err := s.BuildEndpoint(PixCashInDynamicPath, nil, "immediate", *transactionId)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error reading response body")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusNoContent {
		var response *PixDeleteResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("Error decoding JSON response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error decoding JSON response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).WithFields(fields).WithError(err).Error("Pix service error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

func (s *PixsService) GetAddressKey(ctx context.Context, key, currentIdentity, account string, searchDict *bool) (*PixAddressKeyResponse, error) {
	var search_dict = false
	if searchDict == nil {
		searchDict = &search_dict // valor for null, não devemos buscar no DICT pois pode afetar o balde de fichas
	}

	// 1. Validar se todos os parâmetros foram enviados
	if key == "" || currentIdentity == "" || account == "" {
		return nil, fmt.Errorf("missing required parameters: key, currentIdentity, or account")
	}

	response := &PixAddressKeyResponse{}

	// Verificar se o searchDict é true e consultar o método GetExternalPixKey
	if searchDict != nil && *searchDict {
		externalPixResponse, err := s.GetExternalPixKey(ctx, account, key, currentIdentity)
		if err != nil {
			return nil, fmt.Errorf("failed to get external pix key: %v", err)
		}
		holder_type := "CPF"
		if len(externalPixResponse.Body.Owner.DocumentNumber) > 11 {
			holder_type = "CNPJ"
		}
		// Atualizar o response com base na resposta do GetExternalPixKey
		response.Status = "FOUND_IN_DICT"
		response.EndToEndID = externalPixResponse.Body.EndToEndId
		response.CreatedAt = externalPixResponse.Body.CreationDate
		response.OwnedAt = externalPixResponse.Body.Account.CreateDate

		response.AddressingKey.Type = PixType(externalPixResponse.Body.KeyType)
		response.AddressingKey.Value = externalPixResponse.Body.Key

		response.Holder.Type = holder_type
		response.Holder.Name = externalPixResponse.Body.Owner.Name
		response.Holder.Document.Type = PixType(holder_type)
		response.Holder.Document.Value = externalPixResponse.Body.Owner.DocumentNumber
	}

	return response, nil
}

// GetEmvQRCodeImmediate decodifica o QR code e faz uma requisição ao endpoint correspondente.
func (s *PixsService) GetEmvQRCodeImmediate(ctx context.Context, merchanturl *string) (*QRCodeImmediateResponse, error) {
	fields := logrus.Fields{"merchantAccountInformation.url": merchanturl}
	logrus.WithFields(fields).Info("Processing GetEmvQRCodeImmediate request")

	if merchanturl == nil {
		err := fmt.Errorf("decoded QR code does not contain a valid URL")
		logrus.WithFields(fields).WithError(err).Error("Invalid decoded QR code")
		return nil, err
	}

	// Transformar a URL no formato esperado (remover https:// e codificar as barras)
	originalURL := *merchanturl
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error parsing URL")
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	encodedPath := url.PathEscape(strings.TrimPrefix(parsedURL.Host+parsedURL.Path, "https://"))
	fields["encoded_url"] = encodedPath
	logrus.WithFields(fields).Info("Encoded URL created successfully")

	// Construir o endpoint para a requisição
	endpoint, err := s.BuildEndpoint(PixEmvUrl, nil, "immediate", "payload", encodedPath)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	// Realizar a requisição GET no endpoint
	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *QRCodeImmediateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetEmvQRCodeDueDate decodifica o QR code e faz uma requisição ao endpoint correspondente para dueDate.
func (s *PixsService) GetEmvQRCodeDueDate(ctx context.Context, merchanturl *string) (*QRCodeDueDateResponse, error) {
	fields := logrus.Fields{"merchantAccountInformation.url": merchanturl}
	logrus.WithFields(fields).Info("Processing GetEmvQRCodeDueDate request")

	if merchanturl == nil {
		err := fmt.Errorf("decoded QR code does not contain a valid URL")
		logrus.WithFields(fields).WithError(err).Error("Invalid decoded QR code")
		return nil, err
	}

	// Transformar a URL no formato esperado (remover https:// e codificar as barras)
	originalURL := *merchanturl
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error parsing URL")
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	encodedPath := url.PathEscape(strings.TrimPrefix(parsedURL.Host+parsedURL.Path, "https://"))
	fields["encoded_url"] = encodedPath
	logrus.WithFields(fields).Info("Encoded URL created successfully")

	// Construir o endpoint para a requisição
	endpoint, err := s.BuildEndpoint(PixEmvUrl, nil, "duedate", "payload", encodedPath)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	// Realizar a requisição GET no endpoint
	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Erro ao ler o corpo da resposta")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *QRCodeDueDateResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("erro ao decodificar a resposta JSON")
			return nil, ErrDefaultPix
		}

		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("erro ao decodificar a resposta JSON")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("erro no serviço Pix")
		return nil, err
	}

	return nil, ErrDefaultPix
}

func (s *PixsService) CreateQrCodeLocation(ctx context.Context, req PixQrCodeLocationRequest) (*PixQrCodeLocationResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.
		WithFields(fields).
		Info("Get QRCode Location")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixQrCodeLocationPath, nil)
	if err != nil {
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error reading response body")
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var response *PixQrCodeLocationResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("Error decoding JSON response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error decoding JSON error response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).WithFields(fields).WithError(err).Error("Pix service error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// CreatePixClaim cadastra um pedido de portabilidade de chave Pix.
func (s *PixsService) CreatePixClaim(ctx context.Context, req PixClaimRequest) (*PixClaimResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Create Pix Claim")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixClaimPath, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for CreatePixClaim")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling CreatePixClaim")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response PixClaimResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("Error decoding json response")
			return nil, ErrDefaultPix
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).WithFields(fields).WithError(err).Error("Celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// ConfirmPixClaim confirma um pedido de portabilidade de chave Pix.
func (s *PixsService) ConfirmPixClaim(ctx context.Context, req PixClaimActionRequest) (*PixClaimResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Confirm Pix Claim")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixClaimPath, nil, "confirm")
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for ConfirmPixClaim")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling ConfirmPixClaim")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixClaimResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// CancelPixClaim Cancelar pedido de portabilidade recebido
func (s *PixsService) CancelPixClaim(ctx context.Context, req PixClaimActionRequest) (*PixClaimResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Confirm Pix Claim")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint, err := s.BuildEndpoint(PixClaimPath, nil, "cancel")
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for ConfirmPixClaim")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling ConfirmPixClaim")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("accept", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixClaimResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}

// GetPixClaim consulta um pedido de portabilidade de chave Pix.
func (s *PixsService) GetPixClaim(ctx context.Context, claimID string) (*PixClaimResponse, error) {
	fields := logrus.Fields{"account": claimID}
	logrus.WithFields(fields).Info("Get Pix Claim")

	endpoint, err := s.BuildEndpoint(PixClaimPath, nil, claimID)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for GetPixClaim")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling GetPixClaim")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixClaimResponse

		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix

}

// GetPixClaimList consulta a lista de pedidos de portabilidade de chave Pix.
func (s *PixsService) GetPixClaimList(ctx context.Context, dateFrom, dateTo string, limit, page int, status, claimType string) (*PixClaimListResponse, error) {
	fields := logrus.Fields{
		"dateFrom":  dateFrom,
		"dateTo":    dateTo,
		"limit":     limit,
		"page":      page,
		"status":    status,
		"claimType": claimType,
	}
	logrus.WithFields(fields).Info("Get Pix Claim List")

	queryParams := map[string]string{
		"DateFrom":     dateFrom,
		"DateTo":       dateTo,
		"LimitPerPage": fmt.Sprintf("%d", limit),
		"Page":         fmt.Sprintf("%d", page),
		"Status":       status,
		"claimType":    claimType,
	}

	endpoint, err := s.BuildEndpoint(PixClaimPath, queryParams)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for GetPixClaimList")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling GetPixClaimList")

	httpReq, err := http.NewRequestWithContext(ctx, "GET", *endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response *PixClaimListResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).Error("error decoding json response")
			return nil, ErrDefaultPix
		}
		return response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).Error("error decoding json response")
		return nil, ErrDefaultPix
	}

	if errResponse.Error != nil {
		err := FindPixError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).WithFields(fields).WithError(err).Error("celcoin get pix error")
		return nil, err
	}

	return nil, ErrDefaultPix
}
