package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example-project/logger"
	"example-project/weather"
)

// WeatherResponse represents the JSON response for the weather endpoint
type WeatherResponse struct {
	Status string `json:"status"`
	Temp   string `json:"temp,omitempty"`
}

// WeatherHandler handles HTTP requests to the /weather endpoint
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	logger.InfoWithFields("Incoming weather request", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
		"query":  r.URL.RawQuery,
		"ip":     r.RemoteAddr,
	})

	if r.Method != http.MethodGet {
		logger.WarnWithFields("Method not allowed", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
		})
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse latitude and longitude from query parameters
	latStr := r.URL.Query().Get("latitude")
	lonStr := r.URL.Query().Get("longitude")

	if latStr == "" || lonStr == "" {
		logger.WarnWithFields("Missing required query parameters", map[string]interface{}{
			"latitude":  latStr,
			"longitude": lonStr,
		})
		http.Error(w, "latitude and longitude query parameters are required", http.StatusBadRequest)
		return
	}

	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		logger.WarnWithFields("Invalid latitude parameter", map[string]interface{}{
			"latitude": latStr,
			"error":    err.Error(),
		})
		http.Error(w, "invalid latitude parameter", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		logger.WarnWithFields("Invalid longitude parameter", map[string]interface{}{
			"longitude": lonStr,
			"error":     err.Error(),
		})
		http.Error(w, "invalid longitude parameter", http.StatusBadRequest)
		return
	}

	// Validate latitude and longitude ranges
	if latitude < -90 || latitude > 90 {
		logger.WarnWithFields("Latitude out of range", map[string]interface{}{
			"latitude": latitude,
		})
		http.Error(w, "latitude must be between -90 and 90", http.StatusBadRequest)
		return
	}
	if longitude < -180 || longitude > 180 {
		logger.WarnWithFields("Longitude out of range", map[string]interface{}{
			"longitude": longitude,
		})
		http.Error(w, "longitude must be between -180 and 180", http.StatusBadRequest)
		return
	}

	// Get weather data from NOAA
	status, temp, err := weather.GetWeatherFromNOAA(latitude, longitude)
	if err != nil {
		logger.ErrorWithFields("Error fetching weather", map[string]interface{}{
			"error":     err.Error(),
			"latitude":  latitude,
			"longitude": longitude,
		})
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{
		Status: status,
		Temp:   temp,
	}

	logger.InfoWithFields("Weather request successful", map[string]interface{}{
		"latitude":  latitude,
		"longitude": longitude,
		"status":    status,
		"temp":      temp,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=10")
	json.NewEncoder(w).Encode(response)
}
