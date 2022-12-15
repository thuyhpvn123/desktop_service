package common

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrorInvalidConnectionAddress = errors.New("invalid connection address")
)

func SplitConnectionAddress(address string) (ip string, port int, err error) {
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return "", 0, ErrorInvalidConnectionAddress
	}
	intPort, err := strconv.Atoi(splited[1])
	if err != nil {
		return "", 0, err
	}
	return splited[0], intPort, nil
}
