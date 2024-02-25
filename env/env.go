package env

import (
	"os"
	"strconv"

	"git.qowevisa.me/Qowevisa/gotell/errors"
)

func GetHost() (string, error) {
	host := os.Getenv("GOTELL_HOST")
	if host == "" {
		return host, errors.ENV_EMPTY
	}
	return host, nil
}

func GetPort() (int, error) {
	portStr := os.Getenv("GOTELL_PORT")
	if portStr == "" {
		return 0, errors.ENV_EMPTY
	}
	port, err := strconv.ParseInt(portStr, 10, 32)
	return int(port), err
}
