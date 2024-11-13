package handlers

import (
	"2024-2A-T06-ES09-G02/src/API/logging"
	"2024-2A-T06-ES09-G02/src/API/models"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"go.uber.org/zap"
)

var UserDatabaseFile = "users.json"

// RegisterHandler - Função de cadastro de usuário com suporte a timeout e logs estruturados com zap
func RegisterHandler(w http.ResponseWriter, r *http.Request, logger *zap.Logger) {
	startTime := time.Now()

	// Verifica se o logger foi passado
	if logger == nil {
		logger = logging.GetProdLogger() // Use o logger de produção como fallback
	}

	// Verifica se o contexto foi cancelado (timeout)
	select {
	case <-r.Context().Done():
		logger.Error("Request timed out due to context cancellation",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", http.StatusGatewayTimeout),
		)
		http.Error(w, "context deadline exceeded", http.StatusGatewayTimeout)
		return
	default:
		// Continua o processamento se o contexto ainda estiver ativo
	}

	var newUser models.User

	// Decodifica o JSON do corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		logger.Error("Invalid JSON input",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", http.StatusBadRequest),
		)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Valida o formato do email
	if !isValidEmail(newUser.Email) {
		logger.Error("Invalid email format",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", http.StatusBadRequest),
		)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Salva o usuário no "banco de dados"
	userID, err := saveUser(newUser)
	if err != nil {
		logger.Error("Failed to save user",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", http.StatusInternalServerError),
		)
		http.Error(w, "Failed to save user.", http.StatusInternalServerError)
		return
	}

	// Resposta de sucesso
	response := map[string]interface{}{
		"message": "User created successfully",
		"user_id": userID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	end_time := time.Now()

	// Loga a requisição com a duração do processamento
	responseTime := time.Since(startTime)
	logger.Info("Request Info",
		zap.String("method", r.Method),
		zap.String("uri", r.RequestURI),
		zap.String("inicio_requisição", startTime.Format(time.RFC3339)),
		zap.String("fim_requisição", end_time.Format(time.RFC3339)),
		zap.Int("status", http.StatusCreated),
		zap.Duration("duration", responseTime),
		zap.String("client_ip", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
	)
}

// Valida o formato do email usando regex
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Função para salvar o usuário em um arquivo JSON
func saveUser(user models.User) (int, error) {

	logger := logging.GetProdLogger()

	var err error

	// Gera um ID único para o usuário
	userID := rand.Intn(1000)
	user.Id = userID

	err = retry(3, 2*time.Second, func() error{
		return SaveToFile(UserDatabaseFile, user)
	})

	if err != nil {

		logger.Error("Falha ao salvar no banco de dados principal", zap.String("database", UserDatabaseFile), zap.Error(err))

		//Falha no banco de dados principal, tenta o banco de dados redundante
		backupDatabaseFile := "backup_user.json"
		err = retry(3, 2*time.Second, func() error {
			return SaveToFile(backupDatabaseFile, user)
		})
		if err != nil {
			// Log de erro ao salvar no banco de dados redundante
			logger.Error("Falha ao salvar no banco de dados redundante", zap.String("database", backupDatabaseFile), zap.Error(err))
			return 0, errors.New("falha ao salvar no banco de dados principal e redundante")
		}
		// Log de sucesso ao salvar no banco de dados redundante
        logger.Info("Usuário salvo no banco de dados redundante", zap.Int("userID", userID), zap.String("database", backupDatabaseFile))
        
    } else {
        // Log de sucesso ao salvar no banco de dados principal
        logger.Info("Usuário salvo no banco de dados principal", zap.Int("userID", userID), zap.String("database", UserDatabaseFile))
    }

	return userID, nil
}

func SaveToFile(fileName string, user models.User) error {
	//Abre o arquivo JSON
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	//Lê os dados existentes no arquivo JSON
	var users []models.User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil && err.Error() != "EOF" {
		return err
	}

	//Adiciona o novo usuário à lista
	users = append(users, user)

	//Volta para o início do arquivo e trunca o arquivo para gravar os dados atualizados 
	file.Seek(0, 0)
	file.Truncate(0)

	//Grava os dados atualizados no arquivo
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(users); err != nil {
		return err
	}

	return nil
}

func retry(attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			//Delay para espera entre falhas
			time.Sleep(sleep)
			continue
		}
		return nil
	}
	return errors.New("todas as tentativas falaharam")
}
