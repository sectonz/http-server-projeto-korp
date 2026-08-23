package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const nomeServico = "http-server-projeto-korp"

type ProjetoResponse struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Servico   string `json:"servico"`
	Timestamp string `json:"timestamp"`
}

var (
	totalRequisicoes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requisicoes_total",
			Help: "Numero total de requisicoes HTTP processadas",
		},
		[]string{"rota", "metodo", "status"},
	)

	duracaoRequisicao = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duracao_segundos",
			Help:    "Duracao das requisicoes HTTP em segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"rota", "metodo", "status"},
	)

	servicoDisponivel = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "servico_disponivel",
			Help: "Disponibilidade do servico (1 = UP, 0 = DOWN)",
		},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		statusStr := strconv.Itoa(rw.statusCode)
		route := r.URL.Path

		totalRequisicoes.WithLabelValues(route, r.Method, statusStr).Inc()
		duracaoRequisicao.WithLabelValues(route, r.Method, statusStr).Observe(duration)
	})
}

func main() {
	servicoDisponivel.Set(1)

	mux := http.NewServeMux()

	mux.HandleFunc("/projeto-korp", projetoHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"servico":   nomeServico,
			"mensagem":  "Servidor em execucao. Acesse /projeto-korp, /health ou /metrics",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	handlerComMetricas := metricsMiddleware(mux)

	log.Printf("Iniciando %s na porta 8080...", nomeServico)
	if err := http.ListenAndServe(":8080", handlerComMetricas); err != nil {
		servicoDisponivel.Set(0)
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}

func projetoHandler(w http.ResponseWriter, r *http.Request) {
	resp := ProjetoResponse{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format("15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Falha ao codificar resposta", http.StatusInternalServerError)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:    "UP",
		Servico:   nomeServico,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
