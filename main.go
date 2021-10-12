package main

import (
	"log"

	"github.com/suutaku/command-server/pkg/server"
)

func main() {
	log.SetFlags(log.Lshortfile)
	server := server.NewServer("5000")
	server.InitServer()
	server.Serve()
}
