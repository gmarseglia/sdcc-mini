package utils

import (
	"math"
	"time"
)

const Step = 25

func SimulatedCPUIntensiveFunction(baseDurationMillis float64, active *int, multiplier int) {

	var counter float64

	for counter < baseDurationMillis {
		time.Sleep(Step * time.Millisecond)
		counter += Step / (float64(*active) * float64(multiplier))
	}

}

func DummyCPUIntensiveFunction(iterations int) float64 {
	result := 0.0

	for i := 0; i < iterations; i++ {
		result += math.Sqrt(float64(i))
		result *= math.Pow(math.Sin(float64(i)), 2)
		result /= math.Pow(math.Cos(float64(i)), 2)
	}

	return result
}
