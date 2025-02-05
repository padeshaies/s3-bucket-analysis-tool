package helpers

import (
	"math"
)

func CalculateBucketCost(bucketSizeInBytes int) float64 {
	// TODO change the 0.023 depending on the region and bucket tier
	bucketPrice := 0.023 * float64(bucketSizeInBytes) / 1024 / 1024 / 1024
	// Round to 2 decimal places
	return round2(bucketPrice)
}

func round2(n float64) float64 {
	return math.Round(n*100) / 100
}
