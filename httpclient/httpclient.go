package httpclient

import (
	"fmt"
	"io"
	"net/http"
)

// MakeRequest performs an HTTP GET request to the specified URL and returns the response body.
// It handles request creation, setting headers, and error checking.
// userAgent is required by some APIs (like NOAA) and should be provided.
func MakeRequest(url, userAgent string) (io.ReadCloser, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("API returned status: %d %s", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}
