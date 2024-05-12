package calc

import "math"

func Deg(deg float32) float32 {
	return float32(float64(deg) * math.Pi / 180.0)
}
