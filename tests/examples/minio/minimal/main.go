package main

import (
	"fmt"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type MinioConfig struct {
	Addr       string
	Bucket     string
	AccessKey  string
	SecretKey  string
	PublicHost string
}

func main() {
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

	minioAdapter := stm.NewMinioAdapter(
		minioConfig.Addr,
		minioConfig.Bucket,
		minioCredentials,
	)

	sm := stm.NewSitemap()

	sm.GetConfig().
		SetDefaultHost("https://example.com").
		SetSitemapsPath("sitemap").
		SetPublicPath("public").
		SetFilename("index").
		SetCompress(false).
		SetAdp(minioAdapter)

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

	log.Println("Sitemap generated")
}
