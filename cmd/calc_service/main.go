package main

import (
	"github.com/InsafMin/go-web-calculator/internal/db"
	"github.com/InsafMin/go-web-calculator/internal/orchestrator/server"
)

func main() {
	db.InitDB("calc.db")

	go server.StartHTTPServer()
	go server.StartGRPCServer()

	select {}
}
