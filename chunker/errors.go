package chunker

import "fmt"

type chunkerError struct {
	Message string
	Code    int
}

func (e *chunkerError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}
