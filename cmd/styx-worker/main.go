package main

import (
	"github.com/fetzi/styx/worker"
)

func main() {

	worker := worker.NewQueueWorker("test")

	worker.Start()
}
