package tools

import "fmt"

// Mock function to simulate weather data retrieval
func GetWeather(location string) string {
	return fmt.Sprintf("In %s, it is Sunny, 25°C", location)
}
