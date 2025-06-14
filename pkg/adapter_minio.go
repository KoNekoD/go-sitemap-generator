package stm

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/http"
)

type MinioAdapter struct {
	Host        string
	Bucket      string
	Credentials *credentials.Credentials
	Transport   http.RoundTripper
	clientOpts  *minio.Options
}

func NewMinioAdapter(host, bucket string, c *credentials.Credentials, opts ...*minio.Options) *MinioAdapter {
	if c == nil {
		c = credentials.NewEnvMinio()
	}

	clientOpts := &minio.Options{Creds: c}
	if len(opts) > 0 {
		if clientOpts = opts[0]; clientOpts.Creds == nil {
			clientOpts.Creds = c
		}
	}

	return &MinioAdapter{Host: host, Bucket: bucket, Credentials: c, clientOpts: clientOpts}
}

func (a *MinioAdapter) Bytes() [][]byte { return nil }

func (a *MinioAdapter) Write(loc *Location, data []byte) {
	ctx := context.Background()
	buffer := bytes.NewBuffer(data)

	if GzipPtn.MatchString(loc.Filename()) {
		buffer = bytes.NewBuffer(nil)
		gz := gzip.NewWriter(buffer)
		_, _ = io.Copy(gz, bytes.NewReader(data))
		_ = gz.Close()
	}

	if a.Transport != nil {
		a.clientOpts.Transport = a.Transport
	}

	minioClient, err := minio.New(a.Host, a.clientOpts)
	if err != nil {
		log.Fatalln("MinIO connection error:", err)
	}

	objectName, objectSize := loc.Path(), int64(buffer.Len())
	putOpts := minio.PutObjectOptions{ContentType: "application/xml"}

	_, err = minioClient.PutObject(ctx, a.Bucket, objectName, buffer, objectSize, putOpts)
	if err != nil {
		log.Fatal("[F] MinIO Upload file Error:", err)
	}
}
