package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const DepCheckFileName = "dep-check-v2.yaml"

type depCheckConfig struct {
	External []string `yaml:"external"`
	Utility  []string `yaml:"utility"`
	Common   []string `yaml:"common"`
	Service  []string `yaml:"service"`
}

var (
	DepCheckConfig   depCheckConfig
	DebugMode        = false
	SelectedLanguage string
)

func (d *depCheckConfig) IsEmpty() bool {
	return len(d.External) == 0 && len(d.Utility) == 0 && len(d.Common) == 0 && len(d.Service) == 0
}

func init() {
	data, err := os.ReadFile(DepCheckFileName)
	if err != nil {
		panic(fmt.Sprint("No", DepCheckFileName, "found", err))
	}
	err = yaml.Unmarshal(data, &DepCheckConfig)
	if err != nil {
		panic(fmt.Sprint(DepCheckFileName, " is not yaml parsable", err))
	}
}
