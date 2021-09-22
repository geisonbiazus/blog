package assert

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
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

func Contains(t *testing.T, content, expected string) {
	t.Helper()
	if !strings.Contains(content, expected) {
		t.Errorf("\nExpected to contain \"%s\" but it doesn't. \nContent:\n%s", expected, content)
	}
}

func Matches(t *testing.T, pattern, value string) {
	t.Helper()
	matched, _ := regexp.MatchString(pattern, value)

	if !matched {
		t.Errorf("\n\"%s\" did not match the regex: %s", value, pattern)
	}
}

func Error(t *testing.T, expected, actual error) {
	t.Helper()

	if !errors.Is(actual, expected) {
		t.Errorf("\nexpected: %v\n  actual: %v\n", expected, actual)
	}
}
