// Package errors stubs the blueprint's internal errors package.
package errors

type Error struct{ Msg string }

func (e *Error) Error() string { return e.Msg }

func New(msg string) *Error { return &Error{Msg: msg} }
