package stm

import (
	"os"
	"testing"
)

func TestBuilderIndexFile(t *testing.T) {
	c := NewConfig().SetCompress(false)

	l := NewLocation(c)

	b := NewBuilderIndexFile(c, l)

	err := b.Add(&BuilderFile{Config: c, loc: l})
	if err != nil {
		t.Fatalf("Failed to add url in BuilderFile: %s", err)
	}

	content := b.Content()
	if len(content) == 0 {
		t.Fatalf("Should be content")
	}

	xmlContent := b.XMLContent()
	if len(xmlContent) == 0 {
		t.Fatalf("Should be xml content")
	}

	b.Write()

	outFile := "public/sitemaps/sitemap1.xml"

	stat, statErr := os.Stat(outFile)
	if statErr != nil {
		t.Fatalf("Failed to stat file: %s", statErr)
	}
	if stat.Size() == 0 {
		t.Fatalf("Should be content")
	}

	_ = os.Remove(outFile)
	_ = os.RemoveAll("public")
}
