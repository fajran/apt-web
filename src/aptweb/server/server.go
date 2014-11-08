package server

import (
	"net/http"
	"time"

	"aptweb"
)

func NewServer(aptWebConfig *aptweb.Config, config *Config) *http.Server {
	h := NewHandler(aptWebConfig, config)
	s := &http.Server{
		Addr:           config.Address,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}
