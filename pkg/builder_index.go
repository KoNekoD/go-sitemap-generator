package stm

import "bytes"

type BuilderIndexFile struct {
	*Config
	loc        *Location
	content    []byte
	linkCount  int
	totalCount int
}

func NewBuilderIndexFile(c *Config, loc *Location) *BuilderIndexFile {
	return &BuilderIndexFile{Config: c, loc: loc}
}

// Add method joins old bytes with creates bytes by it calls from Sitemap.Finalize method.
func (b *BuilderIndexFile) Add(link any) BuilderError {
	builderFile := link.(*BuilderFile)
	builderFile.Write()
	b.content = append(b.content, NewSitemapIndexURL(b.Config, URL{{"loc", builderFile.loc.URL()}}).XML()...)
	b.totalCount += builderFile.linkCount
	b.linkCount++
	return nil
}

// Content and BuilderFile.Content are almost the same behavior.
func (b *BuilderIndexFile) Content() []byte { return b.content }

// XMLContent and BuilderFile.XMLContent share almost the same behavior.
func (b *BuilderIndexFile) XMLContent() []byte {
	return append(append(bytes.Join(bytes.Fields(IndexXMLHeader), []byte(" ")), b.Content()...), IndexXMLFooter...)
}

// Write and BuilderFile.Write are almost the same behavior.
func (b *BuilderIndexFile) Write() {
	c := b.XMLContent()

	b.loc.Write(c, b.linkCount)
}
