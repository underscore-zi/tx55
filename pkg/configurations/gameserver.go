package configurations

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type MetalGearOnline1 struct {
	Host string
	Port uint16
	// LobbyID is the id in the lobbies table. It doesn't use any content from there but it will update the players count
	LobbyID  uint
	Database DatabaseConfig
	LogLevel LogLevelOptions
}

func LoadTOML(filename string, v interface{}) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return toml.Unmarshal(buf, v)
}
