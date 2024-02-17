package main

import (
	"github.com/ijimiji/pipeline/internal/api/http"
	"github.com/ijimiji/pipeline/internal/services/core"
)

type Config struct {
	API  http.Config
	Core core.Config
}
