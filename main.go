package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"

	"github.com/Scalingo/sclng-backend-test-v1/internal/server"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := NewConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(-1)
	}
	srv, err := server.NewServer(cfg.GithubApiKey)
	if err != nil {
		log.WithError(err).Error("Fail to initialize the server")
		os.Exit(-1)
	}

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", PongHandler)
	// Initialize web server and configure the following routes:
	// GET /repos
	// GET /stats
	srv.RegisterHandler(router)

	log.WithField("port", cfg.Port).Info("Listening...")
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
}

func PongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}
