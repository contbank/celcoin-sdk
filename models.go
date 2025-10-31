package celcoin

import (
	"encoding/json"
	"time"
)

type CompanyType string
type PersonType string

const (
	// CelcoinBankCode ...
	CelcoinBankCode string = "509"
	// CelcoinBankISPB ...
	CelcoinBankISPB string = "13935893"
	// CelcoinBankName ...
	CelcoinBankName string = "Celcoin Instituição De Pagamento S.A."

	// LoginPath ...
	LoginPath string = "v5/token"
	// LoginMtlsPath ...
	LoginMtlsPath string = "v5/token"
	// BalancePath
	BalancePath string = "/baas-walletreports/v1/wallet/balance"
	// CustomersPath ...
	CustomersPath string = "/baas-accountmanager/v1/account/fetch"
	// BusinessPath ...
	BusinessPath string = "/baas-accountmanager/v1/account/fetch-business"
	// ProposalsPath ...
	ProposalsPath string = "/onboarding/v1/onboarding-proposal"
	// ProposalFilesPath ...
	ProposalFilesPath string = "/onboarding/v1/onboarding-proposal/files"
	// CancelAccountPath ...
	CancelAccountPath string = "/baas-accountmanager/v1/account/close"
	// UpdateAccountStatusPath ...
	UpdateAccountStatusPath string = "/baas-accountmanager/v1/account/status"
	//NaturalPersonOnboardingPath ...
	NaturalPersonOnboardingPath string = "/onboarding/v1/onboarding-proposal/natural-person"
	// LegalPersonOnboardingPath ...
	LegalPersonOnboardingPath string = "/onboarding/v1/onboarding-proposal/legal-person"
	// ExternalTransfersPath external transfer (TED)
	ExternalTransfersPath string = "/baas-wallet-transactions-webservice/v1/spb/transfer"
	// InternalTransfersPath internal transfer
	InternalTransfersPath string = "/baas-wallet-transactions-webservice/v1/wallet/internal/transfer"

	// Pix ...
	PixClaimPath string = "/celcoin-baas-pix-dict-webservice/v1/pix/dict/claim"
	PixDictPath  string = "/celcoin-baas-pix-dict-webservice/v1/pix/dict/entry"
	//Deprecated
	PixDictDueDatePathDeprecated string = "/pix/v1/dict/v2/key"
	PixDictDueDatePath           string = "/baas/v2/pix/dict/entry/external"
	PixCashOutPath               string = "/baas-wallet-transactions-webservice/v1/pix/payment"
	PixCashInPath                string = "/pix/v2/receivement/v2"
	PixEmvPath                   string = "/pix/v1/emv"
	PixStaticPath                string = "/pix/v1/brcode/static"
	PixCashInStatusPath          string = "/pix/v2/receivement/v2/devolution/status"
	PixEmvUrl                    string = "/pix/v1/collection"
	PixCashInDynamicPath         string = "/pix/v1/collection"
	PixQrCodeLocationPath        string = "/pix/v1/location"

	// StatementPath ...
	StatementPath string = "/baas-walletreports/v1/wallet/movement"
	// DdaSubscriptionPath ...
	DdaSubscriptionPath string = "/dda-subscription-webservice/v1/subscription/Register"
	// IncomeReportPath ...
	IncomeReportPath string = "/baas-accountmanager/v1/account/income-report"

	// Webhook
	WebhookPath    string = "/baas-webhookmanager/v1/webhook"
	WebhookDdaPath string = "/dda-servicewebhook-webservice/v1/webhook"

	//BaasV2ChargePath ...
	BaasV2ChargePath string = "/baas/v2/charge"

	// ONBOARDING CONSTANTS
	// OnboardingStatusProcessing ...
	OnboardingStatusProcessing string = "PROCESSING"

	// OnboardingStatusApproved ...
	OnboardingStatusApproved string = "APPROVED"

	// OnboardingStatusReproved ...
	OnboardingStatusReproved string = "REPROVED"

	// OnboardingStatusPending ...
	OnboardingStatusPending string = "PENDING"

	//ProposalTypeNaturalPerson ...
	ProposalTypeNaturalPerson string = "NATURAL_PERSON"

	// ProposalTypeLegalPerson ...
	ProposalTypeLegalPerson string = "LEGAL_PERSON"

	// LegalPersonOwnerTypeSocio ...
	LegalPersonOwnerTypeSocio string = "SOCIO"

	// LegalPersonOwnerTypeRepresentante ...
	LegalPersonOwnerTypeRepresentante string = "REPRESENTANTE"

	// LegalPersonOwnerTypeDemaisSocios...
	LegalPersonOwnerTypeDemaisSocios string = "DEMAIS_SOCIOS"

	// ProposalTypePF...
	ProposalTypePF string = "PF"

	// ProposalTypePJ...
	ProposalTypePJ string = "PJ"

	// OnboardingProposalCompanyTypes...
	CompanyTypeME     CompanyType = "ME"
	CompanyTypeMEI    CompanyType = "MEI"
	CompanyTypeEPP    CompanyType = "EPP"
	CompanyTypeLTDA   CompanyType = "LTDA"
	CompanyTypeSA     CompanyType = "SA"
	CompanyTypeEI     CompanyType = "EI"
	CompanyTypeEIRELI CompanyType = "EIRELI"
	CompanyTypePJ     CompanyType = "PJ"

	// DefaultOnboardingType ...
	DefaultOnboardingType string = "BAAS"

	// ChargeDiscountModalityFixed ...
	ChargeDiscountModalityFixed string = "FIXED"
	// ChargeDiscountModalityPercentage ...
	ChargeDiscountModalityPercentage string = "PERCENT"
)

const (
	// NaturalPersonType Pessoa Fisica
	NaturalPersonType PersonType = "F"
	// LegalPersonType Pessoa Juridica
	LegalPersonType PersonType = "J"
)

type AccountType string

const (
	// AccountTypeCC Conta corrente
	AccountTypeCC AccountType = "CC"
	// AccountTypeCI Conta de investimento
	AccountTypeCI AccountType = "CI"
	// AccountTypePG Conta de pagamento
	AccountTypePG AccountType = "PG"
	// AccountTypePP Conta poupança
	AccountTypePP AccountType = "PP"
)

type ClientFinality string

const (
	// TaxesLeviesAndFees 1 - Pagamento de Impostos, Tributos e Taxas
	TaxesLeviesAndFeesClientFinality ClientFinality = "1"
	// Dividends 3 - Pagamentos de Dividendos
	DividendsClientFinality ClientFinality = "3"
	// Salaries 4 - Pagamento de Salários
	SalariesClientFinality ClientFinality = "4"
	// Suppliers 5 - Pagamento de Fornecedores
	SuppliersClientFinality ClientFinality = "5"
	// RentAndCondominiumFees 7 - Pagamento de Aluguéis e Taxas de Condomínio
	RentAndCondominiumFeesClientFinality ClientFinality = "7"
	// SchoolTuition 9 - Pagamento de Mensalidade Escolar
	SchoolTuitionClientFinality ClientFinality = "9"
	// AccountCredit 10 - Crédito em Conta
	AccountCreditClientFinality ClientFinality = "10"
	// JudicialDeposit 100 - Depósito Judicial
	JudicialDepositClientFinality ClientFinality = "100"
	// TransfersBetweenSameOwnership 110 - Transferência entre contas de mesma titularidade
	TransfersBetweenSameOwnershipClientFinality ClientFinality = "110"
	// Others 99999 - Outros
	OthersClientFinality ClientFinality = "99999"
)

type PixType string

const (
	// PixCNPJ ...
	PixCNPJ PixType = "CNPJ"
	// PixCPF ...
	PixCPF PixType = "CPF"
	// PixEMAIL ...
	PixEMAIL PixType = "EMAIL"
	// PixPHONE ...
	PixPHONE PixType = "PHONE"
	//  PixEVP ...
	PixEVP PixType = "EVP"
)

type PixClaimType string

const (
	Portability PixClaimType = "PORTABILITY"
	Ownership   PixClaimType = "OWNERSHIP"
)

type StatusClaim string

const (
	Open              StatusClaim = "OPEN"
	WaitingResolution StatusClaim = "WAITING_RESOLUTION"
	Confirmed         StatusClaim = "CONFIRMED"
	CanceledClaim     StatusClaim = "CANCELED"
	CompletedClaim    StatusClaim = "COMPLETED"
)

type CancelReason string

const (
	UserRequested    CancelReason = "USER_REQUESTED"
	ClaimerRequest   CancelReason = "CLAIMER_REQUEST"
	DonorRequest     CancelReason = "DONOR_REQUEST"
	AccountClosure   CancelReason = "ACCOUNT_CLOSURE"
	Fraud            CancelReason = "FRAUD"
	DefaultOperation CancelReason = "DEFAULT_OPERATION"
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
	ClientSecret     string `json:"client_secret"`
}

// WebhookCredential ...
type WebhookCredential struct {
	Login    string `json:"celcoin_webhook_login"`
	Password string `json:"celcoin_webhook_password"`
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
	Entity     string      `validate:"required" json:"entity"`     // Identificador do evento
	WebhookURL string      `validate:"required" json:"webhookUrl"` // URL do webhook
	Auth       WebhookAuth `json:"auth"`                           // Dados de autenticação do webhook
}
type WebhookSubscriptionDdaRequest struct {
	BasicAuthentication *WebhookBasicAuthentication `json:"basicAuthentication"` // identification + password
	OAuthTwo            *WebhookOAuthTwo            `json:"oAuthTwo,omitempty"`
	TypeEventWebhook    string                      `json:"typeEventWebhook"` // "Subscription", "Deletion" or "Invoice"
	URL                 string                      `json:"url"`
}

