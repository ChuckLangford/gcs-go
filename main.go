package main

// GOOGLE_APPLICATION_CREDENTIALS
// GOOGLE_CLOUD_PROJECT

import (
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
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
		for _, meta := range attrs.Metadata {
			objectNames = append(objectNames, meta)
		}
	}
	return objectNames, nil
}

func write(client *storage.Client, bucket, object string, fileName string) error {
	ctx := context.Background()

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func delete(client *storage.Client, bucket, object string) error {
	ctx := context.Background()
	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	var a string
	flag.StringVar(&a, "a", "", "This is the action to take. Valid values: write or delete.")
	flag.Parse()
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

	bucket := os.Getenv("BUCKET_TO_USE")
	if bucket == "" {
		fmt.Fprintf(os.Stderr, "BUCKET_TO_USE environment variable must be set.\n")
		os.Exit(1)
	}

	fileToUpload := os.Getenv("FILE_TO_UPLOAD")
	if fileToUpload == "" {
		fmt.Fprintf(os.Stderr, "FILE_TO_UPLOAD environment variable must be set.\n")
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

	if a == "write" {
		write(client, bucket, "testobject", fileToUpload)
	} else if a == "delete" {
		delete(client, bucket, "testobject")
	}

	objectNames, err := listObjectsInBucket(client, bucket)
	if err != nil {
		log.Fatal(err)
	}

	for _, objectName := range objectNames {
		fmt.Println(objectName)
	}
}
