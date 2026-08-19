package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
)

const nomeServico = "http-server-projeto-korp"

type ProjetoResponse struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

func main() {

	http.HandleFunc("/projeto-korp", projetoHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"teste")
	})

	log.Printf("starting...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func projetoHandler(w http.ResponseWriter, r *http.Request) {

	resp := ProjetoResponse{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format("15:04:05"),
		// Horario: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
