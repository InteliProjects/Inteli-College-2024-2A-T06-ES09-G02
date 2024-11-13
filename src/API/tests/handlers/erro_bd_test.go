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

func TestRegisterHandler_DatabaseConnectionError(t *testing.T) {
	// Inicia o logger para testes
	err := logging.InitTestLogger()
	if err != nil {
		t.Fatalf("Erro ao inicializar logger de teste: %v", err)
	}
	logger := logging.GetTestLogger() // Obter o logger de testes

	logging.LogInfoTest("Iniciando teste: TestRegisterHandler_DatabaseConnectionError")

	// Modifica o caminho do arquivo de banco de dados para um inválido
	originalDatabaseFile := handlers.UserDatabaseFile
	handlers.UserDatabaseFile = "/invalid/path/users.json"

	// Cria um usuário de teste
	user := models.User{Username: "Test User", Email: "testuser@example.com"}
	body, _ := json.Marshal(user)

	// Cria uma requisição HTTP simulada
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	// Envolve o handler com o logger de teste
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, logger) // Passando o logger de teste
	})

	// Mede o tempo de resposta
	startTime := time.Now()
	handler.ServeHTTP(rr, req)
	responseTime := time.Since(startTime)

	// Log do erro de conexão com o banco de dados
	logging.LogInfoTest("Erro simulado de conexão com o banco de dados")
	logging.LogRequestWithDurationTest(req.Method, req.RequestURI, rr.Code, responseTime, req.RemoteAddr, req.UserAgent())

	// Verifica se o código de resposta está correto
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Failed to save user.\n", rr.Body.String())

	// Reseta o caminho do arquivo de banco de dados
	handlers.UserDatabaseFile = originalDatabaseFile

	// Finaliza o log do teste
	logging.LogInfoTest("Teste finalizado: TestRegisterHandler_DatabaseConnectionError")
}
