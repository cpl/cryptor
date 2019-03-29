package tests

import "testing"

// AssertNil checks that the err is nil. If the error is something other
// than nil, then the test will fail and exit with the error message.
func AssertNil(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// AssertEqual makes sure that the two given "interfaces" are NOT diffrent.
func AssertEqual(t *testing.T, i0 interface{}, i1 interface{}) {
	if i0 != i1 {
		t.Fatalf("not equal, expected %v = %v\n", i0, i1)
	}
}
