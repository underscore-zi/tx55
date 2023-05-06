# TX-55 - Metal Gear Online Private Server

This is primarily a private server for Metal Gear Online 1 (MGO or MGO1), a game released for the Playstation 2 in 2006. Along with some associated services run by SaveMGO.com.

# Konami Server Protocol (pkg/konamiserver)

In doing the reverse engineering of the game's network traffic we realized that several Konami games around the same time used a similar TCP-layer protocol. So the [konamiserver](./pkg/konamiserver) package provides generic implementation of this protocol.

The main requirement to using it is to implement the `GameClient` interface to recieve the connection/disconnection and packet events, and a factory function to create a new client for each connection. The [metalgearonline1](./pkg/metalgearonline1) package is an example implementation of this. Though it has some added complexity to do automated parsing and routing of argument structures to the appropriate handlers.

With the interface implemented, the Konami server will handle the incoming packets verifying their hashes for you, and call your packet receiver with the incoming `Packet` structure, and a channel for an outbound packets. The implementation also supports hooking in and outbound packets to make on-the-fly modifications.

# Packet (pkg/packet)

This is kind of the core object for the konamiserver, it consumes Packet structures in, and you send it packets to send out. It encapsulates the Konami Server Protocol's packet structure and provides some helps for serializing/deserializing the data based on golang structures.

On the serializing aspect, if using a structure to fill the payload, they must be binary serializable, that is, it must be a fixed and determinable sized structure. Using things like `uint32` not `int` and only using fixed size arrays. The only real violation of that is that it supports a custom encoding tag `packet:"truncate"` which can be used to indicate that a packet ends with a variable length array that should be truncated before sending. For example the `GameInfo` structure defined in the types package, and used when sending information about a game from the game list to a client. Ends with:

```go
Players           [8]GamePlayerStats `packet:"truncate"`
```

The client could send NULL bytes to fill the space, which is what the binary serializer will do but this tag lets it know to truncate it. This works by looking for the first item that is a zero value (as defined by the reflection library). So you cannot use this if you intend to send zero filled items. This can also only be used with the last item in a structure/packet payload.

# Metal Gear Online 1 Server (pkg/metalgearonline1)

This package is the core of the game server, specifically the `handlers` package within. On top of the Konami Server GameClient interface, I impelment a sort of handler registration system. Various structs can define what packet type they handle and supported argument structures. Then packets are routed to the appropriate handler automatically. 

 - handlers - Contains all of these handler structures for all supported packets types. They are roughly broken down into user-state when the command would be called. So Authenticating, in-lobby, and hosting.
 - models - Contains the gorm models for all the game state data. Each model may also have associated conversion methods implemented for converting between the model and the game's binary representation of the data. 
 - session - Each handler gets a session pointer, this is the structure used to store the active connection state. So for example the User model is stored in here for reference once the user has logged in. This is also where any functionality that should be available to all the clients, like the `Session.Log` logger is accessible from here. Importantly this object knows nothing about the GameClient/Server to avoid any cyclic dependencies.
 - testclient - honestly, this is unmaintained garbage code used to invoke certain server actions without needing to use the game. It can send packets to the server and do some interactions but its really just intended for development testing.
 - types - Is the a central location for all the data types defined during reverse engineering. They are generally are binary serializable structures using sized integers or sized byte arrays. This are roughly grouped based on what I was working on when I discovered the type, but the organization could be better and is something to do in the future.

While most configuration happens through teh `gameserver` binary and its configuration file (see the `configurations` package) it does read one envrionment value `EVENTS_ENDPOINT` if specified the game will POST JSON objects to this endpoint as various game events happen. This is intended to be used with the `restapi` package's `/api/v1/stream/events/:token` endpoint.

# Metal Gear Online 1 - REST API (pkg/restapi)

This is a REST interface used to expose the game-server state to whoever wants that information, but it is also where regular crons such as updating rankings are run. There is also an optional `/api/v1/stream/events` endpoint that is a websocket endpoint that will stream JSON blobs of game-related events as they are reported by the game servers.

The API's endpoints are documented in [pkg/restapi/types.go](./pkg/restapi/types.go), but every request will get a `ResponseJSON` as the response object with a varying `Data` field depending on the request.


# Running the Server for Development

## Cheats
One of the bigger challenges for this is just getting a copy of the game to communicate with a custom server. Right now this does require the use of "cheats" or memory editing. See [cheats.md](./cheats.md) for information about what each of the cheats does but I will assume here that you are running the appropiate cheats for your game version AND cheat device.

## DNS 
Next when you run the game it will make several DNS lookups to get to the game server. You will need to setup your game so that it points to a DNS you control. You can do this with the built-in network configuration editing, or the PS2 network adapters came with a configuration disc. If using PCSX2 (prefered) you can just do this in the settings. The main domains you need to resolve are:

- mgs3sstun.konamionline.com - this needs to point to some compatible STUN server. I used stund to run one you may be able to find a public STUN server but not all STUN servers are compatible and I'm not entirely sure what the requirements are, but running my own always worked.
- savemgo.com - One of the cheats rewrites a konami domain to savemgo.com. This is so we could get around the HTTPS requirement. The game makes a few web requests to this server:
 - /us/mgs3/text/policy_us_en.txt - This is the terms that you initially have to agree to when you are connecting. There is a different file for each language.
 - chgpswd.html, deluser.html, reguser.html - I don't recall exactly what path these make the request to, but the last segment was /<filename>. The expected response value on each is just a "0".
- mgs3sgate02.konamionline.com - This is the actual game server location. It'll first reach out on :5731 to this server. The gate number changes depending on the game's region; 02 is for North America. It reaches here and makes a request for the game's lobbies. Which can be run anywhere. 

## GameServer Configuration

See [pkg/configuration](./pkg/configurations/gameserver.go) for specific details about the configuration format. But basically you just provide it a host/port to listen on, a database connection string and some database config options, and a lobby id to run as. The LobbyID should reflect the ID of the server in the lobby table it is running as, but it doesn't have to. It never uses the info from the table for any reason, only tries to update its player count based on this id and of course games created on it will use the id.

## `gameserver` binary

This should just be a `go build tx55/cmd/gameserver` to build. It does require CGO be enabled, `CGO_ENABLED=1` but otherwise it should build and run. 


## Database

The first run with the gameserver you should provide a config via `-config <filename>` and use the `-migrate` flag. The Migrate flag will create all the necessary database tables. You can use this any time you update this server, though I don't forsee often needing to change the database schema on such an old game. 

Within the database, you'll need to have atleast three entries in the lobby table:
 - The first entry should be of type `LobbyTypeGate` which is the gate server. I don't thing the game actually uses this value, but it expects it to be there and its best just to point it to the same place as the DNS entry for the gate server and port 5731.
 - The second is the account/auth server (`LobbyTypeGate`). After making the lobby list request, it'll try to do some auth related actions, and eventually login. The metalgearonline1 implementation handles all packets the same regardless of the type of server so the distinction between gate/auth/game doesn't matter and these can all be the same, but the game will split requests across them depending on what it is doing.
 - Lastly is the `LobbyTypeGame` entries, these are the ones actually displayed to the user  and that users can connect to and create games on.


