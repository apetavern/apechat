package chat

import (
	"encoding/json"
	"errors"
	"github.com/lesismal/nbio/nbhttp/websocket"
)

var (
	ErrEventNotRecognized = errors.New("event type not recognized")
	ErrFailedToUnmarshal  = errors.New("failed to unmarshal payload")
	ErrFailedToMarshal    = errors.New("failed to marshal response data")
)

// routeEvent routes a given event based in its type into the correct handler
func (m *Manager) routeEvent(e Event, c *Client) error {
	if handler, ok := m.handlers[e.Type]; ok {
		if err := handler(e, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotRecognized
	}
}

// setupEventHandlers stores all functions for handling event types in the Manager
func (m *Manager) setupEventHandlers() {
	m.handlers[int(Message)] = handleChatMessage
	m.handlers[int(Subscription)] = handleSubscribeMessage
	m.handlers[int(ClientInfo)] = handleClientInfo
}

// handleChatMessage accepts and handles a chat message from a chat.Client
func handleChatMessage(e Event, c *Client) error {
	var chatEvent SendChatMessageEvent
	if err := json.Unmarshal(e.Payload, &chatEvent); err != nil {
		return ErrFailedToUnmarshal
	}

	for channel, _ := range c.Manager.channels {
		if channel.Name == chatEvent.Channel {
			channel.OnMessageReceived(chatEvent)
			break
		}
	}

	return nil
}

// handleSubscribeMessage handles subscribing and unsubscribing from channels
func handleSubscribeMessage(e Event, c *Client) error {
	var subscribeEvent SubscribeEvent
	if err := json.Unmarshal(e.Payload, &subscribeEvent); err != nil {
		return ErrFailedToUnmarshal
	}

	m := c.Manager

	if subscribeEvent.Subscribe {
		for channel, _ := range m.channels {
			if channel.Name == subscribeEvent.Channel {
				channel.Subscribe(*c)
				// Send the updated channel info to all clients
				if eventJson, err := BuildEvent(ChannelInfo, channel.ToChannelInfoMessage()); err == nil {
					for cl, _ := range channel.clients {
						_ = cl.Connection.WriteMessage(1, eventJson)
					}
				}
				return nil
			}
		}

		// If no channel was found, create one.
		ch := NewChannel(subscribeEvent.Channel)
		ch.Subscribe(*c)
		m.channels[ch] = true

		// Send the client the info for this channel.
		if eventJson, err := BuildEvent(ChannelInfo, ch.ToChannelInfoMessage()); err == nil {
			_ = c.Connection.WriteMessage(1, eventJson)
		}
	} else {
		for channel, _ := range m.channels {
			if channel.Name == subscribeEvent.Channel {
				channel.Unsubscribe(*c)
				// Send the updated channel info to all clients
				if eventJson, err := BuildEvent(ChannelInfo, channel.ToChannelInfoMessage()); err == nil {
					for cl, _ := range channel.clients {
						_ = cl.Connection.WriteMessage(1, eventJson)
					}
				}
				return nil
			}
		}
	}

	return nil
}

func handleClientInfo(e Event, c *Client) error {
	var clientInfoEvent ClientInfoEvent
	if err := json.Unmarshal(e.Payload, &clientInfoEvent); err != nil {
		return ErrFailedToUnmarshal
	}

	c.Name = clientInfoEvent.Name
	c.SteamId = clientInfoEvent.SteamId
	return nil
}

func pingMessageHandler(conn *websocket.Conn, _ string) {
	if err := conn.WriteMessage(websocket.PongMessage, []byte("PONG")); err != nil {
		_ = conn.Close()
	}
}

func pongMessageHandler(conn *websocket.Conn, _ string) {
	if err := conn.WriteMessage(websocket.PingMessage, []byte("PING")); err != nil {
		_ = conn.Close()
	}
}
