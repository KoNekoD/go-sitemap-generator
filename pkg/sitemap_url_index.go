package stm

import (
	"github.com/beevik/etree"
	"time"
)

// sitemapIndexURL and sitemapURL are almost the same behavior.
type sitemapIndexURL struct {
	*Config
	data URL
}

func NewSitemapIndexURL(c *Config, url URL) SitemapURL {
	return &sitemapIndexURL{Config: c, data: url}
}

// XML and sitemapIndexURL.XML are almost the same behavior.
func (su *sitemapIndexURL) XML() []byte {
	doc := etree.NewDocument()
	sitemap := doc.CreateElement("sitemap")
	SetBuilderElementValue(sitemap, su.data, "loc")
	if _, ok := SetBuilderElementValue(sitemap, su.data, "lastmod"); !ok {
		lastMod := sitemap.CreateElement("lastmod")
		lastMod.SetText(time.Now().Format(time.RFC3339))
	}
	if su.Config.Pretty {
		doc.Indent(2)
	}
	buf := poolBuffer.Get()
	_, _ = doc.WriteTo(buf)
	bytes := buf.Bytes()
	poolBuffer.Put(buf)

	return bytes
}
