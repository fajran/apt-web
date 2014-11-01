package server

import (
	"aptweb"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const API_PATH_PREFIX = "/_/"

type Handler struct {
	aptWebConfig *aptweb.Config
	config       *Config

	staticHandler http.Handler
}

func NewHandler(aptWebConfig *aptweb.Config, config *Config) *Handler {
	h := &Handler{
		aptWebConfig: aptWebConfig,
		config:       config,
	}

	h.staticHandler = http.FileServer(http.Dir(config.DocumentRoot))

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p := strings.TrimPrefix(r.URL.Path, API_PATH_PREFIX); len(p) < len(r.URL.Path) {
		r.URL.Path = p
		h.HandleAPI(w, r)
	} else {
		h.staticHandler.ServeHTTP(w, r)
	}
}

func (h *Handler) HandleAPI(w http.ResponseWriter, r *http.Request) {
	// TODO
	io.WriteString(w, fmt.Sprintf("API Path: %s\n", r.URL.Path))
}

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
