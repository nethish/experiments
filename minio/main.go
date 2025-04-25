package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// MinIO server connection parameters
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"     // default MinIO access key
	secretAccessKey := "minioadmin" // default MinIO secret key
	useSSL := false                 // set to true if your MinIO uses SSL

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Error creating MinIO client:", err)
	}

	// Define bucket name
	bucketName := "my-test-bucket"

	// Check if bucket exists and create it if it doesn't
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalln("Error checking bucket:", err)
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln("Error creating bucket:", err)
		}
		fmt.Printf("Successfully created bucket %s\n", bucketName)
	} else {
		fmt.Printf("Bucket %s already exists\n", bucketName)
	}

	// Upload an object (write operation)
	objectName := "example-object.txt"
	content := []byte("This is a sample content to upload to MinIO")
	contentType := "text/plain"

	// Upload using PutObject
	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, 
		bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("Error uploading object:", err)
	}
	fmt.Printf("Successfully uploaded %s of size %d\n", objectName, len(content))

	// Download an object (read operation)
	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("Error getting object:", err)
	}
	defer object.Close()

	localFile, err := os.Create("downloaded-" + objectName)
	if err != nil {
		log.Fatalln("Error creating local file:", err)
	}
	defer localFile.Close()

	stat, err := object.Stat()
	if err != nil {
		log.Fatalln("Error getting object info:", err)
	}

	if _, err := io.Copy(localFile, object); err != nil {
		log.Fatalln("Error copying object to file:", err)
	}

	fmt.Printf("Successfully downloaded %s of size %d\n", objectName, stat.Size)

	// List all objects in the bucket
	fmt.Println("Objects in bucket:")
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error listing objects:", object.Err)
			continue
		}
		fmt.Printf("- %s (size: %d)\n", object.Key, object.Size)
	}

	// Delete the object
	err = minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln("Error removing object:", err)
	}
	fmt.Printf("Successfully deleted %s\n", objectName)
}
