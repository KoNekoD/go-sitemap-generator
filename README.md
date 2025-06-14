# Go sitemap generator

The easiest way to generate sitemaps in Go.

[![GoDoc](https://godoc.org/github.com/KoNekoD/go-sitemap-generator/stm?status.svg)](https://godoc.org/github.com/KoNekoD/go-sitemap-generator/pkg)

### Usage

```go
package main

import (
	"fmt"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
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
}
```

### Features

- [x] Supports: generate kind of some sitemaps.
    - [x] [News sitemaps](#news-sitemaps)
    - [x] [Video sitemaps](#video-sitemaps)
    - [x] [Image sitemaps](#image-sitemaps)
    - [x] [Geo sitemaps](#geo-sitemaps)
    - [x] [Mobile sitemaps](#mobile-sitemaps)
    - [x] Alternate Links
- [x] Supports: adapters for sitemap storage.
    - [x] Filesystem (
      see [NewFileAdapter](https://github.com/KoNekoD/go-sitemap-generator/blob/main/pkg/adapter_file.go))
    - [x] S3 (see [NewS3Adapter](https://github.com/KoNekoD/go-sitemap-generator/blob/main/pkg/adapter_s3.go))
    - [x] MinIO (see [NewMinioAdapter](https://github.com/KoNekoD/go-sitemap-generator/blob/main/pkg/adapter_minio.go))
    - [x] Buffer (for extend,
      see [NewBufferAdapter](https://github.com/KoNekoD/go-sitemap-generator/blob/main/pkg/adapter_buffer.go))
- [x] [Notifies search engines (Google, Bing) of new sitemaps](#pinging-search-engines)
- [x] [Gives complete control over sitemap contents and naming scheme](#full-example)

### Examples

- [x] [Filesystem](https://github.com/KoNekoD/go-sitemap-generator/blob/main/tests/examples/filesystem/main.go)
- [x] [MinIO](https://github.com/KoNekoD/go-sitemap-generator/blob/main/tests/examples/minio/minimal/main.go)
- [x] [MinIO with orphans removal](https://github.com/KoNekoD/go-sitemap-generator/blob/main/tests/examples/minio/full/main.go)

## Getting Started

### Pinging Search Engines

PingSearchEngines notifies search engines of changes once a sitemap
has been generated or changed. The library will append Google and Bing to any engines passed in to the function.

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	// ...

	sm.Finalize().PingSearchEngines()
}
```

If a `new search engine` needs to be added, there is an option to pass it to a function:

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	// ...

	sm.Finalize().PingSearchEngines("http://newengine.com/ping?url=%s")
}
```

### Options

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func main() {
	sm := stm.NewSitemap()

	// Website's host name
	sm.SetDefaultHost("http://www.example.com")

	// The remote host where sitemaps will be hosted
	sm.SetSitemapsHost("http://s3.amazonaws.com/sitemap-generator/")

	// The directory to write sitemaps to locally
	sm.SetPublicPath("tmp/")

    // Set this parameter to the directory/path if the upload is not to the `SitemapsHost` root
	sm.SetSitemapsPath("sitemaps/")

	var c *credentials.Credentials

	// Struct of `S3Adapter`
	sm.SetAdp(stm.NewS3Adapter("ap-northeast-1", "app-bucket", "public-read", c))

	// Change the output filename
	sm.SetFilename("new_filename")
}
```

### News sitemaps

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	sm.Add(stm.URL{
		{"loc", "/news"},
		{"news", stm.URL{
			{"publication", stm.URL{
				{"name", "Example"},
				{"language", "en"},
			},
			},
			{"title", "Test article"},
			{"keywords", "article, articles about testing"},
			{"stock_tickers", "SAO:PETR3"},
			{"publication_date", "2011-08-22"},
			{"access", "Subscription"},
			{"genres", "PressRelease"},
		}}},
	)
}
```

Look at [Creating a Google News Sitemap](https://support.google.com/news/publisher/answer/74288) as required.

### Video sitemaps

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	sm.Add(stm.URL{
		{"loc", "/videos"},
		{"video", stm.URL{
			{"thumbnail_loc", "http://www.example.com/video1_thumbnail.png"},
			{"title", "Title"},
			{"description", "Description"},
			{"content_loc", "http://www.example.com/cool_video.mpg"},
			{"category", "Category"},
			{"tag", []string{"one", "two", "three"}},
			{"player_loc", stm.Attrs{"https://example.com/p/flash/moogaloop/6.2.9/moogaloop.swf?clip_id=26", map[string]string{"allow_embed": "Yes", "autoplay": "autoplay=1"}},},
		},
		},
	})
}
```

Look at [Video sitemaps](https://support.google.com/webmasters/answer/80471) as required.

### Image sitemaps

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	sm.Add(stm.URL{
		{"loc", "/images"},
		{"image", []stm.URL{
			{{"loc", "http://www.example.com/image.png"}, {"title", "Image"}},
			{{"loc", "http://www.example.com/image1.png"}, {"title", "Image1"}},
		}},
	})
}
```

Look at [Image sitemaps](https://support.google.com/webmasters/answer/178636) as required.

### Geo sitemaps

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	sm.Add(stm.URL{
		{"loc", "/geos"},
		{"geo", stm.URL{
			{"format", "kml"},
		}},
	})
}
```

Couldn't find Geo sitemaps example, although it's similar to:

```xml

<url>
    <loc>/geos</loc>
    <geo:geo>
        <geo:format>kml</geo:format>
    </geo:geo>
</url>
```

### Mobile sitemaps

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	sm.Add(stm.URL{{"loc", "mobiles"}, {"mobile", true}})
}
```

Look at [Feature phone sitemaps](https://support.google.com/webmasters/answer/6082207) as required.

### Full example

```go
package main

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()
	sm.SetDefaultHost("http://example.com")
	sm.SetSitemapsHost("http://s3.amazonaws.com/sitemaps/")
	sm.SetSitemapsPath("sitemaps/")
	sm.SetFilename("anothername")
	sm.SetCompress(true)
	sm.SetVerbose(true)
	sm.SetAdp(stm.NewS3Adapter("ap-northeast-1", "app-bucket", "", nil))

	sm.Create()

	sm.Add(stm.URL{{"loc", "/home"}, {"changefreq", "daily"}})

	sm.Add(stm.URL{{"loc", "/abouts"}, {"mobile", true}})

	sm.Add(stm.URL{{"loc", "/news"},
		{"news", stm.URL{
			{"publication", stm.URL{
				{"name", "Example"},
				{"language", "en"},
			},
			},
			{"title", "Test article"},
			{"keywords", "article, articles about testing"},
			{"stock_tickers", "SAO:PETR3"},
			{"publication_date", "2011-08-22"},
			{"access", "Subscription"},
			{"genres", "PressRelease"},
		},},
	})

	sm.Add(stm.URL{{"loc", "/images"},
		{"image", []stm.URL{
			{{"loc", "http://www.example.com/image.png"}, {"title", "Image"}},
			{{"loc", "http://www.example.com/image1.png"}, {"title", "Image1"}},
		},},
	})

	sm.Add(stm.URL{{"loc", "/videos"},
		{"video", stm.URL{
			{"thumbnail_loc", "http://www.example.com/video1_thumbnail.png"},
			{"title", "Title"},
			{"description", "Description"},
			{"content_loc", "http://www.example.com/cool_video.mpg"},
			{"category", "Category"},
			{"tag", []string{"one", "two", "three"}},
			{"player_loc", stm.Attrs{"https://example.com/p/flash/moogaloop/6.2.9/moogaloop.swf?clip_id=26", map[string]string{"allow_embed": "Yes", "autoplay": "autoplay=1"}}},
		},},
	})

	sm.Add(stm.URL{{"loc", "/geos"},
		{"geo", stm.URL{
			{"format", "kml"},
		}},
	})

	sm.Finalize().PingSearchEngines("http://newengine.com/ping?url=%s")
}
```

### Webserver example

```go
package main

import (
	"log"
	"net/http"

	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func buildSitemap() *stm.Sitemap {
	sm := stm.NewSitemap()
	sm.SetDefaultHost("http://example.com")

	sm.Create()
	sm.Add(stm.URL{{"loc", "/"}, {"changefreq", "daily"}})

	// Note: Do not call `sm.Finalize()` because it flushes
	// the underlying data structure from memory to disk.

	return sm
}

func main() {
	sm := buildSitemap()

	mux := http.NewServeMux()
	mux.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		// Go's webserver automatically sets the correct `Content-Type` header.
		w.Write(sm.XMLContent())
		return
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Documentation

- [API Reference](https://godoc.org/github.com/KoNekoD/go-sitemap-generator/pkg)

The library is organized as follows: a Sitemap structure is created via NewSitemap, then the user can configure the 
behavior via sm.GetConfig, then the user must call sm.Create to create a builder index into which the previously 
configured sm.GetConfig will be thrown, then the only thing left to do is to call add the page to the sitemap 
index via sm. 
Add by throwing stm.URL into it (example `sm.Add(stm.URL{{"loc", link}, {"changefreq", "weekly"}, {"priority", priority}})`), 
the library itself will slice these links by 50K pages, according to the requirements of Google, after adding the 
required number of pages it is necessary to call `sm.Finalize()` to trigger the writing of links through the adapter.

### How to test.

```shell
make run_tests

# If coverage needed
make run_coverage_test
```
