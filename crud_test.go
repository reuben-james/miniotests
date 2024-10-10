// crud_test.go
package miniotests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestCRUDOperations(t *testing.T) {
    t.Log("Starting TestCRUDOperations")

    // Use the centralized TestConfig
    endpoint := TestConfig.Endpoint
    accessKey := TestConfig.AccessKey
    secretKey := TestConfig.SecretKey
    secure := TestConfig.Secure
    transport := TestConfig.Transport

    client, err := minio.New(endpoint, &minio.Options{
        Creds:     credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure:    secure,
        Transport: transport,
    })
    if err != nil {
        t.Fatalf("Failed to create MinIO client: %v", err)
    }

    // Test parameters
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
    t.Logf("Creating bucket: %s", bucketName)
    err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
    if err != nil {
        t.Fatalf("Bucket creation failed: %v", err)
    }
    defer func() {
        err = client.RemoveBucket(context.Background(), bucketName)
        if err != nil {
            t.Logf("Failed to remove bucket %s: %v", bucketName, err)
        }
    }()

    // Upload Object
    t.Logf("Uploading object: %s", testFileName)
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
    objectCh := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
    for object := range objectCh {
        if object.Err != nil {
            t.Fatalf("Error listing objects: %v", object.Err)
        }
        if object.Key == testFileName {
            t.Fatalf("Object %s was not deleted", testFileName)
        }
    }

    t.Log("TestCRUDOperations completed successfully")
}
