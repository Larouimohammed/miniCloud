package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Containername string `json:"containername"`
	Image         string `json:"image"`
	Subnet        string `json:"subnet"`
	Nunofvms      string `json:"numofvms"`
}

func Goyaml(pathfile string) *Config {
	// read the output.yaml file
	data, err := os.ReadFile(pathfile)

	if err != nil {
		panic(err)
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
