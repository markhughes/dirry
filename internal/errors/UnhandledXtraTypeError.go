package errors

import "fmt"

type UnhandledXtraTypeError struct {
	XtraType string
}

func (e *UnhandledXtraTypeError) Error() string {
	return fmt.Sprintf("unhandled xtra type %s", e.XtraType)
}
