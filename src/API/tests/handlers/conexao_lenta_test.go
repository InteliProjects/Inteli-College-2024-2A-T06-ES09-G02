package handlers_test

import (
	"2024-2A-T06-ES09-G02/src/API/handlers"
	"2024-2A-T06-ES09-G02/src/API/logging"
	"2024-2A-T06-ES09-G02/src/API/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler_SlowConnection(t *testing.T) {
	// Inicializa o logger para testes
	err := logging.InitTestLogger()
	if err != nil {
		t.Fatalf("Erro ao inicializar logger de teste: %v", err)
	}
	logger := logging.GetTestLogger()

	logging.LogInfoTest("Iniciando teste: TestRegisterHandler_SlowConnection")

	newUser := models.User{
		Username: "Test User",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userJSON, _ := json.Marshal(newUser)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJSON))
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)            // Simula conexão lenta
		handlers.RegisterHandler(w, r, logger) // Passa o logger de teste aqui
	})

	startTime := time.Now()
	handler.ServeHTTP(rec, req)
	responseTime := time.Since(startTime)

	// Log do tempo de resposta
	logging.LogRequestWithDurationTest(req.Method, req.RequestURI, rec.Code, responseTime, req.RemoteAddr, req.UserAgent())

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.GreaterOrEqual(t, responseTime.Seconds(), 3.0, "Expected response time to be at least 3 seconds")

	var response map[string]interface{}
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "User created successfully", response["message"])

	// Finalizando o log do teste
	logging.LogInfoTest("Teste finalizado: TestRegisterHandler_SlowConnection")
}
