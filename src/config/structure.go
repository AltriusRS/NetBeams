package config

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
	MinimumAccountAge string `toml:"MinimumAccountAge" comment:"The minimum age of an account to be able to join the server\n Leave empty to disable\n Currently does not work, as the server does not know the age of the account"`

	// The maximum amount of time a player is allowed to be idle before being kicked (in minutes)
	MaxIdleTime int `toml:"MaxIdleTime" comment:"The maximum amount of time a player is allowed to be idle before being kicked (in minutes)\n Set to -1 to disable"`

	// The maximum amount of time a player is allowed to be on the server before being kicked (in minutes)
	MaxOnlineTime int `toml:"MaxOnlineTime" comment:"The maximum amount of time a player is allowed to be on the server before being kicked (in minutes)\n Set to -1 to disable"`

	// The amount of time a player is prevented from joining the server after being kicked (in seconds)
	DefaultKickDuration int `toml:"DefaultKickDuration" comment:"The amount of time a player is prevented from joining the server after being kicked (in seconds)\n Set to -1 to disable"`
}
