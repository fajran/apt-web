package v1

import (
	"net/http"

	"aptweb"
)

type Handler struct {
	aptWebConfig *aptweb.Config

	apt *aptweb.Apt
}

func NewHandler(aptWebConfig *aptweb.Config) *Handler {
	h := &Handler{
		aptWebConfig: aptWebConfig,
	}
	h.apt = aptweb.NewApt(aptWebConfig)
	return h
}

func (h *Handler) GetPathPrefix() string {
	return "/api/v1/"
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	if p == "/api/v1/desc" && r.Method == "GET" {
		h.handleDescription(w, r)

	} else if p == "/api/v1/deps" && r.Method == "GET" {
		h.handleDependencies(w, r)

	} else if p == "/api/v1/info" && r.Method == "GET" {
		h.handleInfo(w, r)

	} else {
		http.NotFound(w, r)
	}
}
