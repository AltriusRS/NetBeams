package config

func LoadDefault() BaseConfig {
	return BaseConfig{
		General: GeneralConfig{
			Name:           "NetBeams Server",
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
			MasterNode: "localhost",
			MasterPort: 30815,
			LogLevel:   "info",
			LogFile:    "/logs/netbeams.log",
			ModServer:  "",
		},
		Auth: AuthenticationConfig{
			AllowGuests:          true,
			AllowContentCreators: true,
			AllowStaff:           true,
			MinimumAccountAge:    "0",
			MaxIdleTime:          -1,
			MaxOnlineTime:        -1,
			DefaultKickDuration:  -1,
		},
	}
}
