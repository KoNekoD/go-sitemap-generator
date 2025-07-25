package tests

import (
	"bytes"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"github.com/beevik/etree"
	"github.com/clbanning/mxj"
	"reflect"
	"testing"
	"time"
)

func TestBlank(t *testing.T) {
	if _, err := stm.NewSitemapURL(&stm.Config{}, stm.URL{}); err == nil {
		t.Errorf(`Failed to validate blank arg ( URL{} ): %v`, err)
	}
}

func TestItHasLocElement(t *testing.T) {
	if _, err := stm.NewSitemapURL(&stm.Config{}, stm.URL{}); err == nil {
		t.Errorf(`Failed to validate about must have loc attribute in URL type ( URL{} ): %v`, err)
	}
}

func TestJustSetLocElement(t *testing.T) {
	smu, err := stm.NewSitemapURL(&stm.Config{}, stm.URL{{"loc", "path"}, {"host", "http://example.com"}})

	if err != nil {
		t.Fatalf(`Fatal to validate! This is a critical error: %v`, err)
	}

	doc := etree.NewDocument()
	_ = doc.ReadFromBytes(smu.XML())

	var elm *etree.Element
	url := doc.SelectElement("url")

	elm = url.SelectElement("loc")
	if elm == nil {
		t.Errorf(`Failed to generate xml that loc element is blank: %v`, elm)
	}
	if elm != nil && elm.Text() != "http://example.com/path" {
		t.Errorf(`Failed to generate xml thats deferrent value in loc element: %v`, elm.Text())
	}
}

func TestJustSetLocElementAndThenItNeedsCompleteValues(t *testing.T) {
	smu, err := stm.NewSitemapURL(&stm.Config{}, stm.URL{{"loc", "path"}, {"host", "http://example.com"}})

	if err != nil {
		t.Fatalf(`Fatal to validate! This is a critical error: %v`, err)
	}

	doc := etree.NewDocument()
	_ = doc.ReadFromBytes(smu.XML())

	var elm *etree.Element
	url := doc.SelectElement("url")

	elm = url.SelectElement("loc")
	if elm == nil {
		t.Errorf(`Failed to generate xml that loc element is blank: %v`, elm)
	}
	if elm != nil && elm.Text() != "http://example.com/path" {
		t.Errorf(`Failed to generate xml thats deferrent value in loc element: %v`, elm.Text())
	}

	elm = url.SelectElement("priority")
	if elm == nil {
		t.Errorf(`Failed to generate xml that priority element is nil: %v`, elm)
	}
	if elm != nil && elm.Text() != "0.5" {
		t.Errorf(`Failed to generate xml thats deferrent value in priority element: %v`, elm.Text())
	}

	elm = url.SelectElement("changefreq")
	if elm == nil {
		t.Errorf(`Failed to generate xml that changefreq element is nil: %v`, elm)
	}
	if elm != nil && elm.Text() != "weekly" {
		t.Errorf(`Failed to generate xml thats deferrent value in changefreq element: %v`, elm.Text())
	}

	elm = url.SelectElement("lastmod")
	if elm == nil {
		t.Errorf(`Failed to generate xml that lastmod element is nil: %v`, elm)
	}
	if elm != nil {
		if _, err := time.Parse(time.RFC3339, elm.Text()); err != nil {
			t.Errorf(`Failed to generate xml thats failed to parse datetime in lastmod element: %v`, err)
		}
	}
}

func TestSetNilValue(t *testing.T) {
	smu, err := stm.NewSitemapURL(
		&stm.Config{},
		stm.URL{
			{"loc", "path"},
			{"priority", nil},
			{"changefreq", nil},
			{"lastmod", nil},
			{"host", "http://example.com"},
		},
	)

	if err != nil {
		t.Fatalf(`Fatal to validate! This is a critical error: %v`, err)
	}

	doc := etree.NewDocument()
	_ = doc.ReadFromBytes(smu.XML())

	var elm *etree.Element
	url := doc.SelectElement("url")

	elm = url.SelectElement("loc")
	if elm == nil {
		t.Errorf(`Failed to generate xml that loc element is blank: %v`, elm)
	}
	if elm != nil && elm.Text() != "http://example.com/path" {
		t.Errorf(`Failed to generate xml thats deferrent value in loc element: %v`, elm.Text())
	}

	elm = url.SelectElement("priority")
	if elm != nil {
		t.Errorf(`Failed to generate xml that priority element must be nil: %v`, elm)
	}

	elm = url.SelectElement("changefreq")
	if elm != nil {
		t.Errorf(`Failed to generate xml that changefreq element must be nil: %v`, elm)
	}

	elm = url.SelectElement("lastmod")
	if elm != nil {
		t.Errorf(`Failed to generate xml that lastmod element must be nil: %v`, elm)
	}
}

