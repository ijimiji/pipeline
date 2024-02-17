sidecar:
	go run ./cmd/sidecar/*.go config/sidecar.json

gateway:
	go run ./cmd/gateway/*.go config/gateway.json

core:
	go run ./cmd/core/*.go config/core.json