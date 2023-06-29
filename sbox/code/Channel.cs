using System.Collections.Generic;

namespace ApeChat;

public class Channel
{
	public string Name { get; set; }
	public int UserCount { get; set; }
	public List<ChatEvent> Messages { get; set; } = new();

	public Channel( string name, int userCount )
	{
		Name = name;
		UserCount = userCount;
	}

	public void Update( int userCount )
	{
		UserCount = userCount;
	}
}
