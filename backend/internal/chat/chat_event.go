package chat

import (
	"encoding/json"
)

// Event is an incoming or outgoing message, with a type defined
// to determine how the payload is handled.
type Event struct {
	Type    int             `json:"messageType"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

type EventType int

const (
	Message      EventType = 0
	Subscription EventType = 1
	ChannelInfo  EventType = 2
	ClientInfo   EventType = 3
	Heartbeat    EventType = 99
)

type SendChatMessageEvent struct {
	Channel string `json:"channel"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type SubscribeEvent struct {
	Subscribe bool   `json:"subscribe"`
	Channel   string `json:"channel"`
}

type ChannelInfoEvent struct {
	Channel   string       `json:"channel"`
	ChannelID int32        `json:"channelId"`
	Clients   []ChatClient `json:"clients"`
}

type ClientInfoEvent struct {
	Name    string `json:"name"`
	SteamId int64  `json:"steamId"`
}

func BuildEvent(t EventType, i interface{}) ([]byte, error) {
	// Marshal the generic interface event
	marshalledEvent, err := json.Marshal(i)
	if err != nil {
		return []byte{}, ErrFailedToMarshal
	}
	// Encode marshalled event to a json.RawMessage
	payload := json.RawMessage(marshalledEvent)

	// Build the event with EventType t
	event := Event{
		Type:    int(t),
		Payload: payload,
	}

	// Marshal the whole event and return as []byte
	if eventJson, err := json.Marshal(event); err == nil {
		return eventJson, nil
	}
	return []byte{}, err
}
