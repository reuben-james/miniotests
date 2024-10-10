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
    secure := os.Getenv("MINIO_SECURE") == "true"

    endpoint := net.JoinHostPort(server, port)

    // Valid Credentials
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: secure,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    })
    if err != nil {
        t.Fatalf("Failed to create MinIO client: %v", err)
    }

    // Attempt to list buckets
    _, err = client.ListBuckets(context.Background())
    if err != nil {
        t.Fatalf("Valid credentials failed: %v", err)
    }

    // Invalid Credentials
    invalidClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4("INVALID_KEY", "INVALID_SECRET", ""),
        Secure: secure,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    })
    if err != nil {
        t.Fatalf("Failed to create MinIO client with invalid credentials: %v", err)
    }

    _, err = invalidClient.ListBuckets(context.Background())
    if err == nil {
        t.Fatalf("Invalid credentials should not be authenticated")
    }
}

func TestAuthorization(t *testing.T) {
    server := os.Getenv("MINIO_SERVER")
    port := os.Getenv("MINIO_PORT")
    if port == "" {
        port = "9000"
    }
    accessKey := os.Getenv("MINIO_READONLY_ACCESS_KEY")
    secretKey := os.Getenv("MINIO_READONLY_SECRET_KEY")
    secure := os.Getenv("MINIO_SECURE") == "true"

    endpoint := net.JoinHostPort(server, port)

    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: secure,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    })
    if err != nil {
        t.Fatalf("Failed to create MinIO client: %v", err)
    }

    bucketName := "test-bucket-auth"

    // Attempt to create a bucket with read-only user
    err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
    if err == nil {
        t.Fatalf("Read-only user should not be able to create buckets")
    }
}
