package utils

import "fmt"

func NewError(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

type SystemError error