package env_test

import (
	"os"
	"testing"

	"github.com/geisonbiazus/blog/pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Run("It returns default value when variable is not present", func(t *testing.T) {
		assert.Equal(t, "default", env.GetString("TEST_VAR", "default"))
	})

	t.Run("It returns an ENV var value when variable is present", func(t *testing.T) {
		os.Setenv("TEST_VAR", "value")
		assert.Equal(t, "value", env.GetString("TEST_VAR", "default"))
		os.Unsetenv("TEST_VAR")
	})
}

func TestGetInt(t *testing.T) {
	t.Run("It returns default value when variable is not present", func(t *testing.T) {
		assert.Equal(t, 200, env.GetInt("TEST_VAR", 200))
	})

	t.Run("It returns an ENV var value when variable is present", func(t *testing.T) {
		os.Setenv("TEST_VAR", "3000")
		assert.Equal(t, 3000, env.GetInt("TEST_VAR", 200))
		os.Unsetenv("TEST_VAR")
	})

	t.Run("It returns default value var value when variable is present but not a number", func(t *testing.T) {
		os.Setenv("TEST_VAR", "asd")
		assert.Equal(t, 200, env.GetInt("TEST_VAR", 200))
		os.Unsetenv("TEST_VAR")
	})
}
