package tests

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"testing"
)

func TestBufferAdapter(t *testing.T) {
	a := stm.NewBufferAdapter()

	a.Write(nil, []byte("Hello world"))

	a.Bytes()
}
