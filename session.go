package celcoin

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"

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

// oauthTransport ... é um transporte customizado que adiciona o token e o renova quando necessário
type oauthTransport struct {
	underlyingTransport http.RoundTripper
	session             *Session
	token               string
	tokenExpiration     time.Time
	mutex               *sync.Mutex
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

	/*if config.Mtls != nil {
		if *config.Mtls {
			config.ClientID = &config.Certificate.ClientID
		}
	}*/

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
func CreateMtlsHTTPClient(cert *Certificate, session *Session) *http.Client {
	// Carregar o certificado e a chave privada
	tlsCert, err := tls.X509KeyPair([]byte(cert.Certificate), []byte(cert.PrivateKey))
	if err != nil {
		panic("Erro ao carregar o certificado e a chave privada")
	}

	// Criar um pool de certificados confiáveis (opcional, para validar o servidor)
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM([]byte(cert.Certificate)) {
		panic("Erro ao adicionar o certificado ao pool de CA")
	}

	mtlsTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            caCertPool,
			Certificates:       []tls.Certificate{tlsCert},
			InsecureSkipVerify: true,
		},
	}

	// Obter o token inicial usando o transporte mTLS
	token, expiration, err := fetchAccessToken(&http.Client{Transport: mtlsTransport}, session)
	if err != nil {
		panic(fmt.Sprintf("Erro ao obter token de acesso: %v", err))
	}

	// Configura o transporte com renovação automática de token
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &oauthTransport{
			underlyingTransport: mtlsTransport,
			session:             session,
			token:               token,
			tokenExpiration:     expiration,
			mutex:               &sync.Mutex{},
		},
	}
	return httpClient
}

// CreateOAuth2HTTPClient ... cria um cliente HTTP autenticado via OAuth2 com renovação automática de token
func CreateOAuth2HTTPClient(session *Session) *http.Client {
	httpClient := &http.Client{}
	httpClient.Timeout = 30 * time.Second

	// Obter token inicial
	token, expiration, err := fetchAccessToken(httpClient, session)
	if err != nil {
		panic(fmt.Sprintf("Erro ao obter token de acesso: %v", err))
	}

	// Configurar transporte com renovação automática do token
	httpClient.Transport = &oauthTransport{
		underlyingTransport: http.DefaultTransport,
		session:             session,
		token:               token,
		tokenExpiration:     expiration,
		mutex:               &sync.Mutex{},
	}

	return httpClient
}

func fetchAccessToken(client *http.Client, session *Session) (string, time.Time, error) {
	var data []byte

	if session.Mtls {
		mtlsData := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials",
			session.ClientID, session.ClientSecret)
		data = []byte(mtlsData)
	} else {
		oauth2Data := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials",
			session.ClientID, session.ClientSecret)
		data = []byte(oauth2Data)
	}

	url, err := url.Parse(session.LoginEndpoint)
	if err != nil {
		return "", time.Time{}, err
	}

	url.Path = path.Join(url.Path, LoginPath)
	endpoint := url.String()

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("erro ao criar a requisição de token: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("erro ao fazer a requisição de token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", time.Time{}, fmt.Errorf("falha ao obter token: %s", body)
	}

	var tokenResponse AuthenticationResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", time.Time{}, fmt.Errorf("erro ao decodificar a resposta do token: %w", err)
	}

	// Calcula a hora de expiração com base no tempo atual e no tempo de expiração do token
	expiration := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	return tokenResponse.AccessToken, expiration, nil
}

// RoundTrip ... adiciona o cabeçalho Authorization e renova o token se necessário
func (t *oauthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if time.Now().After(t.tokenExpiration.Add(-1 * time.Minute)) { // Renova o token antes de expirar
		client := &http.Client{Transport: t.underlyingTransport}
		newToken, newExpiration, err := fetchAccessToken(client, t.session)
		if err != nil {
			return nil, fmt.Errorf("falha ao renovar o token: %w", err)
		}
		t.token = newToken
		t.tokenExpiration = newExpiration
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.underlyingTransport.RoundTrip(req)
}
