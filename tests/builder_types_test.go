package tests

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"reflect"
	"testing"
)

func TestURLType(t *testing.T) {
	url := stm.URL{{"loc", "1"}, {"host", "http://example.com"}}
	expect := stm.URL{{"loc", "http://example.com/1"}, {"host", "http://example.com"}}

	url = url.URLJoinBy("loc", "host", "loc")

	if !reflect.DeepEqual(url, expect) {
		t.Fatalf("Failed to join url in URL type: deferrent URL %v and %v", url, expect)
	}
}
