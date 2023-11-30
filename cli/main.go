package main

import (
	"fmt"
)

func main() {
	var config Config
	config.Goyaml("./config.yaml")
	fmt.Printf("CN = %+v\n", config.Goyaml("./config.yaml"))

}
