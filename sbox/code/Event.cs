using System;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace ApeChat;

public class Event
{
	[JsonPropertyName( "messageType" )] public int MessageType { get; set; }
	[JsonPropertyName( "payload" )] public JsonDocument Payload { get; set; }

	public override string ToString()
	{
		return $"{MessageType}: {Payload}";
	}
}

public class ChatEvent
{
	[JsonPropertyName( "channel" )] public string ChannelName { get; set; }
	[JsonPropertyName( "author" )] public string Author { get; set; }
	[JsonPropertyName( "message" )] public string Message { get; set; }
	public DateTime Time { get; set; }

	public override string ToString()
	{
		return $"{Time.ToString( "h:mm tt" )} <{Author}>: {Message}";
	}
}

public class ChannelInfoEvent
{
	[JsonPropertyName( "channel" )] public string ChannelName { get; set; }
	[JsonPropertyName( "channelId" )] public int ChannelId { get; set; }
	[JsonPropertyName( "clients" )] public ChatClient[] Clients { get; set; }
}

public class SubscriptionEvent
{
	[JsonPropertyName( "subscribe" )] public bool Subscribe { get; set; }
	[JsonPropertyName( "channel" )] public string ChannelName { get; set; }
}

public class ClientInfoEvent
{
	[JsonPropertyName( "name" )] public string Name { get; set; }
	[JsonPropertyName( "steamId" )] public long SteamId { get; set; }
}

public enum EventType
{
	ChatMessage = 0,
	Subscription = 1,
	ChannelInfo = 2,
	ClientInfo = 3,
	Heartbeat = 99,
}