// WebhookAuth representa os dados de autenticação para o webhook
type WebhookAuth struct {
	Login    string `json:"login"` // Login para autenticação
	Password string `json:"pwd"`   // Senha para autenticação
	Type     string `json:"type"`  // Tipo de autenticação (basic)
}

// WebhookSubscriptionResponse representa a resposta da rota de cadastro de webhooks
type WebhookSubscriptionResponse struct {
	Version string                   `json:"version"` // Versão da API
	Status  string                   `json:"status"`  // Status da operação (SUCCESS ou ERROR)
	Body    *WebhookSubscriptionBody `json:"body,omitempty"`
	Error   *WebhookError            `json:"error,omitempty"` // Detalhes do erro, se houver
}

// WebhookSubscriptionDdaResponse representa a resposta da rota DDA (formato específico)
type WebhookSubscriptionDdaResponse struct {
	Version string                      `json:"version"` // Versão da API
	Status  int                         `json:"status"`
	Body    *WebhookSubscriptionDdaBody `json:"body,omitempty"`
	Error   *WebhookError               `json:"error,omitempty"`
}

// WebhookSubscriptionDdaBody representa o corpo do response DDA
type WebhookSubscriptionDdaBody struct {
	TypeEventWebhook    string                      `json:"typeEventWebhook"`
	URL                 string                      `json:"url"`
	BasicAuthentication *WebhookBasicAuthentication `json:"basicAuthentication"`
	OAuthTwo            *WebhookOAuthTwo            `json:"oAuthTwo,omitempty"`
}

// WebhookBasicAuthentication representa basic auth no payload DDA
type WebhookBasicAuthentication struct {
	Identification string `json:"identification"`
	Password       string `json:"password"`
}

// WebhookOAuthTwo representa o objeto oAuthTwo do payload DDA
type WebhookOAuthTwo struct {
	Endpoint     string `json:"endpoint"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
	Code         string `json:"code"`
	RefreshToken string `json:"refreshToken"`
	ContentType  string `json:"contentType"`
	GrantType    string `json:"grantType"`
}

// WebhookSubscriptionBody representa informações do cadastro
type WebhookSubscriptionBody struct {
	SubscriptionId string `json:"subscriptionId"` // Código do webhook
}

// WebhookError representa informações de erro em uma resposta da API
type WebhookError struct {
	ErrorCode string `json:"errorCode"` // Código do erro
	Message   string `json:"message"`   // Mensagem do erro
}

// WebhookSubscription representa os dados retornados na consulta de webhooks cadastrados
type WebhookSubscription struct {
	SubscriptionId string      `json:"subscriptionId"`
	Entity         string      `json:"entity"`
	WebhookURL     string      `json:"webhookUrl"`
	Active         bool        `json:"active"`
	CreateDate     time.Time   `json:"createDate"`
	LastUpdateDate time.Time   `json:"lastUpdateDate"`
	Auth           WebhookAuth `json:"auth"`
}

// WebhookQueryResponse representa a resposta da API ao consultar webhooks cadastrados
type WebhookQueryResponse struct {
	Version string `json:"version"`
	Status  string `json:"status"`
	Body    struct {
		Subscriptions []WebhookSubscription `json:"subscriptions"`
	} `json:"body"`
	Error *WebhookError `json:"error,omitempty"`
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
	ClientRequestId string `json:"clientRequestId,omitempty"` // ID da solicitação do cliente
}

/*WEBHOOK MODELS*/

// BalanceResponse ...
type BalanceResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Body    struct {
		Amount float64 `json:"amount"`
	} `json:"body"`
}

// AccountResponse ...
type AccountResponse struct {
	Balance *BalanceResponse `json:"balance,omitempty"`
	Status  string           `json:"status,omitempty"`
	Branch  string           `json:"branch,omitempty"`
	Number  string           `json:"number,omitempty"`
	Bank    *BankData        `json:"bank,omitempty"`
}

// BankData ...
type BankData struct {
	ISPB string `json:"ispb,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"compe,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Errors    []ErrorModel `json:"errors,omitempty"`
	Title     string       `json:"title,omitempty"`
	Status    int32        `json:"status,omitempty"`
	TraceId   string       `json:"traceId,omitempty"`
	Reference string       `json:"reference,omitempty"`
	CodeMessageErrorResponse
}

// ErrorDefaultResponse ...
type ErrorDefaultResponse struct {
	Status  *string       `json:"status"`
	Version *string       `json:"version"`
	Error   *ErrorDefault `json:"error"`
}

// ErroDefaultResponse necessário para tratar erros no dda
type ErroDefaultResponse struct {
	Status  *int          `json:"status"`
	Version *string       `json:"version"`
	Error   *ErrorDefault `json:"erro"`
}

// ErrorDefault ...
type ErrorDefault struct {
	ErrorCode *string `json:"errorCode"`
	Message   *string `json:"message"`
}

// CustomerResponse ...
type CustomerResponse struct {
	Body    CustomerResponseBody `json:"body"`
	Version string               `json:"version"`
	Status  string               `json:"status"`
}

// CustomerResponseBody ...
type CustomerResponseBody struct {
	StatusAccount              string     `json:"statusAccount"`
	DocumentNumber             string     `json:"documentNumber"`
	PhoneNumber                string     `json:"phoneNumber"`
	Email                      string     `json:"email"`
	ClientCode                 string     `json:"clientCode"`
	MotherName                 string     `json:"motherName"`
	FullName                   string     `json:"fullName"`
	SocialName                 string     `json:"socialName"`
	BirthDate                  string     `json:"birthDate"`
	Address                    Address    `json:"address"`
	IsPoliticallyExposedPerson bool       `json:"isPoliticallyExposedPerson"`
	Account                    Account    `json:"account"`
	CreateDate                 CustomTime `json:"createDate"`
}

// Address ... representa o objeto "address"
type Address struct {
	PostalCode        string  `json:"postalCode"`
	Street            string  `json:"street"`
	Number            string  `json:"number"`
	AddressComplement string  `json:"addressComplement"`
	Neighborhood      string  `json:"neighborhood"`
	City              string  `json:"city"`
	State             string  `json:"state"`
	Longitude         *string `json:"longitude"`
	Latitude          *string `json:"latitude"`
}

// Account..  representa o objeto "account"
type Account struct {
	Branch  string `json:"branch"`
	Account string `json:"account"`
}

// BusinessResponse ...
type BusinessResponse struct {
	Body    BusinessResponseBody `json:"body"`
	Version string               `json:"version"`
	Status  string               `json:"status"`
}

// BusinessResponseBody ...
type BusinessResponseBody struct {
	StatusAccount       string          `json:"statusAccount"`
	DocumentNumber      string          `json:"documentNumber"`
	ClientCode          string          `json:"clientCode"`
	BusinessPhoneNumber string          `json:"businessPhoneNumber"`
	BusinessEmail       string          `json:"businessEmail"`
	CreateDate          CustomTime      `json:"createDate"`
	BusinessName        string          `json:"businessName"`
	TradingName         string          `json:"tradingName"`
	Owners              []Owner         `json:"owners"`
	BusinessAccount     BusinessAccount `json:"businessAccount"`
	BusinessAddress     Address         `json:"businessAddress"`
}

// BusinessOnboardingResponse ... representa a resposta do onboarding de customer
type BusinessOnboardingResponse struct {
	Body    BusinessOnboardingResponseBody `json:"body"`
	Version string                         `json:"version"`
	Status  string                         `json:"status"`
}

// BusinessOnboardingResponseBody ... representa o corpo da resposta do onboarding de customer
type BusinessOnboardingResponseBody struct {
	ProposalID     string `json:"proposalId"`
	ClientCode     string `json:"clientCode"`
	DocumentNumber string `json:"documentNumber"`
}

// BusinessOnboardingRequest ... representa o payload para o onboarding de uma empresa.
type BusinessOnboardingRequest struct {
	ClientCode      string  `json:"clientCode"`
	ContactNumber   string  `json:"contactNumber"`
	DocumentNumber  string  `json:"documentNumber"`
	BusinessEmail   string  `json:"businessEmail"`
	BusinessName    string  `json:"businessName"`
	TradingName     string  `json:"tradingName"`
	CompanyType     string  `json:"companyType"`
	Owner           []Owner `json:"owner"`
	BusinessAddress Address `json:"businessAddress"`
	OnboardingType  string  `json:"onboardingType"`
}

// BusinessOnboardingRequest ... representa o payload para o onboarding de uma empresa na migração.
type BusinessOnboardingMigrationRequest struct {
	ClientCode      string                   `json:"clientCode"`
	ContactNumber   string                   `json:"contactNumber"`
	DocumentNumber  string                   `json:"documentNumber"`
	BusinessEmail   string                   `json:"businessEmail"`
	BusinessName    string                   `json:"businessName"`
	TradingName     string                   `json:"tradingName"`
	CompanyType     string                   `json:"companyType"`
	Owner           []Owner                  `json:"owner"`
	BusinessAddress Address                  `json:"businessAddress"`
	OnboardingType  string                   `json:"onboardingType"`
	Files           []CustomerFilesMigration `json:"files"`
}

// Owner ...
type Owner struct {
	OwnerType                  string  `json:"ownerType"`
	DocumentNumber             string  `json:"documentNumber"`
	PhoneNumber                string  `json:"phoneNumber"`
	Email                      string  `json:"email"`
	FullName                   string  `json:"fullName"`
	SocialName                 string  `json:"socialName"`
	BirthDate                  string  `json:"birthDate"`
	MotherName                 string  `json:"motherName"`
	Address                    Address `json:"address"`
	IsPoliticallyExposedPerson bool    `json:"isPoliticallyExposedPerson"`
}

