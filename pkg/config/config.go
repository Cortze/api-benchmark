package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

type Config struct {
	BenchmarkName string   `json:"benchmark-name"`
	HostEndpoint  string   `json:"host-endpoint"`
	Query         string   `json:"query"`
	Replaces      []string `json:"replaces"`
	RangeValues   []string `json:"range-values"`
	QueryFile     string   `json:"query-file"`
	NumQueries    int      `json:"num-queries"`
	QueryBackup   string   `json:"query-backup"`
	ConcurrentReq int      `json:"concurrent-req"`
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigFromJson(jfile string) (*Config, error) {
	conf := NewConfig()

	// opend the file
	f, err := ioutil.ReadFile(jfile)
	if err != nil {
		return conf, errors.Wrap(err, fmt.Sprintf("unable to open file %s.", jfile))
	}

	// read the bytes
	err = json.Unmarshal([]byte(f), conf)
	if err != nil {
		return conf, errors.Wrap(err, fmt.Sprintf("unable to parse content of json file %s.", jfile))
	}

	return conf, nil
}
