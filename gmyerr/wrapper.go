package gmyerr

import (
	"fmt"
)

func WrapPrefix(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", prefix, err)
}
