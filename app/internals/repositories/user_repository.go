package repositories

import "database/sql"

type userRepository struct {
	db *sql.DB
}

type UserRepository interface {
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
