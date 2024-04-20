package renderer

import (
	"os"

	"github.com/LNMMusic/optional"
)

type Config struct {
	ForceFallbackAdapter optional.Option[bool]
}

func fillDefault(config Config) Config {
	if !config.ForceFallbackAdapter.IsSome() {
		config.ForceFallbackAdapter = optional.Some(os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1")
	}
	return config
}
