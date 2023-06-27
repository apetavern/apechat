package chat

import "github.com/lesismal/nbio/nbhttp/websocket"

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
