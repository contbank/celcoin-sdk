package celcoin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/contbank/grok"
)

var (
	// ErrEntryNotFound ...
	ErrEntryNotFound = grok.NewError(http.StatusNotFound, "NOT_FOUND", "not found")
	// ErrDuplicateCompany ...
	ErrDuplicateCompany = grok.NewError(http.StatusConflict, "DUPLICATE_COMPANY", "duplicate company")
	// ErrInvalidToken ...
	ErrInvalidToken = grok.NewError(http.StatusConflict, "INVALID_TOKEN", "invalid token")
	// ErrInvalidBusinessSize ...
	ErrInvalidBusinessSize = grok.NewError(http.StatusBadRequest, "INVALID_BUSINESS_SIZE", "invalid business size")
	// ErrEmailAlreadyInUse ...
	ErrEmailAlreadyInUse = grok.NewError(http.StatusBadRequest, "EXISTS_EMAIL", "email already in use")
	// ErrPhoneAlreadyInUse ...
	ErrPhoneAlreadyInUse = grok.NewError(http.StatusBadRequest, "EXISTS_PHONE", "phone already in use")
	// ErrCustomerRegistrationCannotBeReplaced ...
	ErrCustomerRegistrationCannotBeReplaced = grok.NewError(http.StatusConflict, "CUSTOMER_CANNOT_BE_REPLACED", "customer registration cannot be replaced")
	// ErrAccountHolderNotExists ...
	ErrAccountHolderNotExists = grok.NewError(http.StatusBadRequest, "NOT_EXISTS_HOLDER", "account holder not exists")
	// ErrHolderAlreadyHaveAAccount ...
	ErrHolderAlreadyHaveAAccount = grok.NewError(http.StatusConflict, "EXISTS_HOLDER", "holder already have a account")
	// ErrInvalidCorrelationID ...
	ErrInvalidCorrelationID = grok.NewError(http.StatusBadRequest, "INVALID_CORRELATION_ID", "invalid correlation id")
	// ErrInvalidAmount ...
	ErrInvalidAmount = grok.NewError(http.StatusBadRequest, "INVALID_AMOUNT", "invalid amount")
	// ErrInsufficientBalance ...
	ErrInsufficientBalance = grok.NewError(http.StatusBadRequest, "INSUFFICIENT_BALANCE", "insufficient balance")
	// ErrInvalidTransferAuthenticationCode ...
	ErrInvalidTransferAuthenticationCode = grok.NewError(http.StatusBadRequest, "INVALID_TRANSFER_AUTHENTICATION_CODE", "invalid authentication code or transfer request id")
	// ErrInvalidAccountNumber ...
	ErrInvalidAccountNumber = grok.NewError(http.StatusBadRequest, "INVALID_ACCOUNT_NUMBER", "invalid account number")
	// ErrOutOfServicePeriod ...
	ErrOutOfServicePeriod = grok.NewError(http.StatusBadRequest, "OUT_SERVICE_PERIOD", "out of service period")
	// ErrCashoutLimitNotEnough ...
	ErrCashoutLimitNotEnough = grok.NewError(http.StatusBadRequest, "CASHOUT_LIMIT_NOT_ENOUGH", "cashout limit not enough")
	// ErrInvalidParameter ...
	ErrInvalidParameter = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER", "invalid parameter")
	// ErrInvalidParameterLength ...
	ErrInvalidParameterLength = grok.NewError(http.StatusBadRequest, "INVALID_PARAMETER_LENGTH", "invalid parameter length")
	// ErrInvalidUF ...
	ErrInvalidUF = grok.NewError(http.StatusBadRequest, "INVALID_UF", "invalid uf")
	// ErrInvalidAddressNumberLength ...
	ErrInvalidAddressNumberLength = grok.NewError(http.StatusBadRequest, "INVALID_ADDRESS_NUMBER_LENGTH", "invalid address number length")
	// ErrInvalidRegisterNameLength ...
	ErrInvalidRegisterNameLength = grok.NewError(http.StatusBadRequest, "INVALID_REGISTER_NAME_LENGHT", "invalid register name length")
	// ErrInvalidParameterSpecialCharacters ...
	ErrInvalidParameterSpecialCharacters = grok.NewError(http.StatusBadRequest, "INVALID_SPECIAL_CHARACTERS", "invalid parameter with special characters")
	// ErrInvalidSocialNameLength ...
	ErrInvalidSocialNameLength = grok.NewError(http.StatusBadRequest, "INVALID_SOCIAL_NAME_LENGTH", "invalid social name length")
	// ErrInvalidEmailLength ...
	ErrInvalidEmailLength = grok.NewError(http.StatusBadRequest, "INVALID_EMAIL_LENGTH", "invalid email length")
	// ErrInvalidAPIEndpoint ...
	ErrInvalidAPIEndpoint = grok.NewError(http.StatusBadRequest, "INVALID_API_ENDPOINT", "invalid api endpoint")
	// ErrInvalidAPIVersion ...
	ErrInvalidAPIVersion = grok.NewError(http.StatusBadRequest, "INVALID_API_VERSION", "invalid api version")
	// ErrMethodNotAllowed ...
	ErrMethodNotAllowed = grok.NewError(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	// ErrSendDocumentAnalysis ...
	ErrSendDocumentAnalysis = grok.NewError(http.StatusMethodNotAllowed, "SEND_DOCUMENT_ANALYSIS_ERROR", "send document analysis error")
	// ErrGetDocumentAnalysis ...
	ErrGetDocumentAnalysis = grok.NewError(http.StatusMethodNotAllowed, "GET_DOCUMENT_ANALYSIS_ERROR", "get document analysis error")
	// ErrBoletoInvalidStatus ...
	ErrBoletoInvalidStatus = grok.NewError(http.StatusUnprocessableEntity, "INVALID_BOLETO_STATUS", "boleto was in an invalid status")
	// ErrInvalidCardProxy ...
	ErrInvalidCardProxy = grok.NewError(http.StatusBadRequest, "INVALID_CARD_PROXY", "invalid card proxy")
	// ErrBarcodeNotFound ...
	ErrBarcodeNotFound = grok.NewError(http.StatusNotFound, "BARCODE_NOT_FOUND", "bar code not found")
	// ErrPaymentInvalidStatus ...
	ErrPaymentInvalidStatus = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PAYMENT_STATUS", "payment was in an invalid status")
	// ErrDefaultTransfers ...
	ErrDefaultTransfers = grok.NewError(http.StatusConflict, "TRANSFERS_ERROR", "error transfers")
	// ErrDefaultFindTransfers ...
	ErrDefaultFindTransfers = grok.NewError(http.StatusConflict, "FIND_TRANSFERS_ERROR", "error find transfers")
	// ErrDefaultPayment ...
	ErrDefaultPayment = grok.NewError(http.StatusInternalServerError, "PAYMENT_ERROR", "error payment")
	// ErrDefaultBusinessAccounts ...
	ErrDefaultBusinessAccounts = grok.NewError(http.StatusConflict, "BUSINESS_ACCOUNTS_ERROR", "error business accounts")
	// ErrDefaultCorporationBusinessAccounts ...
	ErrDefaultCorporationBusinessAccounts = grok.NewError(http.StatusConflict, "CORPORATION_BUSINESS_ACCOUNTS_ERROR", "error corporation business accounts")
	// ErrDefaultCustomersAccounts ...
	ErrDefaultCustomersAccounts = grok.NewError(http.StatusInternalServerError, "CUSTOMERS_ACCOUNTS_ERROR", "error customers accounts")
	// ErrDefaultCustomersOnboarding ...
	ErrDefaultCustomersOnboarding = grok.NewError(http.StatusInternalServerError, "CUSTOMERS_ONBOARDING_ERROR", "error customers onboarding")
	// ErrDefaultCustomersProposal ...
	ErrDefaultCustomersProposal = grok.NewError(http.StatusInternalServerError, "CUSTOMERS_PROPOSAL_ERROR", "error customers proposal")
	// ErrDefaultBalance ...
	ErrDefaultBalance = grok.NewError(http.StatusInternalServerError, "BALANCE_ERROR", "error balance")
	// ErrDefaultLogin ...
	ErrDefaultLogin = grok.NewError(http.StatusInternalServerError, "LOGIN_ERROR", "error login")
	// ErrDefaultBank ...
	ErrDefaultBank = grok.NewError(http.StatusInternalServerError, "BANK_ERROR", "error bank")
	// ErrDefaultBankStatements ...
	ErrDefaultBankStatements = grok.NewError(http.StatusInternalServerError, "BANK_STATEMENTS_ERROR", "error bank statements")
	// ErrDefaultIncomeReport ...
	ErrDefaultIncomeReport = grok.NewError(http.StatusInternalServerError, "INCOME_REPORT_ERROR", "error income report")
	//ErrDefaultBoletos ...
	ErrDefaultBoletos = grok.NewError(http.StatusInternalServerError, "BOLETOS_ERROR", "error boletos celcoin")
	// ErrDefaultFreshDesk ...
	ErrDefaultFreshDesk = grok.NewError(http.StatusInternalServerError, "FRESH_DESK_ERROR", "error in fresh desk api")
	// ErrFreshDeskTicketNotFound ...
	ErrFreshDeskTicketNotFound = grok.NewError(http.StatusNotFound, "FRESH_DESK_TICKET_NOT_FOUND", "error in fresh desk ticket not found")
	// ErrInvalidRecipientBranch ...
	ErrInvalidRecipientBranch = grok.NewError(http.StatusConflict, "INVALID_RECIPIENT_BRANCH", "invalid recipient branch number")
	// ErrInvalidRecipientAccount ...
	ErrInvalidRecipientAccount = grok.NewError(http.StatusConflict, "INVALID_RECIPIENT_ACCOUNT", "invalid recipient account number")
	// ErrDefaultCard ...
	ErrDefaultCard = grok.NewError(http.StatusInternalServerError, "CARD_ERROR", "error card")
	// ErrDefaultPix ...
	ErrDefaultPix = grok.NewError(http.StatusInternalServerError, "PIX_ERROR", "error pix")
	// ErrKeyNotFound ...
	ErrKeyNotFound = grok.NewError(http.StatusNotFound, "KEY_NOT_FOUND", "key not found")
	// ErrBoletoNotFound ...
	ErrBoletoNotFound = grok.NewError(http.StatusNotFound, "BOLETO_NOT_FOUND", "boleto not found")
	// ErrInvalidQrCodePayload ...
	ErrInvalidQrCodePayload = grok.NewError(http.StatusConflict, "INVALID_QRCODE_PAYLOAD", "invalid qrcode payload")
	// ErrInvalidKeyType ...
	ErrInvalidKeyType = grok.NewError(http.StatusUnprocessableEntity, "INVALID_KEY_TYPE", "invalid key type")
	// ErrInvalidParameterPix ...
	ErrInvalidParameterPix = grok.NewError(http.StatusUnprocessableEntity, "INVALID_PARAMENTER", "invalid parameter")
	// ErrInsufficientBalancePix ...
	ErrInsufficientBalancePix = grok.NewError(http.StatusConflict, "INSUFFICIENT_BALANCE", "insufficient balance")
	// ErrInvalidAccountType ...
	ErrInvalidAccountType = grok.NewError(http.StatusUnprocessableEntity, "INVALID_ACCOUNT_TYPE", "invalid account type")
	// ErrCardActivate ...
	ErrCardActivate = grok.NewError(http.StatusNotModified, "CARD_ACTIVATE_ERROR", "error card activate")
	// ErrCardStatusUpdate ...
	ErrCardStatusUpdate = grok.NewError(http.StatusNotModified, "CARD_STATUS_UPDATE_ERROR", "error update status card")
	// ErrCardDuplicate ...
	ErrCardDuplicate = grok.NewError(http.StatusNotModified, "CARD_DUPLICATE_ERROR", "error card duplicate")
	// ErrCardPasswordUpdate ...
	ErrCardPasswordUpdate = grok.NewError(http.StatusNotModified, "CARD_PASSWORD_UPDATE_ERROR", "error update password card")
	// ErrInvalidPassword ...
	ErrInvalidPassword = grok.NewError(http.StatusUnauthorized, "INVALID_PASSWORD", "invalid password")
	// ErrUnauthorized ...
	ErrUnauthorized = grok.NewError(http.StatusUnauthorized, "UNAUTHORIZED", "error unauthorized")
	// ErrBlockedByRiskAnalysis ...
	ErrBlockedByRiskAnalysis = grok.NewError(http.StatusForbidden, "BLOCKED", "blocked")
	// ErrBankslipAlreadyCancelled ...
	ErrBankslipAlreadyCancelled = grok.NewError(http.StatusBadRequest, "BANKSLIP_HAS_ALREADY_BEEN_CANCELED", "bankslip has already been canceled")
	// ErrBankslipLimitQuantityExceeded ...
	ErrBankslipLimitQuantityExceeded = grok.NewError(http.StatusBadRequest, "LIMIT_QUANTITY_EXCEEDED", "maximum quantity limit per month exceeded")
	// ErrBankslipLimitNotEnough ...
	ErrBankslipLimitNotEnough = grok.NewError(http.StatusBadRequest, "LIMIT_NOT_ENOUGH", "limit not enough")
	// ErrAccountWasClosed ...
	ErrAccountWasClosed = grok.NewError(http.StatusBadRequest, "ACCOUNT_WAS_CLOSED", "account was closed")
	// ErrInvalidDocument ...
	ErrInvalidDocument = grok.NewError(http.StatusBadRequest, "ACCOUNT_DOCUMENT_INVALID", "invalid document for this account")
	// ErrInvalidCardName ...
	ErrInvalidCardName = grok.NewError(http.StatusBadRequest, "INVALID_CARD_NAME", "invalid card name")
	// ErrInvalidIdentifier ...
	ErrInvalidIdentifier = grok.NewError(http.StatusBadRequest, "INVALID_IDENTIFIER", "invalid identifier")
	// ErrCardAlreadyActivated ...
	ErrCardAlreadyActivated = grok.NewError(http.StatusConflict, "CARD_ALREADY_ACTIVATED", "card already activated")
	// ErrOperationNotAllowedCardStatus ...
	ErrOperationNotAllowedCardStatus = grok.NewError(http.StatusMethodNotAllowed, "OPERATION_NOT_ALLOWED", "operation not allowed for current card status")
	// ErrInvalidIncomeReportCalendar ...
	ErrInvalidIncomeReportCalendar = grok.NewError(http.StatusBadRequest, "INVALID_INCOME_REPORT_CALENDAR", "invalid income report calendar")
	// ErrInvalidIncomeReportParameter ...
	ErrInvalidIncomeReportParameter = grok.NewError(http.StatusBadRequest, "INVALID_INCOME_REPORT_PARAMETER", "invalid income report parameter")
	// ErrDefaultCancelCustomersAccounts ...
	ErrDefaultCancelCustomersAccounts = grok.NewError(http.StatusConflict, "CANCEL_CUSTOMERS_ACCOUNTS_ERROR", "error cancel customers accounts")
	// ErrAccountNonZeroBalance ...
	ErrAccountNonZeroBalance = grok.NewError(http.StatusConflict, "ACCOUNT_NON_ZERO_BALANCE", "error account non zero balance")
	// ErrAccountAlreadyBeenCanceled ...
	ErrAccountAlreadyBeenCanceled = grok.NewError(http.StatusUnprocessableEntity, "ACCOUNT_ALREADY_BEEN_CANCELED", "error account already been canceled")
	// ErrAccountNotFound ...
	ErrAccountNotFound = grok.NewError(http.StatusNotFound, "ACCOUNT_NOT_FOUND", "error account not found")
	// ErrSenderAccountStatusNotAllowCashOut ...
	ErrSenderAccountStatusNotAllowCashOut = grok.NewError(http.StatusBadRequest, "SENDER_ACCOUNT_STATUS_NOT_ALLOW_CASH_OUT", "error sender account status not allow cash out")
	// ErrRecipientAccountStatusNotAllowCashIn ...
	ErrRecipientAccountStatusNotAllowCashIn = grok.NewError(http.StatusBadRequest, "RECIPIENT_ACCOUNT_STATUS_NOT_ALLOW_CASH_IN", "error recipient account status not allow cash in")
	// ErrSenderAccountNotFound ...
	ErrSenderAccountNotFound = grok.NewError(http.StatusBadRequest, "SENDER_ACCOUNT_NOT_FOUND", "error sender account not found")
	// ErrRecipientAccountNotFound ...
	ErrRecipientAccountNotFound = grok.NewError(http.StatusBadRequest, "RECIPIENT_ACCOUNT_NOT_FOUND", "error recipient account not found")
	// ErrTimeout ...
	ErrTimeout = grok.NewError(http.StatusBadRequest, "TIMEOUT", "timeout")
	// ErrInvalidBankBranch ...
	ErrInvalidBankBranch = grok.NewError(http.StatusBadRequest, "INVALID_BANK_BRANCH", "error invalid bank branch")
	// ErrInvalidBankAccount ...
	ErrInvalidBankAccount = grok.NewError(http.StatusBadRequest, "INVALID_BANK_ACCOUNT", "error invalid bank account number")
	// ErrInvalidBankAccountOrBranch ...
	ErrInvalidBankAccountOrBranch = grok.NewError(http.StatusBadRequest, "INVALID_BANK_ACCOUNT_OR_BRANCH", "error invalid account number or branch")
	// ErrRecipientAccountDoesNotMatchTheDocument ...
	ErrRecipientAccountDoesNotMatchTheDocument = grok.NewError(http.StatusBadRequest, "RECIPIENT_ACCOUNT_DOES_NOT_MATCH_THE_DOCUMENT", "error recipient account does not match the document")
	// ErrSenderAccountDoesNotMatchTheDocument ...
	ErrSenderAccountDoesNotMatchTheDocument = grok.NewError(http.StatusBadRequest, "SENDER_ACCOUNT_DOES_NOT_MATCH_THE_DOCUMENT", "error sender account does not match the document")
	// ErrTransferWasReproved ...
	ErrTransferWasReproved = grok.NewError(http.StatusBadRequest, "TRANSFER_WAS_REPROVED", "error transfer was reproved")
	// ErrTransferAmountNotReserved ...
	ErrTransferAmountNotReserved = grok.NewError(http.StatusBadRequest, "TRANSFER_AMOUNT_NOT_RESERVED", "error transfer amount not reserved")
	// ErrTransferOrderNotProcessed ...
	ErrTransferOrderNotProcessed = grok.NewError(http.StatusBadRequest, "TRANSFER_ORDER_NOT_PROCESSED", "error transfer order not processed")
	// ErrInternalTransferNotCompleted ...
	ErrInternalTransferNotCompleted = grok.NewError(http.StatusBadRequest, "INTERNAL_TRANSFER_NOT_COMPLETED", "error internal transfer not completed")
	// ErrScheduleNotAllowed ...
	ErrScheduleNotAllowed = grok.NewError(http.StatusBadRequest, "SCHEDULE_NOT_ALLOWED", "error schedule not allowed")
	// ErrInvalidEndToEndId ...
	ErrInvalidEndToEndId = grok.NewError(http.StatusBadRequest, "INVALID_END_TO_END_ID", "error invalid end to end id")
	//ErrInvalidIssuerAddress ...
	ErrInvalidIssuerAddress = grok.NewError(http.StatusBadRequest, "INVALID_ISSUER_ADDRESS", "error invalid issuer address")
	// ErrScouterQuantity ...
	ErrScouterQuantity = grok.NewError(http.StatusUnprocessableEntity, "SCOUTER_QUANTITY_ERROR", "error max boleto reached")
	// ErrAmountNotAllowed ...
	ErrAmountNotAllowed = grok.NewError(http.StatusBadRequest, "AMOUNT_NOT_ALLOWED", "error amount not allowed")
	// ErrInvalidName ...
	ErrInvalidName = grok.NewError(http.StatusBadRequest, "INVALID_NAME", "error invalid name")
	// ErrSimpleBusinessNotAllowed ...
	ErrSimpleBusinessNotAllowed = grok.NewError(http.StatusMethodNotAllowed, "SIMPLE_BUSINESS_NOT_ALLOWED", "simple business not allowed")
	// ErrCorporationBusinessNotAllowed ...
	ErrCorporationBusinessNotAllowed = grok.NewError(http.StatusMethodNotAllowed, "CORPORATION_BUSINESS_NOT_ALLOWED", "corporation business not allowed")
)

// CelcoinError ...
type CelcoinError ErrorModel

// Error ..
type Error struct {
	ErrorKey  string
	GrokError *grok.Error
}

type ErrorCard struct {
	Code         string
	Messages     []string
	Metadata     interface{}
	PropertyName string
	Reasons      []interface{}
}

type CelcoinCardError struct {
	ErrorsCard ErrorCard
}

var errorList = []Error{
	{
		ErrorKey:  "INVALID_PERSONAL_BUSINESS_SIZE",
		GrokError: ErrInvalidBusinessSize,
	},
	{
		ErrorKey:  "EMAIL_ALREADY_IN_USE",
		GrokError: ErrEmailAlreadyInUse,
	},
	{
		ErrorKey:  "PHONE_ALREADY_IN_USE",
		GrokError: ErrPhoneAlreadyInUse,
	},
	{
		ErrorKey:  "CUSTOMER_REGISTRATION_CANNOT_BE_REPLACED",
		GrokError: ErrCustomerRegistrationCannotBeReplaced,
	},
	{
		ErrorKey:  "ACCOUNT_HOLDER_NOT_EXISTS",
		GrokError: ErrAccountHolderNotExists,
	},
	{
		ErrorKey:  "HOLDER_ALREADY_HAVE_A_ACCOUNT",
		GrokError: ErrHolderAlreadyHaveAAccount,
	},
	{
		ErrorKey:  "SCOUTER_QUANTITY",
		GrokError: ErrScouterQuantity,
	},
	{
		ErrorKey:  "BANKSLIP_SETTLEMENT_STATUS_VALIDATE",
		GrokError: ErrBoletoInvalidStatus,
	},
	{
		ErrorKey:  "BAR_CODE_NOT_FOUND",
		GrokError: ErrBarcodeNotFound,
	},
	{
		ErrorKey:  "INVALID_PARAMETER",
		GrokError: ErrInvalidParameter,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_LENGTH",
		GrokError: ErrInvalidParameterLength,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_SPECIAL_CHARACTERS",
		GrokError: ErrInvalidParameterSpecialCharacters,
	},
	{
		ErrorKey:  "INVALID_ADDRESS_NUMBER_LENGTH",
		GrokError: ErrInvalidAddressNumberLength,
	},
	{
		ErrorKey:  "INVALID_REGISTER_NAME_LENGTH",
		GrokError: ErrInvalidRegisterNameLength,
	},
	{
		ErrorKey:  "INVALID_SOCIAL_NAME_LENGTH",
		GrokError: ErrInvalidSocialNameLength,
	},
	{
		ErrorKey:  "INVALID_EMAIL_LENGTH",
		GrokError: ErrInvalidEmailLength,
	},
	{
		ErrorKey:  "HOLDER_HAS_SOME_ACCOUNTS_WITH_NON_ZERO_BALANCE",
		GrokError: ErrAccountNonZeroBalance,
	},
	{
		ErrorKey:  "HOLDER_HAS_ALREADY_BEEN_CANCELED",
		GrokError: ErrAccountAlreadyBeenCanceled,
	},
}

// CelcoinTransferError ..
type CelcoinTransferError KeyValueErrorModel

// TransferError ..
type TransferError struct {
	celcoinTransferError CelcoinTransferError
	grokError            *grok.Error
}

var transferErrorList = []TransferError{
	{
		celcoinTransferError: CelcoinTransferError{Key: "x-correlation-id"},
		grokError:            ErrInvalidCorrelationID,
	},
	{
		celcoinTransferError: CelcoinTransferError{Key: "$.amount"},
		grokError:            ErrInvalidAmount,
	},
	{
		celcoinTransferError: CelcoinTransferError{Key: "INSUFFICIENT_BALANCE"},
		grokError:            ErrInsufficientBalance,
	},
	{
		celcoinTransferError: CelcoinTransferError{Key: "CASH_OUT_NOT_ALLOWED_OUT_OF_BUSINESS_PERIOD"},
		grokError:            ErrOutOfServicePeriod,
	},
	{
		celcoinTransferError: CelcoinTransferError{Key: "CASHOUT_LIMIT_NOT_ENOUGH"},
		grokError:            ErrCashoutLimitNotEnough,
	},
	{
		celcoinTransferError: CelcoinTransferError{Key: "Recipient.Branch"},
		grokError:            ErrInvalidRecipientBranch,
	},
	{
		celcoinTransferError: CelcoinTransferError{Key: "Recipient.Account"},
		grokError:            ErrInvalidRecipientAccount,
	},
}

// FindError Find errors.
func FindError(code string, messages ...string) *Error {
	code = verifyInvalidParameter(code, messages)

	for _, v := range errorList {
		if v.ErrorKey == code {
			return &v
		}
	}

	return &Error{
		ErrorKey:  code,
		GrokError: grok.NewError(http.StatusConflict, code, messages...),
	}
}

// FindBalanceError ... find errors for celcoin balance api
func FindBalanceError(code string, messages string) *Error {
	code = mapBalanceErrorCode(code, messages)

	for _, v := range errorList {
		if v.ErrorKey == code {
			return &v
		}
	}

	return &Error{
		ErrorKey:  code,
		GrokError: grok.NewError(http.StatusConflict, code, messages),
	}
}

// FindErrorByErrorModel ..
func FindErrorByErrorModel(response ErrorModel) *Error {
	if response.Code != "" {
		return FindError(response.Code, response.Messages...)
	}
	return &Error{
		ErrorKey:  response.Key,
		GrokError: grok.NewError(http.StatusBadRequest, response.Key, response.Value),
	}
}

// verifyInvalidParameter Find the correspondent error message.
func verifyInvalidParameter(code string, messages []string) string {
	if code == "INVALID_PARAMETER" {
		for _, m := range messages {
			if strings.Contains(strings.ToLower(m), "length of 'building number'") {
				return "INVALID_ADDRESS_NUMBER_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "length of 'register name'") {
				return "INVALID_REGISTER_NAME_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "length of 'social name'") {
				return "INVALID_SOCIAL_NAME_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "length of 'email'") {
				return "INVALID_EMAIL_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "not allowed to include numbers or special characters") {
				return "INVALID_PARAMETER_SPECIAL_CHARACTERS"
			} else if strings.Contains(strings.ToLower(m), "not allowed to include special characters") {
				return "INVALID_PARAMETER_SPECIAL_CHARACTERS"
			} else if strings.Contains(strings.ToLower(m), "length of") {
				return "INVALID_PARAMETER_LENGTH"
			} else if strings.Contains(strings.ToLower(m), "invalid brazilian state acronym") {
				return "INVALID_UF"
			}
		}
	}
	return code
}

// mapBalanceErrorCode ... mapeia os códigos de erro para mensagens específicas para api balance da celcoin.
func mapBalanceErrorCode(code string, message string) string {
	lowerMessage := strings.ToLower(message)
	switch code {
	case "CBE073":
		if strings.Contains(lowerMessage, "informar pelo menos um dos campos") {
			return "MISSING_REQUIRED_FIELDS"
		}
	case "CBE039":
		if strings.Contains(lowerMessage, "account invalido") {
			return "INVALID_ACCOUNT"
		}
	case "CBE040":
		if strings.Contains(lowerMessage, "documentnumber invalido") {
			return "INVALID_DOCUMENT_NUMBER"
		}
	case "CBE041":
		if strings.Contains(lowerMessage, "tamanho maximo de 20 caracteres") {
			return "INVALID_ACCOUNT_LENGTH"
		}
	case "CBE042":
		if strings.Contains(lowerMessage, "tamanho maximo de 14 caracteres") {
			return "INVALID_DOCUMENT_NUMBER_LENGTH"
		}
	case "CBE089":
		if strings.Contains(lowerMessage, "conta esta bloqueada") {
			return "ACCOUNT_BLOCKED"
		}
	case "CBE090":
		if strings.Contains(lowerMessage, "conta esta encerrada") {
			return "ACCOUNT_CLOSED"
		}
	}
	return code // Retorna o código original se nenhuma correspondência for encontrada.
}

// errorIncomeReportList ...
var errorIncomeReportList = []Error{
	{
		ErrorKey:  "INVALID_CALENDAR_FOR_INCOME_REPORT",
		GrokError: ErrInvalidIncomeReportCalendar,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_INCOME_REPORT",
		GrokError: ErrInvalidIncomeReportParameter,
	},
}

// FindIncomeReportError Find income report errors.
func FindIncomeReportError(code string, messages ...string) *grok.Error {
	code = verifyInvalidIncomeReportParameter(code, messages)

	for _, v := range errorCardList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusInternalServerError, code, messages...)
}

