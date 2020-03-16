package main

import (
	"fmt"

	"reflect"
	"xsky-demon/errors"
)

// type Op uint32

var rbdErrorCodePrefix = "02-01-01"

// const
const (
	Create = iota + 1
	// Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod

	Test = 14
)

func main() {
	err := errors.Errorf("xxx")
	fmt.Println(err)
	fmt.Println(Chmod)
	ops := Create | Chmod
	fmt.Println(ops ^ Chmod)
	fmt.Printf("%T\n", Remove)
	fmt.Printf("%s-%04X\n", rbdErrorCodePrefix, Remove)
	fmt.Printf("%s-%04X\n", rbdErrorCodePrefix, Test)
	fmt.Println(test("hello"))
	fmt.Println(test("hello %s %s %d", "world", "word", 1))
	fmt.Println(test("hello %s %s %d", []interface{}{"world", "word", 1}...))

	volErr := newVolumeError(CodeVolumeResourceExistError, nil, "pool1", "volume1", "volume")
	fmt.Printf("*****%T\n", volErr)
	_, ok := volErr.(errors.Coded)
	fmt.Printf("%+v, %v\n", volErr.(errors.Coded), ok)
	volErrVal := reflect.ValueOf(volErr)
	for i := 0; i < volErrVal.NumMethod(); i++ {
		fmt.Printf(">> method=%v\n", volErrVal.Method(i))
		m := volErrVal.MethodByName("Codes")
		if !m.IsValid() {
			fmt.Println("not found")
		} else {
			fmt.Println(m.Call(nil))
			fmt.Println(m.Call(nil)[0].String())
		}
		// if m, ok := volErrVal.MethodByName("Code"); ok {
		// 	fmt.Printf("**Code() = %v\n", m)
		// 	fmt.Println(m.Call())
		// } else {
		// 	fmt.Println("not found")
		// }
	}
}

func test(cmd string, args ...interface{}) string {
	return fmt.Sprintf(cmd, args...)
}

// TaskVolumeErrCodeToMessage is map of error code to their messages
var TaskVolumeErrCodeToMessage = map[int]string{
	CodeVolumeResourceExistError:             "The resource %s already exists",
	CodeVolumeProcessStoppedOrRestartedError: "The %s service is stopped or restarted",
	CodeVolumeInvalidMigrationStatus:         "Invalid volume migration status code %d",
	CodeVolumeGetMigrationStatusError:        "Failed to get volume %s migration status",
	CodeVolumeMigrationError:                 "Volume %s migration failed with status code %d",
}

// list of error codes of volume related tasks
const (
	CodeVolumeResourceExistError = iota + 1
	CodeVolumeProcessStoppedOrRestartedError
	CodeVolumeInvalidMigrationStatus
	CodeVolumeGetMigrationStatusError
	CodeVolumeMigrationError
)

// ErrCodeTaskVolumePrefix is the prefix of task volume error message
var ErrCodeTaskVolumePrefix = "01-02"

type volumeError struct {
	*errors.Err
	pool, volume string
}

func (err *volumeError) Error() string {
	return fmt.Sprintf("Pool: %s, volume: %s\n%s", err.pool, err.volume, err.Err.Error())
}

// newVolumeError returns a task volume error
func newVolumeError(code int, origErr error, pool, volume string, msgArgs ...interface{}) error {
	return &volumeError{
		Err: errors.NewCodedError(1,
			errors.CatenateErrorPrefixCode(ErrCodeTaskVolumePrefix, code),
			origErr,
			TaskVolumeErrCodeToMessage[code],
			msgArgs...).(*errors.Err),
		pool:   pool,
		volume: volume,
	}
}
