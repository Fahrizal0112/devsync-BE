package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSStorage struct {
	client     *storage.Client
	bucketName string
}

func NewGCSStorage(projectID, bucketName, credentialsPath string) (*GCSStorage, error) {
	ctx := context.Background()
	
	var client *storage.Client
	var err error
	
	if credentialsPath != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	} else {
		// Use default credentials (for production with service account)
		client, err = storage.NewClient(ctx)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}

	return &GCSStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (g *GCSStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%s/%d_%s", folder, time.Now().Unix(), file.Filename)
	
	obj := g.client.Bucket(g.bucketName).Object(filename)
	writer := obj.NewWriter(ctx)
	
	// Set content type
	writer.ContentType = file.Header.Get("Content-Type")
	if writer.ContentType == "" {
		writer.ContentType = "application/octet-stream"
	}

	if _, err := io.Copy(writer, src); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	// Return public URL
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, filename), nil
}

func (g *GCSStorage) DeleteFile(ctx context.Context, filename string) error {
	obj := g.client.Bucket(g.bucketName).Object(filename)
	return obj.Delete(ctx)
}

func (g *GCSStorage) GetSignedURL(ctx context.Context, filename string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}

	url, err := g.client.Bucket(g.bucketName).SignedURL(filename, opts)
	if err != nil {
		return "", err
	}
	
	return url, nil
}

func (g *GCSStorage) Close() error {
	return g.client.Close()
}