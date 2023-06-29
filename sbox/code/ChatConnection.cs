using Sandbox;
using System.Text.Json;
using System.Threading.Tasks;

namespace ApeChat;

public class ChatConnection
{
	private readonly WebSocket _webSocket;

	public ChatConnection()
	{
		_webSocket = new WebSocket();

		_ = Connect();
	}

	private async Task Connect()
	{
		_webSocket.OnMessageReceived += OnMessageReceived;

		await _webSocket.Connect( "ws://localhost:80/ws" );
	}

	private void OnMessageReceived( string jsonMessage )
	{
		var message = JsonSerializer.Deserialize<Event>( jsonMessage );
		var messageType = (EventType)message.MessageType;
		var payload = message.Payload;
		switch ( messageType )
		{
			case EventType.ChatMessage:
				var chatMessage = payload.Deserialize<ChatEvent>();
				ChatManager.HandleReceivedChatEvent( chatMessage );
				break;
			case EventType.ChannelInfo:
				var channelInfo = payload.Deserialize<ChannelInfoEvent>();
				if ( ChatManager.Channels.TryGetValue( channelInfo.ChannelName, out var channel ) )
				{
					channel.Update( channelInfo.Clients.Length );
				}
				else
				{
					ChatManager.Channels.Add( channelInfo.ChannelName, new Channel( channelInfo.ChannelName, channelInfo.Clients.Length ) );
				}
				break;
			case EventType.ClientInfo:
				var clientInfo = new ClientInfoEvent
				{
					Name = Game.UserName,
					SteamId = Game.SteamId,
				};
				var p = Json.Serialize( clientInfo );
				var e = new Event
				{
					MessageType = (int)EventType.ClientInfo,
					Payload = JsonDocument.Parse( p )
				};
				_ = SendMessage( Json.Serialize( e ) );
				break;
			case EventType.Heartbeat:
				var heartbeat = new Event
				{
					MessageType = (int)EventType.Heartbeat,
				};
				_ = SendMessage( Json.Serialize( heartbeat ) );
				break;
			default:
				break;
		}
	}

	public async Task SendMessage( string jsonMessage )
	{
		Log.Info( jsonMessage );
		await _webSocket.Send( jsonMessage );
	}
}
