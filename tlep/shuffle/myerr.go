package shuffle

import "errors"

var (
	BENCHMARK_INTERNAL_ERROR = errors.New("Internal error")
	BENCHMARK_LEN_DN_MATCH   = errors.New("Lengths do not match")
)
