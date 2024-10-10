// connectivity_test.go
package miniotests

import (
	"net/http"
	"testing"
	"time"
)

func TestConnectivity(t *testing.T) {
    t.Log("Starting TestConnectivity")

    secure := TestConfig.Secure
    t.Logf("Scheme: %s", TestConfig.Scheme)

    address := TestConfig.Endpoint
    url := TestConfig.Scheme + "://" + address
    t.Logf("Connecting to Endpoint: %s", url)

    httpClient := &http.Client{
        Timeout: 5 * time.Second,
    }

    if secure {
        t.Log("Configuring TLS settings")
        httpClient.Transport = TestConfig.Transport
    }

    // Test HTTP(S) connectivity
    t.Log("Sending HTTP GET request")
    response, err := httpClient.Get(url)
    if err != nil {
        t.Fatalf("HTTP connectivity test failed: %v", err)
    }
    defer response.Body.Close()
    t.Logf("Received response with status code: %d", response.StatusCode)

    if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusForbidden {
        t.Fatalf("Unexpected HTTP status code: %d", response.StatusCode)
    }

    t.Log("TestConnectivity completed successfully")
}
