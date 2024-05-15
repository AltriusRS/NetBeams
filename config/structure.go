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
}

// MiscConfig is the miscellaneous server settings
type MiscConfig struct {
	ImScaredOfUpdates     bool `toml:"ImScaredOfUpdates"`     // Hides the periodic update message which notifies you of a new server version. You should really keep this on and always update as soon as possible. For more information visit https://wiki.beammp.com/en/home/server-maintenance#updating-the-server. An update message will always appear at startup regardless.
	SendErrorsShowMessage bool `toml:"SendErrorsShowMessage"` // You can turn on/off the SendErrors message you get on startup here
	SendErrors            bool `toml:"SendErrors"`            // If SendErrors is `true`, the server will send helpful info about crashes and other issues back to the BeamMP developers. This info may include your config, who is on your server at the time of the error, and similar general information. This kind of data is vital in helping us diagnose and fix issues faster. This has no impact on server performance. You can opt-out of this system by setting this to `false`
}

// NetBeamsConfig is the settings specific to NetBeams
type NetBeamsConfig struct {
	MasterNode string `toml:"MasterNode"`
	MasterPort int    `toml:"MasterPort"`
	LogLevel   string `toml:"LogLevel"`
	LogFile    string `toml:"LogFile"`
	ModServer  string `toml:"ModServer"`
}
