package main

import (
	"context"
	"fmt"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"log"
	"slices"
	"time"
)

func main() {
	err := NewSitemapGenerator().Generate()
	if err != nil {
		log.Fatal(err)
	}
}

type MinioConfig struct {
	Addr       string
	Bucket     string
	AccessKey  string
	SecretKey  string
	PublicHost string
}

type SitemapGenerator struct {
	minioClient  *minio.Client
	minioAdapter *stm.MinioAdapter
	minioConfig  *MinioConfig
}

func NewSitemapGenerator() *SitemapGenerator {
	minioConfig := &MinioConfig{
		Addr:       "localhost:9000",
		Bucket:     "my-bucket",
		AccessKey:  "123",
		SecretKey:  "456",
		PublicHost: "https://example.com",
	}

	minioCredentials := credentials.NewStaticV4(
		minioConfig.AccessKey,
		minioConfig.SecretKey,
		"",
	)

	minioOptions := &minio.Options{Creds: minioCredentials}

	minioClient, err := minio.New(minioConfig.Addr, minioOptions)
	if err != nil {
		log.Fatal("MinIO connection error:", "err", err)
	}

	minioAdapter := stm.NewMinioAdapter(
		minioConfig.Addr,
		minioConfig.Bucket,
		minioCredentials,
	)

	return &SitemapGenerator{
		minioClient:  minioClient,
		minioAdapter: minioAdapter,
		minioConfig:  minioConfig,
	}
}

func (p *SitemapGenerator) getFilesList(modifiedAfter *time.Time) []string {
	ctx := context.Background()
	filenames := make([]string, 0)

	listOpts := minio.ListObjectsOptions{
		Prefix:    "public/sitemap/",
		Recursive: true,
	}

	iter := p.minioClient.ListObjects(ctx, p.minioConfig.Bucket, listOpts)
	for oldFile := range iter {
		if modifiedAfter != nil && oldFile.LastModified.Before(*modifiedAfter) {
			continue
		}

		filenames = append(filenames, oldFile.Key)
	}

	return filenames
}

func (p *SitemapGenerator) removeOrphanFiles(
	oldFilenames []string,
	newFilenames []string,
) error {
	ctx := context.Background()
	for _, oldFilename := range oldFilenames {
		if slices.Contains(newFilenames, oldFilename) {
			continue
		}

		removeOpts := minio.RemoveObjectOptions{}

		err := p.minioClient.RemoveObject(
			ctx,
			p.minioConfig.Bucket,
			oldFilename,
			removeOpts,
		)
		if err != nil {
			return errors.Wrap(err, "Error deleting old file")
		}

		log.Println("Deleted orphan old file", oldFilename)
	}

	return nil
}

func (p *SitemapGenerator) Generate() error {
	oldFilenames := p.getFilesList(nil)
	now := time.Now().UTC()
	sm := stm.NewSitemap()

	sm.GetConfig().
		SetDefaultHost("https://example.com").
		SetSitemapsPath("sitemap").
		SetPublicPath("public").
		SetFilename("index").
		SetCompress(false).
		SetAdp(p.minioAdapter)

	sm.Create()

	getSitemapFilename := func(count int) string {
		return fmt.Sprintf("users-%03d.xml", count)
	}
	sm.GetNamer().Opts.SetBuildName(getSitemapFilename)

	priority := 0.8

	link := "example.com/1/2/3"

	sm.Add(
		stm.URL{
			{"loc", link},
			{"changefreq", "weekly"},
			{"priority", priority},
		},
	)

	sm.Finalize()

	// can be used when getSitemapFilename method will change, reset counter
	sm.GetNamer().Reset()

	newFilenames := p.getFilesList(&now)

	if err := p.removeOrphanFiles(oldFilenames, newFilenames); err != nil {
		return errors.Wrap(err, "Error removing orphan files")
	}

	log.Println("Sitemap generated")

	return nil
}
