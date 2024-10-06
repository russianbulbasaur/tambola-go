package repositories

import (
	"cmd/tambola/models"
	"database/sql"
	"log"
)

type gameRepository struct {
	db *sql.DB
}

type GameRepository interface {
	CreateGame(host *models.User) int64
	UpdateGameActivity(activity string, gameId int64)
}

func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{db: db}
}

func (gr *gameRepository) CreateGame(host *models.User) int64 {
	results, err := gr.db.Query(
		`insert into games(host_id) values($1) returning id`, host.Id)
	if err != nil {
		log.Fatalln(err)
	}
	var gameId int64
	for results.Next() {
		err = results.Scan(&gameId)
		if err != nil {
			log.Fatalln(err)
		}
	}
	results.Close()
	return gameId
}

func (gr *gameRepository) UpdateGameActivity(activity string, gameId int64) {
	_, err := gr.db.Exec(
		`update games set activity=$1 where id=$2`, activity, gameId)
	if err != nil {
		log.Fatalln(err)
	}
}
