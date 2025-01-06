package celcoin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// NewMockServer cria um servidor HTTP com rotas configuradas.
func NewMockServer() *httptest.Server {
	handler := http.NewServeMux()

	// Configurar rotas (pode ser modularizado por funcionalidade)
	RegisterWebhookRoutes(handler)
	RegisterAuthRoutes(handler) // Adicionando mock de login

	// Retorna o servidor configurado
	return httptest.NewServer(handler)
}

// RegisterWebhookRoutes registra as rotas relacionadas a Webhooks no handler.
func RegisterWebhookRoutes(handler *http.ServeMux) {
	handler.HandleFunc("/baas-webhookmanager/v1/webhook/subscription", handleCreateSubscription)
}

// RegisterAuthRoutes registra as rotas relacionadas a autenticação no handler.
func RegisterAuthRoutes(handler *http.ServeMux) {
	handler.HandleFunc("/auth/login", handleLogin) // Mockando a rota de login
}

// handleLogin simula o endpoint de login da API.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Validar método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validar headers
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	// Simulando resposta de login bem-sucedido com token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": "Bearer fake-api-key", // Mock de um token
	})
}

// handleCreateSubscription simula o endpoint de criação de Webhook Subscription.
func handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	// Validar método HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validar headers
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	// Validar payload
	var req WebhookSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Simular resposta de sucesso
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(WebhookSubscriptionResponse{
		Version: "1.0.0",
		Status:  "SUCCESS",
	})
}
