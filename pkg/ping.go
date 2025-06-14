package stm

import (
	"fmt"
	"net/http"
	"time"
)

// PingSearchEngines requests some ping server from it calls Sitemap.PingSearchEngines.
func PingSearchEngines(c *Config, urls ...string) {
	if urls == nil {
		urls = DefaultPingLinks
	}
	urls = append(urls, c.SearchEngines...)
	totalEngines := len(urls)
	does := make(chan string, totalEngines)
	locationUrl := c.GetIndexLocation().URL()
	client := http.Client{Timeout: 5 * time.Second}

	for _, url := range urls {
		pingUrl := fmt.Sprintf(url, locationUrl)
		c.GetOnPingStart()(pingUrl)
		go func(pingUrl string) {
			resp, err := client.Get(pingUrl)
			if err != nil {
				does <- fmt.Sprintf("[E] Ping failed: %s (URL:%s)", err, pingUrl)
				return
			}
			defer func() { _ = resp.Body.Close() }()

			does <- fmt.Sprintf("Successful ping of `%s`", pingUrl)
		}(pingUrl)
	}

	for i := 0; i < totalEngines; i++ {
		c.GetOnPingEnd()(<-does)
	}
}
