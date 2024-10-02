package models

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"runtime"
	"sync"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type User struct {
	Id         int64           `json:"id"`
	Name       string          `json:"name"`
	Conn       *websocket.Conn `json:"-"`
	GameServer GameServer      `json:"-"`
	Send       chan []byte     `json:"-"`
	IsHost     bool            `json:"is_host"`
	Lock       sync.Mutex      `json:"-"`
}

func (player *User) disconnect() {
	player.GameServer.Log(fmt.Sprintf("Killing %s's read thread", player.Name))
	log.Println("Goroutines : ", runtime.NumGoroutine())
	player.GameServer.RemovePlayer(player)
	player.Conn.Close()
}

func (player *User) ReadPump(gameCtx context.Context) {
	defer player.disconnect()

	player.Conn.SetReadLimit(maxMessageSize)
	player.Conn.SetReadDeadline(time.Now().Add(pongWait))
	player.Conn.SetPongHandler(func(string) error { player.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		select {
		case <-gameCtx.Done():
			return
		default:
			_, jsonMessage, err := player.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("unexpected close error: %v", err)
				}

				return
			}
			player.GameServer.BroadcastMessage(jsonMessage)
		}
	}
}

func (player *User) WritePump(gameCtx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		player.GameServer.Log(fmt.Sprintf("Killing %s's write thread", player.Name))
		log.Println("Goroutines : ", runtime.NumGoroutine())
		ticker.Stop()
		player.disconnect()
	}()
	for {
		select {
		case message, ok := <-player.Send:
			player.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := player.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			player.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := player.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-gameCtx.Done():
			return
		}
	}
}
