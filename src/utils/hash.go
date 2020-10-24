package utils

import (
	"errors"
	"strconv"

	"github.com/segmentio/fasthash/fnv1a"
)

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password should not be empty")
	}

	return strconv.FormatUint(fnv1a.HashString64(password), 32), nil
}
