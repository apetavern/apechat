package chat

import (
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log"
	"time"
)

var (
	pongWait     = time.Second * 10
	pingInterval = (pongWait * 9) / 10
)

// ClientList is a map used to help manage clients
type ClientList map[*Client]bool

// Client represents a websocket client
type Client struct {
	Connection *websocket.Conn
	Manager    *Manager
	Name       string
	SteamId    int64
}

type ChatClient struct {
	Name    string `json:"name"`
	SteamId int64  `json:"steamId"`
}

// NewClient initializes a new Client
func NewClient(conn *websocket.Conn, m *Manager) *Client {
	return &Client{
		Connection: conn,
		Manager:    m,
	}
}

func (c *Client) handleHeartbeat() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.Manager.removeClient(c)
	}()

	for {
		select {
		case <-ticker.C:
			jsonEvent, err := BuildEvent(Heartbeat, "")
			if err != nil {
				return
			}
			if err := c.Connection.WriteMessage(1, jsonEvent); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
