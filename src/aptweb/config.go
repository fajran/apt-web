package aptweb

import (
	"encoding/json"
	"errors"
	"io"
)

type Config struct {
	AptGetPath   string `json:"apt-get"`
	AptCachePath string `json:"apt-cache"`
	DistDir      string `json:"dist-dir"`
	DistList     []Dist `json:"dists"`
}

type Dist struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Arch string `json:"arch"`
}

func NewConfigFromJson(r io.Reader) (*Config, error) {
	config := &Config{}

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	if len(config.AptGetPath) == 0 || len(config.AptCachePath) == 0 || len(config.DistDir) == 0 {
		return nil, errors.New("Incomplete config")
	}

	return config, nil
}