// BusinessAccount ...
type BusinessAccount struct {
	Branch  string `json:"branch"`
	Account string `json:"account"`
}

// CustomTime ... é um tipo customizado para lidar com datas no formato específico
type CustomTime struct {
	time.Time
}

// UnmarshalJSON ... método para deserializar CustomTime
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}

	ct.Time = t
	return nil
}

// MarshalJSON ... método para serializar CustomTime
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Time.Format("2006-01-02T15:04:05"))
}

// CustomerAddress ... representa o endereço do cliente
type CustomerAddress struct {
	PostalCode        string `json:"postalCode"`
	Street            string `json:"street"`
	Number            string `json:"number"`
	AddressComplement string `json:"addressComplement"`
	Neighborhood      string `json:"neighborhood"`
	City              string `json:"city"`
	State             string `json:"state"`
}

// Customer ... representa o cliente
type Customer struct {
	ClientCode                 string          `json:"clientCode"`
	DocumentNumber             string          `json:"documentNumber"`
	PhoneNumber                string          `json:"phoneNumber"`
	Email                      string          `json:"email"`
	MotherName                 string          `json:"motherName"`
	FullName                   string          `json:"fullName"`
	SocialName                 string          `json:"socialName"`
	BirthDate                  string          `json:"birthDate"`
	Address                    CustomerAddress `json:"address"`
	IsPoliticallyExposedPerson bool            `json:"isPoliticallyExposedPerson"`
	OnboardingType             string          `json:"onboardingType"`
}

// CustomerMigration ... representa o cliente a ser migrado
type CustomerMigration struct {
	ClientCode                 string                   `json:"clientCode"`
	DocumentNumber             string                   `json:"documentNumber"`
	PhoneNumber                string                   `json:"phoneNumber"`
	Email                      string                   `json:"email"`
	MotherName                 string                   `json:"motherName"`
	FullName                   string                   `json:"fullName"`
	SocialName                 string                   `json:"socialName"`
	BirthDate                  string                   `json:"birthDate"`
	Address                    CustomerAddress          `json:"address"`
	IsPoliticallyExposedPerson bool                     `json:"isPoliticallyExposedPerson"`
	OnboardingType             string                   `json:"onboardingType"`
	Files                      []CustomerFilesMigration `json:"files"`
}

type CustomerFilesMigration struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// CustomerOnboardingResponse ... representa a resposta do onboarding de customer
type CustomerOnboardingResponse struct {
	Body    CustomerOnboardingResponseBody `json:"body"`
	Version string                         `json:"version"`
	Status  string                         `json:"status"`
}

// CustomerOnboardingResponseBody ... representa o corpo da resposta do onboarding de customer
type CustomerOnboardingResponseBody struct {
	ProposalID     string `json:"proposalId"`
	ClientCode     string `json:"clientCode"`
	DocumentNumber string `json:"documentNumber"`
}

// OnboardingProposalResponse ... representa a resposta do método GetOnboardingProposal
type OnboardingProposalResponse struct {
	Body    OnboardingProposalResponseBody `json:"body"`
	Version string                         `json:"version"`
	Status  string                         `json:"status"`
}

// OnboardingProposalResponseBody ... representa o corpo da resposta do método GetOnboardingProposal
type OnboardingProposalResponseBody struct {
	Limit        int        `json:"limit"`
	CurrentPage  int        `json:"currentPage"`
	LimitPerPage int        `json:"limitPerPage"`
	TotalPages   int        `json:"totalPages"`
	TotalItems   int        `json:"totalItems"`
	Proposals    []Proposal `json:"proposal"`
}

// Proposal ... representa uma proposta no corpo da resposta
type Proposal struct {
	ProposalID     string          `json:"proposalId"`
	ClientCode     string          `json:"clientCode"`
	DocumentNumber string          `json:"documentNumber"`
	Status         string          `json:"status"`
	ProposalType   string          `json:"proposalType"`
	CreatedAt      string          `json:"createdAt"`
	UpdatedAt      string          `json:"updatedAt"`
	DocumentsCopys []DocumentsCopy `json:"documentscopys"`
}

// DocumentsCopy ... representa uma cópia de documento na proposta
type DocumentsCopy struct {
	ProposalID      string `json:"proposalId"`
	DocumentNumber  string `json:"documentNumber"`
	DocumentsCopyID string `json:"documentscopyId"`
	Status          string `json:"status"`
	URL             string `json:"url"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updateAt"`
}

// OnboardingProposalFilesResponse ... representa a resposta do método GetOnboardingProposalFiles
type OnboardingProposalFilesResponse struct {
	Body    OnboardingProposalFilesResponseBody `json:"body"`
	Version string                              `json:"version"`
	Status  string                              `json:"status"`
}

// CancelAccountResponse ... representa a resposta do método CancelAccount
type CancelAccountResponse struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

// UpdateAccountStatusResponse ... representa a resposta do método CancelAccount
type UpdateAccountStatusResponse struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

// OnboardingProposalFilesResponseBody ... representa o corpo da resposta do método GetOnboardingProposalFiles
type OnboardingProposalFilesResponseBody struct {
	Files          []OnboardingFile `json:"files"`
	ClientCode     string           `json:"clientCode"`
	DocumentNumber string           `json:"documentNumber"`
	ProposalID     string           `json:"proposalId"`
}

// OnboardingFile ... representa um arquivo de onboarding
type OnboardingFile struct {
	Type           string    `json:"type"`
	URL            string    `json:"url"`
	ExpirationTime time.Time `json:"expirationTime"`
}

