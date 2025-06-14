package stm

import (
	"bytes"
)

type BuilderFile struct {
	*Config
	loc       *Location
	content   []byte
	linkCount int
	newsCount int
}

func NewBuilderFile(c *Config, loc *Location) *BuilderFile {
	return &BuilderFile{Config: c, loc: loc, content: make([]byte, 0, MaxSitemapFilesize)}
}

// Add method joins old bytes with creates bytes by it calls from Sitemap.Add method.
func (b *BuilderFile) Add(url any) BuilderError {
	b.linkCount++
	smu, err := NewSitemapURL(b.Config, MergeMap(url.(URL), URL{{"host", b.loc.Config.DefaultHost}}))
	if err != nil {
		return &builderFileError{error: err, invalidUrl: true}
	}

	xmlBytes := smu.XML()
	if !b.isFileCanFit(xmlBytes) {
		return &builderFileError{error: err, full: true}
	}
	b.content = append(b.content, xmlBytes...)

	return nil
}

// isFileCanFit checks bytes to bigger than const values.
func (b *BuilderFile) isFileCanFit(bytes []byte) bool {
	return len(append(b.content, bytes...)) < MaxSitemapFilesize && b.linkCount < MaxSitemapLinks && b.newsCount < MaxSitemapNews
}

// clear will initialize xml content.
func (b *BuilderFile) clear() {
	b.content = make([]byte, 0, MaxSitemapFilesize)
}

// Content will return pooled bytes on content attribute.
func (b *BuilderFile) Content() []byte {
	return b.content
}

// XMLContent will return an XML of the sitemap built
func (b *BuilderFile) XMLContent() []byte {
	return append(append(bytes.Join(bytes.Fields(XMLHeader), []byte(" ")), b.Content()...), XMLFooter...)
}

// Write will write pooled bytes with header and footer to Location path for output sitemap file.
func (b *BuilderFile) Write() {
	b.loc.ReserveName()
	b.loc.Write(b.XMLContent(), b.linkCount)
	b.clear()
}
