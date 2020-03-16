package main

import (
	"fmt"
)

type rpcClientError struct {
	error
	code, message, details string
}

func newRPCClientError(code, message, details string) *rpcClientError {
	return &rpcClientError{
		code:    code,
		message: message,
		details: details,
	}
}

func (err *rpcClientError) Code() string {
	return err.code
}

func (err *rpcClientError) Message() string {
	return err.message
}

func (err *rpcClientError) Error() string {
	return fmt.Sprintf("[%s] %s\nDetails: %s", err.code, err.message, err.details)
}

func (err *rpcClientError) Details() string {
	return err.Error()
}

type rpcError struct {
	*rpcClientError
	others string
}

func (err *rpcError) Error() string {
	if err.others != "" {
		return fmt.Sprintf("Others=%s\n%s", err.others, err.rpcClientError.Error())
	}
	return err.rpcClientError.Error()
}

// func (err *rpcError) Details() string {
// 	return "hello"
// }

func main() {
	err := &rpcError{
		rpcClientError: &rpcClientError{
			code:    "01-001-0001",
			message: "unspecified error",
			details: "initialization error of rbd",
		},
		others: "hello, others",
	}
	fmt.Println(err)
	fmt.Println("---------------")
	fmt.Println(err.Details())
}
