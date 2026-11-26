package gossa

import (
	"fmt"
	"go/token"
)

type CompactError struct {
	Message string
	Pos     token.Pos
}

type Error struct {
	Message  string
	Position token.Position
}

func (e *CompactError) Error() string {
	return fmt.Sprintf("(compact) %d: %s", int(e.Pos), e.Message)
}

func Errorf(format string, p token.Pos, a ...any) *CompactError {
	return &CompactError{
		Message: fmt.Sprintf(format, a...),
		Pos:     p,
	}
}

func (e *CompactError) Expand(f *token.FileSet) *Error {
	return &Error{
		Message:  e.Message,
		Position: f.Position(e.Pos),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Position, e.Message)
}
