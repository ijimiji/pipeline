package images

import (
	"log/slog"
	"net/http"

	"github.com/ijimiji/pipeline/internal/services/core"
	"github.com/ijimiji/pipeline/proto"
)

func newGenerateHandler(coreClient core.Client) *generateHandler {
	return &generateHandler{
		coreClient: coreClient,
	}
}

type generateHandler struct {
	coreClient core.Client
}

func (h *generateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := h.coreClient.Generate(r.Context(), &proto.GenerateRequest{}); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
