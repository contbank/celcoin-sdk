package celcoin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/contbank/grok"
	"github.com/sirupsen/logrus"
)

// Webhooks ...
type Webhooks struct {
	session        Session
	httpClient     *http.Client
	authentication *Authentication
}

// NewWebhooks ...
func NewWebhooks(httpClient *http.Client, session Session) *Webhooks {
	return &Webhooks{
		session:        session,
		httpClient:     httpClient,
		authentication: NewAuthentication(httpClient, session),
	}
}

// CreateSubscription faz a chamada à API para cadastrar um webhook
func (s *Webhooks) CreateSubscription(ctx context.Context, req WebhookSubscriptionRequest) (*WebhookSubscriptionResponse, error) {
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

	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/subscription", s.session.APIEndpoint)

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error serializing request")

		return nil, fmt.Errorf("error serializing request: %v", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}
	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}

	var response WebhookSubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return &response, nil
}

// GetSubscriptions faz a chamada à API para consultar os webhooks cadastrados
func (s *Webhooks) GetSubscriptions(ctx context.Context, entity string, active *bool) (*WebhookQueryResponse, error) {
	fields := logrus.Fields{
		"entity": entity,
		"active": active,
	}
	logrus.
		WithFields(fields).
		Info("Get Subscription")

	queryParams := url.Values{}
	if entity != "" {
		queryParams.Add("entity", entity)
	}
	if active != nil {
		queryParams.Add("active", fmt.Sprintf("%t", *active))
	}
	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/subscription?%s", s.session.APIEndpoint, queryParams.Encode())

	httpReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}

	var response WebhookQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return &response, nil
}

// UpdateSubscription faz a chamada à API para atualizar um webhook existente
func (s *Webhooks) UpdateSubscription(ctx context.Context, req WebhookUpdateRequest) (*WebhookUpdateResponse, error) {
	fields := logrus.Fields{
		"request": req,
	}
	logrus.
		WithFields(fields).
		Info("Update Subscription")

	err := grok.Validator.Struct(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error validating model")
		return nil, grok.FromValidationErros(err)
	}

	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/subscription/%s", s.session.APIEndpoint, req.SubscriptionID)

	payload, err := json.Marshal(req)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error serializing request")
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	httpReq, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição HTTP: %v", err)
	}

	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}

	var response WebhookUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}
	return &response, nil
}

// DeleteSubscription faz a chamada à API para excluir um webhook existente
func (s *Webhooks) DeleteSubscription(ctx context.Context, subscriptionID string) (*WebhookDeleteResponse, error) {
	fields := logrus.Fields{
		"subscription_id": subscriptionID,
	}
	logrus.
		WithFields(fields).
		Info("Delete Subscription")

	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/subscription/pix-payment-out?SubscriptionId=%s", s.session.APIEndpoint, subscriptionID)

	httpReq, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}
	var response WebhookDeleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return &response, nil
}

// GetWebhookReplayCount realiza a consulta de quantidade de webhooks enviados
func (s *Webhooks) GetWebhookReplayCount(ctx context.Context, entity, dateFrom, dateTo string, optionalParams map[string]string) (*WebhookReplayResponse, error) {
	fields := logrus.Fields{
		"entity":          entity,
		"date_from":       dateFrom,
		"date_to":         dateTo,
		"optional_params": optionalParams,
	}
	logrus.
		WithFields(fields).
		Info("Get Webhook Replay Count")

	if entity == "" || dateFrom == "" || dateTo == "" {
		logrus.
			WithFields(fields).
			Error("the parameters (entity, dateFrom, dateTo) must be provided")
		return nil, errors.New("the parameters (entity, dateFrom, dateTo) must be provided")
	}

	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/replay/%s?DateFrom=%s&DateTo=%s", s.session.APIEndpoint, entity, url.QueryEscape(dateFrom), url.QueryEscape(dateTo))

	if optionalParams != nil {
		queryParams := url.Values{}
		for key, value := range optionalParams {
			queryParams.Add(key, value)
		}
		endpoint += "&" + queryParams.Encode()
	}

	httpReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}
	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}

	var response WebhookReplayResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return &response, nil
}

