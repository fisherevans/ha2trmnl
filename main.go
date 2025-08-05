package main

import (
	"fmt"
	"os"

	"fisherevans.com/ha2trmnl/pkg"
	"gopkg.in/yaml.v3"
)

func main() {
	file := "config.yaml"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	fmt.Println("Config file: " + file)
	configContents, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var instance pkg.Instance
	if err = yaml.Unmarshal(configContents, &instance); err != nil {
		panic(err)
	}

	instance.Run()
}