// verifyInvalidIncomeReportParameter Find the correspondent error message for income reports.
func verifyInvalidIncomeReportParameter(code string, messages []string) string {
	if code == "CALENDAR_NOT_ALLOWED" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "calendar informed is not allowed"):
				return "INVALID_CALENDAR_FOR_INCOME_REPORT"
			default:
				return "INVALID_PARAMETER_INCOME_REPORT"
			}
		}
	}
	return code
}

var errorCardList = []Error{
	{
		ErrorKey:  "INVALID_CARD_PASSWORD",
		GrokError: ErrInvalidPassword,
	},
	{
		ErrorKey:  "OPERATION_NOT_ALLOWED_FOR_CURRENT_CARD_STATUS",
		GrokError: ErrOperationNotAllowedCardStatus,
	},
	{
		ErrorKey:  "CARD_ALREADY_ACTIVATED",
		GrokError: ErrCardAlreadyActivated,
	},
	{
		ErrorKey:  "INVALID_CARD_NAME_EMPTY",
		GrokError: ErrInvalidCardName,
	},
	{
		ErrorKey:  "INVALID_DOCUMENT_NUMBER_EMPTY",
		GrokError: ErrInvalidIdentifier,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_CARD",
		GrokError: ErrInvalidParameter,
	},
}

