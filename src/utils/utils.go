package utils

import (
	"log"
	"math"
	"net"
	"os"
	"time"
)

const Step = 25

func SimulatedCPUIntensiveFunction(baseDurationMillis float64, active *int, multiplier int) {

	var counter float64

	for counter < baseDurationMillis {
		counter += Step / (float64(*active) * float64(multiplier))
		time.Sleep(Step * time.Millisecond)
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

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func SetupFieldOptional(field *string, envName string, defaultValue string) {
	setupField(false, field, envName, defaultValue, nil)
}

func SetupFieldMandatory(field *string, envName string, defaultValue string, callback func()) {
	setupField(true, field, envName, defaultValue, callback)
}

func setupField(mandatory bool, field *string, envName string, defaultValue string, callback func()) {
	if *field == "" {
		*field = os.Getenv(envName)
		if *field == "" {
			if mandatory {
				callback()
				return
			} else {
				*field = defaultValue
			}
		}
	}
}
