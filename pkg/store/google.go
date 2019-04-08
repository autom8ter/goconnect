package store

import (
	gstorage "cloud.google.com/go/storage"
	"context"
	"github.com/autom8ter/goconnect/pkg/store/driver"
)

type GoogleCloudStorage struct {
	driver.Storage
	client     *gstorage.Client
	bucket     *gstorage.BucketHandle
	bucketName string
}

func NewGoogleCloudStorage(client *gstorage.Client, bucket *gstorage.BucketHandle, bucketName string) *GoogleCloudStorage {
	return &GoogleCloudStorage{client: client, bucket: bucket, bucketName: bucketName}
}

func (gcs *GoogleCloudStorage) PublicURL(filename string) string {
	return "https://storage.googleapis.com/" + gcs.bucketName + "/" + filename
}

func (gcs *GoogleCloudStorage) Store(ctx context.Context, filename string, data []byte, metadata map[string]string) error {
	o := gcs.bucket.Object(filename)
	w := o.NewWriter(ctx)

	w.ObjectAttrs = gstorage.ObjectAttrs{
		Name:     filename,
		Metadata: metadata,
	}

	_, err := w.Write(data)
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (gcs *GoogleCloudStorage) Delete(ctx context.Context, filename string) error {
	o := gcs.bucket.Object(filename)
	return o.Delete(ctx)
}
