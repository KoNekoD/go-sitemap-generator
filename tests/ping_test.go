package tests

import (
	"github.com/KoNekoD/go-sitemap-generator/pkg"
	"testing"
)

func TestPingSearchEngines(t *testing.T) {
	stm.PingSearchEngines(stm.NewConfig(), append(stm.DefaultPingLinks, "http://www.example.com")...)

	stm.PingSearchEngines(stm.NewConfig())
}
