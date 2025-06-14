package stm

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestFileAdapter(t *testing.T) {
	l := NewLocation(NewConfig().SetCompress(false))

	a := NewFileAdapter()

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
	l := NewLocation(NewConfig())

	a := NewFileAdapter()

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
