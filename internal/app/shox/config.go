package shox

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Shell string `yaml:"shell"`
	Bar   struct {
		Format  string `yaml:"format"`
		Colours struct {
			Bg string `yaml:"bg"`
			Fg string `yaml:"fg"`
		} `yaml:"colours"`
		Padding uint16 `yaml:"padding"`
	} `yaml:"bar"`
}

func LoadConfig() (*Config, error) {

	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("cannot find current user")
	}

	home := usr.HomeDir
	if home == "" {
		return nil, fmt.Errorf("user has no home directory")
	}

	places := []string{}

	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgHome != "" {
		places = append(places, filepath.Join(xdgHome, "shox/config.yaml"))
	}

	places = append(places, filepath.Join(home, ".config/shox/config.yaml"))
	places = append(places, filepath.Join(home, ".shox.yaml"))

	for _, place := range places {
		if data, err := ioutil.ReadFile(place); err == nil {
			config := Config{}
			err := yaml.Unmarshal([]byte(data), &config)
			return &config, err
		}
	}

	// config doesn't exist
	return nil, nil
}
