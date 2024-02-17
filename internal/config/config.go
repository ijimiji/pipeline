package config

import (
	"encoding/json"
	"os"
)

func New[T any](file string) T {
	var ret T
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bytes, &ret); err != nil {
		panic(err)
	}

	return ret
}
