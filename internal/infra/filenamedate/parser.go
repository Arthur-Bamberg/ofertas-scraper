package filenamedate

import (
	"regexp"
)

var dateRe = regexp.MustCompile(`(20\d{2})[-_]?(\d{2})[-_]?(\d{2})`)

// Parser extracts the first YYYY-MM-DD-like token from a filename.
type Parser struct{}

func (Parser) Parse(filename string) (string, bool) {
	m := dateRe.FindStringSubmatch(filename)
	if m == nil {
		return "", false
	}
	return m[1] + "-" + m[2] + "-" + m[3], true
}
