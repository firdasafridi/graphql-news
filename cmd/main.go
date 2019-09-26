package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/firdasafridi/graphql-news/service/graphserver"
)

func main() {
	graphServer, err := graphserver.NewGraphServer()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		errServer := graphServer.Start()
		if errServer != nil {
			log.Fatal(errServer)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	graphServer.Stop()
}
