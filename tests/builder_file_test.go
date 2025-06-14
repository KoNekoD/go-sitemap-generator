package tests

import (
	"fmt"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"os"
	"testing"
)

func TestBuilderFile(t *testing.T) {
	c := stm.NewConfig().SetCompress(false)

	l := stm.NewLocation(c)

	b := stm.NewBuilderFile(c, l)

	err := b.Add(stm.URL{{"loc", "http://www.example.com"}})
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

	content = b.Content()
	if len(content) != 0 {
		t.Fatalf("Should be empty")
	}

	outFile := "public/sitemaps/sitemap1.xml.gz"
	_ = os.Remove(outFile)
	_ = os.RemoveAll("public")
}

func TestBuilderFileOverflow(t *testing.T) {
	c := stm.NewConfig().SetCompress(false)

	l := stm.NewLocation(c)

	b := stm.NewBuilderFile(c, l)

	for i := range 50001 {
		err := b.Add(stm.URL{{"loc", fmt.Sprintf("http://www.example.com/%d", i)}})
		if err != nil && !(i == 49999 && err.FullError()) {
			t.Fatalf("Failed to add url in BuilderFile: %s", err)
		}
		if err != nil && err.FullError() {
			break
		}
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

	content = b.Content()
	if len(content) != 0 {
		t.Fatalf("Should be empty")
	}

	outFile := "public/sitemaps/sitemap1.xml.gz"
	_ = os.Remove(outFile)
	_ = os.RemoveAll("public")
}

func TestBuilderFileInvalidUrl(t *testing.T) {
	c := stm.NewConfig().SetCompress(false)

	l := stm.NewLocation(c)

	b := stm.NewBuilderFile(c, l)

	err := b.Add(stm.URL{{"aaa", "http://www.example.com/"}})
	if !err.InvalidUrlErr() {
		t.Fatalf("Should be invalid url error")
	}
}
