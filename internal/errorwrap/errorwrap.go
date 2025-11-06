package erwp

import (
	"log"
	"runtime/debug"
)

type Result[T any] struct {
	V   T
	Err error
}

func Try[T any](v T, err error) Result[T] { // pack (v, err) into one value
	return Result[T]{V: v, Err: err}
}

func CleanUp(onErr ...func()) {
	for _, fn := range onErr {
		if fn != nil {
			fn()
		}
	}
}

func MustReturn[T any](r Result[T], onErrs ...func()) T { // unpacks inside
	if r.Err != nil {
		CleanUp(onErrs...)
		debug.PrintStack()
		log.Fatalf("fatal: %v", r.Err)
	}
	return r.V
}

func MustDo(err error, onErr ...func()) {
	if err != nil {
		CleanUp(onErr...)
		debug.PrintStack()
		log.Fatalf("fatal: %v", err)
	}
}

func LetDo(err error, onErr ...func()) {
	if err != nil {
		CleanUp(onErr...)
		debug.PrintStack()
		log.Panicf("panic: %v", err)
	}
}

func LetReturn[T any](r Result[T], onErrs ...func()) T { // unpacks inside
	if r.Err != nil {
		CleanUp(onErrs...)
		debug.PrintStack()
		log.Panicf("panic: %v", r.Err)
	}
	return r.V
}
