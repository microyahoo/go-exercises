package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	pattern = regexp.MustCompile(`\A(-)?((?P<weeks>[\d\.]+)w)?((?P<days>[\d\.]+)d)?((?P<hours>[\d\.]+)h)?((?P<minutes>[\d\.]+)m)?((?P<seconds>[\d\.]+?)s)?\z`)
)

// matches = [-1.2d0.5s -   1.2d 1.2     0.5s 0.5]
// i = 0, name = , value = -1.2d0.5s
// i = 1, name = , value = -
// i = 2, name = , value =
// i = 3, name = weeks, value =
// i = 4, name = , value = 1.2d
// i = 5, name = days, value = 1.2
// i = 6, name = , value =
// i = 7, name = hours, value =
// i = 8, name = , value =
// i = 9, name = minutes, value =
// i = 10, name = , value = 0.5s
// i = 11, name = seconds, value = 0.5
// The pattern is -1.2d0.5s: -28h48m0.5s

func ParseDuration(s string) (time.Duration, error) {
	var (
		matches []string
		prefix  string
	)
	if pattern.MatchString(s) {
		matches = pattern.FindStringSubmatch(s)
	} else {
		return 0, errors.New(fmt.Sprintf("The duration string is invalid."))
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

// func durationFunc(prefix string) func(string, float64) time.Duration {
// 	return func(format string, f float64) time.Duration {
// 		d, err := time.ParseDuration(fmt.Sprintf(prefix+format, f))
// 		if err != nil {
// 			panic(fmt.Sprintf("Failed to parse the duration string."))
// 		}
// 		return d
// 	}
// }

// func durationFromMatchesAndPrefix(matches []string, prefix string) (time.Duration, error) {
// 	d := time.Duration(0)

// 	duration := durationFunc(prefix)

// 	fmt.Printf("matches = %s\n", matches)
// 	for i, name := range pattern.SubexpNames() {
// 		value := matches[i]
// 		fmt.Printf("i = %d, name = %s, value = %s\n", i, name, value)
// 		if i == 0 || name == "" || value == "" {
// 			continue
// 		}

// 		f, err := strconv.ParseFloat(value, 64)
// 		fmt.Println(err)
// 		if err == nil {
// 			switch name {
// 			case "weeks":
// 				d += duration("%fh", f*24*7)
// 			case "days":
// 				fmt.Printf("%fh\n", f*24)
// 				d += duration("%fh", f*24)
// 			case "hours":
// 				d += duration("%fh", f)
// 			case "minutes":
// 				d += duration("%fm", f)
// 			case "seconds":
// 				d += duration("%fs", f)
// 			}
// 		} else {
// 			fmt.Println("Overflow!!!")
// 			return 0, err
// 		}
// 	}

// 	return d, nil
// }

// func ParseDuration(s string) time.Duration {
// 	d, err := time.ParseDuration(s)
// 	if err != nil {
// 		fmt.Println("The duration string is invalid")
// 		panic("The duration string is invalid")
// 	}
// 	return d
// }

var once sync.Once

func main() {
	content := []byte(`
	# comment line
	option1: value1
	option2: value2

	# another comment line
	option3: value3
`)

	// Regex pattern captures "key: value" pair from the content.
	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\s+(?P<value>\w+)$`)

	// Template to convert "key: value" to "key=value" by
	// referencing the values captured by the regex pattern.
	template := []byte("$key=$value\n")

	result := []byte{}

	// For each match of the regex in the content.
	for _, submatches := range pattern.FindAllSubmatchIndex(content, -1) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		fmt.Println(submatches)
		result = pattern.Expand(result, template, content, submatches)
		fmt.Println(string(result))
	}
	fmt.Println(string(result))
	fmt.Println("---------------1---------------")
	duration, err := ParseDuration("-1w2d3h4m5s")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("The pattern is %s: %v\n", "-1w2d3h4m5s", duration)

	duration, err = ParseDuration("2d5s")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("The pattern is %s: %v\n", "2d5s", duration)

	duration, err = ParseDuration("-1.2d0.5s")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("The pattern is %s: %v\n", "-1.2d0.5s", duration)

	duration, err = ParseDuration("1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111129999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999922222222222222222222222222222222222222222222222222222222211111111111111111111111111111111111111111111.2d0.5s")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("The pattern is %s: %v\n", "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111222222222222222222222222222222222222222222222222222222222211111111111111111111111111111111111111111111.2d0.5s", duration)

	fmt.Println("---------------2---------------")
	duration, err = time.ParseDuration("266666666666666669926443103427127881237749279907337016759964514044272624082952483033132846110731183416873414767488447358124724331574097619765385979383727448392489123198636792964885279971501016289902592.000000h")
	if err == nil {
		fmt.Println("y")
		fmt.Println(duration)
	} else {
		fmt.Println("y")
		fmt.Println(err)
	}

	fmt.Println("---------------3---------------")
}

func test_once() {
	once.Do(func() {
	})
}

// https://github.com/peterhellberg/duration/blob/master/duration.go
