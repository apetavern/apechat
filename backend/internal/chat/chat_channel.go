package chat

import (
	"fmt"
)

type ChannelList map[*Channel]bool

type Channel struct {
	Name    string
	clients map[Client]bool
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:    name,
		clients: map[Client]bool{},
	}
}

func (c *Channel) ToChannelInfoMessage() *ChannelInfoEvent {
	var clients []ChatClient
	for client, _ := range c.clients {
		clients = append(clients, ChatClient{Name: client.Name, SteamId: client.SteamId})
	}

	return &ChannelInfoEvent{
		Channel: c.Name,
		Clients: clients,
	}
}

func (c *Channel) Subscribe(client Client) {
	c.clients[client] = true
}

func (c *Channel) Unsubscribe(client Client) {
	delete(c.clients, client)
}

func (c *Channel) OnMessageReceived(msg SendChatMessageEvent) {
	eventJson, err := BuildEvent(Message, msg)
	if err != nil {
		fmt.Printf(ErrFailedToMarshal.Error())
		return
	}

	for client, connected := range c.clients {
		if !connected {
			continue
		}

		err := client.Connection.WriteMessage(
			1,
			eventJson,
		)
		if err != nil {
			fmt.Printf("an error occured while writing: %v\n", err)
		}
	}
}
