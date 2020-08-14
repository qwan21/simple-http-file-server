package config

import (
	"flag"
)

//Config is General Configuration
type Config struct {
	Host      string
	Port      string
	Directory string
	Route     string
}

//Get function returns the configuration
func Get(fs *flag.FlagSet) (*Config, error) {
	cfg, err := initConfig(fs)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func initConfig(fs *flag.FlagSet) (cfg *Config, err error) {

	var mergedCfg Config

	settingsName := []string{"host", "port", "directory", "route"}

	mergedCfg.Host = fs.Lookup(settingsName[0]).Value.String()
	mergedCfg.Port = fs.Lookup(settingsName[1]).Value.String()
	mergedCfg.Directory = fs.Lookup(settingsName[2]).Value.String()
	mergedCfg.Route = fs.Lookup(settingsName[3]).Value.String()

	return &mergedCfg, nil
}
