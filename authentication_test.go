// authentication_test.go
package miniotests

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestAuthentication(t *testing.T) {
    t.Log("Starting TestAuthentication")

    endpoint := TestConfig.Endpoint
    accessKey := TestConfig.AccessKey
    secretKey := TestConfig.SecretKey
    secure := TestConfig.Secure
    transport := TestConfig.Transport

    t.Logf("Using endpoint: %s", endpoint)
    t.Logf("Secure mode: %v", secure)

    // Create client with good credentials
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

    t.Log("TestAuthentication completed successfully")
}