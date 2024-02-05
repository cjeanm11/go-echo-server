package main

import (
	srv "server-template/internal/server"
)

func main() {
	option := srv.Option(srv.WithPort(8081));
	server := srv.NewServer(option)
	server.Start()
}
