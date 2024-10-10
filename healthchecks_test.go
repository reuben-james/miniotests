// health_test.go
package miniotests

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestHealthChecks(t *testing.T) {
    server := os.Getenv("MINIO_SERVER")
    port := os.Getenv("MINIO_PORT")
    if port == "" {
        port = "9000"
    }

    secure := *useTLS
    scheme := "http"
    if secure {
        scheme = "https"
    }
    address := net.JoinHostPort(server, port)
    healthLiveEndpoint := scheme + "://" + address + "/minio/health/live"
    healthReadyEndpoint := scheme + "://" + address + "/minio/health/ready"

    httpClient := &http.Client{
        Timeout: 5 * time.Second,
    }

    if secure {
        httpClient.Transport = &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
    }

    // Liveness Probe
    resp, err := httpClient.Get(healthLiveEndpoint)
    if err != nil {
        t.Fatalf("Liveness probe failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Unexpected status code for liveness probe: %d", resp.StatusCode)
    }

    // Readiness Probe
    resp, err = httpClient.Get(healthReadyEndpoint)
    if err != nil {
        t.Fatalf("Readiness probe failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Unexpected status code for readiness probe: %d", resp.StatusCode)
    }
}
