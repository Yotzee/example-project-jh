package weather

import (
	"encoding/json"
	"fmt"
	"strconv"

	"example-project/httpclient"
	"example-project/logger"
)

const (
	userAgent   = "example-project-weather-service/1.0 (contact: support@example.com)"
	noaaBaseURL = "https://api.weather.gov"
)

// pointsResponse represents the response from the NOAA points API
type pointsResponse struct {
	Properties struct {
		GridId              string `json:"gridId"`
		GridX               int    `json:"gridX"`
		GridY               int    `json:"gridY"`
		ObservationStations string `json:"observationStations"`
	} `json:"properties"`
}

// stationFeature represents a single station in the stations response
type stationFeature struct {
	Properties struct {
		StationIdentifier string `json:"stationIdentifier"`
	} `json:"properties"`
}

// stationsResponse represents the response from the NOAA stations API
type stationsResponse struct {
	Features []stationFeature `json:"features"`
}

// observationResponse represents the response from the NOAA observations API
type observationResponse struct {
	Properties struct {
		Temperature struct {
			Value    *float64 `json:"value"` // Use pointer to handle null values
			UnitCode string   `json:"unitCode"`
		} `json:"temperature"`
		TextDescription string `json:"textDescription,omitempty"`
	} `json:"properties"`
}

// getGridpoint fetches the gridpoint information and observation stations URL for the given latitude and longitude.
func getGridpoint(latitude, longitude float64) (string, string, error) {
	logger.DebugWithFields("Fetching gridpoint from NOAA", map[string]interface{}{
		"latitude":  latitude,
		"longitude": longitude,
	})

	pointsURL := fmt.Sprintf("%s/points/%.4f,%.4f", noaaBaseURL, latitude, longitude)

	body, err := httpclient.MakeRequest(pointsURL, userAgent)
	if err != nil {
		logger.ErrorWithFields("Failed to get gridpoint from NOAA points API", map[string]interface{}{
			"error":     err.Error(),
			"latitude":  latitude,
			"longitude": longitude,
			"url":       pointsURL,
		})
		return "", "", fmt.Errorf("failed to get gridpoint from NOAA points API: %w", err)
	}
	defer body.Close()

	var pointsResp pointsResponse
	if err := json.NewDecoder(body).Decode(&pointsResp); err != nil {
		logger.ErrorWithFields("Failed to decode NOAA points API response", map[string]interface{}{
			"error":     err.Error(),
			"latitude":  latitude,
			"longitude": longitude,
		})
		return "", "", fmt.Errorf("failed to decode NOAA points API response: %w", err)
	}

	logger.DebugWithFields("Successfully retrieved gridpoint", map[string]interface{}{
		"gridId": pointsResp.Properties.GridId,
		"gridX":  pointsResp.Properties.GridX,
		"gridY":  pointsResp.Properties.GridY,
	})

	return pointsResp.Properties.GridId, pointsResp.Properties.ObservationStations, nil
}

// determineStatus determines the weather status based on temperature in Fahrenheit.
// Returns "hot" if >75F, "cold" if <50F, and "moderate" for temperatures between 50-75F.
func determineStatus(tempF float64) string {
	if tempF > 75 {
		return "hot"
	} else if tempF < 50 {
		return "cold"
	}
	return "moderate"
}

// GetWeatherFromNOAA fetches weather data from NOAA API for the given latitude and longitude.
// It returns the weather status (hot/cold/moderate) and temperature in Fahrenheit.
func GetWeatherFromNOAA(latitude, longitude float64) (string, string, error) {
	logger.InfoWithFields("Fetching weather from NOAA", map[string]interface{}{
		"latitude":  latitude,
		"longitude": longitude,
	})

	// Step 1: Get gridpoint information and observation stations URL from lat/lon
	gridId, observationStationsURL, err := getGridpoint(latitude, longitude)
	if err != nil {
		return "", "", err
	}

	// Step 2: Get the list of observation stations
	stationsBody, err := httpclient.MakeRequest(observationStationsURL, userAgent)
	if err != nil {
		logger.ErrorWithFields("Failed to get observation stations", map[string]interface{}{
			"error": err.Error(),
			"url":   observationStationsURL,
		})
		return "", "", fmt.Errorf("failed to get observation stations: %w", err)
	}
	defer stationsBody.Close()

	var stationsResp stationsResponse
	if err := json.NewDecoder(stationsBody).Decode(&stationsResp); err != nil {
		logger.ErrorWithFields("Failed to decode stations response", map[string]interface{}{
			"error": err.Error(),
		})
		return "", "", fmt.Errorf("failed to decode stations response: %w", err)
	}

	if len(stationsResp.Features) == 0 {
		logger.WarnWithFields("No observation stations available", map[string]interface{}{
			"gridId": gridId,
		})
		return "", "", fmt.Errorf("no observation stations available")
	}

	// Use the first (closest) station
	stationId := stationsResp.Features[0].Properties.StationIdentifier
	logger.DebugWithFields("Fetching current observations from NOAA station", map[string]interface{}{
		"stationId": stationId,
		"gridId":    gridId,
	})

	// Step 3: Get latest observation from the station
	observationURL := fmt.Sprintf("%s/stations/%s/observations/latest", noaaBaseURL, stationId)
	obsBody, err := httpclient.MakeRequest(observationURL, userAgent)
	if err != nil {
		logger.ErrorWithFields("Failed to get current observations from NOAA", map[string]interface{}{
			"error":     err.Error(),
			"stationId": stationId,
			"url":       observationURL,
		})
		return "", "", fmt.Errorf("failed to get current observations from NOAA: %w", err)
	}
	defer obsBody.Close()

	var observationResp observationResponse
	if err := json.NewDecoder(obsBody).Decode(&observationResp); err != nil {
		logger.ErrorWithFields("Failed to decode NOAA observations response", map[string]interface{}{
			"error":     err.Error(),
			"stationId": stationId,
		})
		return "", "", fmt.Errorf("failed to decode NOAA observations response: %w", err)
	}

	// Check if temperature value is available
	if observationResp.Properties.Temperature.Value == nil {
		logger.WarnWithFields("Temperature data not available from NOAA observations", map[string]interface{}{
			"stationId": stationId,
		})
		return "", "", fmt.Errorf("temperature data not available from NOAA observations")
	}

	// Convert from Celsius to Fahrenheit (NOAA returns Celsius)
	tempC := *observationResp.Properties.Temperature.Value
	tempF := (tempC * 9.0 / 5.0) + 32.0

	status := determineStatus(tempF)

	logger.InfoWithFields("Successfully retrieved current weather data", map[string]interface{}{
		"latitude":  latitude,
		"longitude": longitude,
		"status":    status,
		"tempF":     tempF,
		"tempC":     tempC,
		"gridId":    gridId,
		"stationId": stationId,
	})

	return status, strconv.FormatFloat(tempF, 'f', 1, 64) + "°F", nil
}
