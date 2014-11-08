package server

import (
	"aptweb"
	"net/http"

	"aptweb/server/v1"
)

type Handler struct {
	aptWebConfig *aptweb.Config
	config       *Config

	apt *aptweb.Apt
}

func NewHandler(aptWebConfig *aptweb.Config, config *Config) http.Handler {
	h := v1.NewHandler(aptWebConfig)

	m := http.NewServeMux()
	m.Handle(h.GetPathPrefix(), h)
	m.Handle("/", http.FileServer(http.Dir(config.DocumentRoot)))

	return m
}