// GetWebhookReplay realiza a consulta para recuperar os detalhes dos webhooks enviados
func (s *Webhooks) GetWebhookReplay(ctx context.Context, entity, dateFrom, dateTo string, onlyPending bool) (*WebhookReplayResponse, error) {
	fields := logrus.Fields{
		"entity":       entity,
		"date_from":    dateFrom,
		"date_to":      dateTo,
		"only_pending": onlyPending,
	}
	logrus.
		WithFields(fields).
		Info("Get Webhook Replay")

	if entity == "" || dateFrom == "" || dateTo == "" {
		logrus.
			WithFields(fields).
			Error("the parameters (entity, dateFrom, dateTo) must be provided")
		return nil, errors.New("the parameters (entity, dateFrom, dateTo) must be provided")
	}

	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/replay/%s?DateFrom=%s&DateTo=%s&OnlyPending=%t", s.session.APIEndpoint, entity, url.QueryEscape(dateFrom), url.QueryEscape(dateTo), onlyPending)

	httpReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}
	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}

	var response WebhookReplayResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return &response, nil
}

// GetWebhookReplaySendCount realiza a consulta para recuperar a quantidade de webhooks enviados
func (s *Webhooks) GetWebhookReplaySendCount(ctx context.Context, entity, dateFrom, dateTo string) (*WebhookReplayCountResponse, error) {
	fields := logrus.Fields{
		"entity":    entity,
		"date_from": dateFrom,
		"date_to":   dateTo,
	}
	logrus.
		WithFields(fields).
		Info("Get Webhook Replay")

	if entity == "" || dateFrom == "" || dateTo == "" {
		logrus.
			WithFields(fields).
			Error("the parameters (entity, dateFrom, dateTo) must be provided")
		return nil, errors.New("the parameters (entity, dateFrom, dateTo) must be provided")
	}

	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/replay/%s/count?DateFrom=%s&DateTo=%s", s.session.APIEndpoint, entity, url.QueryEscape(dateFrom), url.QueryEscape(dateTo))

	httpReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}
	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}
	var response WebhookReplayCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return &response, nil
}

// ReplayMessageFromWebhook reenvia o webhook com base nos parâmetros fornecidos
func (s *Webhooks) ReplayMessageFromWebhook(ctx context.Context, entity, webhookID, dateFrom, dateTo string, onlyPending bool, filter WebhookReplayRequest) (*WebhookReplayResponse, error) {

	fields := logrus.Fields{
		"entity":     entity,
		"webhook_id": webhookID,
		"date_from":  dateFrom,
		"date_to":    dateTo,
	}
	logrus.
		WithFields(fields).
		Info("Replay Message from Webhook")

	if entity == "" || webhookID == "" {
		logrus.
			WithFields(fields).
			Error("the parameters (entity, webhookID) must be provided")
		return nil, errors.New("the parameters (entity, webhookID) must be provided")
	}

	endpoint := fmt.Sprintf("%s/baas-webhookmanager/v1/webhook/replay/%s?webhookId=%s", s.session.APIEndpoint, entity, webhookID)
	if dateFrom != "" {
		endpoint += "&DateFrom=" + url.QueryEscape(dateFrom)
	}
	if dateTo != "" {
		endpoint += "&DateTo=" + url.QueryEscape(dateTo)
	}
	if onlyPending {
		endpoint += "&OnlyPending=true"
	}

	httpReq, err := http.NewRequest("PUT", endpoint, nil)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error creating http request")
		return nil, fmt.Errorf("error creating http request: %v", err)
	}
	token, err := s.authentication.Token(ctx)
	if err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error authentication")
		return nil, err
	}

	httpReq.Header.Add("Authorization", token)
	httpReq.Header.Add("Content-type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error http client status code")
		return nil, errors.New("error http client status code")
	}
	var response WebhookReplayResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logrus.
			WithFields(fields).
			WithError(err).
			Error("error unmarshal")
		return nil, err
	}

	return &response, nil
}