// FindCardError Find cards errors.
func FindCardError(code string, messages ...string) *grok.Error {
	code = verifyInvalidCardParameter(code, messages)

	for _, v := range errorCardList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusConflict, code, messages...)
}

// verifyInvalidCardParameter Find the correspondent error message for Cards.
func verifyInvalidCardParameter(code string, messages []string) string {
	if code == "INVALID_PARAMETER" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "card name"):
				return "INVALID_CARD_NAME_EMPTY"
			case strings.Contains(strings.ToLower(m), "document number"):
				return "INVALID_DOCUMENT_NUMBER_EMPTY"
			default:
				return "INVALID_PARAMETER_CARD"
			}
		}
	} else if code == "009" {
		return "OPERATION_NOT_ALLOWED_FOR_CURRENT_CARD_STATUS"
	} else if code == "011" {
		return "INVALID_CARD_PASSWORD"
	} else if code == "021" {
		return "CARD_ALREADY_ACTIVATED"
	}
	return code
}

var errorPixList = []Error{
	{
		ErrorKey:  "ENTRY_NOT_FOUND",
		GrokError: ErrKeyNotFound,
	},
	{
		ErrorKey:  "INVALID_QRCODE_PAYLOAD_CONTENT_TO_DECODE",
		GrokError: ErrInvalidQrCodePayload,
	},
	{
		ErrorKey:  "INVALID_KEY_TYPE",
		GrokError: ErrInvalidKeyType,
	},
	{
		ErrorKey:  "INVALID_PARAMETER_PIX",
		GrokError: ErrInvalidParameterPix,
	},
	{
		ErrorKey:  "INSUFFICIENT_BALANCE",
		GrokError: ErrInsufficientBalancePix,
	},
	{
		ErrorKey:  "INVALID_ACCOUNT_TYPE",
		GrokError: ErrInvalidAccountType,
	},
	{
		ErrorKey:  "SENDER_ACCOUNT_STATUS_NOT_ALLOW_CASH_OUT",
		GrokError: ErrSenderAccountStatusNotAllowCashOut,
	},
	{
		ErrorKey:  "RECIPIENT_ACCOUNT_STATUS_NOT_ALLOW_CASH_IN",
		GrokError: ErrRecipientAccountStatusNotAllowCashIn,
	},
	{
		ErrorKey:  "INVALID_RECIPIENT_ACCOUNT",
		GrokError: ErrInvalidRecipientAccount,
	},
	{
		ErrorKey:  "SENDER_ACCOUNT_NOT_FOUND",
		GrokError: ErrSenderAccountNotFound,
	},
	{
		ErrorKey:  "RECIPIENT_ACCOUNT_NOT_FOUND",
		GrokError: ErrRecipientAccountNotFound,
	},
	{
		ErrorKey:  "CASHOUT_LIMIT_NOT_ENOUGH",
		GrokError: ErrCashoutLimitNotEnough,
	},
	{
		ErrorKey:  "TIMEOUT",
		GrokError: ErrTimeout,
	},
	{
		ErrorKey:  "INVALID_BANK_BRANCH",
		GrokError: ErrInvalidBankBranch,
	},
	{
		ErrorKey:  "INVALID_BANK_ACCOUNT",
		GrokError: ErrInvalidBankAccount,
	},
	{
		ErrorKey:  "RECIPIENT_ACCOUNT_DOES_NOT_MATCH_THE_DOCUMENT",
		GrokError: ErrRecipientAccountDoesNotMatchTheDocument,
	},
	{
		ErrorKey:  "SENDER_ACCOUNT_DOES_NOT_MATCH_THE_DOCUMENT",
		GrokError: ErrSenderAccountDoesNotMatchTheDocument,
	},
	{
		ErrorKey:  "TRANSFER_WAS_REPROVED",
		GrokError: ErrTransferWasReproved,
	},
	{
		ErrorKey:  "TRANSFER_AMOUNT_NOT_RESERVED",
		GrokError: ErrTransferAmountNotReserved,
	},
	{
		ErrorKey:  "TRANSFER_ORDER_NOT_PROCESSED",
		GrokError: ErrTransferOrderNotProcessed,
	},
	{
		ErrorKey:  "INTERNAL_TRANSFER_NOT_COMPLETED",
		GrokError: ErrInternalTransferNotCompleted,
	},
	{
		ErrorKey:  "ACCOUNT_NOT_FOUND",
		GrokError: ErrAccountNotFound,
	},
	{
		ErrorKey:  "SCHEDULE_NOT_ALLOWED",
		GrokError: ErrScheduleNotAllowed,
	},
	{
		ErrorKey:  "INVALID_END_TO_END_ID",
		GrokError: ErrInvalidEndToEndId,
	},
}

