package trash

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
)

var (
	pattern = regexp.MustCompile(`\A(-)?((?P<weeks>[\d\.]+)w)?((?P<days>[\d\.]+)d)?((?P<hours>[\d\.]+)h)?((?P<minutes>[\d\.]+)m)?((?P<seconds>[\d\.]+?)s)?\z`)
)

// ParseDuration parses the duration string to time.Duration, which can be
// treated as a supplement of time.ParseDuration(string).
// Currently, the following format are supported.
//		w - weeks
//		d - days
//		h - hours
//		m - minutes
//		s - seconds
func ParseDuration(s string) (time.Duration, error) {
	var (
		matches []string
		prefix  string
	)
	if pattern.MatchString(s) {
		matches = pattern.FindStringSubmatch(s)
	} else {
		return 0, errors.Errorf("The duration string is invalid.")
	}

	if strings.HasPrefix(s, "-") {
		prefix = "-"
	}
	return durationFromMatchesAndPrefix(matches, prefix)
}

func durationFunc(prefix string) func(string, float64) time.Duration {
	return func(format string, f float64) time.Duration {
		d, err := time.ParseDuration(fmt.Sprintf(prefix+format, f))
		if err != nil {
			panic(fmt.Sprintf("Failed to parse the duration string."))
		}
		return d
	}
}

func durationFromMatchesAndPrefix(matches []string, prefix string) (time.Duration, error) {
	d := time.Duration(0)

	duration := durationFunc(prefix)

	for i, name := range pattern.SubexpNames() {
		value := matches[i]
		if i == 0 || name == "" || value == "" {
			continue
		}

		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return d, err
		}
		switch name {
		case "weeks":
			d += duration("%fh", f*24*7)
		case "days":
			d += duration("%fh", f*24)
		case "hours":
			d += duration("%fh", f)
		case "minutes":
			d += duration("%fm", f)
		case "seconds":
			d += duration("%fs", f)
		}
	}

	return d, nil
}
