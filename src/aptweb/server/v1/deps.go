package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"aptweb"
	"aptweb/server/util"
)

func (h *Handler) handleDependencies(w http.ResponseWriter, r *http.Request) {
	d, err := getDistIndex(r)
	if err != nil {
		util.BadRequest(w, r)
		return
	}
	pkgs, err := getPackages(r)
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
	ii, err := a.Install(pkgs)
	if err != nil {
		log.Printf("Error showing dependencies: dist=%s dist=%s error: %v", dist.Path, pkgs, err)
		util.InternalServerError(w, r)
		return
	}

	if ii == nil {
		http.NotFound(w, r)
		return
	}

	var data struct {
		Packages    []string `json:"pkgs"`
		Urls        []string `json:"urls"`
		Suggested   []string `json:"suggested"`
		Recommended []string `json:"recommended"`
		Install     []string `json:"install"`
	}

	data.Packages = pkgs
	data.Suggested = ii.Packages[aptweb.GROUP_SUGGESTED]
	data.Recommended = ii.Packages[aptweb.GROUP_RECOMMENDED]
	data.Install = ii.Packages[aptweb.GROUP_INSTALL]

	baseUrl := strings.TrimRight(h.aptWebConfig.RepoBaseUrl, "/")
	for _, u := range ii.Urls {
		data.Urls = append(data.Urls, strings.TrimPrefix(u.Url, baseUrl))
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	e := json.NewEncoder(w)
	e.Encode(data)
}
