package main

import (
	srv "server-template/internal/server"
)

func main() {
	port := srv.GetPortOrDefault(8081)
	server := srv.NewServer(srv.WithPort(port))
	server.Start()
}
