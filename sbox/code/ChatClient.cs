using System.Text.Json.Serialization;

namespace ApeChat;

public class ChatClient
{
	[JsonPropertyName( "name" )]
	public string Name { get; set; }
	public long SteamId { get; set; }

	public ChatClient( string name, long steamId )
	{
		Name = name;
		SteamId = steamId;
	}
}
