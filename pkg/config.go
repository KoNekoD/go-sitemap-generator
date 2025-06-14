package stm

import (
	"fmt"
	"log"
)

type Config struct {
	MaxProc int

	// DefaultHost - your website's host name
	DefaultHost string

	// SitemapsHost - remote host where your sitemaps will be hosted
	SitemapsHost string

	// SitemapsPath - sets this to a directory/path if you don't want to upload to the root of your SitemapsHost
	SitemapsPath string

	// PublicPath - directory to write sitemaps to locally
	PublicPath string

	// Filename - output file name
	Filename string

	// Verbose - verbose output to console
	Verbose bool

	// Compress - compress the output file
	Compress bool

	// Pretty - allows pretty formating to the output files
	Pretty bool

	// Adp output file storage. We have S3Adapter and FileAdapter (default: FileAdapter)
	Adp Adapter
	Nmr *Namer
	Loc *Location

	SearchEngines []string

	OnPingStart func(pingUrl string)
	OnPingEnd   func(msg string)

	OnLocationWriteEnd func(msg string)

	OnInvalidUrl func(err error)
}

func NewConfig() *Config {
	return &Config{
		DefaultHost:  "http://www.example.com",
		SitemapsHost: "", // http://s3.amazonaws.com/sitemap-generator/,
		PublicPath:   "public/",
		SitemapsPath: "sitemaps/",
		Filename:     "sitemap",
		Verbose:      true,
		Compress:     true,
		Pretty:       false,
		Adp:          NewFileAdapter(),
	}
}

func (c *Config) SetMaxProc(v int) *Config         { c.MaxProc = v; return c }
func (c *Config) SetDefaultHost(v string) *Config  { c.DefaultHost = v; return c }
func (c *Config) SetSitemapsHost(v string) *Config { c.SitemapsHost = v; return c }
func (c *Config) SetSitemapsPath(v string) *Config { c.SitemapsPath = v; return c }
func (c *Config) SetPublicPath(v string) *Config   { c.PublicPath = v; return c }
func (c *Config) SetFilename(v string) *Config     { c.Filename = v; return c }
func (c *Config) SetVerbose(v bool) *Config        { c.Verbose = v; return c }
func (c *Config) SetCompress(v bool) *Config       { c.Compress = v; return c }
func (c *Config) SetPretty(v bool) *Config         { c.Pretty = v; return c }
func (c *Config) SetAdp(v Adapter) *Config         { c.Adp = v; return c }
func (c *Config) SetNmr(v *Namer) *Config          { c.Nmr = v; return c }
func (c *Config) SetLoc(v *Location) *Config       { c.Loc = v; return c }
func (c *Config) GetLocation() *Location           { return NewLocation(c) }
func (c *Config) Clone() *Config                   { o := *c; return &o }

func (c *Config) GetSitemapsHost() string {
	if c.SitemapsHost != "" {
		return c.SitemapsHost
	}
	return c.DefaultHost
}

func (c *Config) GetIndexLocation() *Location {
	c2 := c.Clone()
	c2.SetNmr(NewNamer(&NOpts{base: c.Filename}))
	return NewLocation(c2)
}

func (c *Config) GetNamer() *Namer {
	if c.Nmr == nil {
		c.Nmr = NewNamer(&NOpts{base: c.Filename, zero: 1, start: 2})
	}

	return c.Nmr
}

func (c *Config) GetOnPingStart() func(pingUrl string) {
	if c.OnPingStart == nil {
		c.OnPingStart = func(pingUrl string) { fmt.Println("Ping now:", pingUrl) }
	}

	return c.OnPingStart
}

func (c *Config) GetOnPingEnd() func(msg string) {
	if c.OnPingEnd == nil {
		c.OnPingEnd = func(msg string) { fmt.Println(msg) }
	}

	return c.OnPingEnd
}

func (c *Config) GetOnLocationWriteEnd() func(msg string) {
	if c.OnLocationWriteEnd == nil {
		c.OnLocationWriteEnd = func(msg string) { fmt.Println(msg) }
	}
	return c.OnLocationWriteEnd
}

func (c *Config) SetOnInvalidUrl(v func(err error)) *Config { c.OnInvalidUrl = v; return c }

func (c *Config) GetOnInvalidUrl() func(err error) {
	if c.OnInvalidUrl == nil {
		c.OnInvalidUrl = func(err error) { log.Fatalf("[F] Sitemap: %s", err) }
	}
	return c.OnInvalidUrl
}
