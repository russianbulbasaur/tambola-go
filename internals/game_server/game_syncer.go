package game_server

import (
	"cmd/tambola/internals/repositories"
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

type gameSyncer struct {
	gameRepo      repositories.GameRepository
	stripedGameId int64
	gameState     GameState
}

type GameSyncer interface {
	Sync(gameCtx context.Context)
}

func NewGameSyncer(gameId string, gameRepo repositories.GameRepository,
	gameState GameState) GameSyncer {
	stripedGameId, err := strconv.ParseInt(strings.Trim(gameId, "TMB"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return &gameSyncer{
		stripedGameId: stripedGameId,
		gameRepo:      gameRepo,
		gameState:     gameState,
	}
}

func (gs *gameSyncer) Sync(gameCtx context.Context) {
	log.Printf("Launching syncer %d\n", gs.stripedGameId)
	for {
		select {
		case _ = <-gameCtx.Done():
			log.Printf("Killing syncer %d", gs.stripedGameId)
			return
		default:
			gameActivity := gs.gameState.GetStateJson()
			gs.gameRepo.UpdateGameActivity(gameActivity, gs.stripedGameId)
			time.Sleep(time.Second * 10)
		}
	}
}
