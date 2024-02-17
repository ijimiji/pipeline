package images

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ijimiji/pipeline/internal/services/core"
	"github.com/ijimiji/pipeline/proto"
)

func newDiscardHandler(coreClient core.Client) *discardHandler {
	return &discardHandler{
		coreClient: coreClient,
	}
}

type discardHandler struct {
	coreClient core.Client
}

// @Summary Discard image generation
// @Description Discard image generation
// @Tags Image generation flow
// @Accept  json
// @Produce  json
// @Router /images/{id} [delete]
// @Param id path string true "Image group ID"
func (h *discardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.coreClient.Discard(r.Context(), &proto.DiscardRequest{
		ID: id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
