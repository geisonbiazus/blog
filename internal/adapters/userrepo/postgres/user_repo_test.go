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

type testUserRepoFixture struct {
	repo    *postgres.UserRepo
	uuidGen *uuid.Generator
	user    auth.User
}

func TestUserRepo(t *testing.T) {
	setup := func(db *sql.DB) *testUserRepoFixture {
		repo := postgres.NewUserRepo(db)
		uuidGen := uuid.NewGenerator()

		user := auth.User{
			ID:             uuidGen.Generate(),
			Name:           "Name",
			Email:          "user@example.com",
			ProviderUserID: "provider_user_id",
			AvatarURL:      "http://example.com/avatar",
		}

		return &testUserRepoFixture{
			repo:    repo,
			uuidGen: uuidGen,
			user:    user,
		}
	}

	t.Run("CreateUser", func(t *testing.T) {
		t.Run("It creates a new user", func(t *testing.T) {
			dbTest(func(ctx context.Context, db *sql.DB) {
				f := setup(db)

				err := f.repo.CreateUser(ctx, f.user)

				assert.Nil(t, err)

				createdUser, err := f.repo.FindUserByID(ctx, f.user.ID)

				assert.Nil(t, err)
				assert.Equal(t, createdUser, f.user)
			})
		})

		t.Run("It returns error when user already exists", func(t *testing.T) {
			dbTest(func(ctx context.Context, db *sql.DB) {
				f := setup(db)

				f.repo.CreateUser(ctx, f.user)
				err := f.repo.CreateUser(ctx, f.user)

				assert.NotNil(t, err)
			})
		})
	})

	t.Run("UpdateUser", func(t *testing.T) {
		t.Run("It updates the user values", func(t *testing.T) {
			dbTest(func(ctx context.Context, db *sql.DB) {
				f := setup(db)
				user := f.user

				f.repo.CreateUser(ctx, user)

				user.AvatarURL = "https://example.com/new-avatar"
				user.Email = "new-email@example.com"
				user.Name = "new name"
				user.ProviderUserID = "new_provider_user_id"

				err := f.repo.UpdateUser(ctx, user)

				assert.Nil(t, err)

				updatedUser, err := f.repo.FindUserByID(ctx, user.ID)

				assert.Nil(t, err)
				assert.Equal(t, user, updatedUser)
			})
		})

		t.Run("It returns error when user does not exist", func(t *testing.T) {
			dbTest(func(ctx context.Context, db *sql.DB) {
				f := setup(db)

				err := f.repo.UpdateUser(ctx, f.user)

				assert.Equal(t, auth.ErrUserNotFound, err)
			})
		})

		t.Run("It returns error when update is not possible", func(t *testing.T) {
			dbTest(func(ctx context.Context, db *sql.DB) {
				f := setup(db)

				user := f.user
				f.repo.CreateUser(ctx, user)

				user2 := f.user

				user2.ID = f.uuidGen.Generate()
				user2.AvatarURL = "https://example.com/another-avatar"
				user2.Email = "another-email@example.com"
				user2.Name = "another name"
				user2.ProviderUserID = "another_provider_user_id"

				f.repo.CreateUser(ctx, user2)

				user.Email = user2.Email

				err := f.repo.UpdateUser(ctx, user)

				assert.NotNil(t, err)
			})
		})
	})
}
