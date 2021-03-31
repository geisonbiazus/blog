package internal_test

import "testing"

func TestNothing(t *testing.T) {
	t.Run("Dummy test", func(t *testing.T) {
		if false != true {
			t.Error("Test failed.")
		}
	})
}
