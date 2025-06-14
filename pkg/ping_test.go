package stm

import "testing"

func TestPingSearchEngines(t *testing.T) {
	PingSearchEngines(NewConfig(), append(DefaultPingLinks, "http://www.example.com")...)

	PingSearchEngines(NewConfig())
}
