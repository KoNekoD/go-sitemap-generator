package stm

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

// Location provides sitemap's path and filename on file systems, and it provides proxy for Adapter interface also.
type Location struct {
	*Config
	nmr      *Namer
	filename string
}

func NewLocation(c *Config) *Location { return &Location{Config: c} }

func (loc *Location) Namer() *Namer { return loc.Config.GetNamer() }

// IsReservedName confirms that keeps filename on Location.filename.
func (loc *Location) IsReservedName() bool { return loc.filename != "" }

// IsVerbose returns boolean about verbosed summary.
func (loc *Location) IsVerbose() bool { return loc.Config.Verbose }

// Directory indicates where sitemap files are
func (loc *Location) Directory() string {
	return filepath.Join(loc.PublicPath, loc.SitemapsPath)
}

// Path indicates where sitemap name is.
func (loc *Location) Path() string {
	return filepath.Join(loc.PublicPath, loc.SitemapsPath, loc.Filename())
}

// PathInPublic indicates where url file path is.
func (loc *Location) PathInPublic() string {
	return filepath.Join(loc.Config.SitemapsPath, loc.Filename())
}

// URL returns path to combine SitemapsHost, sitemapsPath and Filename on website with it uses ResolveReference.
func (loc *Location) URL() string {
	base, _ := url.Parse(loc.GetSitemapsHost())

	for _, ref := range []string{loc.Config.SitemapsPath + "/", loc.Filename()} {
		base, _ = base.Parse(ref)
	}

	return base.String()
}

// Filesize returns file size this struct has.
func (loc *Location) Filesize() int64 {
	f, _ := os.Open(loc.Path())
	defer func() { _ = f.Close() }()

	fi, err := f.Stat()
	if err != nil {
		return 0
	}

	return fi.Size()
}

// Filename returns sitemap filename.
func (loc *Location) Filename() string {
	nmr := loc.Namer()
	if loc.filename == "" && nmr == nil {
		log.Fatal("[F] No filename or namer set")
	}

	if loc.filename == "" {
		loc.filename = nmr.String()

		if !loc.Config.Compress {
			newName := reGzip.ReplaceAllString(loc.filename, "")
			loc.filename = newName
		}
	}
	return loc.filename
}

// ReserveName returns that sets filename if this struct didn't keep filename,
// and it returns reserved filename if this struct keeps filename also.
func (loc *Location) ReserveName() string {
	nmr := loc.Namer()
	if nmr != nil {
		loc.Filename()
		nmr.Next()
	}

	return loc.filename
}

// Write writes sitemap and index files that used from Adapter interface.
func (loc *Location) Write(data []byte, linkCount int) {
	loc.Config.Adp.Write(loc, data)
	if !loc.IsVerbose() {
		return
	}

	if output := loc.Summary(linkCount); output != "" {
		loc.Config.GetOnLocationWriteEnd()(output)
	}
}

// Summary outputs to generated file summary for console.
func (loc *Location) Summary(linkCount int) string {
	nmr := loc.Namer()
	if nmr.IsStart() {
		return ""
	}

	out := fmt.Sprintf("%s '%d' links", loc.PathInPublic(), linkCount)

	size := loc.Filesize()
	if size <= 0 {
		return out
	}

	return fmt.Sprintf("%s / %d bytes", out, size)
}
