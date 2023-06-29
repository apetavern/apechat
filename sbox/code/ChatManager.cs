using Sandbox;
using System.Collections.Generic;
using System.Text.Json;

namespace ApeChat;

public class ChatManager
{
	private ChatConnection Chat { get; set; }

	public static ChatClient LocalClient => new( Game.UserName, Game.SteamId );
	public static Dictionary<string, Channel> Channels { get; set; } = new();
	public static Channel ActiveChannel { get; set; }

	public static ChatManager Instance;

	public ChatManager()
	{
		Chat = new ChatConnection();
		Instance = this;
	}

	public static void SetActiveChannel( string channelName )
	{
		var channel = Channels[channelName];
		ActiveChannel = channel;
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

	public static void SendMessage( string message )
	{
		var chatMsg = new ChatEvent
		{
			ChannelName = ActiveChannel.Name,
			Author = LocalClient.Name,
			Message = message,
		};

		var payload = Json.Serialize( chatMsg );
		var msg = new Event
		{
			MessageType = (int)EventType.ChatMessage,
			Payload = JsonDocument.Parse( payload )
		};

		_ = Instance.Chat.SendMessage( Json.Serialize( msg ) );
	}

	public static void HandleReceivedChatEvent( ChatEvent chatEvent )
	{
		if ( Channels.TryGetValue( chatEvent.ChannelName, out var channel ) )
		{
			channel.Messages.Add( chatEvent );
		}
		else
		{
			Log.Warning( "Received a message for a channel we do not know." );
		}
	}
}
