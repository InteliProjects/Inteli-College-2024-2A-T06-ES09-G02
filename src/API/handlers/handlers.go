package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
	"2024-2A-T06-ES09-G02/src/API/models"
)

func NormalHandler(w http.ResponseWriter, r *http.Request) {
	response := models.Response{ Message: "Normal Response", Status: 200}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RandomDelayHandler(w http.ResponseWriter, r *http.Request) {
	delay := time.Duration(rand.Intn(500)) * time.Millisecond
	time.Sleep(delay)
	response := models.Response{ Message: "Response with random delay", Status: 200}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RandomFailureHandler(w http.ResponseWriter, r *http.Request) {
	if rand.Float32() < 0.3 {
		response := models.Response{ Message: "Random Failure", Status: 500}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	} else {
		response := models.Response {Message: "Success", Status: 200}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}