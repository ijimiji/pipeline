package models

type ImageGroup struct {
	ID     string  `json:"id"`
	Images []Image `json:"images"`
}