func verifyInvalidPixParameter(code string, messages []string) string {
	if code == "INVALID_PARAMETER" {
		for _, m := range messages {
			switch {
			case strings.Contains(strings.ToLower(m), "addressing key value does not match with addressing key type"):
				return "INVALID_KEY_TYPE"
			case strings.Contains(strings.ToLower(m), "sender.account.type"):
				return "INVALID_ACCOUNT_TYPE"
			default:
				return "INVALID_PARAMETER_PIX"
			}
		}
	}
	return code
}

// FindPixError
func FindPixError(code string, messages ...string) *grok.Error {
	code = verifyInvalidPixParameter(code, messages)

	for _, v := range errorPixList {
		if v.ErrorKey == code {
			return v.GrokError
		}
	}

	return grok.NewError(http.StatusConflict, code, messages...)
}

// ParseErr ..
func ParseErr(err error) (*Error, bool) {
	celcoinErr, ok := err.(*Error)
	return celcoinErr, ok
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"Key: %s - Messages: %s",
		e.ErrorKey,
		strings.Join(e.GrokError.Messages, "\n"),
	)
}

// errorBoletoList ...
var errorBoletoList = []Error{
	{
		ErrorKey:  "ACCOUNT_VALIDATE",
		GrokError: ErrInvalidBankAccountOrBranch,
	},
	{
		ErrorKey:  "REGISTEREDNAME_INVALID",
		GrokError: ErrInvalidName,
	},
	{
		ErrorKey:  "ADDRESS_INVALID",
		GrokError: ErrInvalidIssuerAddress,
	},
	{
		ErrorKey:  "ACCOUNT_INTERNAL_ERROR",
		GrokError: ErrDefaultBoletos,
	},
	{
		ErrorKey:  "SCOUTER_QUANTITY",
		GrokError: ErrScouterQuantity,
	},
	{
		ErrorKey:  "SCOUTER_MAXIMUM_AMOUNT",
		GrokError: ErrAmountNotAllowed,
	},
	{
		ErrorKey:  "SCOUTER_MINIMUM_AMOUNT",
		GrokError: ErrAmountNotAllowed,
	},
	{
		ErrorKey:  "INVALID_PARAMETER",
		GrokError: ErrInvalidParameter,
	},
	{
		ErrorKey:  "BANKSLIP_UNAUTHORIZED",
		GrokError: ErrUnauthorized,
	},
	{
		ErrorKey:  "BLOCKED_BY_RISK_ANALYSIS",
		GrokError: ErrBlockedByRiskAnalysis,
	},
	{
		ErrorKey:  "BANKSLIP_HAS_ALREADY_BEEN_CANCELED",
		GrokError: ErrBankslipAlreadyCancelled,
	},
	{
		ErrorKey:  "LIMIT_QUANTITY_EXCEEDED",
		GrokError: ErrBankslipLimitQuantityExceeded,
	},
	{
		ErrorKey:  "LIMIT_NOT_ENOUGH",
		GrokError: ErrBankslipLimitNotEnough,
	},
	{
		ErrorKey:  "ACCOUNT_WAS_CLOSED",
		GrokError: ErrAccountWasClosed,
	},
	{
		ErrorKey:  "ACCOUNT_DOCUMENT_INVALID",
		GrokError: ErrInvalidDocument,
	},
}

