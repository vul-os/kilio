package middleware

import (
	"fmt"
	"runtime"
)

func NewError(err error) error {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("ERROR: [%s][%d] : %s", file, line, err)
}
