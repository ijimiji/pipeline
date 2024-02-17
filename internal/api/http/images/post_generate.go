package images

import (
	"encoding/json"
	"net/http"

	"github.com/ijimiji/pipeline/internal/services/core"
	"github.com/ijimiji/pipeline/proto"
)

type generateRequest struct {
	// @Description What should be drawn
	Prompt string `example:"very fat cat" json:"prompt"`
	// @Description Amount of images to be drawn for your request
	ImagesCount int `example:"4" json:"imagesCount"`
}

type generateResponse struct {
	// @Description ID of image group
	ID string `json:"id" example:"1234"`
}

func newGenerateHandler(coreClient core.Client) *generateHandler {
	return &generateHandler{
		coreClient: coreClient,
	}
}

type generateHandler struct {
	coreClient core.Client
}

// @Summary Generate new image
// @Description Generate new image by given prompt.
// @Tags Image generation flow
// @Accept  json
// @Produce  json
// @Success 200 {object} generateResponse
// @Router /images/generate [post]
// @Param body body generateRequest true "Request body"
func (h *generateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	out, err := h.coreClient.Generate(r.Context(), &proto.GenerateRequest{
		Prompt: req.Prompt,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder.Encode(generateResponse{
		ID: out.ID,
	})
}
