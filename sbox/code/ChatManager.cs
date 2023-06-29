using Sandbox;
using System.Collections.Generic;
using System.Text.Json;

namespace ApeChat;

public class ChatManager
{
	private ChatConnection Chat { get; set; }

	public static ChatClient LocalClient => new( Game.UserName, Game.SteamId );
	public static Dictionary<string, Channel> Channels { get; set; } = new();

	public static ChatManager Instance;

	public ChatManager()
	{
		Chat = new ChatConnection();
		Instance = this;
	}

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

		_ = Instance.Chat.SendMessage( Json.Serialize( msg ) );
	}
}
