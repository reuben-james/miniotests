// health_test.go
package miniotests

import (
	"net/http"
	"testing"
	"time"
)

func TestHealthChecks(t *testing.T) {
    t.Log("Starting TestHealthChecks")

    // Ensure that TestConfig.Endpoint is not empty
    if TestConfig.Endpoint == "" {
        t.Fatal("TestConfig.Endpoint is empty. Ensure that MINIO_SERVER and MINIO_PORT environment variables are set.")
    }

    address := TestConfig.Endpoint
    healthLiveURL := TestConfig.Scheme + "://" + address + "/minio/health/live"
    healthReadyURL := TestConfig.Scheme + "://" + address + "/minio/health/ready"

    t.Logf("Health live endpoint: %s", healthLiveURL)
    t.Logf("Health ready endpoint: %s", healthReadyURL)



    httpClient := &http.Client{
        Timeout: 5 * time.Second,
    }

    // Configure TLS settings if secure
    if TestConfig.Secure {
        t.Log("Configuring TLS settings")
        httpClient.Transport = TestConfig.Transport
    }

    // Liveness Probe
    t.Run("Liveness Probe", func(t *testing.T) {
        t.Log("Sending GET request to liveness endpoint")
        resp, err := httpClient.Get(healthLiveURL)
        if err != nil {
            t.Fatalf("Liveness probe failed: %v", err)
        }
        defer resp.Body.Close()

        t.Logf("Received response with status code: %d", resp.StatusCode)
        if resp.StatusCode != http.StatusOK {
            t.Fatalf("Unexpected status code for liveness probe: %d", resp.StatusCode)
        }
    })

    // Readiness Probe
    t.Run("Readiness Probe", func(t *testing.T) {
        t.Log("Sending GET request to readiness endpoint")
        resp, err := httpClient.Get(healthReadyURL)
        if err != nil {
            t.Fatalf("Readiness probe failed: %v", err)
        }
        defer resp.Body.Close()

        t.Logf("Received response with status code: %d", resp.StatusCode)
        if resp.StatusCode != http.StatusOK {
            t.Fatalf("Unexpected status code for readiness probe: %d", resp.StatusCode)
        }
    })

    t.Log("TestHealthChecks completed successfully")
}
