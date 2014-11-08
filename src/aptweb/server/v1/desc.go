package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"aptweb"
	"aptweb/server/util"
)

func (h *Handler) handleDescription(w http.ResponseWriter, r *http.Request) {
	d, err := getDistIndex(r)
	if err != nil {
		util.BadRequest(w, r)
		return
	}
	pkg, err := getPackage(r)
	if err != nil {
		util.BadRequest(w, r)
		return
	}

	if len(h.aptWebConfig.DistList) <= d {
		http.NotFound(w, r)
		return
	}

	dist := h.aptWebConfig.DistList[d]
	a := aptweb.NewAction(dist, h.apt)
	di, err := a.Show(pkg)
	if err != nil {
		log.Printf("Error showing description: pkg=%s dist=%s error: %v", pkg, dist.Path, err)
		util.InternalServerError(w, r)
		return
	}

	if di == nil {
		http.NotFound(w, r)
		return
	}

	var desc struct {
		Pkg         string            `json:"pkg"`
		Description map[string]string `json:"description"`
	}

	desc.Pkg = pkg
	desc.Description = di

	w.Header()["Content-Type"] = []string{"application/json"}
	e := json.NewEncoder(w)
	e.Encode(desc)
}