func TestAutoGenerateSitemapHost(t *testing.T) {
	smu, err := stm.NewSitemapURL(&stm.Config{}, stm.URL{{"loc", "path"}, {"host", "http://example.com"}})

	if err != nil {
		t.Fatalf(`Fatal to validate! This is a critical error: %v`, err)
	}

	doc := etree.NewDocument()
	_ = doc.ReadFromBytes(smu.XML())

	var elm *etree.Element
	url := doc.SelectElement("url")

	elm = url.SelectElement("loc")
	if elm == nil {
		t.Errorf(`Failed to generate xml that loc element is blank: %v`, elm)
	}
	if elm != nil && elm.Text() != "http://example.com/path" {
		t.Errorf(`Failed to generate xml thats deferrent value in loc element: %v`, elm.Text())
	}
}

func TestNewsSitemaps(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	data := stm.URL{
		{"loc", "/news"}, {
			"news", stm.URL{
				{
					"publication", stm.URL{
						{"name", "Example"},
						{"language", "en"},
					},
				},
				{"title", "My Article"},
				{"keywords", "my article, articles about myself"},
				{"stock_tickers", "SAO:PETR3"},
				{"publication_date", "2011-08-22"},
				{"access", "Subscription"},
				{"genres", "PressRelease"},
			},
		},
	}
	expect := []byte(`
	<root>
		<news:news>
			<news:keywords>my article, articles about myself</news:keywords>
			<news:stock_tickers>SAO:PETR3</news:stock_tickers>
			<news:publication_date>2011-08-22</news:publication_date>
			<news:access>Subscription</news:access>
			<news:genres>PressRelease</news:genres>
			<news:publication>
				<news:name>Example</news:name>
				<news:language>en</news:language>
			</news:publication>
			<news:title>My Article</news:title>
		</news:news>
	</root>`)

	stm.SetBuilderElementValue(root, data, "news")
	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func TestImageSitemaps(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	data := stm.URL{
		{"loc", "/images"}, {
			"image", []stm.URL{
				{{"loc", "http://www.example.com/image.png"}, {"title", "Image"}},
				{{"loc", "http://www.example.com/image1.png"}, {"title", "Image1"}},
			},
		},
	}
	expect := []byte(`
	<root>
		<image:image>
			<image:loc>http://www.example.com/image.png</image:loc>
			<image:title>Image</image:title>
		</image:image>
		<image:image>
			<image:loc>http://www.example.com/image1.png</image:loc>
			<image:title>Image1</image:title>
		</image:image>
	</root>`)

	stm.SetBuilderElementValue(root, data, "image")
	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func TestVideoSitemaps(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	data := stm.URL{
		{"loc", "/videos"}, {
			"video", stm.URL{
				{"thumbnail_loc", "http://www.example.com/video1_thumbnail.png"},
				{"title", "Title"},
				{"description", "Description"},
				{"content_loc", "http://www.example.com/cool_video.mpg"},
				{"category", "Category"},
				{"tag", []string{"one", "two", "three"}},
			},
		},
	}

	expect := []byte(`
	<root>
		<video:video>
			<video:thumbnail_loc>http://www.example.com/video1_thumbnail.png</video:thumbnail_loc>
			<video:title>Title</video:title>
			<video:description>Description</video:description>
			<video:content_loc>http://www.example.com/cool_video.mpg</video:content_loc>
			<video:tag>one</video:tag>
			<video:tag>two</video:tag>
			<video:tag>three</video:tag>
			<video:category>Category</video:category>
		</video:video>
	</root>`)

	stm.SetBuilderElementValue(root, data, "video")
	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func TestGeoSitemaps(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	data := stm.URL{{"loc", "/geos"}, {"geo", stm.URL{{"format", "kml"}}}}

	expect := []byte(`
	<root>
		<geo:geo>
			<geo:format>kml</geo:format>
		</geo:geo>
	</root>`)

	stm.SetBuilderElementValue(root, data, "geo")
	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func TestMobileSitemaps(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	data := stm.URL{{"loc", "/mobile"}, {"mobile", true}}

	expect := []byte(`
	<root>
	  <loc>/mobile</loc>
	  <mobile:mobile/>
	</root>`)

	stm.SetBuilderElementValue(root, data.URLJoinBy("loc", "host", "loc"), "loc")
	stm.SetBuilderElementValue(root, data, "mobile")

	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func TestAlternateLinks(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	loc := "/alternates"
	data := stm.URL{
		{"loc", loc}, {
			"xhtml:link", []stm.Attr{
				{
					"rel":      "alternate",
					"hreflang": "zh-tw",
					"href":     loc + "?locale=zh-tw",
				},
				{
					"rel":      "alternate",
					"hreflang": "en-us",
					"href":     loc + "?locale=en-us",
				},
			},
		},
	}

	expect := []byte(`
       <root>
         <loc>/alternates</loc>
         <xhtml:link rel="alternate" hreflang="zh-tw" href="/alternates?locale=zh-tw"/>
         <xhtml:link rel="alternate" hreflang="en-us" href="/alternates?locale=en-us"/>
       </root>`)

	stm.SetBuilderElementValue(root, data.URLJoinBy("loc", "host", "loc"), "loc")
	stm.SetBuilderElementValue(root, data, "xhtml:link")

	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func TestAttr(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	data := stm.URL{
		{"loc", "/videos"}, {
			"video", stm.URL{
				{"thumbnail_loc", "http://www.example.com/video1_thumbnail.png"},
				{"title", "Title"},
				{"description", "Description"},
				{"content_loc", "http://www.example.com/cool_video.mpg"},
				{"category", "Category"},
				{"tag", []string{"one", "two", "three"}},
				{
					"player_loc",
					stm.Attrs{
						"https://f.vimeocdn.com/p/flash/moogaloop/6.2.9/moogaloop.swf?clip_id=26",
						stm.Attr{"allow_embed": "Yes", "autoplay": "autoplay=1"},
					},
				},
			},
		},
	}

	expect := []byte(`
	<root>
		<video:video>
			<video:thumbnail_loc>http://www.example.com/video1_thumbnail.png</video:thumbnail_loc>
			<video:title>Title</video:title>
			<video:description>Description</video:description>
			<video:content_loc>http://www.example.com/cool_video.mpg</video:content_loc>
			<video:tag>one</video:tag>
			<video:tag>two</video:tag>
			<video:tag>three</video:tag>
			<video:category>Category</video:category>
			<video:player_loc allow_embed="Yes" autoplay="autoplay=1">https://f.vimeocdn.com/p/flash/moogaloop/6.2.9/moogaloop.swf?clip_id=26</video:player_loc>
		</video:video>
	</root>`)

	stm.SetBuilderElementValue(root, data, "video")

	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func TestAttrWithoutTypedef(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")

	data := stm.URL{
		{"loc", "/videos"}, {
			"video", stm.URL{
				{"thumbnail_loc", "http://www.example.com/video1_thumbnail.png"},
				{"title", "Title"},
				{"description", "Description"},
				{"content_loc", "http://www.example.com/cool_video.mpg"},
				{"category", "Category"},
				{"tag", []string{"one", "two", "three"}},
				{
					"player_loc",
					stm.Attrs{
						"https://f.vimeocdn.com/p/flash/moogaloop/6.2.9/moogaloop.swf?clip_id=26",
						map[string]string{"allow_embed": "Yes", "autoplay": "autoplay=1"},
					},
				},
			},
		},
	}

	expect := []byte(`
	<root>
		<video:video>
			<video:thumbnail_loc>http://www.example.com/video1_thumbnail.png</video:thumbnail_loc>
			<video:title>Title</video:title>
			<video:description>Description</video:description>
			<video:content_loc>http://www.example.com/cool_video.mpg</video:content_loc>
			<video:tag>one</video:tag>
			<video:tag>two</video:tag>
			<video:tag>three</video:tag>
			<video:category>Category</video:category>
			<video:player_loc allow_embed="Yes" autoplay="autoplay=1">https://f.vimeocdn.com/p/flash/moogaloop/6.2.9/moogaloop.swf?clip_id=26</video:player_loc>
		</video:video>
	</root>`)

	stm.SetBuilderElementValue(root, data, "video")

	buf := &bytes.Buffer{}
	_, _ = doc.WriteTo(buf)

	mdata, _ := mxj.NewMapXml(buf.Bytes())
	mexpect, _ := mxj.NewMapXml(expect)

	if !reflect.DeepEqual(mdata, mexpect) {
		t.Error(`Failed to generate sitemap xml thats deferrent output value in URL type`)
	}
}

func BenchmarkGenerateXML(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	forPerformance := 500

	for k := 0; k <= forPerformance; k++ {
		for i := 1; i <= forPerformance; i++ {

			var smu stm.SitemapURL
			var data stm.URL

			data = stm.URL{{"loc", "/mobile"}, {"mobile", true}}
			smu, _ = stm.NewSitemapURL(&stm.Config{}, data)
			smu.XML()

			data = stm.URL{
				{"loc", "/images"}, {
					"image", []stm.URL{
						{{"loc", "http://www.example.com/image.png"}, {"title", "Image"}},
						{{"loc", "http://www.example.com/image1.png"}, {"title", "Image1"}},
					},
				},
			}
			smu, _ = stm.NewSitemapURL(&stm.Config{}, data)
			smu.XML()

			data = stm.URL{
				{"loc", "/videos"}, {
					"video", stm.URL{
						{"thumbnail_loc", "http://www.example.com/video1_thumbnail.png"},
						{"title", "Title"},
						{"description", "Description"},
						{"content_loc", "http://www.example.com/cool_video.mpg"},
						{"category", "Category"},
						{"tag", []string{"one", "two", "three"}},
					},
				},
			}
			smu, _ = stm.NewSitemapURL(&stm.Config{}, data)
			smu.XML()

			data = stm.URL{
				{"loc", "/news"}, {
					"news", stm.URL{
						{
							"publication", stm.URL{
								{"name", "Example"},
								{"language", "en"},
							},
						},
						{"title", "My Article"},
						{"keywords", "my article, articles about myself"},
						{"stock_tickers", "SAO:PETR3"},
						{"publication_date", "2011-08-22"},
						{"access", "Subscription"},
						{"genres", "PressRelease"},
					},
				},
			}

			smu, _ = stm.NewSitemapURL(&stm.Config{}, data)
			smu.XML()
		}
	}
}
