package baseutil

import "github.com/GaryBoone/GoStats/stats"

func StandardDeviationInt(inList []int) float64 {
	floatList := make([]float64, len(inList))
	for i, value := range inList {
		floatList[i] = float64(value)
	}
	stdDev := stats.StatsSampleStandardDeviation(floatList)
	return stdDev
}
