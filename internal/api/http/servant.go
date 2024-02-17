package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ijimiji/pipeline/internal/api/http/images"
	"github.com/ijimiji/pipeline/internal/services/core"
	swagger "github.com/swaggo/http-swagger/v2"
)

func NewServant(config Config, coreClient core.Client) *servant {
	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer, SwaggerCORS)

	apiRouter := chi.NewRouter()
	apiRouter.Mount("/images", images.New(coreClient))
	router.Get("/swagger/*", swagger.Handler(
		swagger.URL("swagger/doc.json"),
	))

	router.Mount("/api/v1", apiRouter)

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
