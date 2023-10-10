package errors

import "fmt"

type UnhandledCastTypeError struct {
	CastType     int32
	CastTypeName string
}

func (e *UnhandledCastTypeError) Error() string {
	return fmt.Sprintf("unhandled cast type %d %s", e.CastType, e.CastTypeName)
}
