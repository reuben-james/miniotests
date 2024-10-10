package miniotests

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestConnectivity(t *testing.T) {
    server := os.Getenv("MINIO_SERVER")
    port := os.Getenv("MINIO_PORT") // Default is 9000
    if port == "" {
        port = "9000"
    }
    address := net.JoinHostPort(server, port)

    // DNS Resolution
    _, err := net.LookupHost(server)
    if err != nil {
        t.Fatalf("DNS resolution failed for %s: %v", server, err)
    }

    // Port Availability
    conn, err := net.DialTimeout("tcp", address, 5*time.Second)
    if err != nil {
        t.Fatalf("Port %s is not open on %s: %v", port, server, err)
    }
    conn.Close()

    // HTTP Connectivity
    secure := os.Getenv("MINIO_SECURE")
    scheme := "http"
    if secure == "true" {
        scheme = "https"
    }
    url := scheme + "://" + address
    httpClient := &http.Client{
        Timeout: 5 * time.Second,
    }

    // Skip SSL verification if using self-signed certificates
    if secure == "true" {
        httpClient.Transport = &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
    }

    resp, err := httpClient.Get(url)
    if err != nil {
        t.Fatalf("HTTP connectivity test failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
        t.Fatalf("Unexpected HTTP status code: %d", resp.StatusCode)
    }
}
