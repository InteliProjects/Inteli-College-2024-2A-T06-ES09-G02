package main

import (
	"2024-2A-T06-ES09-G02/src/API/logging"
	"2024-2A-T06-ES09-G02/src/API/router"
	"log"
	"net/http"
)

func main() {
	log.Println("API rodando...")

	// Inicializa o logger de produção
	err := logging.InitProdLogger()
	if err != nil {
		log.Fatalf("Erro ao inicializar logger de produção: %v", err)
	}

	// Inicia as rotas
	r := router.InitializeRoutes()

	// Executa o servidor
	log.Println("API está rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
