package aptweb

import (
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

type Action struct {
	Dist Dist
	Apt  *Apt
}

func NewAction(dist Dist, apt *Apt) *Action {
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

func GetReturnCode(err error) (int, error) {
	if err == nil {
		return 0, nil
	}

	if xerr, ok := err.(*exec.ExitError); ok {
		s := xerr.Sys().(syscall.WaitStatus).ExitStatus()
		return s, nil
	} else {
		return 0, err
	}
}

func (a *Action) Show(pkg string) (DetailInfo, error) {
	name := SanitizePackage(pkg)
	if len(name) == 0 {
		return nil, errors.New("Invalid package name")
	}

	args := make([]string, 0)
	args = append(args, "-c=apt.conf", "show", name)

	out, err := a.Apt.Cache(a.Dist, args)
	if ret, _ := GetReturnCode(err); ret == 100 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	r := strings.NewReader(out)
	di := ParseDetail(r)

	return di, nil
}
