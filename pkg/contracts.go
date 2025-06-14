package stm

// Adapter provides interface for writes some kind of sitemap.
type Adapter interface {
	Write(loc *Location, data []byte)
	Bytes() [][]byte
}

// Builder provides interface for adds some kind of url sitemap.
type Builder interface {
	XMLContent() []byte
	Content() []byte
	Add(any) BuilderError
	Write()
}

// SitemapURL provides generated xml interface.
type SitemapURL interface {
	XML() []byte
}
