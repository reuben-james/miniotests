package miniotests

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand" // Alias to avoid conflict with crypto/rand
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func init() {
    mrand.Seed(time.Now().UnixNano())
}

func BenchmarkUploadObject(b *testing.B) {
    // Setup MinIO client
    client, err := minio.New(TestConfig.Endpoint, &minio.Options{
        Creds:     credentials.NewStaticV4(TestConfig.AccessKey, TestConfig.SecretKey, ""),
        Secure:    TestConfig.Secure,
        Transport: TestConfig.Transport,
    })
    if err != nil {
        b.Fatalf("Failed to create MinIO client: %v", err)
    }

    // Create a bucket for benchmarking
    bucketName := "benchmark-bucket-create"
    ctx := context.Background()
    err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
    if err != nil {
        // Ignore if the bucket already exists
        exists, errBucketExists := client.BucketExists(ctx, bucketName)
        if errBucketExists == nil && exists {
            // Bucket exists
        } else {
            b.Fatalf("Failed to create bucket: %v", err)
        }
    }
    defer client.RemoveBucket(ctx, bucketName)

    // Prepare data to upload
    objectName := "test-object"
    dataSize := int64(1024 * 1024) // 1MB
    data := rand.Reader

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err = client.PutObject(ctx, bucketName, fmt.Sprintf("%s-%d", objectName, i), io.LimitReader(data, dataSize), dataSize, minio.PutObjectOptions{})
        if err != nil {
            b.Fatalf("Failed to upload object: %v", err)
        }
    }
    b.StopTimer()

    // Optionally, clean up objects
    for i := 0; i < b.N; i++ {
        err = client.RemoveObject(ctx, bucketName, fmt.Sprintf("%s-%d", objectName, i), minio.RemoveObjectOptions{})
        if err != nil {
            b.Logf("Failed to remove object: %v", err)
        }
    }
}

func BenchmarkDownloadObject(b *testing.B) {
    client, err := minio.New(TestConfig.Endpoint, &minio.Options{
        Creds:     credentials.NewStaticV4(TestConfig.AccessKey, TestConfig.SecretKey, ""),
        Secure:    TestConfig.Secure,
        Transport: TestConfig.Transport,
    })
    if err != nil {
        b.Fatalf("Failed to create MinIO client: %v", err)
    }

    // Create a bucket and upload an object for benchmarking
    bucketName := "benchmark-bucket-read"
    ctx := context.Background()
    err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
    if err != nil {
        exists, errBucketExists := client.BucketExists(ctx, bucketName)
        if errBucketExists == nil && exists {
            // Bucket exists
        } else {
            b.Fatalf("Failed to create bucket: %v", err)
        }
    }
    defer client.RemoveBucket(ctx, bucketName)

    objectName := "test-object"
    dataSize := int64(1024 * 1024) // 1MB
    data := rand.Reader

    _, err = client.PutObject(ctx, bucketName, objectName, io.LimitReader(data, dataSize), dataSize, minio.PutObjectOptions{})
    if err != nil {
        b.Fatalf("Failed to upload object: %v", err)
    }
    defer client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        obj, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
        if err != nil {
            b.Fatalf("Failed to get object: %v", err)
        }
        _, err = io.Copy(io.Discard, obj)
        if err != nil {
            b.Fatalf("Failed to read object data: %v", err)
        }
        obj.Close()
    }
    b.StopTimer()
}

func BenchmarkUploadObjectParallel(b *testing.B) {
    // Setup MinIO client
    client, err := minio.New(TestConfig.Endpoint, &minio.Options{
        Creds:     credentials.NewStaticV4(TestConfig.AccessKey, TestConfig.SecretKey, ""),
        Secure:    TestConfig.Secure,
        Transport: TestConfig.Transport,
    })
    if err != nil {
        b.Fatalf("Failed to create MinIO client: %v", err)
    }

    // Create a bucket for benchmarking
    bucketName := "benchmark-bucket-upload-parallel"
    ctx := context.Background()
    err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
    if err != nil {
        exists, errBucketExists := client.BucketExists(ctx, bucketName)
        if errBucketExists == nil && exists {
            // Bucket exists
        } else {
            b.Fatalf("Failed to create bucket: %v", err)
        }
    }
    defer client.RemoveBucket(ctx, bucketName)

    dataSize := int64(1 * 1024 * 1024) // 1MB

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            objectName := fmt.Sprintf("test-object-%d-%d", time.Now().UnixNano(), mrand.Int())
            _, err := client.PutObject(ctx, bucketName, objectName, io.LimitReader(rand.Reader, dataSize), dataSize, minio.PutObjectOptions{})
            if err != nil {
                b.Errorf("Failed to upload object: %v", err)
                continue
            }
            // Optionally, remove the object
            err = client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
            if err != nil {
                b.Logf("Failed to remove object: %v", err)
            }
        }
    })
    b.StopTimer()
}