package stm

import (
	"log"
	"runtime"
)

// Sitemap provides interface for create sitemap xml file and that has convenient interface.
// And also needs to use first this struct if it wants to use this package.
type Sitemap struct {
	*Config
	builderMain  Builder
	builderIndex Builder
}

// NewSitemap returns the created the Sitemap's pointer
func NewSitemap(c ...*Config) *Sitemap {
	var config *Config
	if len(c) > 0 {
		config = c[0]
	} else {
		config = NewConfig()
	}

	if config.MaxProc > 0 {
		maxProc := config.MaxProc
		log.SetFlags(log.LstdFlags | log.Llongfile)
		if maxProc < 1 || maxProc > runtime.NumCPU() {
			maxProc = runtime.NumCPU()
		}
		log.Printf("Max processors %d\n", maxProc)
		runtime.GOMAXPROCS(maxProc)
	}

	return &Sitemap{Config: config}
}

// Create method must be that calls first this method in that before call to Add method on this struct.
func (sm *Sitemap) Create() *Sitemap {
	sm.builderIndex = NewBuilderIndexFile(sm.Config, sm.Config.GetIndexLocation())
	return sm
}

// Add Should call this after call to Create method on this struct.
func (sm *Sitemap) Add(url any) *Sitemap {
	if sm.builderMain == nil {
		sm.builderMain = NewBuilderFile(sm.Config, sm.Config.GetLocation())
	}

	err := sm.builderMain.Add(url)
	if err != nil {
		if err.FullError() {
			sm.Finalize()
			return sm.Add(url)
		}
		if err.InvalidUrlErr() {
			sm.Config.GetOnInvalidUrl()(err)
		}
	}

	return sm
}

// XMLContent returns the XML content of the sitemap
func (sm *Sitemap) XMLContent() []byte {
	return sm.builderMain.XMLContent()
}

// Finalize writes sitemap and index files if it had some
// specific condition in BuilderFile struct.
func (sm *Sitemap) Finalize() *Sitemap {
	_ = sm.builderIndex.Add(sm.builderMain)
	sm.builderIndex.Write()
	sm.builderMain = nil
	return sm
}

func (sm *Sitemap) PingSearchEngines(urls ...string) { PingSearchEngines(sm.Config, urls...) }

func (sm *Sitemap) GetConfig() *Config {
	return sm.Config
}
