package celcoin

const (
	// LoginPath ...
	LoginPath string = "v5/token"
	// LoginMtlsPath ...
	LoginMtlsPath string = "v5/token"
	// BalancePath
	BalancePath string = "/baas-walletreports/v1/wallet/balance"
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

// BalanceResponse ...
type BalanceResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Body    struct {
		Amount float64 `json:"amount"`
	} `json:"body"`
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
