package handlers_test

import (
	"2024-2A-T06-ES09-G02/src/API/handlers"
	"2024-2A-T06-ES09-G02/src/API/logging"
	"2024-2A-T06-ES09-G02/src/API/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler_Timeout(t *testing.T) {
	// Inicia o logger para testes
	err := logging.InitTestLogger()
	if err != nil {
		t.Fatalf("Erro ao inicializar logger de teste: %v", err)
	}
	logger := logging.GetTestLogger()

	logging.LogInfoTest("Iniciando teste: TestRegisterHandler_Timeout")

	newUser := models.User{
		Username: "Test User",
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userJSON, _ := json.Marshal(newUser)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJSON))

	// Cria um contexto com timeout de 3 segundos
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(4 * time.Second) // Simula resposta lenta
		handlers.RegisterHandler(w, r, logger)
	})

	startTime := time.Now()
	handler.ServeHTTP(rec, req)
	responseTime := time.Since(startTime)

	assert.Equal(t, http.StatusGatewayTimeout, rec.Code)

	// Verifica o conteúdo da resposta, se aplicável
	var response map[string]interface{}
	err = json.NewDecoder(rec.Body).Decode(&response)
	if err == nil {
		assert.Equal(t, "context deadline exceeded", response["message"])
	}

	logging.LogRequestWithDurationTest(req.Method, req.RequestURI, rec.Code, responseTime, req.RemoteAddr, req.UserAgent())

	// Finalizando o log do teste
	logging.LogInfoTest("Teste finalizado: TestRegisterHandler_Timeout")
}
