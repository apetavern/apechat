package chat

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	errBeforeUpgrade = flag.Bool("error-before-upgrade", false, "return an error on upgrade with body")
	keepAliveTime    = time.Second * 30
)

// Manager holds references to all clients, channels, and handlers
type Manager struct {
	clients  ClientList
	channels ChannelList
	handlers map[int]EventHandler
	sync.RWMutex
}

// NewManager initializes all the values inside the Manager
func NewManager() *Manager {
	channels := make(ChannelList)

	m := &Manager{
		clients:  ClientList{},
		channels: channels,
		handlers: map[int]EventHandler{},
	}
	m.setupEventHandlers()
	return m
}

// addClient adds a client to the Manager ClientList
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		err := client.Connection.Close()
		if err != nil {
			fmt.Printf("error occured while closing connection: %v\n", err)
		}
		delete(m.clients, client)
	}
}

func (m *Manager) newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.KeepaliveTime = keepAliveTime
	u.CheckOrigin = func(r *http.Request) bool {
		checkReferrer := r.Header.Get("Referer") == "https://sbox.facepunch.com/"
		checkAgent := r.Header.Get("User-Agent") == "facepunch-sbox"
		return checkReferrer && checkAgent
	}

	var client *Client
	u.OnMessage(func(c *websocket.Conn, messageType websocket.MessageType, data []byte) {
		if messageType == 1 {
			// Unmarshal incoming data into chat.Event
			var request Event
			if err := json.Unmarshal(data, &request); err != nil {
				log.Printf("error unmarshalling message: %v\n", err)
				return
			}

			// Route the event through the manager
			if err := m.routeEvent(request, client); err != nil {
				log.Printf("error handling message: %v\n", err)
				return
			}
		}
	})

	u.OnOpen(func(conn *websocket.Conn) {
		client = NewClient(conn, m)
		m.addClient(client)

		go client.handleHeartbeat()

		// Request info from the client
		if jsonEvent, err := BuildEvent(ClientInfo, ""); err == nil {
			_ = conn.WriteMessage(1, jsonEvent)
		}
	})

	u.OnClose(func(c *websocket.Conn, err error) {
		fmt.Println("OnClose:", c.RemoteAddr().String(), err)
		m.removeClient(client)
	})

	return u
}

func (m *Manager) onWebsocket(w http.ResponseWriter, r *http.Request) {
	if *errBeforeUpgrade {
		w.WriteHeader(http.StatusForbidden)
		_, err := w.Write([]byte("returning an error"))
		if err != nil {
			return
		}
		return
	}

	upgrader := m.newUpgrader()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		return
	}

	fmt.Println("OnOpen:", conn.RemoteAddr().String())
}

func (m *Manager) serve() error {
	mux := &http.ServeMux{}
	mux.HandleFunc("/ws", m.onWebsocket)

	svr := nbhttp.NewServer(nbhttp.Config{
		Network: "tcp",
		Addrs:   []string{"localhost:80"},
		Handler: mux,
	})

	err := svr.Start()
	if err != nil {
		return fmt.Errorf("nbio.Start failed: %v\n", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = svr.Shutdown(ctx)
	return err
}

func StartServer() error {
	manager := NewManager()
	err := manager.serve()
	return err
}
