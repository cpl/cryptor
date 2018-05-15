package utils

import "testing"

// CheckErr ...
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// CheckErrTest ...
func CheckErrTest(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
