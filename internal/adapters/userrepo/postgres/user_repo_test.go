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

func dbTest(cb func(ctx context.Context, db *sql.DB)) {
	ctx := context.Background()
	db, _ := sql.Open("pgx", "postgres://postgres:postgres@localhost:5433/blog_test?sslmode=disable")
	defer db.Close()
	tx, _ := db.BeginTx(ctx, nil)
	ctx = context.WithValue(ctx, postgres.TxKey, tx)
	cb(ctx, db)
	tx.Rollback()
}

func TestUserRepo(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		t.Run("It creates a new user", func(t *testing.T) {
			dbTest(func(ctx context.Context, db *sql.DB) {
				repo := postgres.NewUserRepo(db)
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

				createdUser, err := repo.FindUserByID(ctx, user.ID)

				assert.Nil(t, err)
				assert.Equal(t, createdUser, user)
			})
		})

		t.Run("It returns error when user already exists", func(t *testing.T) {
			dbTest(func(ctx context.Context, db *sql.DB) {
				repo := postgres.NewUserRepo(db)
				uuidGen := uuid.NewGenerator()

				user := auth.User{
					ID:             uuidGen.Generate(),
					Name:           "Name",
					Email:          "user@example.com",
					ProviderUserID: "provider_user_id",
					AvatarURL:      "http://example.com/avatar",
				}

				repo.CreateUser(ctx, user)
				err := repo.CreateUser(ctx, user)

				assert.NotNil(t, err)
			})
		})
	})
}
