# Flow

This document describes the general flow of the protocol.
All the steps will include their protocol type

> ## Editor's Note
>
> This documentation is for the VC2.0 protocol, and is not yet complete.
> The BeamMP team are currently working on a new protocol, which is slated to be released sometime in 2024, however I am unable to find any information on the new protocol, nor am I aware of any implementations of it from which to base any documentation. Once this new protocol is released, I will rename this document to `vc2.0.md` and update the links to the new protocol.

## Packet Definition

### TCP

```md
[][][][] [][][][][]...
 Header    Data
```

Header is a 4 byte signed integer, and the data is of an undefined length, as specified by the header. The maximum size of the data is 100 MB. A negative header value should cause the server to close the connection, as it is not a part of the protocol, and indicates an improperly implemented client.

### UDP

> Not currently documented

## Normal Connection

The client will connect to the server, and immediately send a TCP packet with a single byte, to define the expected behaviour of the server. The server will then respond appropriately following the normal control flow based on the byte received.

### First Byte

| Value | Description                                                                               |
|:-----:|-------------------------------------------------------------------------------------------|
|  `C`  | The client is requesting to connect to the server.                                        |
|  `D`  | The client is requesting to download mods from the server.                                |
|  `P`  | This option is not known by me, but I believe it is a health check by the listing server. |

### C

> The client is requesting to connect to the server.

The client will open a TCP connection to the server, and send a single byte, `C`, to indicate that it is ready to connect.

The table below describes the standard control flow to continue the connection process.

> The length constraints are applied to only the data portion of the given packet. Headers are not included in the length calculation.

| Direction        | Protocol | Header | Length Constraints | Data         | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
|------------------|:--------:|--------|--------------------|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Client -> Server |   TCP    | Yes    | >=5                | `VC<SEMVER>` | The client sends its version to the server.                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| Server -> Client |   TCP    | No     | 1                  | `A`          | The server responds with a single byte to indicate it has accepted the client version. If the server does not accept the version, it should send a [kick](#kick) packet                                                                                                                                                                                                                                                                                                                                                |
| Client -> Server |   TCP    | Yes    | 36                 | `UUID`       | the client sends a public key to the server, which is used to obtain the user's details from the BeamMP Authentication API. The server should then present the profile data to any plugins which are requesting access via the plugin API (if one is present). Should the server, or one of its plugins decide to not authorize the connection of the client, the server should close the connection using the [kick](#kick) packet, with preference as to providing a reason to the client to prevent user confusion. |

> At this point, we come to a logical fork in the protocol. If the server has a configured `Password` in the `ServerConfig.toml`, then the server should implement the following, otherwise, it should continue from the [Mod Syncing](#mod-syncing) section.

| Direction        | Protocol | Header | Length Constraints | Data       | Description                                                       |
|------------------|:--------:|--------|--------------------|------------|-------------------------------------------------------------------|
| Server -> Client |   TCP    | Yes    | 1                  | `S`        | The server sends a single byte to indicate it has a password set. |
| Client -> Server |   TCP    | Yes    | N/A                | `PASSWORD` | The client sends the password to the server.                      |

> If the server accepts the provided password, it should continue from the [Mod Syncing](#mod-syncing) section. Otherwise, it should close the connection using the [kick](#kick) packet, sending the reason as something referencing an invalid password.

### Mod Syncing

> Please note that this section is incomplete because this server does not implement the mod syncing feature at this time. This may be added in the future.
> If you know more information about the mod syncing feature, please inform me, or open a pull request. Thanks!

| Direction        | Protocol | Header | Length Constraints | Data                                               | Description                                                                                                                                    |
|------------------|:--------:|--------|--------------------|----------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| Server -> Client |   TCP    | Yes    | >1                 | `P<player ID number>`                              | The server sends the client the ID number it has been assigned. This must be sent as a text character, or it will not be receive by the client |
| Client -> Server |   TCP    | Yes    | 2                  | `SR`                                               | The client sends the server the request to synchronize the mod list.                                                                           |
| Server -> Client |   TCP    | Yes    | N/A                | **Empty**: `-` <br/> **Modded**: Currently Unknown | The server will send the mod list to the client, or if no mods are present, it will send the first example packet.                             |

> We now reach another logical fork in the protocol. If the server sent a modlist to the client, the client will request the mods from the server, however, if the server did not send a modlist, then the client should continue from the [Map Loading](#map-loading) section.

#### Map Loading

| Direction        | Protocol | Header | Length Constraints | Data                      | Description                                                                           |
|------------------|:--------:|--------|--------------------|---------------------------|---------------------------------------------------------------------------------------|
| client -> Server |   TCP    | Yes    | 4                  | `Done`                    | The client sends the server a notification that it has finished syncing the mod list. |
| Server -> Client |   TCP    | Yes    | >11                | `M/path/to/map/info.json` | The server sends the path to the map information file to the client.                  |

> At this point, there will be a delay whilst the client loads the mods, and the map for the player.

| Direction        | Protocol | Header | Length Constraints | Data           | Description                                                                       |
|------------------|:--------:|--------|--------------------|----------------|-----------------------------------------------------------------------------------|
| Client -> Server |   TCP    | Yes    | 1                  | `H`            | The client sends the server confirmation it has loaded the map successfully.      |
| Server -> Client |   TCP    | Yes    | >2                 | `Sn<username>` | The server sends the client the username of the user account that loaded the map. |

> At this point, the client should open a UDP connection to the server. The details of this protocol are elsewhere in this document.
> You should also continually send packets to the client, the purpose of which is populating the player list. The frequency of these packets appears to be approximately once per second.

| Direction        | Protocol | Header | Length Constraints | Data                                                                      | Description                                         |
|------------------|:--------:|--------|--------------------|---------------------------------------------------------------------------|-----------------------------------------------------|
| Server -> client |   TCP    | Yes    | >2                 | `Ss<player #>/<player limit>:<player name>` <br> Eg: `Ss1/8:guest1004808` | Set the leaderboard for `player #` to `player name` |

### B

> The server is ready to accept a connection.

### D

> Not currently documented, but is used to download mods from the server.

### P

> This option is not known by me, but I believe it is a health check by the listing server.

## Miscellaneous Packets

### Kick

> The client is being disconnected from the server by force.

```md
HEADER | K<DATA>
```

For example...

```md
 | KUnauthorized
```

This packet is sent by the server to the client when the client is being disconnected by the server.
