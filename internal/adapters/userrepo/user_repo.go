package userrepo

import (
	"database/sql"

	"github.com/geisonbiazus/blog/internal/adapters/userrepo/memory"
	"github.com/geisonbiazus/blog/internal/adapters/userrepo/postgres"
)

func NewMemoryUserRepo() *memory.UserRepo {
	return memory.NewUserRepo()
}

func NewPostgresUserRepo(db *sql.DB) *postgres.UserRepo {
	return postgres.NewUserRepo(db)
}
