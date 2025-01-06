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
	// ErrInvalidAuthenticationCodeOrAccount ...
	ErrInvalidAuthenticationCodeOrAccount = grok.NewError(http.StatusBadRequest, "INVALID_AUTHENTICATION_CODE_OR_ACCOUNT_NUMBER", "invalid authentication code or account number")
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

// FindTransferError ..
func FindTransferError(transferErrorResponse TransferErrorResponse) *grok.Error {
	// get the error code if errors list is null
	if len(transferErrorResponse.Errors) == 0 && transferErrorResponse.Code != "" {
		transferErrorResponse.Errors = []KeyValueErrorModel{
			{
				Key: transferErrorResponse.Code,
			},
		}
	}
	// checking the errors list
	errorModel := transferErrorResponse.Errors[0]
	for _, v := range transferErrorList {
		if v.celcoinTransferError.Key == errorModel.Key {
			return v.grokError
		}
	}
	return grok.NewError(http.StatusBadRequest, errorModel.Key, errorModel.Key+" - "+errorModel.Value)
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