// OnboardingErrorMessages ... mapeia os códigos de erro para suas respectivas mensagens
var OnboardingErrorMessages = map[string]string{
	"OBE001": "Token de autorização não enviado.",
	"OBE002": "Token enviado está no formato incorreto.",
	"OBE003": "Token inválido.",
	"OBE004": "Token expirado.",
	"OBE005": "Usuario não encontrado.",
	"OBE006": "Cliente não possui produto Onboarding ativo.",
	"OBE007": "O campo clientCode é obrigatório.",
	"OBE008": "O campo documentNumber é obrigatório e deve ser um CPF válido.",
	"OBE009": "O campo documentNumber é obrigatório e deve ser um CNPJ válido.",
	"OBE010": "O campo phoneNumber ou contactNumber é obrigatório e deve ser um telefone válido.",
	"OBE011": "O campo email é obrigatório e deve ser um email válido.",
	"OBE012": "O campo motherName é obrigatório e deve ser completo.",
	"OBE013": "O campo fullName é obrigatório e deve ser completo.",
	"OBE014": "O campo fullName possui tamanho máximo de 120 caracteres.",
	"OBE015": "socialName inválido.",
	"OBE016": "O campo birthDate é obrigatório e deve ser no formato (DD-MM-YYYY).",
	"OBE017": "O campo address é obrigatório.",
	"OBE018": "O campo onboardingType é obrigatório e deve conter um tipo válido.",
	"OBE019": "O campo postalCode é obrigatório e deve ser um CEP existente.",
	"OBE020": "O campo street é obrigatório deve respeitar o limite de caracteres e conter um formato de texto válido.",
	"OBE021": "Number inválido.",
	"OBE022": "AddressComplement inválido.",
	"OBE023": "O campo neighborhood é obrigatório e deve conter um formato de texto válido.",
	"OBE024": "O campo city é obrigatório e deve conter um formato de texto válido.",
	"OBE025": "O campo state é obrigatório e deve ser uma estado valido.",
	"OBE026": "O campo businessEmail é obrigatório e deve ser um email válido.",
	"OBE027": "O campo businessName é obrigatório e deve conter um formato de texto válido.",
	"OBE028": "O campo tradingName é obrigatório deve respeitar o limite de caracteres e conter um formato de texto válido.",
	"OBE029": "O campo owner.documentNumber é obrigatório e deve ser um CPF ou CNPJ válido.",
	"OBE030": "O campo owner.name é obrigatório e deve ser completo.",
	"OBE031": "O campo owner.email é obrigatório e deve ser um email válido.",
	"OBE032": "O campo owner.address é obrigatório.",
	"OBE033": "Cadastro não permitido para menores de idade.",
	"OBE034": "Formato do JSON esta fora do padrão. Verifique a documentação.",
	"OBE035": "Não foi possivel realizar essa operação. Tente novamente mais tarde.",
	"OBE036": "CompanyType inválido.",
	"OBE037": "O campo Owners deve conter um array de no mínimo 1 e máximo 10.",
	"OBE038": "Owners não podem ser duplicados.",
	"OBE039": "O campo businessAddress é obrigatório.",
	"OBE040": "O campo ownerType é obrigatório e deve conter um valor válido.",
	"OBE041": "BackgroundCheck não encontrado ou com status diferente de pendente.",
	"OBE042": "Erro ao atualizar backgroundCheck.",
	"OBE043": "Documentscopy não encontrado ou com status diferente de pendente.",
	"OBE044": "Erro ao atualizar documentscopy.",
	"OBE045": "Status da proposta inexistente verifique a documentação por favor.",
	"OBE046": "Data inválida.",
	"OBE047": "Limite inserido inválido. Os campos limit ou limitPerPage devem ter valores entre 1 e 200.",
	"OBE048": "O campo documentNumber deve ser um CPF ou CNPJ válido.",
	"OBE049": "Não foi encontrada nenhuma proposta referente aos dados informados.",
	"OBE050": "A data inicial não pode ser maior que a data final.",
	"OBE051": "Ao não enviar o proposalId os campos data inicial e a data final são obrigatórios.",
	"OBE052": "O intervalo de dias entre a data inicial e a data final não deve ser maior que {0} dias.",
	"OBE053": "O campo ownerType deve conter pelo menos um sócio ou representante",
	"OBE054": "ProposalId e clientCode não enviados. Ao menos um desses parametros deve ser enviado.",
	"OBE055": "Não foram encontrados arquivos para o proposalId ou clientCode informado(s).",
	"OBE056": "Não foram encontradas documentoscopias referentes ao proposalId ou clientCode enviado.",
	"OBE057": "Ocorreu um erro ao buscar documentos.",
	"OBE058": "ClientType inválido.",
	"OBE059": "SourceType inválido.",
	"OBE060": "O campo clientId é obrigatório.",
	"OBE061": "Source inválido.",
	"OBE062": "ClientCode já vinculado a outra proposta, esse campo deve ser único por proposta.",
	"OBE063": "Não foram encontrados registros para a sua requisição.",
	"OBE064": "Já existe uma proposta em aberto para esse documentNumber.",
	"OBE065": "O campo dateFrom é obrigatório.",
	"OBE066": "O campo dateTo é obrigatório.",
	"OBE067": "O campo partner.partnerName deve conter um valor válido.",
	"OBE068": "O campo partner.parameter deve ser preenchido.",
	"OBE069": "Ao enviar os campos partner.parameters, o campo partner.partnerName deve ser obrigatório.",
	"OBE070": "Não foram encontrados dados, pois o usuário ainda não iniciou a jornada webview. Tente novamente mais tarde.",
	"OBE071": "Ocorreu um erro ao consultar parceiro. Favor tentar novamente mais tarde.",
	"OIE999": "Ocorreu um erro interno durante a chamada da api",
}

