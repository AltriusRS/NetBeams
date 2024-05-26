package config

import "time"

// BaseConfig is the main config struct for the server
type BaseConfig struct {
	// General server settings
	General GeneralConfig `toml:"General"  comment:"This is the BeamMP-Server config file.\nHelp & Documentation: 'https://docs.beammp.com/server/server-maintenance/'\nIMPORTANT: Fill in the AuthKey with the key you got from 'https://keymaster.beammp.com/' on the left under 'Keys'\n This section stores general server settings"`

	// Miscellaneous server settings
	Misc MiscConfig `toml:"Misc"     comment:"This is for miscellaneous server settings provided by the official implementation of BeamMP"`

	// NetBeams specific settings
	NetBeams NetBeamsConfig `toml:"NetBeams" comment:"Configuration options specific to the NetBeams project"`

	// Authentication settings (also NetBeams specific, but seperated for organizational purposes)
	Auth AuthenticationConfig `toml:"Auth"     comment:"These are the authentication settings for the NetBeams project, which allows you to define which players are allowed to join the server"`
}

// GeneralConfig is the general server settings
type GeneralConfig struct {
	// Name of the server
	Name string `toml:"Name" comment:"The name of your server. See https://wiki.beammp.com/en/home/server-maintenance#customize-the-look-of-your-server-name for more information"`

	// Port to listen on
	Port int `toml:"Port" comment:"The port to use for connections. This should be a port that is not in use by any other application. Make sure to port forward to this port on your router."`

	// Authentication key
	AuthKey string `toml:"AuthKey" comment:"This is the authentication key for your server. Please use the BeamMP Keymaster to generate a key, see https://keymaster.beammp.com/"`

	// Whether to log chat messages in the console / log
	LogChat bool `toml:"LogChat" comment:"If the output log should contain chat messages"`

	// Add custom identifying tags to your server to make it easier to find. Format should be TagA,TagB,TagC. Note the comma seperation.
	Tags string `toml:"Tags" comment:""`

	// Whether to log debug messages
	Debug bool `toml:"Debug" comment:""`

	// Whether the server is private or public
	Private bool `toml:"Private" comment:""`

	// Maximum number of cars on the server
	MaxCars int `toml:"MaxCars" comment:""`

	// Maximum number of players on the server
	MaxPlayers int `toml:"MaxPlayers" comment:""`

	// Map to use
	Map string `toml:"Map" comment:""`

	// Description of the server
	Description string `toml:"Description" comment:""`

	// Folder to load resources from
	ResourceFolder string `toml:"ResourceFolder" comment:""`

	// Password to use for the server
	Password string `toml:"Password" comment:""`
}

// MiscConfig is the miscellaneous server settings
type MiscConfig struct {

	// Hides the periodic update message which notifies you of a new server version. You should really keep this on and always update as soon as possible. For more information visit https://wiki.beammp.com/en/home/server-maintenance#updating-the-server. An update message will always appear at startup regardless.
	ImScaredOfUpdates bool `toml:"ImScaredOfUpdates"`

	// You can turn on/off the SendErrors message you get on startup here
	SendErrorsShowMessage bool `toml:"SendErrorsShowMessage"`

	// If SendErrors is `true`, the server will send helpful info about crashes and other issues back to the BeamMP developers. This info may include your config, who is on your server at the time of the error, and similar general information. This kind of data is vital in helping us diagnose and fix issues faster. This has no impact on server performance. You can opt-out of this system by setting this to `false`
	SendErrors bool `toml:"SendErrors"`
}

// NetBeamsConfig is the settings specific to NetBeams
type NetBeamsConfig struct {

	// The IP address of the master node
	MasterNode string `toml:"MasterNode"`

	// The port of the master node
	MasterPort int `toml:"MasterPort"`

	// The log level of the server
	LogLevel string `toml:"LogLevel"`

	// The log file of the server
	LogFile string `toml:"LogFile"`

	// The IP address of the mod server (This can be grabbed from the master node at startup)
	ModServer string `toml:"ModServer"`
}

// AuthenticationConfig is the authentication settings specific to NetBeams
type AuthenticationConfig struct {

	// Whether guests are allowed to join the server
	AllowGuests bool `toml:"AllowGuests" comment:"Whether guest accounts are allowed to join the server\nThis will automatically prevent guest accounts from joining the server in the authentication step."`

	// Whether content creators are allowed to join the server
	AllowContentCreators bool `toml:"AllowContentCreators" comment:"Whether content creators are allowed to join the server\nThis will automatically prevent content creators from joining the server in the authentication step."`

	// Whether BeamMP staff are allowed to join the server
	AllowStaff bool `toml:"AllowStaff" comment:"Whether BeamMP staff are allowed to join the server\nThis will automatically prevent BeamMP staff from joining the server in the authentication step."`

	// The minimum age of an account to be able to join the server
	// MinimumAccountAge string `toml:"MinimumAccountAge" comment:"The minimum age of an account to be able to join the server\n Leave empty to disable\n Currently does not work, as the server does not know the age of the account"`

	// Idle player detection settings
	Idle AuthIdleConfig `toml:"Idle" comment:"Idle player detection settings"`

	// Online player detection settings
	Online AuthOnlineConfig `toml:"Online" comment:"Online player detection settings"`

	// VPN detection settings
	VPN AuthVPNConfig `toml:"VPN" comment:"VPN detection settings"`

	// Proxy detection settings
	Proxy AuthProxyConfig `toml:"Proxy" comment:"Proxy detection settings"`

	// Kick player detection settings
	Kick AuthKickConfig `toml:"Kick" comment:"Kick player detection settings"`

	// Admin player detection settings
	Admin AuthAdminConfig `toml:"Admin" comment:"Admin player detection settings"`
}

