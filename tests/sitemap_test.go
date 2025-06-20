package tests

import (
	"bytes"
	"fmt"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"os"
	"reflect"
	"testing"

	"github.com/clbanning/mxj"
)

func TestSitemap(t *testing.T) {
	t.Run(
		"MaxProc", func(t *testing.T) {
			stm.NewSitemap(stm.Config{MaxProc: 100})
		},
	)

	t.Run(
		"FullOnAddErr", func(t *testing.T) {
			c := stm.NewConfig()
			sm := stm.NewSitemap(*c).Create()

			for i := range 50000 {
				sm.Add(stm.URL{{"loc", fmt.Sprintf("http://www.example.com/%d", i)}})
			}

			_ = os.RemoveAll("public")
		},
	)

	t.Run(
		"InvalidUrl", func(t *testing.T) {
			c := stm.NewConfig().SetOnInvalidUrl(func(err error) {})
			sm := stm.NewSitemap(*c).Create()

			sm.Add(stm.URL{{"aaa", "http://www.example.com"}})
		},
	)

	t.Run(
		"XMLContent", func(t *testing.T) {
			out := stm.NewSitemap(*stm.NewConfig()).Create().Add(
				stm.URL{
					{
						"loc",
						"http://www.example.com",
					},
				},
			).XMLContent()
			if len(out) == 0 {
				t.Fatalf("Should be content")
			}
		},
	)

	t.Run(
		"PingSearchEngines", func(t *testing.T) {
			stm.NewSitemap(*stm.NewConfig()).PingSearchEngines()
		},
	)
}

func TestSitemapGenerator(t *testing.T) {
	buf := stm.BufferAdapter{}

	sm := stm.NewSitemap()
	sm.SetPretty(true)
	sm.SetVerbose(false)
	sm.SetAdp(&buf)

	sm.Create()
	for i := 1; i <= 10; i++ {
		sm.Add(
			stm.URL{
				{"loc", "home"},
				{"changefreq", "always"},
				{"mobile", true},
				{"lastmod", "2018-10-28T17:56:02+09:00"},
			},
		)
		sm.Add(stm.URL{{"loc", "readme"}, {"lastmod", "2018-10-28T17:56:02+09:00"}})
		sm.Add(stm.URL{{"loc", "aboutme"}, {"priority", 0.1}, {"lastmod", "2018-10-28T17:56:02+09:00"}})
	}
	sm.Finalize()

	buffers := buf.Bytes()

	data := buffers[len(buffers)-1]
	expect := []byte(`
	<?xml version="1.0" encoding="UTF-8"?>
	<sitemapindex xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/siteindex.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	  <sitemap>
		<loc>http://www.example.com/sitemaps//sitemap1.xml.gz</loc>
		<lastmod>2018-10-28T17:37:21+09:00</lastmod>
	  </sitemap>
	</sitemapindex>`)

	mdata, _ := mxj.NewMapXml(data)
	mexpect, _ := mxj.NewMapXml(expect)
	mdata.Remove("sitemapindex.sitemap.lastmod")
	mexpect.Remove("sitemapindex.sitemap.lastmod")

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemapindex`)
	}

	bufs := bytes.Buffer{}
	for _, buf := range buffers[:len(buffers)-1] {
		bufs.Write(buf)
	}
	data = bufs.Bytes()
	expect = []byte(`
	<?xml version="1.0" encoding="UTF-8"?> <urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1" xmlns:video="http://www.google.com/schemas/sitemap-video/1.1" xmlns:geo="http://www.google.com/geo/schemas/sitemap/1.0" xmlns:news="http://www.google.com/schemas/sitemap-news/0.9" xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0" xmlns:pagemap="http://www.google.com/schemas/sitemap-pagemap/1.0" xmlns:xhtml="http://www.w3.org/1999/xhtml" ><url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	<url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	</urlset>
	`)

	mdata, _ = mxj.NewMapXml(data)
	mexpect, _ = mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate dataindex`)
	}

}
