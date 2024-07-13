package utils

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"time"

	pb "mini/proto"
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

func SetupFieldMandatory(field *string, envName string, callback func()) {
	setupField(true, field, envName, "", callback)
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

func PrettyPrint(title string, matrix [][]int32) {
	maxWidth := 0
	for _, row := range matrix {
		for _, value := range row {
			width := len(fmt.Sprintf("%d", value))
			if width > maxWidth {
				maxWidth = width
			}
		}
	}

	// Print the matrix with pretty formatting
	fmt.Println(title)
	for _, row := range matrix {
		for _, value := range row {
			fmt.Printf("%*d ", maxWidth, value)
		}
		fmt.Println()
	}
}

func MatrixToProto(matrix [][]int32) *pb.Matrix {
	protoMatrix := &pb.Matrix{}
	for _, row := range matrix {
		protoRow := &pb.Row{}
		for _, value := range row {
			protoRow.Values = append(protoRow.Values, value)
		}
		protoMatrix.Rows = append(protoMatrix.Rows, protoRow)
	}
	return protoMatrix
}

func ProtoToMatrix(protoMatrix *pb.Matrix) [][]int32 {
	matrix := [][]int32{}
	for _, protoRow := range protoMatrix.Rows {
		row := []int32{}
		for _, value := range protoRow.Values {
			row = append(row, value)
		}
		matrix = append(matrix, row)
	}
	return matrix
}

func GenerateEmptyMatrix(height int, width int) [][]int32 {
	result := make([][]int32, height)
	for i := 0; i < height; i++ {
		result[i] = make([]int32, width)
	}
	return result
}

func GenerateRandomMatrix(targetSize int) [][]int32 {
	result := make([][]int32, targetSize)
	for i := 0; i < targetSize; i++ {
		result[i] = make([]int32, targetSize)
		for j := 0; j < targetSize; j++ {
			result[i][j] = int32(rand.Intn(256))
		}
	}
	return result
}
