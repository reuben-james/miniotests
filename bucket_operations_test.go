package miniotests

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestCRUDOperations(t *testing.T) {
    server := os.Getenv("MINIO_SERVER")
    port := os.Getenv("MINIO_PORT")
    if port == "" {
        port = "9000"
    }
    accessKey := os.Getenv("MINIO_ACCESS_KEY")
    secretKey := os.Getenv("MINIO_SECRET_KEY")
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

    bucketName := "test-bucket-" + time.Now().Format("20060102150405")
    testFileName := "testfile.txt"
    testFileContent := "This is a test file."

    // Create test file
    err = os.WriteFile(testFileName, []byte(testFileContent), 0644)
    if err != nil {
        t.Fatalf("Failed to create test file: %v", err)
    }

    defer os.Remove(testFileName)

    // Create Bucket
    err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
    if err != nil {
        t.Fatalf("Bucket creation failed: %v", err)
    }

    defer client.RemoveBucket(context.Background(), bucketName)

    // Upload Object
    _, err = client.FPutObject(context.Background(), bucketName, testFileName, testFileName, minio.PutObjectOptions{})
    if err != nil {
        t.Fatalf("File upload failed: %v", err)
    }

    // Download Object
    downloadedFileName := "downloaded_" + testFileName
    err = client.FGetObject(context.Background(), bucketName, testFileName, downloadedFileName, minio.GetObjectOptions{})
    if err != nil {
        t.Fatalf("File download failed: %v", err)
    }

    defer os.Remove(downloadedFileName)

    // Verify File Integrity
    originalContent, err := os.ReadFile(testFileName)
    if err != nil {
        t.Fatalf("Failed to read original file: %v", err)
    }

    downloadedContent, err := os.ReadFile(downloadedFileName)
    if err != nil {
        t.Fatalf("Failed to read downloaded file: %v", err)
    }

    if string(originalContent) != string(downloadedContent) {
        t.Fatalf("File content mismatch")
    }

    // Update Object
    updatedContent := testFileContent + "\nAdding a new line."
    err = os.WriteFile(testFileName, []byte(updatedContent), 0644)
    if err != nil {
        t.Fatalf("Failed to update test file: %v", err)
    }

    _, err = client.FPutObject(context.Background(), bucketName, testFileName, testFileName, minio.PutObjectOptions{})
    if err != nil {
        t.Fatalf("File update failed: %v", err)
    }

    // Delete Object
    err = client.RemoveObject(context.Background(), bucketName, testFileName, minio.RemoveObjectOptions{})
    if err != nil {
        t.Fatalf("File deletion failed: %v", err)
    }

    // Verify Deletion
    objectsCh := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
    for obj := range objectsCh {
        if obj.Err != nil {
            t.Fatalf("Error listing objects: %v", obj.Err)
        }
        if obj.Key == testFileName {
            t.Fatalf("Object %s was not deleted", testFileName)
        }
    }
}
