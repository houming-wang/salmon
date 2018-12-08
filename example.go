package main

import (
	"log"
	"sync"
	"time"

	"salmon/pkg/election"

	"github.com/pborman/uuid"
)

func main() {
	name := uuid.New()
	e, err := election.NewElection(name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("I'm:", name)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.Start()
	}()

	time.Sleep(120 * time.Second)
	e.Stop()
	wg.Wait()
}
