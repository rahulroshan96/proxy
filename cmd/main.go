package main

import "github.com/rahulroshan96/proxy/server"

func main() {
	server := server.NewServer()
	server.Run()
}
