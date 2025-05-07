package main

import (
	"context"
	"log"

	"github.com/brotich/go-log-forwarder/config"
	"github.com/brotich/go-log-forwarder/internal"
)

func main() {
	ctx := context.Background()
	srv, err := internal.NewServer(ctx, config.Config{
		ListenAddr: ":8090",
	})
	if err != nil {
		log.Fatalln("got error starting server ", err)
	}

	log.Println("started server")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln("error running server ", err)
	}
}
