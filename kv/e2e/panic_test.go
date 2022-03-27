package e2e_test

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	// "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/v2"
)

// SkipPanic is the value that will be panicked from Skip.
type SkipPanic struct {
	Message        string // The failure message passed to Fail
	Filename       string // The filename that is the source of the failure
	Line           int    // The line number of the filename that is the source of the failure
	FullStackTrace string // A full stack trace starting at the source of the failure
}

const GINKGO_PANIC = `
Your test failed.
Ginkgo panics to prevent subsequent assertions from running.
Normally Ginkgo rescues this panic so you shouldn't see it.

But, if you make an assertion in a goroutine, Ginkgo can't capture the panic.
To circumvent this, you should call

	defer GinkgoRecover()

at the top of the goroutine that caused this panic.
`

// String makes SkipPanic look like the old Ginkgo panic when printed.
func (SkipPanic) String() string { return GINKGO_PANIC }

// Skip wraps ginkgo.Skip so that it panics with more useful
// information about why the test is being skipped. This function will
// panic with a SkipPanic.
func skip(message string, callerSkip ...int) {
	skip := 1
	if len(callerSkip) > 0 {
		skip += callerSkip[0]
	}

	_, file, line, _ := runtime.Caller(skip)
	sp := SkipPanic{
		Message:        message,
		Filename:       file,
		Line:           line,
		FullStackTrace: pruneStack(skip),
	}

	defer func() {
		e := recover()
		fmt.Printf("************(e != nil): %t, %T\n", e != nil, e)
		if e != nil {
			panic(sp)
		}
	}()

	ginkgo.Skip(message, skip)
}

// ginkgo adds a lot of test running infrastructure to the stack, so
// we filter those out
var stackSkipPattern = regexp.MustCompile(`onsi/ginkgo`)

func pruneStack(skip int) string {
	skip += 2 // one for pruneStack and one for debug.Stack
	stack := debug.Stack()
	scanner := bufio.NewScanner(bytes.NewBuffer(stack))
	var prunedStack []string

	// skip the top of the stack
	for i := 0; i < 2*skip+1; i++ {
		scanner.Scan()
	}

	for scanner.Scan() {
		if stackSkipPattern.Match(scanner.Bytes()) {
			scanner.Scan() // these come in pairs
		} else {
			prunedStack = append(prunedStack, scanner.Text())
			scanner.Scan() // these come in pairs
			prunedStack = append(prunedStack, scanner.Text())
		}
	}

	return strings.Join(prunedStack, "\n")
}

func skipInternalf(caller int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Logf(msg)
	skip(msg, caller+1)
}

func nowStamp() string {
	return time.Now().Format(time.StampMilli)
}

func log(level string, format string, args ...interface{}) {
	fmt.Fprintf(ginkgo.GinkgoWriter, nowStamp()+": "+level+": "+format+"\n", args...)
}

// Logf logs the info.
func Logf(format string, args ...interface{}) {
	log("INFO", format, args...)
}

// SIGDescribe annotates the test with the SIG label.
func SIGDescribe(text string, body func()) bool {
	return ginkgo.Describe("[sig-storage] "+text, body)
}

var _ = SIGDescribe("PersistentVolumes", func() {
	ginkgo.Context("", func() {
		ginkgo.It("", func() {
			gate := "HonorPVReclaimPolicy"
			fmt.Printf("======hello world: %s\n", gate)
			// {
			skipInternalf(1, "Only supported when %v feature is enabled", gate)
			fmt.Println("******hello world")
			// }
			fmt.Println("******hello world again")
			ginkgo.By("delete pv", func() {
				fmt.Println("******hello world again again")
			})
			fmt.Println("******hello world again 2")
		})
		// }, 10)
		fmt.Println("******hello world again again again")
		ginkgo.AfterEach(func() {
			time.Sleep(5 * time.Second)
			fmt.Println("******hello world again again again 2")
		})
	})
})
