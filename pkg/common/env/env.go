package env

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var ErrEnvNotSet = errors.New("environment variable not set")

func GetBool(key string) bool {
	val := os.Getenv(key)
	return val != ""
}

func GetInt(key string) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return -1, ErrEnvNotSet
	}
	i, err := strconv.ParseInt(val, 10, 32)
	return int(i), errors.Wrap(err, "failed to parse number")
}

func GetString(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", ErrEnvNotSet
	}
	return val, nil
}

func GetDuration(key string) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		return 0, ErrEnvNotSet
	}
	t, err := time.ParseDuration(val)
	return t, errors.Wrap(err, "failed to parse duration")
}
