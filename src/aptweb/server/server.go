package server

import (
	"aptweb"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const API_PATH_PREFIX = "/_/"

type Handler struct {
	aptWebConfig *aptweb.Config
	config       *Config

	apt *aptweb.Apt

	staticHandler http.Handler
}

func NewHandler(aptWebConfig *aptweb.Config, config *Config) *Handler {
	h := &Handler{
		aptWebConfig: aptWebConfig,
		config:       config,
	}

	h.apt = aptweb.NewApt(aptWebConfig)
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
	p := r.URL.Path
	if p == "description" && r.Method == "GET" {
		h.HandleDescription(w, r)
	} else if p == "dependencies" && r.Method == "GET" {
		h.HandleDependencies(w, r)
	} else if p == "info" && r.Method == "GET" {
		h.HandleInfo(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *Handler) HandleDescription(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	if len(qs["pkg"]) == 0 || len(qs["d"]) == 0 {
		BadRequest(w, r)
		return
	}

	pkg := strings.TrimSpace(qs["pkg"][0])
	d, err := strconv.ParseInt(qs["d"][0], 10, 32)
	if err != nil || len(pkg) == 0 || d < 0 {
		BadRequest(w, r)
		return
	}

	h.ShowDescription(w, r, int(d), pkg)
}

func (h *Handler) ShowDescription(w http.ResponseWriter, r *http.Request, d int, pkg string) {
	if len(h.aptWebConfig.DistList) <= d {
		http.NotFound(w, r)
		return
	}

	dist := h.aptWebConfig.DistList[d]
	a := aptweb.NewAction(dist, h.apt)
	di, err := a.Show(pkg)
	if err != nil {
		log.Printf("Error showing description: pkg=%s dist=%s error: %v", pkg, dist.Path, err)
		InternalServerError(w, r)
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

func unique(list []string) []string {
	m := make(map[string]bool)
	for _, s := range list {
		m[s] = true
	}

	res := make([]string, 0)
	for k, _ := range m {
		res = append(res, k)
	}
	return res
}

func splitPackages(list string) []string {
	re := regexp.MustCompile("\\s+")
	s := re.Split(strings.TrimSpace(list), -1)
	pkgs := unique(s)
	sort.Strings(pkgs)
	return pkgs
}

func (h *Handler) HandleDependencies(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	if len(qs["pkgs"]) == 0 || len(qs["d"]) == 0 {
		BadRequest(w, r)
		return
	}

	pkgs := splitPackages(qs["pkgs"][0])

	d, err := strconv.ParseInt(qs["d"][0], 10, 32)
	if err != nil || len(pkgs) == 0 || d < 0 {
		BadRequest(w, r)
		return
	}

	h.ShowDependencies(w, r, int(d), pkgs)
}

func (h *Handler) ShowDependencies(w http.ResponseWriter, r *http.Request, d int, pkgs []string) {
	if len(h.aptWebConfig.DistList) <= d {
		http.NotFound(w, r)
		return
	}

	dist := h.aptWebConfig.DistList[d]
	a := aptweb.NewAction(dist, h.apt)
	ii, err := a.Install(pkgs)
	if err != nil {
		log.Printf("Error showing dependencies: dist=%s dist=%s error: %v", dist.Path, pkgs, err)
		InternalServerError(w, r)
		return
	}

	if ii == nil {
		http.NotFound(w, r)
		return
	}

	var data struct {
		Packages    []string `json:"pkgs"`
		Deps        []string `json:"deps"`
		Suggested   []string `json:"suggested"`
		Recommended []string `json:"recommended"`
	}

	data.Packages = pkgs
	data.Suggested = ii.Packages[aptweb.GROUP_SUGGESTED]
	data.Recommended = ii.Packages[aptweb.GROUP_RECOMMENDED]

	baseUrl := strings.TrimRight(h.aptWebConfig.RepoBaseUrl, "/")
	for _, u := range ii.Urls {
		data.Deps = append(data.Deps, strings.TrimPrefix(u.Url, baseUrl))
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	e := json.NewEncoder(w)
	e.Encode(data)
}

func (h *Handler) HandleInfo(w http.ResponseWriter, r *http.Request) {
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
