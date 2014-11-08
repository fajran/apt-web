package v1

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) handleInfo(w http.ResponseWriter, r *http.Request) {
	type RepoInfo struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	type DistInfo struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	var info struct {
		Repos []RepoInfo `json:"repos"`
		Dists []DistInfo `json:"dists"`
	}

	for _, r := range h.aptWebConfig.RepoList {
		repo := RepoInfo{r.Name, strings.TrimRight(r.Url, "/")}
		info.Repos = append(info.Repos, repo)
	}

	for index, dist := range h.aptWebConfig.DistList {
		di := DistInfo{
			Id:   index,
			Name: dist.Name,
		}
		info.Dists = append(info.Dists, di)
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	e := json.NewEncoder(w)
	e.Encode(info)
}
