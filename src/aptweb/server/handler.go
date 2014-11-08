package server

import (
	"aptweb"
	"net/http"
	"strings"
)

const API_PATH_PREFIX = "/_/"

type Handler struct {
	aptWebConfig *aptweb.Config
	config       *Config

	apt *aptweb.Apt
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func NewHandler(aptWebConfig *aptweb.Config, config *Config) http.Handler {
	h := &Handler{
		aptWebConfig: aptWebConfig,
		config:       config,
	}
	h.apt = aptweb.NewApt(aptWebConfig)

	m := http.NewServeMux()
	m.HandleFunc(API_PATH_PREFIX, stripPath(API_PATH_PREFIX, h.HandleAPI))
	m.Handle("/", http.FileServer(http.Dir(config.DocumentRoot)))

	return m
}

func stripPath(prefix string, h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
			r.URL.Path = p
			h(w, r)
		} else {
			// Should not happen
			BadRequest(w, r)
		}
	}
}
