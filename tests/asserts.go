package tests

import "testing"

// AssertNil checks that the err is nil. If the "object" is something other
// than nil, then the test will fail and exit with the object.
func AssertNil(t *testing.T, obj interface{}) {
	if obj != nil {
		t.Fatal(obj)
	}
}

// AssertEqual makes sure that the two given "objects" are NOT different (!=).
func AssertEqual(t *testing.T, obj0 interface{}, obj1 interface{}) {
	if obj0 != obj1 {
		t.Fatalf("not equal, expected %v = %v\n", obj0, obj1)
	}
}