// UnmarshalJSON ... customizado para OnboardingFile para lidar com o formato de tempo
func (f *OnboardingFile) UnmarshalJSON(data []byte) error {
	type Alias OnboardingFile
	aux := &struct {
		ExpirationTime string `json:"expirationTime"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	expirationTime, err := time.Parse(time.RFC3339, aux.ExpirationTime)
	if err != nil {
		return err
	}

	f.ExpirationTime = expirationTime
	return nil
}

/* PIX */
// PixKeyRequest representa o payload para criação de uma chave Pix.
type PixKeyRequest struct {
	Account string `json:"account" validate:"required"`
	KeyType string `json:"keyType" validate:"required,oneof=EVP CPF CNPJ EMAIL PHONE"`
	Key     string `json:"key,omitempty"`
}

// PixKeyResponse representa a resposta principal para uma operação relacionada a Pix Keys.
type PixKeyResponse struct {
	Body    PixKeyResponseBody `json:"body"`
	Version string             `json:"version"`
	Status  string             `json:"status"`
}

// PixKeyResponseBody representa o corpo da resposta para uma operação relacionada a Pix Keys.
type PixKeyResponseBody struct {
	KeyType string        `json:"keyType"`
	Key     string        `json:"key"`
	Account PixKeyAccount `json:"account"`
	Owner   PixKeyOwner   `json:"owner"`
}

// PixKeyAccount representa os detalhes da conta vinculada a uma Pix Key.
type PixKeyAccount struct {
	Participant string    `json:"participant"`
	Branch      string    `json:"branch"`
	Account     string    `json:"account"`
	AccountType string    `json:"accountType"`
	CreateDate  time.Time `json:"createDate"`
}
type PixExternalKeyAccount struct {
	Participant string    `json:"participant"`
	Branch      int       `json:"branch"`
	Account     string    `json:"accountNumber"`
	AccountType string    `json:"accountType"`
	OpeningDate time.Time `json:"openingDate"`
}

type PixExternalKeyAccountDueDate struct {
	Participant string    `json:"participant"`
	Branch      string    `json:"branch"`
	Account     string    `json:"accountNumber"`
	AccountType string    `json:"accountType"`
	OpeningDate time.Time `json:"openingDate"`
}

// PixKeyOwner representa as informações do proprietário da Pix Key.
type PixKeyOwner struct {
	Type           string `json:"type"` // NATURAL_PERSON ou LEGAL_PERSON
	DocumentNumber string `json:"documentNumber"`
	Name           string `json:"name"`
}

// PixKeyListResponse representa a resposta ao consultar todas as chaves Pix de uma conta.
type PixKeyListResponse struct {
	Version string                 `json:"version"`
	Status  string                 `json:"status"`
	Body    PixKeyListResponseBody `json:"body"`
}

// PixKeyListResponseBody representa o corpo da resposta ao consultar todas as chaves Pix de uma conta.
type PixKeyListResponseBody struct {
	ListKeys []PixKeyListItem `json:"listKeys"`
}

// PixKeyListItem representa uma chave Pix na lista de chaves retornada.
type PixKeyListItem struct {
	KeyType string        `json:"keyType"`
	Key     string        `json:"key"`
	Account PixKeyAccount `json:"account"`
	Owner   PixKeyOwner   `json:"owner"`
}

// PixExternalKeyRequest representa os parâmetros para consulta de uma chave Pix externa (DICT).
type PixExternalKeyRequest struct {
	Key        string `json:"key" validate:"required"`
	OwnerTaxID string `json:"ownerTaxId"`
}

// PixExternalKeyResponse representa a resposta de uma consulta de chave Pix externa (DICT).
type PixExternalKeyResponse struct {
	Status  string                       `json:"status"`
	Version string                       `json:"version"`
	Body    PixExternalKeyResponseBody   `json:"body,omitempty"`
	Error   *PixExternalKeyErrorResponse `json:"error,omitempty"`
}

// PixExternalKeyResponseBody representa o corpo da resposta para a consulta de chave Pix externa.
type PixExternalKeyResponseBody struct {
	KeyType          string        `json:"keyType"`
	Key              string        `json:"key"`
	Account          PixKeyAccount `json:"account"`
	Owner            PixKeyOwner   `json:"owner"`
	EndToEndId       string        `json:"endtoEndId"`
	CreationDate     time.Time     `json:"creationDate"`
	KeyOwnershipDate time.Time     `json:"keyOwnershipDate"`
	IsSameTaxId      bool          `json:"isSameTaxId"`
}

// PixExternalKeyErrorResponse representa os detalhes de erro em caso de falha na consulta de chave Pix externa.
type PixExternalKeyErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

type PixExternalKeyDueDateResponse struct {
	Status  string                    `json:"status"`
	Version string                    `json:"version"`
	Body    PixExternalKeyDueDateBody `json:"body,omitempty"`
}

// PixExternalKeyDueDateBody representa a resposta de uma consulta de chave Pix externa (COBV - DUEDATE).
type PixExternalKeyDueDateBody struct {
	Key              string                       `json:"key"`
	KeyType          string                       `json:"keyType"`
	Account          PixExternalKeyAccountDueDate `json:"account"`
	Owner            PixKeyOwner                  `json:"owner"`
	EndToEndId       string                       `json:"endtoendid"`
	CreationDate     time.Time                    `json:"creationDate"`
	KeyOwnershipDate time.Time                    `json:"keyOwnershipDate"`
}

// PixCashOutRequest representa os dados para realizar um Pix Cash-Out.
type PixCashOutRequest struct {
	Amount float64 `json:"amount" description:"O valor da transação (required)"`
	//VlcpAmount                float64     `json:"vlcpAmount" description:"O valor da compra (Pix Troco)"`
	//VldnAmount                float64     `json:"vldnAmount" description:"O valor em dinheiro disponibilizado (Pix Troco)"`
	//WithdrawalServiceProvider string      `json:"withdrawalServiceProvider" description:"O Identificador ISPB do serviço de saque (Pix Saque/Troco)"`
	//WithdrawalAgentMode       string      `json:"withdrawalAgentMode" description:"Modo do agente de retirada. AGTEC: Estabelecimento Comercial, AGTOT: Entidade Jurídica cuja atividade é a prestação de serviços auxiliares de serviços financeiros, AGPSS: Participante Pix que presta diretamente o serviço de saque."`
	ClientCode                string      `json:"clientCode" description:"A identificação única da transacção dada pelo lado do cliente. Este valor não pode ser repetido (required)"`
	TransactionIdentification string      `json:"transactionIdentification" description:"Identificador do QRCode a ser pago (ver regras de preenchimento)"`
	EndToEndId                string      `json:"endToEndId" description:"Identificador de ponta a ponta associado a este pedido de iniciação de pagamento. Deve ser o mesmo da consulta ao DICT, quando aplicável."`
	DebitParty                DebitParty  `json:"debitParty" description:"Dados bancários da conta do pagador na Celcoin"`
	CreditParty               CreditParty `json:"creditParty" description:"Dados bancários da conta do recebedor"`
	InitiationType            string      `json:"initiationType" description:"Representa o tipo de pagamento que será iniciado (required)"`
	TaxIdPaymentInitiator     string      `json:"taxIdPaymentInitiator" description:"CNPJ do iniciador de pagamentos. Utilizado apenas se o campo 'initiationType' for igual a 'PAYMENT_INITIATOR'."`
	RemittanceInformation     string      `json:"remittanceInformation" description:"Texto a ser apresentado ao pagador para informação correlacionada, em formato livre."`
	PaymentType               string      `json:"paymentType" description:"Representa o tipo de pagamento: IMMEDIATE (padrão), FRAUD (suspeita de fraude), SCHEDULED (programado)."`
	Urgency                   string      `json:"urgency" description:"Define a urgência do pagamento: HIGH (padrão), NORMAL (programado)."`
	TransactionType           string      `json:"transactionType" description:"Tipo de transação: TRANSFER (padrão), CHANGE (Pix Troco), WITHDRAWAL (Pix Saque)."`
}

// DebitParty representa os dados do pagador.
type DebitParty struct {
	Account     string `json:"account" description:"Conta bancária do pagador"`
	Bank        string `json:"bank" description:"Banco do pagador"`
	Branch      string `json:"branch" description:"Agência do pagador"`
	PersonType  string `json:"personType" description:"Tipo de pessoa do pagador (Física/Jurídica)"`
	TaxId       string `json:"taxId" description:"CPF/CNPJ do pagador"`
	AccountType string `json:"accountType" description:"Tipo de conta do pagador CACC, TRAN, SLRY, SVGS"`
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
}

// CreditParty representa os dados do recebedor.
type CreditParty struct {
	Account     string `json:"account" description:"Conta bancária do recebedor"`
	Bank        string `json:"bank" description:"Banco do recebedor"`
	Branch      string `json:"branch" description:"Agência do recebedor"`
	PersonType  string `json:"personType" description:"Tipo de pessoa do recebedor (Física/Jurídica)"`
	TaxId       string `json:"taxId" description:"CPF/CNPJ do recebedor"`
	AccountType string `json:"accountType" description:"Tipo de conta do recebedor CACC, TRAN, SLRY, SVGS "`
	Name        string `json:"name" description:"Nome do recebedor"`
	Key         string `json:"key" description:"Chave Pix do recebedor"`
}

// PixCashOutResponse ...
type PixCashOutResponse struct {
	Status  string                   `json:"status"`
	Version string                   `json:"version"`
	Body    PixCashOutResponseBody   `json:"body,omitempty"`
	Error   *PixCashOutErrorResponse `json:"error,omitempty"`
}

// PixCashOutErrorResponse...
type PixCashOutErrorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

// PixCashOutResponseBody...
type PixCashOutResponseBody struct {
	ID                        string      `json:"id"` // id transação retornado para ser usado em get posteriormente
	Amount                    float64     `json:"amount"`
	ClientCode                string      `json:"clientCode"` // identificador unico fornecido pelo cliente na requisição de payment.
	TransactionIdentification *string     `json:"transactionIdentification,omitempty"`
	EndToEndID                string      `json:"endToEndId"` // identificador ponta-a-ponta da transação. O mesmo retornado na consulta DICT e no retorno do endpoint de pagamento (/baas-wallet-transactions-webservice/v1/pix/payment).
	InitiationType            string      `json:"initiationType"`
	PaymentType               string      `json:"paymentType"`
	Urgency                   string      `json:"urgency"`
	TransactionType           string      `json:"transactionType"`
	DebitParty                DebitParty  `json:"debitParty"`
	CreditParty               CreditParty `json:"creditParty"`
	RemittanceInformation     string      `json:"remittanceInformation"`
}

// ErrorDetails ...
type ErrorDetails struct {
	Code        string `json:"code" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// QRCodeResponse representa a resposta decodificada do QR Code.
type QRCodeResponse struct {
	Type                       string                     `json:"type"`
	Collection                 interface{}                `json:"collection"`
	PayloadFormatIndicator     string                     `json:"payloadFormatIndicator"`
	MerchantAccountInformation MerchantAccountInformation `json:"merchantAccountInformation"`
	MerchantCategoryCode       int                        `json:"merchantCategoryCode"`
	TransactionCurrency        int                        `json:"transactionCurrency"`
	TransactionAmount          float64                    `json:"transactionAmount"`
	CountryCode                string                     `json:"countryCode"`
	MerchantName               string                     `json:"merchantName"`
	MerchantCity               string                     `json:"merchantCity"`
	PostalCode                 string                     `json:"postalCode"`
	InitiationMethod           interface{}                `json:"initiationMethod"`
	TransactionIdentification  string                     `json:"transactionIdentification"`
}

// MerchantAccountInformation representa as informações da conta do comerciante.
type MerchantAccountInformation struct {
	URL                       interface{} `json:"url"`
	GUI                       string      `json:"gui"`
	Key                       string      `json:"key"`
	AdditionalInformation     string      `json:"additionalInformation"`
	WithdrawalServiceProvider interface{} `json:"withdrawalServiceProvider"`
}

// PixTransactionResponse representa os dados para uma transação Pix.
type PixCashoutStatusTransactionResponse struct {
	Status  string                          `json:"status"`
	Version string                          `json:"version"`
	Body    PixCashoutStatusTransactionBody `json:"body"`
	Error   *PixCashOutErrorResponse        `json:"error,omitempty"`
}

// PixTransactionBody representa o corpo da transação Pix.
type PixCashoutStatusTransactionBody struct {
	ID                        string        `json:"id"`
	Amount                    float64       `json:"amount"`
	ClientCode                string        `json:"clientCode"`
	TransactionIdentification *string       `json:"transactionIdentification,omitempty"`
	EndToEndID                string        `json:"endToEndId"`
	InitiationType            string        `json:"initiationType"`
	PaymentType               string        `json:"paymentType"`
	Urgency                   string        `json:"urgency"`
	TransactionType           string        `json:"transactionType"`
	DebitParty                DebitParty    `json:"debitParty"`
	CreditParty               CreditParty   `json:"creditParty"`
	RemittanceInformation     string        `json:"remittanceInformation"`
	Error                     *ErrorDetails `json:"error,omitempty"`
}

// PixCashInTransactionResponse representa a resposta para a consulta de uma devolução de pagamento ou transferência Pix (Cash-In).
type PixCashinStatusTransactionResponse struct {
	Status               string  `json:"status"`               // Status da transação
	ReturnIdentification string  `json:"returnIdentification"` // Identificação de devolução
	TransactionId        int64   `json:"transactionId"`        // Identificador da transação
	TransactionIdPayment int64   `json:"transactionIdPayment"` // Identificador do pagamento
	TransactionType      string  `json:"transactionType"`      // Tipo da transação (e.g., REVERTED)
	Amount               float64 `json:"amount"`               // Valor da transação
	Reason               string  `json:"reason"`               // Motivo da devolução
	ReversalDescription  string  `json:"reversalDescription"`  // Descrição da reversão (se houver)
	CreatedAt            string  `json:"createdAt"`            // Data de criação da transação
}

// PixCashInStaticRequest representa o request necessário para solicitar um pix estático
type PixCashInStaticRequest struct {
	Key                       string      `json:"key" validate:"required"`
	Amount                    float64     `json:"amount" validate:"required"`
	TransactionIdentification string      `json:"transactionIdentification" validate:"required"`
	Merchant                  PixMerchant `json:"merchant" validate:"required"`
	Tags                      []string    `json:"tags,omitempty"`
	AdditionalInformation     string      `json:"additionalInformation,omitempty"`
	Withdrawal                bool        `json:"withdrawal"`
}

// PixCashInStaticResponse representa a resposta ao realizar um Pix Cash-in por Cobrança Estática.
type PixCashInStaticResponse struct {
	TransactionId             int    `json:"transactionId"`
	EMVQRCode                 string `json:"emvqrcps"`
	TransactionIdentification string `json:"transactionIdentification"`
}

// PixAddressKeyResponse...
type PixAddressKeyResponse struct {
	EndToEndID    string       `json:"endToEndId"`
	AddressingKey PixTypeValue `json:"addressingKey"`
	Holder        PixHolder    `json:"holder"`
	Status        string       `json:"status"`
	CreatedAt     time.Time    `json:"createdAt"`
	OwnedAt       time.Time    `json:"ownedAt"`
}

type PixTypeValue struct {
	Type  PixType `json:"type"`
	Value string  `json:"value"`
}

type PixHolder struct {
	Type       string       `json:"type"`
	Name       string       `json:"name"`
	SocialName string       `json:"socialName,omitempty"`
	Document   PixTypeValue `json:"document"`
}

type QRCodeImmediateResponse struct {
	Status             string           `json:"status"`
	InfoAdicionais     *string          `json:"infoAdicionais,omitempty"`
	TxID               string           `json:"txid"`
	Chave              string           `json:"chave"`
	SolicitacaoPagador *string          `json:"solicitacaoPagador,omitempty"`
	Valor              QRCodeValor      `json:"valor"`
	Calendario         QRCodeCalendario `json:"calendario"`
	Revisao            int              `json:"revisao"`
}

type QRCodeValor struct {
	Original            string  `json:"original"`
	Abatimento          string  `json:"abatimento"`
	Desconto            string  `json:"desconto"`
	Multa               string  `json:"multa"`
	Juros               string  `json:"juros"`
	Final               string  `json:"final"`
	ModalidadeAlteracao int     `json:"modalidadeAlteracao"`
	Retirada            *string `json:"retirada,omitempty"`
}

type QRCodeCalendario struct {
	Criacao                string `json:"criacao"`
	Expiracao              int    `json:"expiracao"`
	Apresentacao           string `json:"apresentacao"`
	ValidadeAposVencimento int    `json:"validadeAposVencimento"`
}

// QRCodeDueDateResponse representa a resposta do endpoint de dueDate.
type QRCodeDueDateResponse struct {
	Calendar              PixCalendar         `json:"calendar"`
	Debtor                PixDebtor           `json:"debtor"`
	Receiver              PixReceiver         `json:"receiver"`
	TransactionID         string              `json:"transactionIdentification"`
	Revision              string              `json:"revision"`
	Status                string              `json:"status"`
	Key                   string              `json:"key"`
	Amount                PixQrCodeAmount     `json:"amount"`
	AdditionalInformation []PixAdditionalInfo `json:"additionalInformation"`
}

// Receiver representa os detalhes do recebedor.
type PixReceiver struct {
	CNPJ        string `json:"cnpj,omitempty"`
	CPF         string `json:"cpf,omitempty"`
	FantasyName string `json:"fantasyName,omitempty"`
	PublicArea  string `json:"publicArea,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postalCode,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
}

// Amount representa os valores relacionados à transação.
type PixAmount struct {
	Original   *string            `json:"original,omitempty"`
	Abatement  *PixChargeFee      `json:"abatement,omitempty"`
	Discount   *PixAmountDiscount `json:"discount,omitempty"`
	Interest   *PixChargeFee      `json:"interest,omitempty"`
	Fine       *PixChargeFee      `json:"fine,omitempty"`
	Final      *string            `json:"final,omitempty"`
	ChangeType int                `json:"changeType"`
	Withdrawal *PixWithdrawal     `json:"withdrawal,omitempty"`
	Change     *PixChangeDetails  `json:"change,omitempty"`
}

// Amount representa os valores relacionados à transação.
type PixAmountCashIn struct {
	Original   *float64           `json:"original,omitempty"`
	Abatement  *PixChargeFee      `json:"abatement,omitempty"`
	Discount   *PixAmountDiscount `json:"discount,omitempty"`
	Interest   *PixChargeFee      `json:"interest,omitempty"`
	Fine       *PixChargeFee      `json:"fine,omitempty"`
	Final      *float64           `json:"final,omitempty"`
	ChangeType int                `json:"changeType"`
	Withdrawal *PixWithdrawal     `json:"withdrawal,omitempty"`
	Change     *PixChangeDetails  `json:"change,omitempty"`
}

// Amount representa os valores relacionados à transação.
type PixQrCodeAmount struct {
	Original  *string `json:"original,omitempty"`
	Abatement *string `json:"abatement,omitempty"`
	Discount  *string `json:"discount,omitempty"`
	Interest  *string `json:"interest,omitempty"`
	Fine      *string `json:"fine,omitempty"`
	Final     *string `json:"final,omitempty"`
}
type PixWithdrawal struct {
	VldnAmount                *float64 `json:"vldnAmount"`
	AgentMode                 *string  `json:"agentMode,omitempty"`
	WithdrawalServiceProvider *string  `json:"withdrawalServiceProvider,omitempty"`
	ChangeType                int      `json:"changeType"`
}

type PixChangeDetails struct {
	VldnAmount                *float64 `json:"vldnAmount"`
	VlcpAmount                *float64 `json:"vlcpAmount"`
	AgentMode                 *string  `json:"agentMode,omitempty"`
	WithdrawalServiceProvider *string  `json:"withdrawalServiceProvider,omitempty"`
	ChangeType                int      `json:"changeType"`
}

// AdditionalInfo representa informações adicionais incluídas na transação.
type PixAdditionalInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PixCashInDueDateResponse ...
type PixCashInDueDateResponse struct {
	TransactionIdentification string              `json:"transactionIdentification"`
	TransactionID             int64               `json:"transactionId"`
	ClientRequestID           string              `json:"clientRequestId"`
	Status                    string              `json:"status"`
	LastUpdate                *time.Time          `json:"lastUpdate,omitempty"`
	PayerQuestion             string              `json:"payerQuestion"`
	AdditionalInformation     []PixAdditionalInfo `json:"additionalInformation"`
	Debtor                    PixDebtor           `json:"debtor"`
	Amount                    PixAmountCashIn     `json:"amount"`
	Location                  *PixLocation        `json:"location,omitempty"` // Campo opcional
	Key                       string              `json:"key"`
	Receiver                  PixReceiver         `json:"receiver"`
	Calendar                  PixCalendar         `json:"calendar"`
	CreateAt                  time.Time           `json:"createAt"`
}

type PixDebtor struct {
	Name       string  `json:"name,omitempty"`
	CPF        *string `json:"cpf,omitempty"`
	CNPJ       string  `json:"cnpj,omitempty"`
	City       string  `json:"city,omitempty"`
	PublicArea string  `json:"publicArea,omitempty"`
	State      string  `json:"state,omitempty"`
	PostalCode string  `json:"postalCode,omitempty"`
	Email      string  `json:"email,omitempty"`
}
type PixChargeFee struct {
	AmountPerc *string `json:"amountPerc,omitempty"` // Aceita string ou número
	Modality   *string `json:"modality"`
}

type PixLocation struct {
	Merchant   PixMerchant `json:"merchant"`
	URL        string      `json:"url"`
	EMV        string      `json:"emv"`
	Type       string      `json:"type"`
	LocationID string      `json:"locationId"`
	ID         *string     `json:"id,omitempty"`
}

// PixMerchant representa as informações do comerciante no Pix Cash-in.
type PixMerchant struct {
	PostalCode           string `json:"postalCode" validate:"required"`
	City                 string `json:"city" validate:"required"`
	MerchantCategoryCode string `json:"merchantCategoryCode"`
	Name                 string `json:"name" validate:"required"`
}

// Calendar representa os detalhes do calendário da transação.
type PixCalendar struct {
	ExpirationAfterPayment  string `json:"expirationAfterPayment,omitempty"`
	CreatedAt               string `json:"createdAt,omitempty"`
	DueDate                 string `json:"dueDate,omitempty"`
	ValidateAfterExpiration int    `json:"validateAfterExpiration,omitempty"`
	Presentation            string `json:"presentation,omitempty"`
	Expiration              int    `json:"expiration,omitempty"`
}

type PixCashInDueDateRequest struct {
	PayerQuestion          string              `json:"payerQuestion"`
	ClientRequestID        string              `json:"clientRequestId"`
	ExpirationAfterPayment int                 `json:"expirationAfterPayment"`
	DueDate                string              `json:"duedate"`
	Debtor                 PixDebtor           `json:"debtor"`
	Receiver               PixReceiver         `json:"receiver"`
	LocationID             int64               `json:"locationId"`
	Amount                 float64             `json:"amount"`
	AdditionalInformation  []PixAdditionalInfo `json:"additionalInformation"`
	AmountInterest         PixAmountInterest   `json:"amountInterest"`
	AmountDiscount         PixAmountDiscount   `json:"amountDicount"`
	AmountAbatement        PixAmountAbatement  `json:"amountAbatement"`
	AmountFine             PixAmountFine       `json:"amountFine"`
	Key                    string              `json:"key"`
}

type PixAmountInterest struct {
	HasCondition bool   `json:"hasInterest"`
	AmountPerc   string `json:"amountPerc"`
	Modality     string `json:"modality"`
}

type PixAmountDiscount struct {
	HasCondition      bool                   `json:"hasDicount"`
	AmountPerc        string                 `json:"amountPerc"`
	Modality          string                 `json:"modality"`
	DiscountDateFixed []PixDiscountDateFixed `json:"discountDateFixed"`
}

type PixDiscountDateFixed struct {
	Date       string `json:"date"`
	AmountPerc string `json:"amountPerc"`
}
type PixAmountAbatement struct {
	HasCondition bool   `json:"hasAbatement"`
	AmountPerc   string `json:"amountPerc"`
	Modality     string `json:"modality"`
}
type PixAmountFine struct {
	HasCondition bool   `json:"hasFine"`
	AmountPerc   string `json:"amountPerc"`
	Modality     string `json:"modality"`
}

type PixCashInImmediateRequest struct {
	ClientRequestID       string              `json:"clientRequestId"`
	PayerQuestion         string              `json:"payerQuestion"`
	Key                   string              `json:"key"`
	LocationID            int64               `json:"locationId"`
	Debtor                PixDebtor           `json:"debtor"`
	Amount                PixAmount           `json:"amount"`
	Calendar              PixCalendar         `json:"calendar"`
	AdditionalInformation []PixAdditionalInfo `json:"additionalInformation"`
}

type PixCashInImmediateResponse struct {
	Revision                  int64               `json:"revision"`
	TransactionID             int64               `json:"transactionId"`
	ClientRequestID           string              `json:"clientRequestId"`
	Status                    string              `json:"status"`
	LastUpdate                *time.Time          `json:"lastUpdate,omitempty"`
	PayerQuestion             string              `json:"payerQuestion"`
	AdditionalInformation     []PixAdditionalInfo `json:"additionalInformation"`
	Debtor                    PixDebtor           `json:"debtor"`
	Amount                    PixAmount           `json:"amount"`
	Location                  PixLocation         `json:"location"`
	Key                       string              `json:"key"`
	Calendar                  PixCalendar         `json:"calendar"`
	CreatedAt                 time.Time           `json:"createAt"`
	TransactionIdentification string              `json:"transactionIdentification"`
}

// PixDeleteResponse representa a resposta da API para a exclusão do Pix.
type PixDeleteResponse struct {
	Message       string `json:"message"`
	TransactionID int64  `json:"transactionId"`
	Status        int    `json:"status"`
}

type PixQrCodeLocationRequest struct {
	ClientRequestID string            `json:"clientRequestId" validate:"required,uuid4"`
	Type            string            `json:"type" validate:"required,oneof=COB COBV"`
	Merchant        PixQrCodeMerchant `json:"merchant" validate:"required"`
}

// QrCodeMerchant representa as informações do comerciante no request
type PixQrCodeMerchant struct {
	MerchantCategoryCode string `json:"merchantCategoryCode" validate:"required"`
	PostalCode           string `json:"postalCode" validate:"required"`
	City                 string `json:"city" validate:"required"`
	Name                 string `json:"name" validate:"required"`
}

// QrCodeLocationResponse representa a resposta da criação do QR Code Location
type PixQrCodeLocationResponse struct {
	LocationID      int64             `json:"locationId"`
	Status          string            `json:"status"`
	ClientRequestID string            `json:"clientRequestId"`
	URL             string            `json:"url"`
	EMV             string            `json:"emv"`
	Type            string            `json:"type"`
	Merchant        PixQrCodeMerchant `json:"merchant"`
}

// PixClaimRequest representa o payload para requisições de portabilidade de chave Pix.
type PixClaimRequest struct {
	Key       string `json:"key" validate:"required"`
	KeyType   string `json:"keyType" validate:"required,oneof=EMAIL CPF CNPJ PHONE EVP"`
	Account   string `json:"account" validate:"required"`
	ClaimType string `json:"claimType" validate:"required,oneof=OWNERSHIP"`
}

// PixClaimResponse representa a resposta de operações individuais de portabilidade de chave Pix.
type PixClaimResponse struct {
	Version string               `json:"version"`
	Status  string               `json:"status"`
	Body    PixClaimResponseBody `json:"body"`
}

// PixClaimResponseBody representa o corpo da resposta de uma única portabilidade de chave Pix.
type PixClaimResponseBody struct {
	ID                  string             `json:"id"`
	ClaimType           string             `json:"claimType"`
	Key                 string             `json:"key"`
	KeyType             string             `json:"keyType"`
	ClaimerAccount      PixClaimKeyAccount `json:"claimerAccount"`
	Claimer             PixClaimKeyOwner   `json:"claimer"`
	DonorParticipant    string             `json:"donorParticipant"`
	Status              string             `json:"status"`
	CreateTimestamp     string             `json:"createTimestamp"`
	CompletionPeriodEnd string             `json:"completionPeriodEnd"`
	ResolutionPeriodEnd string             `json:"resolutionPeriodEnd"`
	LastModified        string             `json:"lastModified"`
	ConfirmReason       string             `json:"confirmReason,omitempty"`
	CancelReason        string             `json:"cancelReason,omitempty"`
	CancelledBy         string             `json:"cancelledBy,omitempty"`
	DonorAccount        PixClaimKeyAccount `json:"donorAccount,omitempty"`
}

// PixClaimListResponse representa a resposta da consulta de lista de reivindicações de chaves Pix.
type PixClaimListResponse struct {
	Version string                   `json:"version"`
	Status  string                   `json:"status"`
	Body    PixClaimListResponseBody `json:"body"`
}

// PixClaimListResponseBody representa o corpo da resposta da consulta de lista de reivindicações de chaves Pix.
type PixClaimListResponseBody struct {
	Claims []PixClaimResponseBody `json:"claims"`
}

// PixClaimActionRequest representa requisições para confirmação ou cancelamento de portabilidade.
type PixClaimActionRequest struct {
	ID     string `json:"id" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}

// PixKeyAccount representa os detalhes da conta bancária associada a uma chave Pix.
type PixClaimKeyAccount struct {
	Participant string `json:"participant,omitempty"`
	Branch      string `json:"branch"`
	Account     string `json:"account"`
	AccountType string `json:"accountType,omitempty"`
	TaxID       string `json:"taxId,omitempty"`
	Name        string `json:"name,omitempty"`
}

// PixKeyOwner representa o proprietário da chave Pix.
type PixClaimKeyOwner struct {
	PersonType string `json:"personType"`
	TaxID      string `json:"taxId"`
	Name       string `json:"name"`
}

/* TRANSFERS */
// TransfersRequest ...
type TransfersRequest struct {
	Amount          float64                     `validate:"required" json:"amount"`
	ClientCode      string                      `validate:"required" json:"clientCode"`
	ClientRequestId string                      `json:"clientRequestId"`
	DebitParty      TransfersDebitPartyRequest  `validate:"required,dive" json:"debitParty"`
	CreditParty     TransfersCreditPartyRequest `validate:"required,dive" json:"creditParty"`
	ClientFinality  ClientFinality              `json:"clientFinality"`
	Description     string                      `json:"description"`
}

// TransfersDebitPartyRequest ...
type TransfersDebitPartyRequest struct {
	AccountNumber string `validate:"required" json:"account"`
	BankISPB      string `json:"bank"`
}

// TransfersDebitPartyResponse ...
type TransfersDebitPartyResponse struct {
	AccountNumber string      `json:"account"`
	AccountBranch string      `json:"branch"`
	Identifier    string      `json:"taxId"`
	AccountName   string      `json:"name"`
	AccountType   AccountType `json:"accountType"`
	PersonType    PersonType  `json:"personType"`
	BankISPB      string      `json:"bank"`
}

// TransfersCreditPartyRequest ...
type TransfersCreditPartyRequest struct {
	BankISPB      string      `validate:"required" json:"bank"`
	AccountNumber string      `validate:"required" json:"account"`
	AccountBranch string      `validate:"required" json:"branch"`
	Identifier    string      `validate:"required,cnpjcpf" json:"taxId"`
	AccountName   string      `validate:"required" json:"name"`
	AccountType   AccountType `validate:"required" json:"accountType"`
	PersonType    PersonType  `validate:"required" json:"personType"`
}

// TransfersCreditPartyResponse ...
type TransfersCreditPartyResponse struct {
	TransfersCreditPartyRequest
}

// TransfersBodyResponse ...
type TransfersBodyResponse struct {
	ID              string                       `json:"id"`
	Amount          float64                      `json:"amount"`
	ClientCode      string                       `json:"clientCode"`
	ClientRequestId string                       `json:"clientRequestId"`
	DebitParty      TransfersDebitPartyResponse  `json:"debitParty"`
	CreditParty     TransfersCreditPartyResponse `json:"creditParty"`
	EndToEndId      string                       `json:"endToEndId"`
	Description     string                       `json:"description"`
}

// TransfersResponse representa a resposta da rota de cadastro de webhooks
type TransfersResponse struct {
	Version string                 `json:"version"` // Versão da API
	Status  string                 `json:"status"`  // Status da operação (SUCCESS ou ERROR)
	Body    *TransfersBodyResponse `json:"body,omitempty"`
	Error   *TransfersError        `json:"error,omitempty"` // Detalhes do erro, se houver
}

// TransfersError representa informações de erro em uma resposta da API
type TransfersError struct {
	ErrorCode string `json:"errorCode"` // Código do erro
	Message   string `json:"message"`   // Mensagem do erro
}

/* BOLETO */

// CreateBoletoRequest is the payload to create a new Celcoin charge/boleto.
type CreateBoletoRequest struct {
	ExternalID             *string       `json:"externalId,omitempty"`
	ExpirationAfterPayment *int          `json:"expirationAfterPayment,omitempty"`
	DueDate                *string       `json:"dueDate,omitempty"`
	Amount                 *float64      `json:"amount,omitempty"`
	Key                    *string       `json:"key,omitempty"` // optional
	Debtor                 *Debtor       `json:"debtor,omitempty"`
	Receiver               *Receiver     `json:"receiver,omitempty"`
	Instructions           *Instructions `json:"instructions,omitempty"`
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
	Fine     *float64  `json:"fine,omitempty"`
	Interest *float64  `json:"interest,omitempty"`
	Discount *Discount `json:"discount,omitempty"`
}

type Discount struct {
	Amount    *float64 `json:"amount,omitempty"`
	Modality  *string  `json:"modality,omitempty"`  // "fixed" or "percent"
	LimitDate *string  `json:"limitDate,omitempty"` // e.g. "2025-01-20T00:00:00.0000000"
}

// CreateBoletoResponse is the simplified response from POST /charge.
type CreateBoletoResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

// QueryBoletoResponse is the simplified struct for GET /charge?TransactionId=...
type QueryBoletoResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

// CancelInput is the JSON body for DELETE /charge/:id requests.
type CancelInput struct {
	Reason string `json:"reason"`
}

// StatementResponse ... define a estrutura da resposta da API.
type StatementResponse struct {
	Status       string        `json:"status"`
	Version      string        `json:"version"`
	TotalItems   int           `json:"totalItems"`
	CurrentPage  int           `json:"currentPage"`
	LimitPerPage int           `json:"limitPerPage"`
	TotalPages   int           `json:"totalPages"`
	DateFrom     string        `json:"dateFrom"`
	DateTo       string        `json:"dateTo"`
	Body         StatementBody `json:"body"`
}

// StatementBody ... define a estrutura do corpo da resposta.
type StatementBody struct {
	Account        string              `json:"account"`
	DocumentNumber string              `json:"documentNumber"`
	Movements      []StatementMovement `json:"movements"`
}

// StatementMovement ... define a estrutura de cada movimento.
type StatementMovement struct {
	ID             string  `json:"id"`
	ClientCode     string  `json:"clientCode"`
	Description    string  `json:"description"`
	CreateDate     string  `json:"createDate"`
	LastUpdateDate string  `json:"lastUpdateDate"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	BalanceType    string  `json:"balanceType"`
	MovementType   string  `json:"movementType"`
}

// StatementRequest ... define a estrutura da requisição da API.
type StatementRequest struct {
	Account        *string `json:"Account"`
	DateFrom       *string `json:"DateFrom"`
	DateTo         *string `json:"DateTo"`
	DocumentNumber *string `json:"DocumentNumber"`
	LimitPerPage   *int64  `json:"LimitPerPage"`
	Page           *int64  `json:"Page"`
}

// IncomeReportPayerSource ... define a estrutura da fonte pagadora.
type IncomeReportPayerSource struct {
	Name           string `json:"name"`
	DocumentNumber string `json:"documentNumber"`
}

// IncomeReportOwner ... define a estrutura do proprietário.
type IncomeReportOwner struct {
	DocumentNumber string `json:"documentNumber"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	CreateDate     string `json:"createDate"`
}

// IncomeReportAccount ... define a estrutura da conta.
type IncomeReportAccount struct {
	Branch  string `json:"branch"`
	Account string `json:"account"`
}

// IncomeReportBalance ... define a estrutura do saldo.
type IncomeReportBalance struct {
	CalendarYear string  `json:"calendarYear"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Type         string  `json:"type"`
}

// IncomeReportBody ... define a estrutura do corpo do relatório de rendimentos.
type IncomeReportBody struct {
	PayerSource IncomeReportPayerSource `json:"payerSource"`
	Owner       IncomeReportOwner       `json:"owner"`
	Account     IncomeReportAccount     `json:"account"`
	Balances    []IncomeReportBalance   `json:"balances"`
	IncomeFile  string                  `json:"incomeFile"`
	FileType    string                  `json:"fileType"`
}

// IncomeReportResponse ... define a estrutura da resposta da API de relatório de rendimentos.
type IncomeReportResponse struct {
	Version string           `json:"version"`
	Status  string           `json:"status"`
	Body    IncomeReportBody `json:"body"`
}

///PAYMENT

// ValidatePaymentRequest representa a requisição para validar um pagamento.
// No SDK do Bankly este modelo possui o campo Code; adapte-o se o payload da Celcoin exigir outro nome.
type ValidatePaymentRequest struct {
	BarCode *BarcodeData `json:"barCode,omitempty"`
}

// BarcodeData representa os dados do código de barras.
type BarcodeData struct {
	DigitableLine *string `json:"digitable,omitempty"`
	BarCode       *string `json:"barCode,omitempty"`
}

// ValidatePaymentResponse..
type ValidatePaymentResponse struct {
	ID                *string  `json:"id,omitempty"`
	Assignor          *string  `json:"assignor,omitempty"`
	Code              *string  `json:"code,omitempty"`
	Digitable         *string  `json:"digitable,omitempty"`
	Amount            *float64 `json:"amount,omitempty"`
	OriginalAmount    *float64 `json:"originalAmount,omitempty"`
	MinAmount         *float64 `json:"minAmount,omitempty"`
	MaxAmount         *float64 `json:"maxAmount,omitempty"`
	AllowChangeAmount *bool    `json:"allowChangeAmount,omitempty"`
	DueDate           *string  `json:"dueDate,omitempty"`
	SettleDate        *string  `json:"settleDate,omitempty"`
	NextSettle        *bool    `json:"nextSettle,omitempty"`
}

// ConfirmPaymentRequest..
type ConfirmPaymentRequest struct {
	ID          *string  `validate:"required" json:"id,omitempty"`
	Amount      *float64 `validate:"required" json:"amount,omitempty"`
	Description *string  `json:"description,omitempty"`
	BankBranch  *string  `validate:"required" json:"bankBranch,omitempty"`
	BankAccount *string  `validate:"required" json:"bankAccount,omitempty"`
}

// ConfirmPaymentResponse..
type ConfirmPaymentResponse struct {
	AuthenticationCode *string    `json:"authenticationCode,omitempty"`
	SettledDate        *time.Time `json:"settledDate,omitempty"`
}

// FilterPaymentsRequest..
type FilterPaymentsRequest struct {
	BankBranch  *string `validate:"required" json:"bankBranch"`
	BankAccount *string `validate:"required" json:"bankAccount"`
	PageSize    *int    `validate:"required" json:"pageSize"`
	PageToken   *string `json:"pageToken,omitempty"`
}

// FilterPaymentsResponse..
type FilterPaymentsResponse struct {
	NextPageToken *string            `json:"nextPage,omitempty"`
	Data          []*PaymentResponse `json:"data,omitempty"`
}

// DetailPaymentRequest..
type DetailPaymentRequest struct {
	BankBranch         string `validate:"required" json:"bankBranch"`
	BankAccount        string `validate:"required" json:"bankAccount"`
	AuthenticationCode string `validate:"required" json:"authenticationCode"`
}

// PaymentPayer, BusinessHours e Charges (auxiliares, seguindo padrão existente)
type PaymentPayer struct {
	Name           string `json:"name,omitempty"`
	DocumentNumber string `json:"documentNumber,omitempty"`
}

type BusinessHours struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type Charges struct {
	InterestAmountCalculated float64 `json:"interestAmountCalculated,omitempty"`
	FineAmountCalculated     float64 `json:"fineAmountCalculated,omitempty"`
	DiscountAmount           float64 `json:"discountAmount,omitempty"`
}

// ChargeRequest ... define a estrutura da requisição de cobrança.
type ChargeRequest struct {
	TransactionID *string `json:"transactionId"`
	ExternalID    *string `json:"externalId"`
}

// ChargeDebtor ... define a estrutura do devedor.
type ChargeDebtor struct {
	Name         string `json:"name"`
	Document     string `json:"document"`
	PostalCode   string `json:"postalCode"`
	PublicArea   string `json:"publicArea"`
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

// ChargeReceiver ... define a estrutura do recebedor.
type ChargeReceiver struct {
	Name       string `json:"name"`
	Document   string `json:"document"`
	PostalCode string `json:"postalCode"`
	PublicArea string `json:"publicArea"`
	City       string `json:"city"`
	State      string `json:"state"`
	Account    string `json:"account"`
}

// ChargeDiscount ... define a estrutura do desconto.
type ChargeDiscount struct {
	Amount    float64 `json:"amount"`
	Modality  string  `json:"modality"`
	LimitDate string  `json:"limitDate"`
}

// ChargeInstructions ... define a estrutura das instruções.
type ChargeInstructions struct {
	Fine     float64        `json:"fine"`
	Interest float64        `json:"interest"`
	Discount ChargeDiscount `json:"discount"`
}

// ChargeDetails ... define a estrutura dos detalhes do boleto.
type ChargeDetails struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	BankEmissor   string `json:"bankEmissor"`
	BankNumber    string `json:"bankNumber"`
	BankAgency    string `json:"bankAgency"`
	BankAccount   string `json:"bankAccount"`
	BarCode       string `json:"barCode"`
	BankLine      string `json:"bankLine"`
	BankAssignor  string `json:"bankAssignor"`
}

// ChargePix ... define a estrutura do Pix.
type ChargePix struct {
	TransactionID             string `json:"transactionId"`
	TransactionIdentification string `json:"transactionIdentification"`
	Status                    string `json:"status"`
	Key                       string `json:"key"`
	Emv                       string `json:"emv"`
}

// ChargeBody ... define a estrutura do corpo do boleto.
type ChargeBody struct {
	TransactionID string             `json:"transactionId"`
	ExternalID    string             `json:"externalId"`
	Amount        float64            `json:"amount"`
	DueDate       string             `json:"duedate"`
	Status        string             `json:"status"`
	Debtor        ChargeDebtor       `json:"debtor"`
	Receiver      ChargeReceiver     `json:"receiver"`
	Instructions  ChargeInstructions `json:"instructions"`
	Boleto        ChargeDetails      `json:"boleto"`
	Pix           ChargePix          `json:"pix"`
	Split         []interface{}      `json:"split"` // Defina os campos necessários para a estrutura Split, se houver
}

// ChargeResponse ... define a estrutura da resposta do boleto.
type ChargeResponse struct {
	Body    ChargeBody `json:"body"`
	Version string     `json:"version"`
	Status  string     `json:"status"`
}

// =================================================================================
// PAYMENT SERVICE MODELS DEFINITIONS
// =================================================================================

type PaymentCategory int

const (

	// BillPaymentConfirmBasePath ...
	BillPaymentConfirmBasePath = "/baas/v2/billpayment"
	// BillPaymentAuthorizeBasePath ...
	BillPaymentAuthorizeBasePath = "/v5/transactions/billpayments"
	// BillPaymenStatusBasePath ...
	BillPaymenStatusBasePath = "/baas/v2/billpayment"
	// BillPaymentAuthorizePath ...
	BillPaymentAuthorizePath = "authorize"
	// BillPaymentStatusPath ...
	BillPaymentStatusPath = "status"

	// PaymentCategoryConcessionaireAndTaxes
	PaymentCategoryConcessionaireAndTaxes PaymentCategory = 1
	// PaymentCategoryCompensationForm ...
	PaymentCategoryCompensationForm PaymentCategory = 2
)

// GetPaymentRequest ... define a estrutura da requisição para executar consulta de uam cobrança.
type GetPaymentRequest struct {
	ClientRequestID string `json:"clientRequestId"`
	TransactionID   string `json:"id"`
}

// GetPaymentTag ... define a estrutura das tags.
type GetPaymentTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetPaymentBarCodeInfo ... define a estrutura das informações do código de barras.
type GetPaymentBarCodeInfo struct {
	Digitable string `json:"digitable"`
}

// GetPaymentError ... define a estrutura do erro.
type GetPaymentError struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

// GetPaymentResponseBody ... define a estrutura do corpo da resposta de obtenção de pagamento.
type GetPaymentResponseBody struct {
	ID                     string                `json:"id"`
	ClientRequestID        string                `json:"clientRequestId"`
	Account                int                   `json:"account"`
	Amount                 float64               `json:"amount"`
	TransactionIDAuthorize int                   `json:"transactionIdAuthorize"`
	HasOccurrence          bool                  `json:"hasOccurrence"`
	Tags                   []GetPaymentTag       `json:"tags"`
	BarCodeInfo            GetPaymentBarCodeInfo `json:"barCodeInfo"`
	Error                  GetPaymentError       `json:"error"`
	PaymentDate            string                `json:"paymentDate"`
}

// GetPaymentResponse define a estrutura da resposta de obtenção de pagamento.
type GetPaymentResponse struct {
	Body    GetPaymentResponseBody `json:"body"`
	Status  string                 `json:"status"`
	Version string                 `json:"version"`
}

// PaymentCode ... define a estrutura do código de barras.
type PaymentCode struct {
	Type      PaymentCategory `json:"type"`
	Digitable string          `json:"digitable"`
	BarCode   string          `json:"barCode"`
}

// PaymentAuthorizeRequest ... define a estrutura da requisição para autorizar um pagamento.
type PaymentAuthorizeRequest struct {
	ExternalTerminal string      `json:"externalTerminal"` // Terminal de identificação externa do sistema do cliente, Ex: CPF
	ExternalNSU      int         `json:"externalNSU"`      // Identificador da transação do sistema cliente
	BarCode          PaymentCode `json:"barCode"`
}

// PaymentRegisterData ... define a estrutura dos dados de registro do pagamento.
type PaymentRegisterData struct {
	DocumentRecipient       string  `json:"documentRecipient"`
	DocumentPayer           string  `json:"documentPayer"`
	PayDueDate              string  `json:"payDueDate"`
	NextBusinessDay         *string `json:"nextBusinessDay"`
	DueDateRegister         string  `json:"dueDateRegister"`
	AllowChangeValue        bool    `json:"allowChangeValue"`
	Recipient               string  `json:"recipient"`
	Payer                   string  `json:"payer"`
	DiscountValue           float64 `json:"discountValue"`
	InterestValueCalculated float64 `json:"interestValueCalculated"`
	MaxValue                float64 `json:"maxValue"`
	MinValue                float64 `json:"minValue"`
	FineValueCalculated     float64 `json:"fineValueCalculated"`
	OriginalValue           float64 `json:"originalValue"`
	TotalUpdated            float64 `json:"totalUpdated"`
	TotalWithDiscount       float64 `json:"totalWithDiscount"`
	TotalWithAdditional     float64 `json:"totalWithAdditional"`
	TotalPaymentPaid        int     `json:"totalPaymentPaid"`
	TotalValuePaid          float64 `json:"totalValuePaid"`
	MaxPartialsAccepts      int     `json:"maxPartialsAccepts"`
	PaymentSpecies          int     `json:"paymentSpecies"`
	DocumentFinalRecipient  *string `json:"documentFinalRecipient"`
	FinalRecipient          *string `json:"finalRecipient"`
}

// PaymentResponse ... define a estrutura da resposta do pagamento.
type PaymentResponse struct {
	Assignor      *string              `json:"assignor,omitempty"`
	RegisterData  *PaymentRegisterData `json:"registerData,omitempty"`
	SettleDate    *string              `json:"settleDate,omitempty"`
	DueDate       *time.Time           `json:"dueDate,omitempty"`
	EndHour       *string              `json:"endHour,omitempty"`
	IniteHour     *string              `json:"initeHour,omitempty"`
	NextSettle    *string              `json:"nextSettle,omitempty"`
	Digitable     *string              `json:"digitable,omitempty"`
	TransactionID *int                 `json:"transactionId,omitempty"`
	Type          *int                 `json:"type,omitempty"`
	Value         *float64             `json:"value,omitempty"`
	MaxValue      *float64             `json:"maxValue,omitempty"`
	MinValue      *float64             `json:"minValue,omitempty"`
	ErrorCode     *string              `json:"errorCode,omitempty"`
	Message       *string              `json:"message,omitempty"`
	Status        *int                 `json:"status,omitempty"`
}

// ExecPaymentTag define a estrutura das tags.
type ExecPaymentTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ExecPaymentBarCodeInfo define a estrutura das informações do código de barras.
type ExecPaymentBarCodeInfo struct {
	Digitable string `json:"digitable"`
	BarCode   string `json:"barCode"`
}

// ExecPaymentRequest define a estrutura da requisição para executar um pagamento.
type ExecPaymentRequest struct {
	ClientRequestID        string                 `json:"clientRequestId"`
	Amount                 float64                `json:"amount"`
	Account                string                 `json:"account"`
	TransactionIDAuthorize int                    `json:"transactionIdAuthorize"`
	Tags                   []ExecPaymentTag       `json:"tags"`
	BarCodeInfo            ExecPaymentBarCodeInfo `json:"barCodeInfo"`
}

// ExecPaymentResponseBody define a estrutura do corpo da resposta de execução de pagamento.
type ExecPaymentResponseBody struct {
	ID                     string                 `json:"id"`
	ClientRequestID        string                 `json:"clientRequestId"`
	Amount                 float64                `json:"amount"`
	TransactionIDAuthorize int                    `json:"transactionIdAuthorize"`
	Tags                   []ExecPaymentTag       `json:"tags"`
	BarCodeInfo            ExecPaymentBarCodeInfo `json:"barCodeInfo"`
}

// ExecPaymentResponse define a estrutura da resposta de execução de pagamento.
type ExecPaymentResponse struct {
	Body    ExecPaymentResponseBody `json:"body"`
	Status  string                  `json:"status"`
	Version string                  `json:"version"`
}

/* DDA */
// DdaRegisterUserRequest user request for register on DDA

type DdaRegisterUserRequest struct {
	Document        string `validate:"required" json:"document"`
	ClientName      string `validate:"required" json:"clientName"`
	ClientRequestId string `json:"clientRequestId"`
}

type DdaDeleteUserRequest struct {
	Document        string `validate:"required" json:"document"`
	ClientName      string `json:"clientName"`
	ClientRequestId string `json:"clientRequestId"`
}

// DdaRegisterUserResponse response from DDA register user
type DdaRegisterUserResponse struct {
	Status int                          `json:"status,omitempty"` //201,400
	Error  *WebhookError                `json:"error,omitempty"`  // Detalhes do erro, se houver
	Body   *DdaRegisterUserBodyResponse `json:"body,omitempty"`
}

// DdaRegisterUserBodyResponse ...
type DdaRegisterUserBodyResponse struct {
	Document        string `json:"document"`
	ClientRequestId string `json:"clientRequestId"`
	ResponseDate    string `json:"responseDate"`
	Status          string `json:"status"`
	SubscriptionId  string `json:"subscriptionId"`
}
