package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Webhooks define a interface para operações relacionadas a webhooks.
type Webhooks interface {
	CreateSubscription(ctx context.Context, req WebhookSubscriptionRequest) (*WebhookSubscriptionResponse, error)
	CreateSubscriptionDda(ctx context.Context, req WebhookSubscriptionDdaRequest) (*WebhookSubscriptionDdaResponse, error)
	GetSubscriptions(ctx context.Context, entity string, active *bool) (*WebhookQueryResponse, error)
	UpdateSubscription(ctx context.Context, entity string, req WebhookUpdateRequest) (*WebhookUpdateResponse, error)
	DeleteSubscription(ctx context.Context, entity string, subscriptionID string) (*WebhookDeleteResponse, error)
	GetWebhookReplayCount(ctx context.Context, entity, dateFrom, dateTo string, optionalParams map[string]string) (*WebhookReplayResponse, error)
	GetWebhookReplay(ctx context.Context, entity, dateFrom, dateTo string, onlyPending bool) (*WebhookReplayResponse, error)
	GetWebhookReplaySendCount(ctx context.Context, entity, dateFrom, dateTo string) (*WebhookReplayCountResponse, error)
	ReplayMessageFromWebhook(ctx context.Context, entity, webhookID, dateFrom, dateTo string, onlyPending bool, filter WebhookReplayRequest) (*WebhookReplayResponse, error)
}

// Webhooks ...
type WebhooksService struct {
	session        Session
	httpClient     *LoggingHTTPClient
	authentication *Authentication
}

// NewWebhooks cria uma nova instância de WebhooksService.
func NewWebhooks(httpClient *http.Client, session Session) Webhooks {
	return &WebhooksService{
		session:        session,
		httpClient:     NewLoggingHTTPClient(httpClient),
		authentication: NewAuthentication(httpClient, session),
	}
}

// Método genérico para construir URLs
func (s *WebhooksService) buildEndpoint(basePath string, queryParams map[string]string, pathParams ...string) (*string, error) {
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

// CreateSubscription faz a chamada à API para cadastrar um webhook
func (s *WebhooksService) CreateSubscription(ctx context.Context, req WebhookSubscriptionRequest) (*WebhookSubscriptionResponse, error) {
	fields := logrus.Fields{
		"request": req,
	}
	logrus.
		WithFields(fields).
		Info("Create Subscription")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}
	endpoint, err := s.buildEndpoint(WebhookPath, nil, "subscription")
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for CreateSubscription")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling CreatePixKey")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error serializing request")

		return nil, fmt.Errorf("error serializing request: %v", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("Content-Type", "application/json")
	//httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookSubscriptionResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook

}

// CreateSubscriptionDda faz a chamada à API para cadastrar um webhook dda
func (s *WebhooksService) CreateSubscriptionDda(ctx context.Context, req WebhookSubscriptionDdaRequest) (*WebhookSubscriptionDdaResponse, error) {
	fields := logrus.Fields{
		"request": req,
	}
	logrus.
		WithFields(fields).
		Info("Create Subscription")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}
	endpoint, err := s.buildEndpoint(WebhookDdaPath, nil, "register")
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint for CreateSubscriptionDda")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Calling CreateSubscriptionDda")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error serializing request")

		return nil, fmt.Errorf("error serializing request: %v", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", *endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("Content-Type", "application/json")
	//httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusCreated {
		var response WebhookSubscriptionDdaResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook

}

// GetSubscriptions faz a chamada à API para consultar os webhooks cadastrados
func (s *WebhooksService) GetSubscriptions(ctx context.Context, entity string, active *bool) (*WebhookQueryResponse, error) {
	fields := logrus.Fields{
		"entity": entity,
		"active": active,
	}
	logrus.WithFields(fields).Info("Get Subscription")

	// Configuração dos parâmetros da query string
	queryParams := map[string]string{}
	if entity != "" {
		queryParams["entity"] = entity
	}
	if active != nil {
		queryParams["active"] = fmt.Sprintf("%t", *active)
	}

	// Construção do endpoint com o método buildEndpoint
	endpoint, err := s.buildEndpoint(WebhookPath, queryParams, "subscription")
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookQueryResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin get balance error")
		return nil, err
	}

	return nil, ErrDefaultWebhook
}

