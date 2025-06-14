package tests

import (
	"bytes"
	"compress/gzip"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestFileAdapter(t *testing.T) {
	l := stm.NewLocation(stm.NewConfig().SetCompress(false))

	a := stm.NewFileAdapter()

	a.Write(l, []byte("Hello world"))

	a.Bytes()

	outFile := "public/sitemaps/sitemap1.xml"

	outFileContent, _ := os.ReadFile(outFile)

	_ = os.Remove(outFile)
	_ = os.RemoveAll("public")

	if !reflect.DeepEqual(outFileContent, []byte("Hello world")) {
		t.Fatalf("unexcepted file content, got %s", outFileContent)
	}
}

func TestFileAdapterGzip(t *testing.T) {
	l := stm.NewLocation(stm.NewConfig())

	a := stm.NewFileAdapter()

	a.Write(l, []byte("Hello world"))

	a.Bytes()

	outFile := "public/sitemaps/sitemap1.xml.gz"

	outFileContent, _ := os.ReadFile(outFile)

	_ = os.Remove(outFile)
	_ = os.RemoveAll("public")

	g, _ := gzip.NewReader(bytes.NewReader(outFileContent))
	outFileContentDecoded, _ := io.ReadAll(g)

	if !reflect.DeepEqual(outFileContentDecoded, []byte("Hello world")) {
		t.Fatalf("unexcepted file content, got %b", outFileContent)
	}
}
