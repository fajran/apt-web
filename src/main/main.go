package main

import (
	"aptweb"
	"aptweb/server"
	"os"
)

func main() {
	f, _ := os.Open("config.json")
	aptWebConfig, _ := aptweb.NewConfigFromJson(f)
	serverConfig := &server.Config{
		Address:      ":8080",
		DocumentRoot: "www/",
	}

	s := server.NewServer(aptWebConfig, serverConfig)
	s.ListenAndServe()
}
