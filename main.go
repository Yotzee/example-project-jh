package main

import (
	"net/http"

	"example-project/api"
	"example-project/logger"
)

func main() {
	http.HandleFunc("/weather", api.WeatherHandler)

	port := ":8000"
	logger.Infof("Weather service starting on port %s", port)
	logger.Infof("Endpoint available at http://localhost%s/weather", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Errorf("Failed to start server: %v", err)
	}
}
