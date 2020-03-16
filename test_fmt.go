package main

import (
	"fmt"
)

func main() {
	x := "The %s operation failed due to %s with status code %d"
	fmt.Println(newTaskVolumeError(x, "migration", "volume1", 1))
	fmt.Println(task(x, "migration")("volume1", 1))
	fmt.Println(fmt.Sprintf("ABC", nil))
	fmt.Println(0x12)
	fmt.Println(0x0E)

	var at int64
	at = 64

	t := &atTest{at: &at}
	fmt.Println(t)
	fmt.Println((*atTest)(nil))
	fmt.Println(&atTest{at: nil})
}

type atTest struct {
	at *int64
}

func (a *atTest) String() string {
	if a != nil && a.at != nil {
		return fmt.Sprintf("AT: %d", *a.at)
	}
	return "unknown"
}

func newTaskVolumeError(desc, op string, msgArgs ...interface{}) string {
	fmt.Println(msgArgs)
	// return fmt.Sprintf(desc, op, msgArgs...)
	return fmt.Sprintf(desc, []interface{}{op, msgArgs}...)
}

// APIURL returns api url generator with version
func APIURL(op string) func(string) func(string, ...interface{}) string {
	return func(resource string) func(string, ...interface{}) string {
		return func(format string, args ...interface{}) string {
			return fmt.Sprintf(fmt.Sprintf("/%s%s%s", op, resource, format), args...)
		}
	}
}

func task(msg, op string) func(string, ...interface{}) string {
	return func(resource string, args ...interface{}) string {
		xx := fmt.Sprintf(msg, op, resource)
		fmt.Println(xx)
		return fmt.Sprintf(xx, args...)
	}
}
