package handlers_test

import (
	"2024-2A-T06-ES09-G02/src/API/handlers"
	"2024-2A-T06-ES09-G02/src/API/logging"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHandler_InvalidInput(t *testing.T) {
	// Inicia o logger para testes
	err := logging.InitTestLogger()
	if err != nil {
		t.Fatalf("Erro ao inicializar logger de teste: %v", err)
	}
	logger := logging.GetTestLogger() // Obter o logger de testes

	logging.LogInfoTest("Iniciando teste: TestRegisterHandler_InvalidInput")

	startTime := time.Now()

	// JSON inválido para a requisição
	invalidJSON := []byte(`{"Username": "Test User", "Email": "invalid email}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(invalidJSON))
	rec := httptest.NewRecorder()

	// Envolve o handler com o logger de teste
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, logger) // Passando o logger de teste
	})

	handler.ServeHTTP(rec, req)
	responseTime := time.Since(startTime)

	// Loga o tempo de resposta
	logging.LogRequestWithDurationTest(req.Method, req.RequestURI, rec.Code, responseTime, req.RemoteAddr, req.UserAgent())

	// Verifica se o código de resposta é 400 BadRequest
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Finalizando o log do teste
	logging.LogInfoTest("Teste finalizado: TestRegisterHandler_InvalidInput")
}
