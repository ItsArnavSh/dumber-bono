package utility

import (
	"errors"
	"math/rand"
)

func RandomString(items []string) (string, error) {
	if len(items) == 0 {
		return "", errors.New("cannot select from empty slice")
	}
	return items[rand.Intn(len(items))], nil
}
