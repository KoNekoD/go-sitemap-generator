package main

import (
	"fmt"
	"github.com/KoNekoD/go-sitemap-generator/pkg"
)

func main() {
	sm := stm.NewSitemap()

	bufferAdapter := stm.NewBufferAdapter()

	sm.GetConfig().
		SetAdp(bufferAdapter)

	sm.Create()

	link := "example.com/1/2/3"
	priority := 0.8

	sm.Add(stm.URL{{"loc", link}, {"changefreq", "weekly"}, {"priority", priority}})

	sm.Finalize()

	out := bufferAdapter.Bytes()

	for _, outBytes := range out {
		fmt.Println(string(outBytes))
	}
}
