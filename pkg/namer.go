package stm

import (
	"fmt"
	"log"
)

func NewNamer(opts *NOpts) *Namer {
	if opts.extension == "" {
		opts.extension = ".xml.gz"
	}

	namer := &Namer{Opts: opts}
	namer.Reset()
	return namer
}

type NOpts struct {
	base      string // filename base
	zero      int
	extension string
	start     int
	buildName func(count int) string
}

func (o *NOpts) SetBuildName(f func(count int) string) *NOpts { o.buildName = f; return o }
func (o *NOpts) GetBuildName() func(count int) string {
	if o.buildName == nil {
		o.buildName = func(count int) string {
			if count == 0 {
				return fmt.Sprintf("%s%s", o.base, o.extension)
			}
			return fmt.Sprintf("%s%d%s", o.base, count, o.extension)
		}
	}
	return o.buildName
}

// Namer provides sitemap's filename as a file number counter.
type Namer struct {
	count int
	Opts  *NOpts
}

// String returns that combines filename base and file extension.
func (n *Namer) String() string { return n.Opts.GetBuildName()(n.count) }

// Reset will initialize to zero value on Namer's counter.
func (n *Namer) Reset() { n.count = n.Opts.zero }

// IsStart confirms that this struct has zero value.
func (n *Namer) IsStart() bool { return n.count == n.Opts.zero }

// Next is going to go to next index for filename.
func (n *Namer) Next() *Namer {
	if n.IsStart() {
		n.count = n.Opts.start
	} else {
		n.count++
	}
	return n
}

// Previous is going to go to previous index for filename.
func (n *Namer) Previous() *Namer {
	if n.IsStart() {
		log.Fatal("[F] Already at the start of the series")
	}
	if n.count <= n.Opts.start {
		n.count = n.Opts.zero
	} else {
		n.count--
	}
	return n
}
