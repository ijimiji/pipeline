package models

type Image struct {
	ID               string `json:"id"`
	URL              string `json:"url"`
	Prompt           string `json:"prompt"`
	GenerationStatus string `json:"generationStatus"`
}

const (
	GenerationStatusEnqueued = "enqueued"
	GenerationStatusReady    = "ready"
)
