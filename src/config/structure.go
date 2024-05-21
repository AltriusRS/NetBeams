package config

// BaseConfig is the main config struct for the server
type BaseConfig struct {
	General  GeneralConfig  `toml:"General"`  // General server settings
	Misc     MiscConfig     `toml:"Misc"`     // Miscellaneous server settings
	NetBeams NetBeamsConfig `toml:"NetBeams"` // Settings specific to NetBeams
}

// GeneralConfig is the general server settings
type GeneralConfig struct {
	Name           string `toml:"Name"`           // Name of the server
	Port           int    `toml:"Port"`           // Port to listen on
	AuthKey        string `toml:"AuthKey"`        // Authentication key
	LogChat        bool   `toml:"LogChat"`        // Whether to log chat messages in the console / log
	Tags           string `toml:"Tags"`           // Add custom identifying tags to your server to make it easier to find. Format should be TagA,TagB,TagC. Note the comma seperation.
	Debug          bool   `toml:"Debug"`          // Whether to log debug messages
	Private        bool   `toml:"Private"`        // Whether the server is private or public
	MaxCars        int    `toml:"MaxCars"`        // Maximum number of cars on the server
	MaxPlayers     int    `toml:"MaxPlayers"`     // Maximum number of players on the server
	Map            string `toml:"Map"`            // Map to use
	Description    string `toml:"Description"`    // Description of the server
	ResourceFolder string `toml:"ResourceFolder"` // Folder to load resources from
	Password       string `toml:"Password"`       // Password to use for the server
}

// MiscConfig is the miscellaneous server settings
type MiscConfig struct {
	ImScaredOfUpdates     bool `toml:"ImScaredOfUpdates"`     // Hides the periodic update message which notifies you of a new server version. You should really keep this on and always update as soon as possible. For more information visit https://wiki.beammp.com/en/home/server-maintenance#updating-the-server. An update message will always appear at startup regardless.
	SendErrorsShowMessage bool `toml:"SendErrorsShowMessage"` // You can turn on/off the SendErrors message you get on startup here
	SendErrors            bool `toml:"SendErrors"`            // If SendErrors is `true`, the server will send helpful info about crashes and other issues back to the BeamMP developers. This info may include your config, who is on your server at the time of the error, and similar general information. This kind of data is vital in helping us diagnose and fix issues faster. This has no impact on server performance. You can opt-out of this system by setting this to `false`
}

// NetBeamsConfig is the settings specific to NetBeams
type NetBeamsConfig struct {
	MasterNode string `toml:"MasterNode"` // The IP address of the master node
	MasterPort int    `toml:"MasterPort"` // The port of the master node
	LogLevel   string `toml:"LogLevel"`   // The log level of the server
	LogFile    string `toml:"LogFile"`    // The log file of the server
	ModServer  string `toml:"ModServer"`  // The IP address of the mod server (This can be grabbed from the master node at startup)
}

type AuthenticationConfig struct {
	AllowGuests          bool   `toml:"AllowGuests"`          // Whether guests are allowed to join the server
	AllowContentCreators bool   `toml:"AllowContentCreators"` // Whether content creators are allowed to join the server
	AllowStaff           bool   `toml:"AllowStaff"`           // Whether BeamMP staff are allowed to join the server
	MinimumAccountAge    string `toml:"MinimumAccountAge"`    // The minimum age of an account to be able to join the server
}
