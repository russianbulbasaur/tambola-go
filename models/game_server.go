package models

import (
	json2 "encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
)

type GameServer struct {
	Id        int32
	Join      chan *User
	Leave     chan *User
	Broadcast chan []byte
	Stop      chan int64
	State     *GameState
	Lock      sync.Mutex
}

func NewGameServer(gameID int32, host *User) *GameServer {
	log.Println(fmt.Sprintf("Making new game server with game id %d", gameID))
	return &GameServer{
		Id:        gameID,
		Join:      make(chan *User),
		Leave:     make(chan *User),
		Broadcast: make(chan []byte),
		Stop:      make(chan int64),
		State:     NewGameState(host),
	}
}

func (gs *GameServer) StartGameServer() {
	log.Println("Starting game server")
	for {
		select {
		case user := <-gs.Join:
			gs.registerPlayer(user)
		case user := <-gs.Leave:
			log.Printf("User Leaving : %#v", user)
			gs.unregisterPlayer(user)
		case message := <-gs.Broadcast:
			gs.broadcast(message)
		case <-gs.Stop:
			log.Println("Stopping game server")
			break
		}
	}
}

func (gs *GameServer) registerPlayer(user *User) {
	gs.addPlayerToState(user)
	userJoinedPayload := &UserJoinedPayload{User: user}
	message := Message{
		UserJoinedPayload: userJoinedPayload,
		Id:                -1,
		Event:             UserJoinedEvent,
		Sender: &User{
			Id:   -1,
			Name: "Server",
		},
	}
	encodedMessage, err := json2.Marshal(message)
	println(string(encodedMessage))
	if err != nil {
		log.Println(err)
		return
	}
	gs.broadcast(encodedMessage)
	gs.sendGameStateToJoinee(user)
}

func (gs *GameServer) sendGameStateToJoinee(player *User) {
	var players []*User
	for memberPlayer := range gs.State.Players {
		players = append(players, memberPlayer)
	}
	playersAlreadyInLobbyPayload := &PlayersAlreadyInLobbyPayload{Players: players, GameId: gs.Id}
	message := &Message{Id: rand.Int64(),
		PlayersAlreadyInLobbyPayload: playersAlreadyInLobbyPayload,
		Sender: &User{
			Id:   -1,
			Name: "Server",
		}, Event: PlayersInLobbyEvent}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	player.Send <- encodedMessage
}

func (gs *GameServer) unregisterPlayer(player *User) {
	if player.IsHost {
		gs.killGameServer()
		return
	}
	gs.removePlayerFromState(player)
	userLeftPayload := &UserLeftPayload{User: player}
	message := Message{
		UserLeftPayload: userLeftPayload,
		Id:              -1,
		Event:           UserLeftEvent,
	}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	gs.broadcast(encodedMessage)
}

func (gs *GameServer) broadcast(data []byte) {
	isForHost := gs.State.updateGameState(data)
	if !isForHost {
		gs.State.Host.Send <- data
	} else {
		for player := range gs.State.Players {
			log.Println(fmt.Sprintf("Sending to player %s", player.Name))
			player.Lock.Lock()
			player.Send <- data
			player.Lock.Unlock()
		}
	}

	if gs.State.playerCount == 0 {
		gs.killGameServer()
	}
}

func (gs *GameServer) removePlayerFromState(player *User) {
	if gs.State.Players[player] {
		delete(gs.State.Players, player)
		gs.State.playerCount--
	}
}

func (gs *GameServer) addPlayerToState(player *User) {
	gs.State.Players[player] = true
	gs.State.playerCount++
}

func (gs *GameServer) killGameServer() {
	gs.Stop <- 1
}
