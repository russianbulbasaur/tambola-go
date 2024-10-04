package repositories

import (
	"cmd/tambola/models"
	"database/sql"
	"log"
)

type authRepository struct {
	db *sql.DB
}

type AuthRepository interface {
	Signup(string, string) (*models.User, error)
	FindUser(string) (*models.User, error)
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (ar *authRepository) FindUser(phone string) (*models.User, error) {
	results, err := ar.db.Query(`select id,name,phone from users where phone=?`,
		phone)
	defer results.Close()
	if err != nil {
		return nil, err
	}
	var user models.User
	for results.Next() {
		err := results.Scan(&user.Id, &user.Name, &user.Phone)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (ar *authRepository) Signup(name string, phone string) (*models.User, error) {
	results, err := ar.db.Query(
		`insert into users(name,phone) values(?,?) returning id,name,phone`,
		name, phone)
	defer results.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var user models.User
	for results.Next() {
		err := results.Scan(&user.Id, &user.Name, &user.Phone)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}
