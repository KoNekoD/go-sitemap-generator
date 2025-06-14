package stm

import (
	"compress/gzip"
	"log"
	"os"
)

type FileAdapter struct{}

func NewFileAdapter() *FileAdapter {
	return &FileAdapter{}
}

func (a *FileAdapter) Bytes() [][]byte { return nil }

func (a *FileAdapter) Write(loc *Location, data []byte) {
	dir := loc.Directory()
	fi, err := os.Stat(dir)
	if err != nil {
		_ = os.MkdirAll(dir, 0755)
	} else if !fi.IsDir() {
		log.Fatalf("[F] %s should be a directory", dir)
	}

	file, _ := os.OpenFile(loc.Path(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	fi, err = file.Stat()
	if err != nil {
		log.Fatalf("[F] %s file not exists", loc.Path())
	} else if !fi.Mode().IsRegular() {
		log.Fatalf("[F] %s should be a filename", loc.Path())
	}

	if GzipPtn.MatchString(loc.Path()) {
		a.gzip(file, data)
	} else {
		a.plain(file, data)
	}
}

// gzip will create sitemap file as a gzip.
func (a *FileAdapter) gzip(file *os.File, data []byte) {
	gz := gzip.NewWriter(file)
	defer func() { _ = gz.Close() }()
	_, _ = gz.Write(data)
}

// plain will create uncompressed file.
func (a *FileAdapter) plain(file *os.File, data []byte) {
	_, _ = file.Write(data)
	defer func() { _ = file.Close() }()
}
