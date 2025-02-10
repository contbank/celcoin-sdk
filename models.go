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
	CelcoinBankName string = "Contbank S.A. (Celcoin Instituição De Pagamento S.A.)"

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
	//NaturalPersonOnboardingPath ...
	NaturalPersonOnboardingPath string = "/onboarding/v1/onboarding-proposal/natural-person"
	// LegalPersonOnboardingPath ...
	LegalPersonOnboardingPath string = "/onboarding/v1/onboarding-proposal/legal-person"
	// ExternalTransfersPath external transfer (TED)
	ExternalTransfersPath string = "/baas-wallet-transactions-webservice/v1/spb/transfer"
	// InternalTransfersPath internal transfer
	InternalTransfersPath string = "/baas-wallet-transactions-webservice/v1/wallet/internal/transfer"

	// Pix ...
	PixDictPath         string = "/celcoin-baas-pix-dict-webservice/v1/pix/dict/entry"
	PixCashOutPath      string = "/baas-wallet-transactions-webservice/v1/pix/payment"
	PixCashInPath       string = "/pix/v2/receivement/v2"
	PixEmvPath          string = "/pix/v1/emv"
	PixStaticPath       string = "/pix/v1/brcode/static"
	PixCashInStatusPath string = "/pix/v2/receivement/v2/devolution/status"

	// StatementPath ...
	StatementPath string = "/baas-walletreports/v1/wallet/movement"

	// Webhook
	WebhookPath string = "/baas-webhookmanager/v1/webhook"

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

// PixCashOutRequest representa os dados para realizar um Pix Cash-Out.
type PixCashOutRequest struct {
	Amount                    float64     `json:"amount" description:"O valor da transação (required)"`
	VlcpAmount                float64     `json:"vlcpAmount" description:"O valor da compra (Pix Troco)"`
	VldnAmount                float64     `json:"vldnAmount" description:"O valor em dinheiro disponibilizado (Pix Troco)"`
	WithdrawalServiceProvider string      `json:"withdrawalServiceProvider" description:"O Identificador ISPB do serviço de saque (Pix Saque/Troco)"`
	WithdrawalAgentMode       string      `json:"withdrawalAgentMode" description:"Modo do agente de retirada. AGTEC: Estabelecimento Comercial, AGTOT: Entidade Jurídica cuja atividade é a prestação de serviços auxiliares de serviços financeiros, AGPSS: Participante Pix que presta diretamente o serviço de saque."`
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
	InitiationType            string                 `json:"initiationType" `
	TransactionIdentification *string                `json:"transactionIdentification,omitempty" `
	Status                    string                 `json:"status"`
	Amount                    float64                `json:"amount"`
	Currency                  string                 `json:"currency"`
	CreationDate              string                 `json:"creationDate"`
	CompletionDate            *string                `json:"completionDate,omitempty"`
	ErrorDetails              *ErrorDetails          `json:"errorDetails,omitempty"`
	AdditionalInfo            map[string]interface{} `json:"additionalInfo,omitempty"`
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

// PixMerchant representa as informações do comerciante no Pix Cash-in.
type PixMerchant struct {
	PostalCode           string `json:"postalCode" validate:"required"`
	City                 string `json:"city" validate:"required"`
	MerchantCategoryCode int    `json:"merchantCategoryCode"`
	Name                 string `json:"name" validate:"required"`
}

// PixCashInStaticResponse representa a resposta ao realizar um Pix Cash-in por Cobrança Estática.
type PixCashInStaticResponse struct {
	TransactionId             int    `json:"transactionId"`
	EMVQRCode                 string `json:"emvqrcps"`
	TransactionIdentification string `json:"transactionIdentification"`
}

/* TRANSFERS */
// TransfersRequest ...
type TransfersRequest struct {
	Amount          int64                       `validate:"required" json:"amount"`
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
	Amount          int64                        `json:"amount"`
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
	ExternalID             string       `json:"externalId"`
	ExpirationAfterPayment int          `json:"expirationAfterPayment"`
	DueDate                string       `json:"dueDate"`
	Amount                 float64      `json:"amount"`
	Key                    string       `json:"key,omitempty"` // optional
	Debtor                 Debtor       `json:"debtor"`
	Receiver               Receiver     `json:"receiver"`
	Instructions           Instructions `json:"instructions"`
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
	LimitPerPage   *string `json:"LimitPerPage"`
}
