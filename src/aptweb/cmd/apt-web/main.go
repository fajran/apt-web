package main

import (
	"fmt"
	"os"

	"aptweb"
	"aptweb/server"
)

func main() {
	address := ":8080"

	f, _ := os.Open("config.json")
	aptWebConfig, _ := aptweb.NewConfigFromJson(f)
	serverConfig := &server.Config{
		Address:      address,
		DocumentRoot: "www/",
	}

	s := server.NewServer(aptWebConfig, serverConfig)

	fmt.Printf("Starting apt-web server. Listening to %s\n", address)
	s.ListenAndServe()
}
