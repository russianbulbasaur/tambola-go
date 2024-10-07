package game_server

import (
	"cmd/tambola/internals/repositories"
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type gameSyncer struct {
	gameRepo      repositories.GameRepository
	stripedGameId int64
	gameState     GameState
	wg            *sync.WaitGroup
}

type GameSyncer interface {
	Sync(gameCtx context.Context)
}

func NewGameSyncer(gameId string, gameRepo repositories.GameRepository,
	gameState GameState, wg *sync.WaitGroup) GameSyncer {
	stripedGameId, err := strconv.ParseInt(strings.Trim(gameId, "TMB"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return &gameSyncer{
		stripedGameId: stripedGameId,
		gameRepo:      gameRepo,
		gameState:     gameState,
		wg:            wg,
	}
}

func (gs *gameSyncer) close() {
	log.Printf("TMB%[1]d : Killing syncer %[1]d", gs.stripedGameId)
	gs.wg.Done()
}

func (gs *gameSyncer) Sync(gameCtx context.Context) {
	defer gs.close()
	log.Printf("TMB%d : Launching syncer %d\n", gs.stripedGameId, gs.stripedGameId)
	for {
		select {
		case <-gameCtx.Done():
			return
		default:
			gameActivity := gs.gameState.GetStateJson()
			gs.gameRepo.UpdateGameActivity(gameActivity, gs.stripedGameId)
			time.Sleep(time.Second * 10)
		}
	}
}
