package image

type GenerateRequest struct {
	ID     string
	Prompt string
}

type GenerateResponse struct {
	ID string
}
