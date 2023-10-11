package parse

import "fmt"

type (
	Error struct {
		Code   string
		Value  any
		Index  int
		Offset int
		Type   string
	}
)

var _ error = &Error{}

func (e *Error) Error() string {
	return fmt.Sprintf("code=%s, value=%v, index=%d, offset=%d, type=%s", e.Code, e.Value, e.Index, e.Offset, e.Type)
}
