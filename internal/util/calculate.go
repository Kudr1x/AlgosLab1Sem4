package util

import "math"

func CalculateCompressionRatio(original, compressed []byte) float64 {
	return math.Round(float64(len(compressed))/float64(len(original))*1000) / 1000
}
