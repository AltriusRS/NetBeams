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
			UseUPnP:    true,
		},
		Auth: AuthenticationConfig{
			AllowGuests:          true,
			AllowContentCreators: true,
			AllowStaff:           true,

			Idle: AuthIdleConfig{
				Enable:      true,
				MaxTime:     "10 minutes",
				MinDistance: 100,
			},

			Online: AuthOnlineConfig{
				Enable: false,
				Quota:  "2 hours",
			},

			VPN: AuthVPNConfig{
				Enable:           true,
				DefaultBehaviour: "Allow",
				ACL: map[string]bool{
					"MyVPNProvider": true,
				},
			},

			Proxy: AuthProxyConfig{
				Enable:           true,
				DefaultBehaviour: "Allow",
				ACL: map[string]bool{
					"MyProxyProvider": true,
				},
			},

			Kick: AuthKickConfig{
				AdminDuration:  "1 hour",
				IdleDuration:   "5 minutes",
				OnlineDuration: "30 minutes",
			},
		},
	}
}
