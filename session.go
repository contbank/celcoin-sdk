package celcoin

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/contbank/grok"
	"github.com/patrickmn/go-cache"
)

const (
	ApiEndpoint       string = "https://sandbox.openfinance.celcoin.dev"
	LoginEndpoint     string = "https://sandbox.openfinance.celcoin.dev"
	CelcoinEnvSandbox string = "SANDBOX"
	CelcoinEnvProd    string = "PRODUCTION"
)

// Config ...
type Config struct {
	LoginEndpoint *string
	APIEndpoint   *string
	ClientID      *string
	ClientSecret  *string
	APIVersion    *string
	Scopes        *string
	Cache         *cache.Cache
	Mtls          *bool
	CompanyKey    *string
	Certificate   *Certificate
	Environment   *string
}

// Session ...
type Session struct {
	LoginEndpoint string
	APIEndpoint   string
	ClientID      string
	ClientSecret  string
	APIVersion    string
	Cache         cache.Cache
	Scopes        string
	Mtls          bool
	Environment   string
}

// NewSession ...
func NewSession(config Config) (*Session, error) {
	if config.APIEndpoint == nil {
		config.APIEndpoint = String(ApiEndpoint)
	}

	if config.LoginEndpoint == nil {
		config.LoginEndpoint = String(LoginEndpoint)
	}

	if config.APIVersion == nil {
		config.APIVersion = aws.String("1.0")
	}

	if config.ClientID == nil {
		config.ClientID = String(os.Getenv("CELCOIN_CLIENT_ID"))
	}

	if config.ClientSecret == nil {
		config.ClientID = String(os.Getenv("CELCOIN_CLIENT_SECRET"))
	}

	if config.Cache == nil {
		config.Cache = cache.New(10*time.Minute, 1*time.Second)
	}

	if config.Scopes == nil {
		config.Scopes = String("")
	}

	if config.Mtls != nil {
		if *config.Mtls {
			config.ClientID = &config.Certificate.ClientID
		}
	}

	if config.Environment == nil {
		config.Environment = String(CelcoinEnvSandbox)
	}

	var session = &Session{
		LoginEndpoint: *config.LoginEndpoint,
		APIEndpoint:   *config.APIEndpoint,
		ClientID:      *config.ClientID,
		ClientSecret:  *config.ClientSecret,
		APIVersion:    *config.APIVersion,
		Cache:         *config.Cache,
		Scopes:        *config.Scopes,
		Mtls:          *config.Mtls,
		Environment:   *config.Environment,
	}

	return session, nil
}

// CreateMtlsHTTPClient ...
func CreateMtlsHTTPClient(cert *Certificate) *http.Client {
	httpClient := &http.Client{}
	httpClient.Timeout = 30 * time.Second

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(cert.CertificateChain))

	certificate, err := grok.LoadCertificate([]byte(cert.Certificate), []byte(cert.PrivateKey), cert.Passphrase)
	if err != nil {
		panic(err)
	}

	httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            caCertPool,
			Certificates:       []tls.Certificate{*certificate},
			InsecureSkipVerify: true,
		},
	}
	return httpClient
}

// CreateOAuth2HTTPClient ... cria um cliente HTTP autenticado via OAuth2
func CreateOAuth2HTTPClient(session *Session) *http.Client {
	httpClient := &http.Client{}
	httpClient.Timeout = 30 * time.Second

	// Obter token de acesso
	token, err := fetchAccessToken(httpClient, session)
	if err != nil {
		panic(fmt.Sprintf("Erro ao obter token de acesso: %v", err))
	}

	// Configurar transporte com cabeçalho de autenticação
	httpClient.Transport = &oauthTransport{
		underlyingTransport: http.DefaultTransport,
		accessToken:         token,
	}

	return httpClient
}

// fetchAccessToken ... realiza a requisição para obter o token de acesso
func fetchAccessToken(client *http.Client, session *Session) (string, error) {
	data := []byte(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials",
		session.ClientID, session.ClientSecret))

	req, err := http.NewRequest("POST", session.LoginEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("erro ao criar a requisição de token: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer a requisição de token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("falha ao obter token: %s", body)
	}

	var tokenResponse AuthenticationResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("erro ao decodificar a resposta do token: %w", err)
	}

	return tokenResponse.AccessToken, nil
}

// oauthTransport ... é um transporte customizado que adiciona o token de autenticação
type oauthTransport struct {
	underlyingTransport http.RoundTripper
	accessToken         string
}

// RoundTrip ... adiciona o cabeçalho Authorization com o token OAuth2
func (t *oauthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.accessToken))
	return t.underlyingTransport.RoundTrip(req)
}
