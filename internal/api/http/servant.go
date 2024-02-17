package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ijimiji/pipeline/internal/api/http/images"
	"github.com/ijimiji/pipeline/internal/services/core"
)

func NewServant(config Config, coreClient core.Client) *servant {
	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer)

	router.Mount("/images", images.New(coreClient))

	return &servant{
		config: config,
		router: router,
	}
}

type servant struct {
	router chi.Router
	config Config
}

func (s *servant) ListenAndServe() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.router)
}
