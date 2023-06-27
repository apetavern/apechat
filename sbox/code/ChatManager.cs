using Sandbox;
using System.Collections.Generic;
using System.Text.Json;

namespace ApeChat;

public static class ChatManager
{
	private static readonly ChatConnection _chat = new();

	public static ChatClient LocalClient => new( Game.UserName, Game.SteamId );
	public static Dictionary<string, Channel> Channels { get; set; } = new();

	public static void ChannelCreate( string channelName )
	{
		var subMsg = new SubscriptionEvent
		{
			ChannelName = channelName,
			Subscribe = true
		};

		var payload = Json.Serialize( subMsg );
		var msg = new Event
		{
			MessageType = (int)EventType.Subscription,
			Payload = JsonDocument.Parse( payload )
		};

		_ = _chat.SendMessage( Json.Serialize( msg ) );
	}
}
