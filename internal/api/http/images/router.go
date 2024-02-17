package images

import (
	"github.com/go-chi/chi/v5"
	"github.com/ijimiji/pipeline/internal/services/core"
)

func New(coreClient core.Client) chi.Router {
	router := chi.NewRouter()

	router.Post("/generate", newGenerateHandler(coreClient).ServeHTTP)
	router.Get("/{id}", newStatusHandler(coreClient).ServeHTTP)
	router.Delete("/{id}", newDiscardHandler(coreClient).ServeHTTP)

	return router
}
