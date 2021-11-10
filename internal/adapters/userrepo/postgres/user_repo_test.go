package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/geisonbiazus/blog/internal/adapters/idgenerator/uuid"
	"github.com/geisonbiazus/blog/internal/adapters/userrepo/postgres"
	"github.com/geisonbiazus/blog/internal/core/auth"
	"github.com/geisonbiazus/blog/pkg/assert"
)

func TestUserRepo(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		ctx := context.Background()
		db, _ := sql.Open("pgx", "postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable")
		tx, _ := db.BeginTx(ctx, nil)
		defer tx.Rollback()
		defer db.Close()

		repo := postgres.NewUserRepo(tx)
		uuidGen := uuid.NewGenerator()

		user := auth.User{
			ID:             uuidGen.Generate(),
			Name:           "Name",
			Email:          "user@example.com",
			ProviderUserID: "provider_user_id",
			AvatarURL:      "http://example.com/avatar",
		}

		err := repo.CreateUser(ctx, user)

		assert.Nil(t, err)

		createdUser, err := repo.FindById(ctx, user.ID)

		assert.Nil(t, err)
		assert.Equal(t, createdUser, user)
	})
}
