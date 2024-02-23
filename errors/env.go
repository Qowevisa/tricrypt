package errors

import "errors"

var (
	ENV_EMPTY = errors.New("Environment variable was empty")
)
