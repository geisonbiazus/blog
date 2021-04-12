package assert

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if actual != expected {
		t.Errorf("\nexpected: %v\n  actual: %v\n", expected, actual)
	}
}

func DeepEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nexpected: %v\n  actual: %v\n", expected, actual)
	}
}

func Nil(t *testing.T, actual interface{}) {
	t.Helper()
	if actual != nil {
		t.Errorf("\n%v should be nil", actual)
	}
}

func NotNil(t *testing.T, actual interface{}) {
	t.Helper()
	if actual == nil {
		t.Errorf("\n%v should not be nil", actual)
	}
}

func True(t *testing.T, actual bool) {
	t.Helper()
	if !actual {
		t.Errorf("\n%v should not be true", actual)
	}
}

func False(t *testing.T, actual bool) {
	t.Helper()
	if actual {
		t.Errorf("\n%v should not be false", actual)
	}
}
