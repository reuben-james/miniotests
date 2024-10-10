// authentication_test.go
package miniotests

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestAuthentication(t *testing.T) {
    server := os.Getenv("MINIO_SERVER")
    port := os.Getenv("MINIO_PORT")
    if port == "" {
        port = "9000"
    }
    accessKey := os.Getenv("MINIO_ACCESS_KEY")
    secretKey := os.Getenv("MINIO_SECRET_KEY")

    secure := *useTLS

    endpoint := net.JoinHostPort(server, port)

    var transport *http.Transport
    if secure {
        transport = &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
    }

    // Create MinIO client with valid credentials
    client, err := minio.New(endpoint, &minio.Options{
        Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure:    secure,
        Transport: transport,
    })
    if err != nil {
        t.Fatalf("Failed to create MinIO client: %v", err)
    }

    // Attempt to list buckets
    _, err = client.ListBuckets(context.Background())
    if err != nil {
        t.Fatalf("Valid credentials failed: %v", err)
    }

    // Create MinIO client with invalid credentials
    invalidClient, err := minio.New(endpoint, &minio.Options{
        Creds:     credentials.NewStaticV4("INVALID_KEY", "INVALID_SECRET", ""),
        Secure:    secure,
        Transport: transport,
    })
    if err != nil {
        t.Fatalf("Failed to create MinIO client with invalid credentials: %v", err)
    }

    // Attempt to list buckets with invalid credentials
    _, err = invalidClient.ListBuckets(context.Background())
    if err == nil {
        t.Fatalf("Invalid credentials should not be authenticated")
    }
}