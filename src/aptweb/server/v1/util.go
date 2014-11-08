package v1

import (
	"errors"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func getDistIndex(r *http.Request) (int, error) {
	qs := r.URL.Query()
	if len(qs["d"]) == 0 {
		return 0, errors.New("Invalid dist index")
	}

	d, err := strconv.ParseInt(qs["d"][0], 10, 32)
	if err != nil || d < 0 {
		return 0, errors.New("Invalid dist index")
	}

	return int(d), nil
}

func getPackage(r *http.Request) (string, error) {
	qs := r.URL.Query()
	if len(qs["pkg"]) == 0 {
		return "", errors.New("Invalid package")
	}

	pkg := strings.TrimSpace(qs["pkg"][0])
	if len(pkg) == 0 {
		return "", errors.New("Invalid package")
	}

	return pkg, nil
}

func getPackages(r *http.Request) ([]string, error) {
	qs := r.URL.Query()
	if len(qs["pkgs"]) == 0 {
		return nil, errors.New("Invalid packages")
	}

	pkgs := strings.TrimSpace(qs["pkgs"][0])
	if len(pkgs) == 0 {
		return nil, errors.New("Invalid packages")
	}

	return splitPackages(pkgs), nil
}

func splitPackages(list string) []string {
	re := regexp.MustCompile("\\s+")
	s := re.Split(strings.TrimSpace(list), -1)
	pkgs := unique(s)
	sort.Strings(pkgs)
	return pkgs
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
