package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMultipleURLRedirects(t *testing.T) {
	// Save original flag values and restore after test
	originalStsEnable := *stsEnable
	defer func() { *stsEnable = originalStsEnable }()
	*stsEnable = true

	tests := []struct {
		name           string
		requestURL     string
		requestHost    string
		expectedURL    string
		expectedStatus int
	}{
		{
			name:           "Root path redirect",
			requestURL:     "http://example.com/",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with single segment",
			requestURL:     "http://example.com/about",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/about",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with multiple segments",
			requestURL:     "http://example.com/api/v1/users",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/api/v1/users",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with query parameters",
			requestURL:     "http://example.com/search?q=test&page=1",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/search?q=test&page=1",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with complex query string",
			requestURL:     "http://example.com/api?filter=active&sort=name&limit=10",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/api?filter=active&sort=name&limit=10",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Subdomain redirect",
			requestURL:     "http://api.example.com/v1/data",
			requestHost:    "api.example.com",
			expectedURL:    "https://api.example.com/v1/data",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Host with port stripped",
			requestURL:     "http://example.com:8080/page",
			requestHost:    "example.com:8080",
			expectedURL:    "https://example.com/page",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with trailing slash",
			requestURL:     "http://example.com/path/to/resource/",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/path/to/resource/",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with special characters",
			requestURL:     "http://example.com/path/with-dashes_and_underscores",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/path/with-dashes_and_underscores",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with encoded characters in query",
			requestURL:     "http://example.com/search?q=hello%20world&category=tech",
			requestHost:    "example.com",
			expectedURL:    "https://example.com/search?q=hello%20world&category=tech",
			expectedStatus: http.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest("GET", tt.requestURL, nil)
			req.Host = tt.requestHost

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the redirect handler
			redirect(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check Location header
			location := rr.Header().Get("Location")
			if location != tt.expectedURL {
				t.Errorf("handler returned wrong redirect URL: got %v want %v",
					location, tt.expectedURL)
			}

			// Check Server header
			serverHeader := rr.Header().Get("Server")
			expectedServer := "Https-echo/" + Version + " (+" + Homepage + ")"
			if serverHeader != expectedServer {
				t.Errorf("handler returned wrong Server header: got %v want %v",
					serverHeader, expectedServer)
			}

			// Check Strict-Transport-Security header
			stsHeader := rr.Header().Get("Strict-Transport-Security")
			if stsHeader != "max-age=31536000" {
				t.Errorf("handler returned wrong STS header: got %v want %v",
					stsHeader, "max-age=31536000")
			}
		})
	}
}

func TestMultipleURLRedirectsWithCustomPort(t *testing.T) {
	// Save original flag values and restore after test
	originalDestPort := *destPort
	originalStsEnable := *stsEnable
	defer func() {
		*destPort = originalDestPort
		*stsEnable = originalStsEnable
	}()
	*destPort = 8443
	*stsEnable = true

	tests := []struct {
		name           string
		requestURL     string
		requestHost    string
		expectedURL    string
		expectedStatus int
	}{
		{
			name:           "Root path with custom port",
			requestURL:     "http://example.com/",
			requestHost:    "example.com",
			expectedURL:    "https://example.com:8443/",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Path with custom port",
			requestURL:     "http://example.com/api/data",
			requestHost:    "example.com",
			expectedURL:    "https://example.com:8443/api/data",
			expectedStatus: http.StatusTemporaryRedirect,
		},
		{
			name:           "Query string with custom port",
			requestURL:     "http://example.com/search?q=test",
			requestHost:    "example.com",
			expectedURL:    "https://example.com:8443/search?q=test",
			expectedStatus: http.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.requestURL, nil)
			req.Host = tt.requestHost
			rr := httptest.NewRecorder()

			redirect(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			location := rr.Header().Get("Location")
			if location != tt.expectedURL {
				t.Errorf("handler returned wrong redirect URL: got %v want %v",
					location, tt.expectedURL)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/health", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	health(rr, req)

	// Health endpoint should return 200 OK, not redirect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("health endpoint returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Healthy.\n"
	if rr.Body.String() != expected {
		t.Errorf("health endpoint returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestSTSHeaderDisabled(t *testing.T) {
	// Save original flag value and restore after test
	originalStsEnable := *stsEnable
	defer func() { *stsEnable = originalStsEnable }()
	*stsEnable = false

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	redirect(rr, req)

	// Check that STS header is not present
	stsHeader := rr.Header().Get("Strict-Transport-Security")
	if stsHeader != "" {
		t.Errorf("STS header should not be present when disabled, got: %v", stsHeader)
	}
}
