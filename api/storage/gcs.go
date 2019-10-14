package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type GCS struct {
	client   *storage.BucketHandle
	BucketID string
}

func NewGCS(bucketID string) Storage {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	return &GCS{
		client:   client.Bucket(bucketID),
		BucketID: bucketID,
	}
}

func (s GCS) Store(filename string, contents []byte) (string, error) {
	ctx := context.Background()
	w := s.client.Object(filename).NewWriter(ctx)
	defer w.Close()

	w.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	w.CacheControl = "public, max-age=86400"
	if _, err := io.Copy(w, bytes.NewReader(contents)); err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.BucketID, filename)
	return publicURL, nil
}

func (s GCS) Get(filename string) ([]byte, error) {
	ctx := context.Background()
	reader, err := s.client.Object(filename).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "error in reading file to GCS")
	}
	return contents, nil
}

func (s GCS) Delete(filename string) error {
	return nil
}
