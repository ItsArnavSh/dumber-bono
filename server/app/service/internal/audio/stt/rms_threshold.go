package stt

import (
	"math"
)

var Threshold = 7000

func BackgroundTuner() {

}

func RMS(samples []int16) float64 {
	var sum float64

	for _, s := range samples {
		x := float64(s)
		sum += x * x
	}

	return math.Sqrt(sum / float64(len(samples)))
}
