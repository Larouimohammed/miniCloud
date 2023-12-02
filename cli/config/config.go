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
	Nunofinstance int32  `yaml:"numofinstance"`
}

func (C *Config) Goyaml(pathfile string) *Config {
	// read the output.yaml file
	data, err := os.ReadFile(pathfile)
	if err != nil {
		if pathfile == "" {
			data, err = os.ReadFile(DefaultPathFile)
		}
		log.Panicf("err")
	}
	// fmt.Printf("members = %#v\n", string(data))
	var config Config
	databayte := []byte(data)
	// dataconfig := make([]Config, 0)
	err = yaml.Unmarshal([]byte(databayte), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// fmt.Printf("members = %#v\n", config)
	return &config
}
