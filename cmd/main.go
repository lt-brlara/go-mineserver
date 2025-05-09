package main

import (
	"log"

	"github.com/blara/go-mineserver/internal/server"
	"github.com/hashicorp/go-hclog"
)

func main() {
	hclog.Default().SetLevel(hclog.Trace)

	s := server.NewServer(":25565")
	log.Fatal(s.Start())
}
