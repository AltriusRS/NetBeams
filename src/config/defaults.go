package config

// type Tags string

// func (t Tags) String() string {
// 	return string(t)
// }

// func (t Tags) Set(value string) error {

func LoadDefault() BaseConfig {
	return BaseConfig{
		General: GeneralConfig{
			Name:           "^l^NetBeams Server",
			Port:           30814,
			AuthKey:        "",
			LogChat:        true,
			Tags:           "Freeroam, NetBeams",
			Debug:          true,
			Private:        true,
			MaxCars:        2,
			MaxPlayers:     10,
			Map:            "/levels/gridmap_v2/info.json",
			Description:    "A BeamMP server written in Go - https://github.com/altriusrs/NetBeams",
			ResourceFolder: "Resources",
		},
		Misc: MiscConfig{
			ImScaredOfUpdates:     true,
			SendErrorsShowMessage: true,
			SendErrors:            true,
		},
		NetBeams: NetBeamsConfig{
			MasterNode: "",
			MasterPort: 30815,
			LogLevel:   "info",
			LogFile:    "/logs/netbeams.log",
			ModServer:  "",
		},
	}
}
