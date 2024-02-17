package images

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ijimiji/pipeline/internal/models"
	"github.com/ijimiji/pipeline/internal/services/core"
	"github.com/ijimiji/pipeline/internal/slices"
	"github.com/ijimiji/pipeline/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

type statusResponse struct {
	ImageGroup models.ImageGroup `json:"imageGroup"`
}

func newStatusHandler(coreClient core.Client) *StatusHandler {
	return &StatusHandler{
		coreClient: coreClient,
	}
}

type StatusHandler struct {
	coreClient core.Client
}

// @Summary Get image generation status
// @Description Get image generation status
// @Tags Image generation flow
// @Accept  json
// @Produce  json
// @Success 200 {object} statusResponse
// @Router /images/{id} [get]
// @Param id path string true "Image group ID"
func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opsProcessed.Inc()

	encoder := json.NewEncoder(w)
	id := chi.URLParam(r, "id")

	out, err := h.coreClient.Status(r.Context(), &proto.StatusRequest{
		ID: id,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder.Encode(statusResponse{
		ImageGroup: models.ImageGroup{
			ID: out.GetImageGroup().GetID(),
			Images: slices.Map(out.GetImageGroup().GetImages(), func(dto *proto.Image) models.Image {
				return models.Image{
					ID:               dto.GetID(),
					URL:              dto.GetURL(),
					Prompt:           dto.GetPrompt(),
					GenerationStatus: dto.GetStatus(),
				}
			}),
		},
	})
}
