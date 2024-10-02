package repositories

import (
	"database/sql"
	"log"
)

type authRepository struct {
	db *sql.DB
}

type AuthRepository interface {
	Signup(string, string) (int64, error)
	UserExists(string) (bool, error)
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (ar *authRepository) UserExists(phone string) (bool, error) {
	results, err := ar.db.Query(`select count(*) as count from users where phone=?`,
		phone)
	defer results.Close()
	if err != nil {
		return false, err
	}
	var count int32
	for results.Next() {
		err := results.Scan(&count)
		if err != nil {
			return false, err
		}
	}
	return count != 0, nil
}

func (ar *authRepository) Signup(name string, phone string) (int64, error) {
	results, err := ar.db.Query(
		`insert into users(name,phone) values(?,?) returning id`,
		name, phone)
	defer results.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var userId int64
	for results.Next() {
		err := results.Scan(&userId)
		if err != nil {
			return -1, err
		}
	}
	return userId, nil
}
