package myerr

import (
	"errors"
)

var (
	INTERNAL_ERROR   = errors.New("Internal error")
	UNIQUENESS_ERROR = errors.New("Are not unique")
)
