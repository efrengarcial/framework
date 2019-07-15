package util

import "github.com/Masterminds/goutils"

func GeneratePassword() (string, error) {
	return goutils.RandomAlphaNumeric(20)
}

func GenerateResetKey() (string, error) {
	return goutils.RandomNumeric(20)
}
