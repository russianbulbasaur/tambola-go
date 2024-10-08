package game_server

import (
	"cmd/tambola/models"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const (
	writeWait = 10 * time.Second

	pongWait = 20 * time.Second

	pingPeriod = 10 * time.Second

	maxMessageSize = 10000
)

type Player struct {
	*models.User
	Conn       *websocket.Conn `json:"-"`
	GameServer GameServer      `json:"-"`
	Send       chan []byte     `json:"-"`
	Lock       sync.Mutex      `json:"-"`
}

func (player *Player) disconnect() {
	log.Printf(fmt.Sprintf("%s : Killing %s's write thread", player.GameServer.GetGameId(),
		player.GetName()))
	player.GameServer.RemovePlayer(player)
	err := player.Conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (player *Player) ReadPump(gameCtx context.Context) {
	defer func() {
		log.Printf(fmt.Sprintf("%s : Killing %s's read thread", player.GameServer.GetGameId(),
			player.GetName()))
	}()

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

func (player *Player) WritePump(gameCtx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		player.disconnect()
	}()
	for {
		select {
		case message, ok := <-player.Send:
			player.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel
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
