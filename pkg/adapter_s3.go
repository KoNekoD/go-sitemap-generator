package stm

import (
	"bytes"
	"compress/gzip"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"net/http"
)

type S3Adapter struct {
	Region      string
	Bucket      string
	ACL         string
	Credentials *credentials.Credentials
	HttpClient  *http.Client
}

func NewS3Adapter(region, bucket, acl string, c *credentials.Credentials) *S3Adapter {
	return &S3Adapter{Region: region, Bucket: bucket, ACL: acl, Credentials: c, HttpClient: http.DefaultClient}
}

func (a *S3Adapter) Bytes() [][]byte { return nil }

func (a *S3Adapter) Write(loc *Location, data []byte) {
	var reader io.Reader = bytes.NewReader(data)

	if GzipPtn.MatchString(loc.Filename()) {
		var writer *io.PipeWriter

		reader, writer = io.Pipe()
		go func() {
			gz := gzip.NewWriter(writer)
			_, _ = io.Copy(gz, bytes.NewReader(data))
			_ = gz.Close()
			_ = writer.Close()
		}()
	}

	if a.Credentials == nil {
		a.Credentials = credentials.NewEnvCredentials()
	}
	_, _ = a.Credentials.Get()

	s, _ := session.NewSession(&aws.Config{Credentials: a.Credentials, Region: &a.Region, HTTPClient: a.HttpClient})

	uploader := s3manager.NewUploader(s)

	bucket, key, acl := aws.String(a.Bucket), aws.String(loc.Path()), aws.String(a.ACL)

	_, err := uploader.Upload(&s3manager.UploadInput{Bucket: bucket, Key: key, ACL: acl, Body: reader})

	if err != nil {
		log.Fatal("[F] S3 Upload file Error:", err)
	}
}
