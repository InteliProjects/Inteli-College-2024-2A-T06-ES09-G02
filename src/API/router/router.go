package router

import (
	"2024-2A-T06-ES09-G02/src/API/handlers"
	"2024-2A-T06-ES09-G02/src/API/logging"
	"github.com/gorilla/mux"
	"net/http"
)

func InitializeRoutes() *mux.Router {
	r := mux.NewRouter()

	// Logger de produção
	logger := logging.GetProdLogger()

	// Passando o logger de produção para as rotas
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, logger)
	}).Methods("POST")

	r.HandleFunc("/normal", handlers.NormalHandler).Methods("GET")
	r.HandleFunc("/random-delay", handlers.RandomDelayHandler).Methods("GET")
	r.HandleFunc("/random-failure", handlers.RandomFailureHandler).Methods("GET")

	return r
}
