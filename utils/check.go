package utils

import (
	"log"
	"testing"
)

// CheckErr ...
func CheckErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// CheckErrTest ...
func CheckErrTest(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
