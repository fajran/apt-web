package aptweb

import (
	"errors"
	"strings"
)

type Dist struct {
	Name string
	Path string
	Arch string
}

type Action struct {
	Dist *Dist
	Apt  *Apt
}

func NewAction(dist *Dist, apt *Apt) *Action {
	return &Action{
		Dist: dist,
		Apt:  apt,
	}
}

func (a *Action) Install(packages []string) (*InstallInfo, error) {
	names := SanitizePackages(packages)
	if len(names) == 0 {
		return nil, errors.New("Invalid package names")
	}

	args := make([]string, 0)
	args = append(args, "-c=apt.conf", "-y", "--print-uris", "install")
	args = append(args, names...)

	out, err := a.Apt.Get(a.Dist, args)
	if err != nil {
		return nil, err
	}

	r := strings.NewReader(out)
	ii := ParseInstall(r)

	return ii, nil
}

func (a *Action) Show(pkg string) (DetailInfo, error) {
	name := SanitizePackage(pkg)
	if len(name) == 0 {
		return nil, errors.New("Invalid package name")
	}

	args := make([]string, 0)
	args = append(args, "-c=apt.conf", "show", name)

	out, err := a.Apt.Cache(a.Dist, args)
	if err != nil {
		return nil, err
	}

	r := strings.NewReader(out)
	di := ParseDetail(r)

	return di, nil
}
