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

func transactional(ctx context.Context, db *sql.DB, cb func(ctx context.Context)) {
	tx, _ := db.BeginTx(ctx, nil)
	ctx = context.WithValue(ctx, postgres.TxKey, tx)
	cb(ctx)
	tx.Rollback()
}

func TestUserRepo(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		ctx := context.Background()
		db, _ := sql.Open("pgx", "postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable")
		defer db.Close()
		repo := postgres.NewUserRepo(db)
		uuidGen := uuid.NewGenerator()

		transactional(ctx, db, func(ctx context.Context) {
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

	})
}
