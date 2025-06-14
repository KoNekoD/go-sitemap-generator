package stm

import (
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"time"
)

// sitemapURL provides xml validator and xml builder.
type sitemapURL struct {
	*Config
	data URL
}

// NewSitemapURL returns the created the SitemapURL's pointer, and it validates URL types error.
func NewSitemapURL(c *Config, url URL) (SitemapURL, error) {
	smu := &sitemapURL{Config: c, data: url}
	err := smu.validate()
	return smu, err
}

// validate is checking correct keys and checks the existence.
func (su *sitemapURL) validate() error {
	var key string
	var invalid bool
	var locOk, hostOk bool

	for _, value := range su.data {
		key = value[0].(string)
		switch key {
		case "loc":
			locOk = true
		case "host":
			hostOk = true
		}

		invalid = true
		for _, name := range urlModelFieldNames {
			if key == name {
				invalid = false
				break
			}
		}
		if invalid {
			break
		}
	}

	if invalid {
		msg := fmt.Sprintf("Unknown map's key `%s` in URL type", key)
		return errors.New(msg)
	}
	if !locOk {
		msg := fmt.Sprintf("URL type must have `loc` map's key")
		return errors.New(msg)
	}
	if !hostOk {
		msg := fmt.Sprintf("URL type must have `host` map's key")
		return errors.New(msg)
	}
	return nil
}

// XML is building xml.
func (su *sitemapURL) XML() []byte {
	doc := etree.NewDocument()
	url := doc.CreateElement("url")

	SetBuilderElementValue(url, su.data.URLJoinBy("loc", "host", "loc"), "loc")
	if _, ok := SetBuilderElementValue(url, su.data, "lastmod"); !ok {
		lastmod := url.CreateElement("lastmod")
		lastmod.SetText(time.Now().Format(time.RFC3339))
	}
	if _, ok := SetBuilderElementValue(url, su.data, "changefreq"); !ok {
		changefreq := url.CreateElement("changefreq")
		changefreq.SetText("weekly")
	}
	if _, ok := SetBuilderElementValue(url, su.data, "priority"); !ok {
		priority := url.CreateElement("priority")
		priority.SetText("0.5")
	}
	SetBuilderElementValue(url, su.data, "expires")
	SetBuilderElementValue(url, su.data, "mobile")
	SetBuilderElementValue(url, su.data, "news")
	SetBuilderElementValue(url, su.data, "video")
	SetBuilderElementValue(url, su.data, "image")
	SetBuilderElementValue(url, su.data, "geo")

	if su.Config.Pretty {
		doc.Indent(2)
	}
	buf := poolBuffer.Get()
	_, _ = doc.WriteTo(buf)

	bytes := buf.Bytes()
	poolBuffer.Put(buf)

	return bytes
}
