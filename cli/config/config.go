package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	DefaultPathFile = "./config.yaml"
)

type Config struct {
	Containername string `yaml:"containername"`
	Image         string `yaml:"image"`
	Subnet        string `yaml:"subnet"`
	Replicas      int32  `yaml:"replicas"`
}

func (C *Config) Goyaml(pathfile string) *Config {
	// read the output.yaml file
	data, err := os.ReadFile(pathfile)
	if err != nil {
		if pathfile == "" {
			return C.Goyaml(DefaultPathFile)
		}
		log.Printf(" marshal error maybe your path config is invalid  path : %v", err)
	}
	var config Config
	databayte := []byte(data)
	err = yaml.Unmarshal([]byte(databayte), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &config
}