// UpdateSubscription faz a chamada à API para atualizar um webhook existente
func (s *WebhooksService) UpdateSubscription(ctx context.Context, entity string, req WebhookUpdateRequest) (*WebhookUpdateResponse, error) {
	fields := logrus.Fields{
		"entity":  entity,
		"request": req,
	}
	logrus.WithFields(fields).Info("Update Subscription")

	// Validação do modelo de requisição
	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("Error validating model")
		return nil, grok.FromValidationErros(err)
	}

	// Construção do endpoint com buildEndpoint
	endpoint, err := s.buildEndpoint(WebhookPath, nil, "subscription", entity)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequest("PUT", *endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição HTTP: %v", err)
	}

	httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookUpdateResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook
}

// DeleteSubscription faz a chamada à API para excluir um webhook existente
func (s *WebhooksService) DeleteSubscription(ctx context.Context, entity string, subscriptionID string) (*WebhookDeleteResponse, error) {
	fields := logrus.Fields{
		"entity":          entity,
		"subscription_id": subscriptionID,
	}
	logrus.WithFields(fields).Info("Delete Subscription")

	// Parâmetros de consulta (query params)
	queryParams := map[string]string{
		"SubscriptionId": subscriptionID,
	}

	// Construção do endpoint com buildEndpoint
	endpoint, err := s.buildEndpoint(WebhookPath, queryParams, "subscription", entity)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequest("DELETE", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookDeleteResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook

}

// GetWebhookReplayCount realiza a consulta de quantidade de webhooks enviados
func (s *WebhooksService) GetWebhookReplayCount(ctx context.Context, entity, dateFrom, dateTo string, optionalParams map[string]string) (*WebhookReplayResponse, error) {
	fields := logrus.Fields{
		"entity":          entity,
		"date_from":       dateFrom,
		"date_to":         dateTo,
		"optional_params": optionalParams,
	}
	logrus.WithFields(fields).Info("Get Webhook Replay Count")

	// Validação dos parâmetros obrigatórios
	if entity == "" || dateFrom == "" || dateTo == "" {
		logrus.
			WithFields(fields).
			Error("the parameters (entity, dateFrom, dateTo) must be provided")
		return nil, errors.New("the parameters (entity, dateFrom, dateTo) must be provided")
	}

	// Adiciona os parâmetros obrigatórios e opcionais como query params
	queryParams := map[string]string{
		"DateFrom": url.QueryEscape(dateFrom),
		"DateTo":   url.QueryEscape(dateTo),
	}

	// Inclui parâmetros opcionais, se presentes
	for key, value := range optionalParams {
		if value != "" {
			queryParams[key] = value
		}
	}

	// Construção do endpoint utilizando buildEndpoint
	endpoint, err := s.buildEndpoint(WebhookPath, queryParams, "replay", entity)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookReplayResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook
}

// GetWebhookReplay realiza a consulta para recuperar os detalhes dos webhooks enviados
func (s *WebhooksService) GetWebhookReplay(ctx context.Context, entity, dateFrom, dateTo string, onlyPending bool) (*WebhookReplayResponse, error) {
	fields := logrus.Fields{
		"entity":       entity,
		"date_from":    dateFrom,
		"date_to":      dateTo,
		"only_pending": onlyPending,
	}
	logrus.WithFields(fields).Info("Get Webhook Replay")

	// Validação dos parâmetros obrigatórios
	if entity == "" || dateFrom == "" || dateTo == "" {
		logrus.
			WithFields(fields).
			Error("the parameters (entity, dateFrom, dateTo) must be provided")
		return nil, errors.New("the parameters (entity, dateFrom, dateTo) must be provided")
	}

	// Construção dos parâmetros de consulta
	queryParams := map[string]string{
		"DateFrom":    url.QueryEscape(dateFrom),
		"DateTo":      url.QueryEscape(dateTo),
		"OnlyPending": fmt.Sprintf("%t", onlyPending),
	}

	// Construção do endpoint usando buildEndpoint
	endpoint, err := s.buildEndpoint(WebhookPath, queryParams, "replay", entity)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookReplayResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook
}

// GetWebhookReplaySendCount realiza a consulta para recuperar a quantidade de webhooks enviados
func (s *WebhooksService) GetWebhookReplaySendCount(ctx context.Context, entity, dateFrom, dateTo string) (*WebhookReplayCountResponse, error) {
	fields := logrus.Fields{
		"entity":    entity,
		"date_from": dateFrom,
		"date_to":   dateTo,
	}
	logrus.WithFields(fields).Info("Get Webhook Replay Send Count")

	// Validação dos parâmetros obrigatórios
	if entity == "" || dateFrom == "" || dateTo == "" {
		logrus.
			WithFields(fields).
			Error("the parameters (entity, dateFrom, dateTo) must be provided")
		return nil, errors.New("the parameters (entity, dateFrom, dateTo) must be provided")
	}

	// Construção dos parâmetros de consulta
	queryParams := map[string]string{
		"DateFrom": url.QueryEscape(dateFrom),
		"DateTo":   url.QueryEscape(dateTo),
	}

	// Construção do endpoint usando buildEndpoint
	endpoint, err := s.buildEndpoint(WebhookPath, queryParams, "replay", entity, "count")
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequest("GET", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookReplayCountResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook
}

// ReplayMessageFromWebhook reenvia o webhook com base nos parâmetros fornecidos
func (s *WebhooksService) ReplayMessageFromWebhook(ctx context.Context, entity, webhookID, dateFrom, dateTo string, onlyPending bool, filter WebhookReplayRequest) (*WebhookReplayResponse, error) {
	fields := logrus.Fields{
		"entity":       entity,
		"webhook_id":   webhookID,
		"date_from":    dateFrom,
		"date_to":      dateTo,
		"only_pending": onlyPending,
	}
	logrus.WithFields(fields).Info("Replay Message from Webhook")

	// Validação dos parâmetros obrigatórios
	if entity == "" || webhookID == "" {
		logrus.WithFields(fields).Error("the parameters (entity, webhookID) must be provided")
		return nil, errors.New("the parameters (entity, webhookID) must be provided")
	}

	// Criação dos parâmetros de consulta
	queryParams := map[string]string{
		"webhookId":   webhookID,
		"DateFrom":    url.QueryEscape(dateFrom),
		"DateTo":      url.QueryEscape(dateTo),
		"OnlyPending": fmt.Sprintf("%t", onlyPending),
	}

	// Construção do endpoint utilizando buildEndpoint
	endpoint, err := s.buildEndpoint(WebhookPath, queryParams, "replay", entity)
	if err != nil {
		logrus.WithFields(fields).WithError(err).Error("Error building endpoint")
		return nil, err
	}

	logrus.WithField("endpoint", *endpoint).Info("Endpoint built successfully")

	httpReq, err := http.NewRequest("PUT", *endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	httpReq.Header.Add("api-version", s.session.APIVersion)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client")
		return nil, err
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		var response WebhookReplayResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			logrus.WithFields(fields).WithError(err).
				Error("error decoding json response")
			return nil, ErrDefaultWebhook
		}
		return &response, nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrEntryNotFound
	}

	var errResponse *ErrorDefaultResponse
	if err := json.Unmarshal(respBody, &errResponse); err != nil {
		logrus.WithFields(fields).WithError(err).
			Error("error decoding json response")
		return nil, ErrDefaultWebhook
	}

	if errResponse.Error != nil {
		err := FindWebhookError(*errResponse.Error.ErrorCode, &resp.StatusCode)
		logrus.WithField("celcoin_error", errResponse.Error).
			WithFields(fields).WithError(err).
			Error("celcoin create subscription error")
		return nil, err
	}

	return nil, ErrDefaultWebhook
}
