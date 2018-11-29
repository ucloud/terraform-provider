package utils

import (
	"regexp"
)

// Patch is the patch object to provider a converter function
type Patch interface {
	Patch([]byte) []byte
}

// RegexpPatcher a patch object to provider a converter function from regular expression
type RegexpPatcher struct {
	pattern     *regexp.Regexp
	replacement string
}

// NewRegexpPatcher will return a patch object to provider a converter function from regular expression
func NewRegexpPatcher(regex string, repl string) *RegexpPatcher {
	return &RegexpPatcher{
		pattern:     regexp.MustCompile(regex),
		replacement: repl,
	}
}

// Patch will convert a bytes to another bytes with patch rules
func (p *RegexpPatcher) Patch(body []byte) []byte {
	return p.pattern.ReplaceAll(body, []byte(p.replacement))
}

// RetCodePatcher will convert `RetCode` as integer
var RetCodePatcher = NewRegexpPatcher(`"RetCode":\s*"(\d+)"`, `"RetCode": $1`)

// PortPatcher will convert `RetCode` as integer
var PortPatcher = NewRegexpPatcher(`"Port":\s*"(\d+)"`, `"Port": $1`)
