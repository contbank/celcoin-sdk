package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Pixs define a interface para operações relacionadas ao serviço de Pix.
type Pixs interface {
	CreatePixKey(ctx context.Context, req PixKeyRequest) (*PixKeyResponse, error)
	GetPixKeys(ctx context.Context, account string) (*PixKeyListResponse, error)
	GetExternalPixKey(ctx context.Context, key string, ownerTaxId string) (*PixExternalKeyResponse, error)
	DeletePixKey(ctx context.Context, account, key string) error
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

// CreatePixKey cadastra uma nova chave Pix.
func (s *PixsService) CreatePixKey(ctx context.Context, req PixKeyRequest) (*PixKeyResponse, error) {
	fields := logrus.Fields{"request": req}
	logrus.WithFields(fields).Info("Create Pix Key")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint := fmt.Sprintf("%s/celcoin-baas-pix-dict-webservice/v1/pix/dict/entry", s.session.APIEndpoint)

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
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

	// Construir o endpoint da API
	endpoint := fmt.Sprintf("%s/celcoin-baas-pix-dict-webservice/v1/pix/dict/entry/%s", s.session.APIEndpoint, account)

	// Criar a requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Obter o token de autenticação
	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in authentication")
		return nil, fmt.Errorf("error obtaining authentication token: %v", err)
	}

	// Adicionar cabeçalhos
	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")

	// Fazer a requisição HTTP
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Ler o corpo da resposta
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
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

	endpoint := fmt.Sprintf("%s/celcoin-baas-pix-dict-webservice/v1/pix/dict/entry/%s", s.session.APIEndpoint, key)
	payload := map[string]string{
		"account": account,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error serializing payload")
		return fmt.Errorf("error serializing payload: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in authentication")
		return err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
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
func (s *PixsService) GetExternalPixKey(ctx context.Context, key string, ownerTaxId string) (*PixExternalKeyResponse, error) {
	fields := logrus.Fields{
		"key":        key,
		"ownerTaxId": ownerTaxId,
	}
	logrus.WithFields(fields).Info("Get External Pix Key")

	// Endpoint de consulta de chaves Pix externas
	endpoint := fmt.Sprintf("%s/celcoin-baas-pix-dict-webservice/v1/pix/dict/entry/external/%s?key=%s&ownerTaxId=%s", s.session.APIEndpoint, ownerTaxId, key, ownerTaxId)

	// Configuração da requisição HTTP
	httpReq, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error creating HTTP request")
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Adicionando o token de autenticação
	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in authentication")
		return nil, err
	}
	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("accept", "application/json")
	httpReq.Header.Add("Content-Type", "application/json")

	// Executando a requisição
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error in HTTP client")
		return nil, err
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
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