type AuthIdleConfig struct {

	// Whether idle player detection is enabled
	Enable bool `toml:"Enable" comment:"Whether idle player detection is enabled"`

	// The maximum amount of time a player is allowed to be idle before being kicked (in minutes)
	MaxTime string `toml:"MaxTime" comment:"The maximum amount of time a player is allowed to be idle before being kicked (in minutes)\n Set to -1 to disable"`

	// The max time in Go Time format
	MaxTimeTime time.Duration

	// The minimum distance a player must have moved to not be considered idle
	MinDistance int `toml:"MinDistance" comment:"The minimum distance a player must have moved to not be considered idle"`
}

type AuthKickConfig struct {

	// The minimum amount of time a player is prevented from joining the server after being kicked by an admin (in seconds)
	AdminDuration string `toml:"MinDuration" comment:"The minimum amount of time a player is prevented from joining the server after being kicked by an admin (in seconds)\n Set to -1 to disable"`

	// The admin duration time in Go Time format
	AdminDurationTime time.Duration

	// The amount of time a player is prevented from joining the server after being kicked for being idle (in seconds)
	IdleDuration string `toml:"IdleDuration" comment:"The amount of time a player is prevented from joining the server after being kicked for being idle (in seconds)\n Set to -1 to disable"`

	//	The idle duration time in Go Time format
	IdleDurationTime time.Duration

	// The amount of time a player is prevented from joining the server after being kicked for reaching their online time quota limit (in seconds)
	OnlineDuration string `toml:"OnlineDuration" comment:"The amount of time a player is prevented from joining the server after being kicked for reaching their online time quota limit (in seconds)\n Set to -1 to disable"`

	// The online duration time in Go Time format
	OnlineDurationTime time.Duration
}

type AuthOnlineConfig struct {

	// Whether players should be kicked for reaching their quota limit
	Enable bool `toml:"Enable" comment:"Whether players should be kicked for reaching their quota limit"`

	// The maximum amount of time a player is allowed to be on the server before being kicked (in minutes)
	Quota string `toml:"Quota" comment:"The maximum amount of time a player is allowed to be on the server before being kicked (in minutes)\n Set to -1 to disable"`

	// The quota time in Go Time format
	QuotaTime time.Duration
}

type AuthVPNConfig struct {

	// Whether VPN detection is enabled
	Enable bool `toml:"Enable" comment:"Whether VPN detection is enabled"`

	// The default behaviour for VPN connections
	DefaultBehaviour string `toml:"DefaultBehaviour" comment:"The default behaviour for VPN connections\n Valid values are 'Allow', 'Block', 'Kick', 'Ban'"`

	// The ACL for VPN providers
	ACL map[string]bool `toml:"ACL" comment:"A list of VPN providers and whether they are allowed to join the server"`
}

type AuthProxyConfig struct {

	// Whether proxy detection is enabled
	Enable bool `toml:"Enable" comment:"Whether proxy detection is enabled"`

	// The default behaviour for proxy connections
	DefaultBehaviour string `toml:"DefaultBehaviour" comment:"The default behaviour for proxy connections\n Valid values are 'Allow', 'Block', 'Kick', 'Ban'"`

	// The ACL for proxy providers
	ACL map[string]bool `toml:"ACL" comment:"A list of proxy providers and whether they are allowed to join the server"`
}

// An array of permissions for each player (use only for admins)
type AuthAdminConfig map[string]PlayerPermissionsConfig

// A struct representing the permissions for each player
type PlayerPermissionsConfig struct {

	// If the player may bypass VPN connection filtering
	BypassVpn bool `toml:"BypassVpn" comment:"Whether VPN connections are allowed to join the server"`

	// If the player may bypass proxy connection filtering
	BypassProxy bool `toml:"BypassProxy" comment:"Whether proxy connections are allowed to join the server"`

	// If the player may bypass idle timeouts
	BypassIdle bool `toml:"BypassIdle" comment:"Whether idle connections are allowed to join the server"`

	// If the player may bypass online quota limits
	BypassOnline bool `toml:"BypassOnline" comment:"Whether online connections are allowed to join the server"`

	// If the player may bypass vehicle limiting
	BypassVehicles bool `toml:"BypassVehicles" comment:"Whether vehicles are allowed to join the server"`

	// If the player may hide their name in the leaderboard
	HideName bool `toml:"HideName" comment:"Whether the player's name is hidden in the leaderboard"`

	// If the player may kick other players
	KickPlayers bool `toml:"KickPlayers" comment:"Whether the user can kick other players"`

	// If the player may ban other players
	BanPlayers bool `toml:"BanPlayers" comment:"Whether the user can ban other players"`

	// If the player may mute other players
	MutePlayers bool `toml:"MutePlayers" comment:"Whether the user can mute other players"`
}

type AllowList struct {

	// A list of players that are allowed to join the server - These players will be able to join the server only if they pass all other authentication checks
	Players []string `toml:"Players" comment:"A list of players that are allowed to join the server - These players will be able to join the server only if they pass all other authentication checks"`
}

type BlockList struct {

	// A list of players that are blocked from joining the server - This is effectively a perma-ban
	Players []string `toml:"Players" comment:"A list of players that are blocked from joining the server - This is effectively a perma-ban"`
}
