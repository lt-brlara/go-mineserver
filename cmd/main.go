package main

import (
	"log"

	"github.com/blara/go-mineserver/internal/server"
)

func main() {
	log.Fatal(server.Run())
}
