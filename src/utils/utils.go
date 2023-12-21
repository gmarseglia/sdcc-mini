package utils

import "math"

func DummyCPUIntensiveFunction(iterations int) float64 {
	result := 0.0

	for i := 0; i < iterations; i++ {
		result += math.Sqrt(float64(i))
		result *= math.Pow(math.Sin(float64(i)), 2)
		result /= math.Pow(math.Cos(float64(i)), 2)
	}

	return result
}
