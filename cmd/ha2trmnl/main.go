package main

import (
	"log"
	"os"

	"fisherevans.com/ha2trmnl/internal/config"
	"fisherevans.com/ha2trmnl/internal/runner"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[ha2trmnl] ")
}

func main() {
	file := "./config.yaml"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	c, err := config.FromFile(file)
	if err != nil {
		panic(err)
	}

	err = runner.Run(c)
	if err != nil {
		panic(err)
	}
}
