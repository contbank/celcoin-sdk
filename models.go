package celcoin

import "time"

const (
	// LoginPath ...
	LoginPath string = "v5/token"
	// LoginMtlsPath ...
	LoginMtlsPath string = "v5/token"
)

// AuthenticationResponse ...
type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Certificate ...
type Certificate struct {
	Certificate      string `json:"certificate"`
	CertificateChain string `json:"certificateChain"`
	SubjectDn        string `json:"subjectDn"`
	PrivateKey       string `json:"privateKey"`
	Passphrase       string `json:"passphrase"`
	UUID             string `json:"uuid"`
	ClientID         string `json:"client_id"`
}

// ErrorLoginResponse ...
type ErrorLoginResponse struct {
	Message string `json:"error"`
}

// ErrorModel ...
type ErrorModel struct {
	Code         string   `json:"code,omitempty"`
	PropertyName string   `json:"propertyName,omitempty"`
	Messages     []string `json:"messages,omitempty"`
	KeyValueErrorModel
}

// KeyValueErrorModel ...
type KeyValueErrorModel struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// CodeMessageErrorResponse ...
type CodeMessageErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// TransferErrorResponse ...
type TransferErrorResponse struct {
	Layer           string               `json:"layer,omitempty"`
	ApplicationName string               `json:"applicationName,omitempty"`
	Errors          []KeyValueErrorModel `json:"errors,omitempty"`
	CodeMessageErrorResponse
}

/* WEBHOOK MODELS */
// WebhookSubscriptionRequest representa o payload para cadastrar e gerenciar webhooks
type WebhookSubscriptionRequest struct {
	Entity     string      `json:"entity"`     // Identificador do evento
	WebhookURL string      `json:"webhookUrl"` // URL do webhook
	Auth       WebhookAuth `json:"auth"`       // Dados de autenticação do webhook
}

// WebhookAuth representa os dados de autenticação para o webhook
type WebhookAuth struct {
	Login    string `json:"login"` // Login para autenticação
	Password string `json:"pwd"`   // Senha para autenticação
	Type     string `json:"type"`  // Tipo de autenticação (basic)
}

// WebhookSubscriptionResponse representa a resposta da rota de cadastro de webhooks
type WebhookSubscriptionResponse struct {
	Version string        `json:"version"`         // Versão da API
	Status  string        `json:"status"`          // Status da operação (SUCCESS ou ERROR)
	Error   *WebhookError `json:"error,omitempty"` // Detalhes do erro, se houver
}

// WebhookError representa informações de erro em uma resposta da API
type WebhookError struct {
	ErrorCode string `json:"errorCode"` // Código do erro
	Message   string `json:"message"`   // Mensagem do erro
}

// WebhookSubscription representa os dados retornados na consulta de webhooks cadastrados
type WebhookSubscription struct {
	Entity         string      `json:"entity"`
	WebhookURL     string      `json:"webhookUrl"`
	Active         bool        `json:"active"`
	CreateDate     time.Time   `json:"createDate"`
	LastUpdateDate time.Time   `json:"lastUpdateDate"`
	Auth           WebhookAuth `json:"auth"`
}

// WebhookQueryResponse representa a resposta da API ao consultar webhooks cadastrados
type WebhookQueryResponse struct {
	Version string                `json:"version"`
	Status  string                `json:"status"`
	Body    []WebhookSubscription `json:"body,omitempty"`
	Error   *WebhookError         `json:"error,omitempty"`
}

// WebhookUpdateRequest representa o payload para atualizar um webhook existente
type WebhookUpdateRequest struct {
	WebhookURL     string      `json:"webhookUrl"`     // URL do webhook
	Auth           WebhookAuth `json:"auth"`           // Dados de autenticação do webhook
	Active         bool        `json:"active"`         // Status de ativação do webhook
	SubscriptionID string      `json:"subscriptionId"` // ID da assinatura a ser atualizada
}

// WebhookUpdateResponse representa a resposta da API ao atualizar um webhook
type WebhookUpdateResponse struct {
	Version string        `json:"version"`         // Versão da API
	Status  string        `json:"status"`          // Status da operação (SUCCESS ou ERROR)
	Error   *WebhookError `json:"error,omitempty"` // Detalhes do erro, se houver
}

// WebhookDeleteResponse representa a resposta da API ao excluir um webhook
type WebhookDeleteResponse struct {
	Version string        `json:"version"`         // Versão da API
	Status  string        `json:"status"`          // Status da operação (SUCCESS ou ERROR)
	Error   *WebhookError `json:"error,omitempty"` // Detalhes do erro, se houver
}

// WebhookReplayResponseBody representa o corpo da resposta da consulta de webhooks enviados
type WebhookReplayResponseBody struct {
	OnlyPending bool   `json:"onlyPending,omitempty"` // Indica se somente os pendentes foram consultados
	Entity      string `json:"entity"`                // Identificador do evento
	DateFrom    string `json:"dateFrom"`              // Data inicial da consulta
	DateTo      string `json:"dateTo"`                // Data final da consulta
	TotalItems  int    `json:"totalItems"`            // Total de itens encontrados
}

// WebhookReplayResponse representa a resposta completa da API para consulta de webhooks
type WebhookReplayResponse struct {
	Body    WebhookReplayResponseBody `json:"body"`            // Corpo da resposta
	Status  string                    `json:"status"`          // Status da operação (SUCCESS ou ERROR)
	Version string                    `json:"version"`         // Versão da API
	Error   *WebhookError             `json:"error,omitempty"` // Detalhes do erro, se houver
}

// WebhookReplayCountRequest representa os parâmetros para consulta da quantidade de webhooks enviados
type WebhookReplayCountRequest struct {
	Entity   string `json:"entity"`   // Identificador do evento
	DateFrom string `json:"dateFrom"` // Data inicial da consulta
	DateTo   string `json:"dateTo"`   // Data final da consulta
}

// WebhookReplayCountResponse representa a resposta da consulta da quantidade de webhooks enviados
type WebhookReplayCountResponse struct {
	TotalItems int           `json:"totalItems"`      // Total de itens encontrados
	Status     string        `json:"status"`          // Status da operação (SUCCESS ou ERROR)
	Version    string        `json:"version"`         // Versão da API
	Error      *WebhookError `json:"error,omitempty"` // Detalhes do erro, se houver
}

// WebhookReplayRequest representa a estrutura do corpo da requisição para reenvio do webhook
type WebhookReplayRequest struct {
	Filter WebhookReplayFilter `json:"filter"`
}

// WebhookReplayFilter representa os filtros que podem ser aplicados ao reenvio do webhook
type WebhookReplayFilter struct {
	DocumentNumber  string `json:"documentNumber,omitempty"`  // Número do documento
	Account         string `json:"account,omitempty"`         // Conta associada ao webhook
	ID              string `json:"id,omitempty"`              // ID do webhook
	ClientRequestID string `json:"clientRequestId,omitempty"` // ID da solicitação do cliente
}
