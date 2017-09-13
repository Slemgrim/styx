package main

import (
	"log"

	"github.com/Slemgrim/styx"
	"github.com/Slemgrim/styx/config"
	"github.com/Slemgrim/styx/worker"
)

func main() {
	config, err := config.ReadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	styx := styx.New(config)
	defer styx.Close()

	worker := worker.New(&styx)
	worker.Start()
	defer worker.Stop()
}
