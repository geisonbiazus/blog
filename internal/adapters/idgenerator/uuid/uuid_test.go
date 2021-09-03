package uuid_test

import (
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
	"github.com/geisonbiazus/blog/pkg/assert"
)

const UUIDRegex = `^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`

func TestUUIDGenerator(t *testing.T) {
	t.Run("It generates a v4 uuid", func(t *testing.T) {
		idgen := uuid.NewGenerator()
		assert.Matches(t, idgen.Generate(), UUIDRegex)
	})
}
