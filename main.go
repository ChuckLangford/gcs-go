package main

// GOOGLE_APPLICATION_CREDENTIALS
// GOOGLE_CLOUD_PROJECT

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"os"
)

func listBuckets(client *storage.Client, projectID string) ([]string, error) {
	ctx := context.Background()
	var buckets []string
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}
	return buckets, nil
}

func listObjectsInBucket(client *storage.Client, bucket string) ([]string, error) {
	ctx := context.Background()
	var objectNames []string
	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objectNames = append(objectNames, attrs.Name)
	}
	return objectNames, nil
}

func main() {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	googleAuth := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if googleAuth == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_APPLICATION_CREDENTIALS environment variable must be set.\n")
		os.Exit(1)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	buckets, err := listBuckets(client, projectID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("buckets: %+v\n", buckets)

	for _, bucket := range buckets {
		fmt.Println()
		fmt.Println(bucket)
		objectNames, err := listObjectsInBucket(client, bucket)
		if err != nil {
			log.Fatal(err)
		}

		for _, objectName := range objectNames {
			fmt.Println(objectName)
		}
	}
}