// OnboardingErrorMappings ... mapeia os códigos de erro do parceiro Celcoin para os códigos de erro do Contbank com descrição
var OnboardingErrorMappings = map[string]struct {
	ContbankCode string
	Description  string
}{
	"OBE001": {"AUTH_TOKEN_NOT_SENT", "Token de autorização não enviado."},
	"OBE002": {"INVALID_AUTH_TOKEN_FORMAT", "Token enviado está no formato incorreto."},
	"OBE003": {"INVALID_AUTH_TOKEN", "Token inválido."},
	"OBE004": {"EXPIRED_AUTH_TOKEN", "Token expirado."},
	"OBE005": {"USER_NOT_FOUND", "Usuario não encontrado."},
	"OBE006": {"NO_ACTIVE_ONBOARDING_PRODUCT", "Cliente não possui produto Onboarding ativo."},
	"OBE007": {"CLIENT_CODE_REQUIRED", "O campo clientCode é obrigatório."},
	"OBE008": {"INVALID_CPF", "O campo documentNumber é obrigatório e deve ser um CPF válido."},
	"OBE009": {"INVALID_CNPJ", "O campo documentNumber é obrigatório e deve ser um CNPJ válido."},
	"OBE010": {"INVALID_PHONE_NUMBER", "O campo phoneNumber ou contactNumber é obrigatório e deve ser um telefone válido."},
	"OBE011": {"INVALID_EMAIL", "O campo email é obrigatório e deve ser um email válido."},
	"OBE012": {"MOTHER_NAME_REQUIRED", "O campo motherName é obrigatório e deve ser completo."},
	"OBE013": {"FULL_NAME_REQUIRED", "O campo fullName é obrigatório e deve ser completo."},
	"OBE014": {"FULL_NAME_TOO_LONG", "O campo fullName possui tamanho máximo de 120 caracteres."},
	"OBE015": {"INVALID_SOCIAL_NAME", "socialName inválido."},
	"OBE016": {"INVALID_BIRTH_DATE", "O campo birthDate é obrigatório e deve ser no formato (DD-MM-YYYY)."},
	"OBE017": {"ADDRESS_REQUIRED", "O campo address é obrigatório."},
	"OBE018": {"INVALID_ONBOARDING_TYPE", "O campo onboardingType é obrigatório e deve conter um tipo válido."},
	"OBE019": {"INVALID_POSTAL_CODE", "O campo postalCode é obrigatório e deve ser um CEP existente."},
	"OBE020": {"INVALID_STREET", "O campo street é obrigatório deve respeitar o limite de caracteres e conter um formato de texto válido."},
	"OBE021": {"INVALID_NUMBER", "Number inválido."},
	"OBE022": {"INVALID_ADDRESS_COMPLEMENT", "AddressComplement inválido."},
	"OBE023": {"INVALID_NEIGHBORHOOD", "O campo neighborhood é obrigatório e deve conter um formato de texto válido."},
	"OBE024": {"INVALID_CITY", "O campo city é obrigatório e deve conter um formato de texto válido."},
	"OBE025": {"INVALID_STATE", "O campo state é obrigatório e deve ser uma estado valido."},
	"OBE026": {"INVALID_BUSINESS_EMAIL", "O campo businessEmail é obrigatório e deve ser um email válido."},
	"OBE027": {"INVALID_BUSINESS_NAME", "O campo businessName é obrigatório e deve conter um formato de texto válido."},
	"OBE028": {"INVALID_TRADING_NAME", "O campo tradingName é obrigatório deve respeitar o limite de caracteres e conter um formato de texto válido."},
	"OBE029": {"INVALID_OWNER_DOCUMENT_NUMBER", "O campo owner.documentNumber é obrigatório e deve ser um CPF ou CNPJ válido."},
	"OBE030": {"OWNER_NAME_REQUIRED", "O campo owner.name é obrigatório e deve ser completo."},
	"OBE031": {"INVALID_OWNER_EMAIL", "O campo owner.email é obrigatório e deve ser um email válido."},
	"OBE032": {"OWNER_ADDRESS_REQUIRED", "O campo owner.address é obrigatório."},
	"OBE033": {"UNDERAGE_REGISTRATION_NOT_ALLOWED", "Cadastro não permitido para menores de idade."},
	"OBE034": {"INVALID_JSON_FORMAT", "Formato do JSON esta fora do padrão. Verifique a documentação."},
	"OBE035": {"OPERATION_FAILED", "Não foi possivel realizar essa operação. Tente novamente mais tarde."},
	"OBE036": {"INVALID_COMPANY_TYPE", "CompanyType inválido."},
	"OBE037": {"INVALID_OWNERS_ARRAY", "O campo Owners deve conter um array de no mínimo 1 e máximo 10."},
	"OBE038": {"DUPLICATE_OWNERS", "Owners não podem ser duplicados."},
	"OBE039": {"BUSINESS_ADDRESS_REQUIRED", "O campo businessAddress é obrigatório."},
	"OBE040": {"INVALID_OWNER_TYPE", "O campo ownerType é obrigatório e deve conter um valor válido."},
	"OBE041": {"BACKGROUND_CHECK_NOT_FOUND", "BackgroundCheck não encontrado ou com status diferente de pendente."},
	"OBE042": {"BACKGROUND_CHECK_UPDATE_FAILED", "Erro ao atualizar backgroundCheck."},
	"OBE043": {"DOCUMENTSCOPY_NOT_FOUND", "Documentscopy não encontrado ou com status diferente de pendente."},
	"OBE044": {"DOCUMENTSCOPY_UPDATE_FAILED", "Erro ao atualizar documentscopy."},
	"OBE045": {"INVALID_PROPOSAL_STATUS", "Status da proposta inexistente verifique a documentação por favor."},
	"OBE046": {"INVALID_DATE", "Data inválida."},
	"OBE047": {"INVALID_LIMIT", "Limite inserido inválido. Os campos limit ou limitPerPage devem ter valores entre 1 e 200."},
	"OBE048": {"INVALID_DOCUMENT_NUMBER", "O campo documentNumber deve ser um CPF ou CNPJ válido."},
	"OBE049": {"PROPOSAL_NOT_FOUND", "Não foi encontrada nenhuma proposta referente aos dados informados."},
	"OBE050": {"INVALID_DATE_RANGE", "A data inicial não pode ser maior que a data final."},
	"OBE051": {"MISSING_REQUIRED_FIELDS", "Ao não enviar o proposalId os campos data inicial e a data final são obrigatórios."},
	"OBE052": {"DATE_RANGE_TOO_LARGE", "O intervalo de dias entre a data inicial e a data final não deve ser maior que {0} dias."},
	"OBE053": {"MISSING_OWNER_TYPE", "O campo ownerType deve conter pelo menos um sócio ou representante"},
	"OBE054": {"MISSING_PROPOSAL_ID_OR_CLIENT_CODE", "ProposalId e clientCode não enviados. Ao menos um desses parametros deve ser enviado."},
	"OBE055": {"FILES_NOT_FOUND", "Não foram encontrados arquivos para o proposalId ou clientCode informado(s)."},
	"OBE056": {"DOCUMENTSCOPIES_NOT_FOUND", "Não foram encontradas documentoscopias referentes ao proposalId ou clientCode enviado."},
	"OBE057": {"DOCUMENT_FETCH_FAILED", "Ocorreu um erro ao buscar documentos."},
	"OBE058": {"INVALID_CLIENT_TYPE", "ClientType inválido."},
	"OBE059": {"INVALID_SOURCE_TYPE", "SourceType inválido."},
	"OBE060": {"CLIENT_ID_REQUIRED", "O campo clientId é obrigatório."},
	"OBE061": {"INVALID_SOURCE", "Source inválido."},
	"OBE062": {"DUPLICATE_CLIENT_CODE", "ClientCode já vinculado a outra proposta, esse campo deve ser único por proposta."},
	"OBE063": {"NO_RECORDS_FOUND", "Não foram encontrados registros para a sua requisição."},
	"OBE064": {"DUPLICATE_PROPOSAL", "Já existe uma proposta em aberto para esse documentNumber."},
	"OBE065": {"DATE_FROM_REQUIRED", "O campo dateFrom é obrigatório."},
	"OBE066": {"DATE_TO_REQUIRED", "O campo dateTo é obrigatório."},
	"OBE067": {"INVALID_PARTNER_NAME", "O campo partner.partnerName deve conter um valor válido."},
	"OBE068": {"PARTNER_PARAMETER_REQUIRED", "O campo partner.parameter deve ser preenchido."},
	"OBE069": {"PARTNER_NAME_REQUIRED", "Ao enviar os campos partner.parameters, o campo partner.partnerName deve ser obrigatório."},
	"OBE070": {"USER_NOT_STARTED_WEBVIEW", "Não foram encontrados dados, pois o usuário ainda não iniciou a jornada webview. Tente novamente mais tarde."},
	"OBE071": {"PARTNER_QUERY_FAILED", "Ocorreu um erro ao consultar parceiro. Favor tentar novamente mais tarde."},
	"OIE999": {"INTERNAL_API_ERROR", "Ocorreu um erro interno durante a chamada da api"},
}

// FindOnboardingError ... retorna a mensagem de erro correspondente ao código de erro de Onboarding
func FindOnboardingError(code string, responseStatus *int) *grok.Error {
	if mapping, exists := OnboardingErrorMappings[code]; exists {
		return grok.NewError(*responseStatus, mapping.ContbankCode, mapping.Description)
	}
	return grok.NewError(http.StatusInternalServerError, "UNKNOWN_ERROR", "unknown error")
}
