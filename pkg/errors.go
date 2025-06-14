package stm

// BuilderError provides interface for it can confirm the error in some difference.
type BuilderError interface {
	error
	FullError() bool
	InvalidUrlErr() bool
}

type builderFileError struct {
	error
	full       bool
	invalidUrl bool
}

// FullError returns true if a sitemap xml had been limit size.
func (e *builderFileError) FullError() bool { return e.full }

// InvalidUrlErr returns true if sitemap tried to create with invalid URL
func (e *builderFileError) InvalidUrlErr() bool { return e.invalidUrl }
