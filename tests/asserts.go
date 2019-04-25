package tests

import "testing"

// AssertNil checks that the "object" is nil. If the "object" is something other
// than nil, then the test will fail and exit with the object.
func AssertNil(t *testing.T, obj interface{}) {
	if obj != nil {
		t.Fatal(obj)
	}
}

// AssertNotNil checks that the "object" is not nil. If the "object" is nil then
// the test will fail and exit with the object.
func AssertNotNil(t *testing.T, obj interface{}, msg string) {
	if obj == nil {
		t.Fatal(msg)
	}
}

// AssertEqual makes sure that the two given "objects" are NOT different (!=).
func AssertEqual(t *testing.T, obj0 interface{}, obj1 interface{}, msg string) {
	if obj0 != obj1 {
		t.Fatalf("%s, expected %v, got %v\n", msg, obj0, obj1)
	}
}
