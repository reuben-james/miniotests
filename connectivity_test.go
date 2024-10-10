// connectivity_test.go
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
    url := scheme + "://" + address

    // Set up HTTP client
    httpClient := &http.Client{
        Timeout: 5 * time.Second,
    }

    if secure {
        httpClient.Transport = &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
    }

    response, err := httpClient.Get(url)
    if err != nil {
        t.Fatalf("HTTP connectivity test failed: %v", err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusForbidden {
        t.Fatalf("Unexpected HTTP status code: %d", response.StatusCode)
    }
}
