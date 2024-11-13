package logging

import (
	"go.uber.org/zap"
	"log"
	"time"
)

// Loggers separados para produção e teste
var prodLogger *zap.Logger
var testLogger *zap.Logger

// InitProdLogger inicializa o logger para produção
func InitProdLogger() error {
	if prodLogger != nil {
		return nil
	}

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"api.json", // Log de produção
		"stdout",
	}
	cfg.ErrorOutputPaths = []string{
		"stderr",
	}

	var err error
	prodLogger, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

// Função para obter o logger de produção
func GetProdLogger() *zap.Logger {
	if prodLogger == nil {
		log.Println("ProdLogger não inicializado. Chamando InitProdLogger().")
		if err := InitProdLogger(); err != nil {
			log.Fatalf("Erro ao inicializar ProdLogger: %v", err)
		}
	}
	return prodLogger
}

// InitTestLogger inicializa o logger para os testes
func InitTestLogger() error {
	if testLogger != nil {
		return nil
	}

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"../test.json", // Log de teste
		"stdout",
	}
	cfg.ErrorOutputPaths = []string{
		"stderr",
	}

	var err error
	testLogger, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

// Função para obter o logger de testes
func GetTestLogger() *zap.Logger {
	return testLogger
}

// LogInfoTest logs general information (para testes)
func LogInfoTest(message string) {
	if testLogger == nil {
		log.Println("TestLogger não inicializado. Certifique-se de chamar InitTestLogger() antes de usar o logger.")
		return
	}
	testLogger.Info(message)
}

// LogErrorTest logs errors for tests
func LogErrorTest(method, uri string, status int, message string) {
	if testLogger == nil {
		log.Println("TestLogger não inicializado. Certifique-se de chamar InitTestLogger() antes de usar o logger.")
		return
	}
	testLogger.Error("Request Error (Test)",
		zap.String("method", method),
		zap.String("uri", uri),
		zap.Int("status", status),
		zap.String("error", message),
	)
}

// LogRequestWithDurationTest logs a request during test execution (para testes)
func LogRequestWithDurationTest(method, uri string, status int, duration time.Duration, clientIP, userAgent string) {
	if testLogger == nil {
		log.Println("TestLogger não inicializado. Certifique-se de chamar InitTestLogger() antes de usar o logger.")
		return
	}
	testLogger.Info("Request Info (Test)",
		zap.String("method", method),
		zap.String("uri", uri),
		zap.Int("status", status),
		zap.Duration("duration", duration),
		zap.String("client_ip", clientIP),
		zap.String("user_agent", userAgent),
	)
}
