package renderer

import (
	"os"
)

type Config struct {
	ForceFallbackAdapter *bool
}

func fillDefault(config Config) Config {
	if config.ForceFallbackAdapter == nil {
		enabled := os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1"
		config.ForceFallbackAdapter = &enabled
	}
	return config
}
