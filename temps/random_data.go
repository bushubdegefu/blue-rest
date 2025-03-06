package temps

import (
	"math/rand"
	"text/template"

	"fmt"
	"strconv"
	"time"
)

// Utility function to generate random strings
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// Generates a random email address
func randomEmail() string {
	return randomString(10) + "@example.com"
}

// Generates a random UUID
func randomUUID() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32())
}

// Generates a random ID (could be a large number, or a UUID)
func randomID() string {
	return randomUUID() // You can replace this with a custom ID generation method
}

// Generates a random uint (unsigned integer)
func randomUInt() uint {
	return uint(rand.Intn(1000000))
}

// Generates a random float64
func randomFloat64() float64 {
	return rand.Float64() * 1000.0
}

// Generates a random float32
func randomFloat32() float32 {
	return rand.Float32() * 1000.0
}

// Generates a random int32
func randomInt32() int32 {
	return int32(rand.Intn(1000000))
}

// Generates a random int64
func randomInt64() int64 {
	return int64(rand.Intn(1000000))
}

// Generates a random string
func randomGenericString() string {
	return randomString(8) // Adjust length if needed
}

// Generates a random time
func randomTime() time.Time {
	return time.Now().Add(time.Duration(rand.Intn(1000000)) * time.Second)
}

// Parsing functions
func parseTime(dateStr string) time.Time {
	// You can use the time.RFC3339 format for your specific date string
	result, _ := time.Parse(time.RFC3339, dateStr)
	return result
}

func parseInt(intString string) int {
	// You can use the time.RFC3339 format for your specific date string
	result, _ := strconv.Atoi(intString)
	return result
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

var FuncMap = template.FuncMap{
	"parseTime":     parseTime,           // Register custom function
	"parseInt":      parseInt,            // Register custom function
	"randomEmail":   randomEmail,         // Register random email function
	"randomUUID":    randomUUID,          // Register random UUID function
	"randomID":      randomID,            // Register random ID function
	"randomUInt":    randomUInt,          // Register random uint function
	"randomFloat64": randomFloat64,       // Register random float64 function
	"randomFloat32": randomFloat32,       // Register random float32 function
	"randomInt32":   randomInt32,         // Register random int32 function
	"randomInt64":   randomInt64,         // Register random int64 function
	"randomString":  randomGenericString, // Register random string function
	"randomTime":    randomTime,          // Register random time function
	"randomBool":    randomBool,          // Register random bool function
}
