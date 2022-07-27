package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"

	"github.com/Scalingo/sclng-backend-test-v1/internal/service/github"
)

type Server struct {
	githubService *github.Service
}

func NewServer(GithubAPIKey string) (*Server, error) {
	githubService, err := github.NewService(GithubAPIKey)
	if err != nil {
		return nil, err
	}
	return &Server{
		githubService,
	}, nil
}

func (s *Server) RegisterHandler(router *handlers.Router) {
	router.HandleFunc("/repos", s.listRepositories)
	router.HandleFunc("/stat", s.getStatistics)
}

func handleError(ctx context.Context, w http.ResponseWriter, httpCode int, message string, err error) error {
	log := logger.Get(ctx)
	log.WithError(err).Error(message)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"code": httpCode, "message": message})
	w.WriteHeader(httpCode)
	return err
}

func (s *Server) listRepositories(w http.ResponseWriter, r *http.Request, req map[string]string) error {
	ctx := r.Context()
	w.Header().Add("Content-Type", "application/json")

	if r.Method != "GET" {
		return handleError(ctx, w, http.StatusMethodNotAllowed, "Method not Allowed", errors.New("invalid HTTP method"))
	}

	filters, err := github.NewListingFilters(r.URL.Query())
	if err != nil {
		return handleError(ctx, w, http.StatusBadRequest, "Invalid filters", err)
	}

	repos, err := s.githubService.ListPublicRepositories(ctx, filters)
	if err != nil {
		return handleError(ctx, w, http.StatusInternalServerError, "Failed to list repositories", err)
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{"repositories": repos})
	if err != nil {
		return handleError(ctx, w, http.StatusInternalServerError, "Fail to encode JSON", err)
	}
	return nil
}

func (s *Server) getStatistics(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	ctx := r.Context()
	w.Header().Add("Content-Type", "application/json")

	if r.Method != "GET" {
		return handleError(ctx, w, http.StatusMethodNotAllowed, "Method not Allowed", errors.New("invalid HTTP method"))
	}

	filters, err := github.NewListingFilters(r.URL.Query())
	if err != nil {
		return handleError(ctx, w, http.StatusBadRequest, "Invalid filters", err)
	}

	stats, err := s.githubService.GetPublicRepositoriesStats(ctx, filters)
	if err != nil {
		return handleError(ctx, w, http.StatusInternalServerError, "Failed to stats repositories", err)
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{"statistics": stats})
	if err != nil {
		return handleError(ctx, w, http.StatusInternalServerError, "Fail to encode JSON", err)
	}
	return nil
}