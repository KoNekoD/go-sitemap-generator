package main

import (
	"fmt"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"log"
)

func main() {
	fileAdapter := stm.NewFileAdapter()

	sm := stm.NewSitemap()

	sm.GetConfig().
		SetDefaultHost("https://example.com").
		SetSitemapsPath("sitemap").
		SetPublicPath("public").
		SetFilename("index").
		SetCompress(false).
		SetAdp(fileAdapter)

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
