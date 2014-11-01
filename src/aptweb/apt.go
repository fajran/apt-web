package aptweb

import (
	"bytes"
	"os/exec"
	"path"
)

type Apt struct {
	Config Config
}

func NewApt(config Config) *Apt {
	return &Apt{
		Config: config,
	}
}

func (apt *Apt) Get(dist *Dist, args []string) (string, error) {
	return apt.run(apt.Config.AptGetPath, dist, args)
}

func (apt *Apt) Cache(dist *Dist, args []string) (string, error) {
	return apt.run(apt.Config.AptCachePath, dist, args)
}

func (apt *Apt) run(cmdPath string, dist *Dist, args []string) (string, error) {
	cmd := exec.Command(cmdPath, args...)
	cmd.Dir = apt.getDir(dist)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func (apt *Apt) getDir(dist *Dist) string {
	return path.Join(apt.Config.DistPath, dist.Path)
}
