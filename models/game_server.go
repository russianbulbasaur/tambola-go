package models

import (
	json2 "encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
)

type GameServer struct {
	Id        int32
	Join      chan *User
	Leave     chan *User
	Broadcast chan []byte
	State     *GameState
}

func NewGameServer(gameID int32, host *User) *GameServer {
	log.Println(fmt.Sprintf("Making new game server with game id %d", gameID))
	return &GameServer{
		Id:        gameID,
		Join:      make(chan *User),
		Leave:     make(chan *User),
		Broadcast: make(chan []byte, 256),
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
			gs.unregisterPlayer(user)
		case message := <-gs.Broadcast:
			gs.broadcast(message)
		}
	}
}

func (gs *GameServer) registerPlayer(user *User) {
	gs.State.Players[user] = true
	userJoinedPayload := UserJoinedPayload{User: user}
	encodedPayload, err := json2.Marshal(userJoinedPayload)
	if err != nil {
		log.Println(err)
		return
	}
	message := Message{
		Payload: string(encodedPayload),
		Id:      -1,
		Event:   UserJoinedEvent,
		Sender: &User{
			Id:   -1,
			Name: "Server",
		},
	}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	gs.broadcast(encodedMessage)
	sendGameStateToJoinee(user, gs)
}

func sendGameStateToJoinee(player *User, gs *GameServer) {
	var players []*User
	for memberPlayer := range gs.State.Players {
		players = append(players, memberPlayer)
	}
	playersAlreadyInLobbyPayload := PlayersAlreadyInLobbyPayload{Players: players, GameId: gs.Id}
	encodedPayload, err := json2.Marshal(playersAlreadyInLobbyPayload)
	if err != nil {
		log.Println(err)
		return
	}
	message := &Message{Id: rand.Int64(),
		Payload: string(encodedPayload)}
	encodedMessage, err := json2.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	player.Send <- encodedMessage
}

func (gs *GameServer) unregisterPlayer(player *User) {
	if gs.State.Players[player] {
		delete(gs.State.Players, player)
	}
	userLeftPayload := UserLeftPayload{User: player}
	encodedPayload, err := json2.Marshal(userLeftPayload)
	message := Message{
		Payload: string(encodedPayload),
		Id:      -1,
		Event:   UserLeftEvent,
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
			player.Send <- data
		}
	}
}
