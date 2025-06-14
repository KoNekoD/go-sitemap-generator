package tests

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"reflect"
	"testing"
)

func TestMergeMap(t *testing.T) {
	var src, dst, expect [][]any
	src = [][]any{{"loc", "1"}, {"changefreq", "2"}, {"mobile", true}, {"host", "http://google.com"}}
	dst = [][]any{{"host", "http://example.com"}}
	expect = [][]any{{"loc", "1"}, {"changefreq", "2"}, {"mobile", true}, {"host", "http://google.com"}}

	src = stm.MergeMap(src, dst)

	if !reflect.DeepEqual(src, expect) {
		t.Fatalf("Failed to maps merge: deferrent map \n%#v\n and \n%#v\n", src, expect)
	}
}
