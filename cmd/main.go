package main

import (
	"log"

	"github.com/blara/go-mineserver/internal/server"
)

func main() {
	log.Println("hello world!")

	s := server.New()

	log.Fatal(s.Run())
}
