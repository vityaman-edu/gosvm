package gossa

import (
	"fmt"
	"go/token"
)

type CompactError struct {
	Pos     token.Pos
	Message string
}

type Error struct {
	Position token.Position
	Message  string
}

func (e *CompactError) Error() string {
	return fmt.Sprintf("(compact) %d: %s", int(e.Pos), e.Message)
}

func Errorf(format string, p token.Pos, a ...any) *CompactError {
	return &CompactError{
		Pos:     p,
		Message: fmt.Sprintf(format, a...),
	}
}

func (e *CompactError) Expand(f *token.FileSet) *Error {
	return &Error{
		Position: f.Position(e.Pos),
		Message:  e.Message,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Position, e.Message)
}
