package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type DetailInfo map[string]string

func ParseDetail(inp io.Reader) DetailInfo {
	var key string

	data := make(DetailInfo)

	s := bufio.NewScanner(inp)
	for s.Scan() {
		line := s.Text()
		if len(strings.Trim(line, " ")) == 0 {
			continue
		}

		if strings.HasPrefix(line, " ") {
			data[key] = fmt.Sprintf("%s\n%s", data[key], strings.Trim(line, " "))

		} else {
			p := strings.SplitN(line, ":", 2)
			key = p[0]
			data[key] = strings.Trim(p[1], " ")
		}
	}

	return data
}

type PackageUrl struct {
	Url  string
	Name string
	Size string
	Hash string
}

type InstallInfo struct {
	Packages map[int][]string
	Urls     []PackageUrl
}

const (
	GROUP_NONE        = 0
	GROUP_EXTRA       = 1
	GROUP_SUGGESTED   = 2
	GROUP_RECOMMENDED = 3
	GROUP_INSTALL     = 4
	GROUP_UPGRADE     = 5
)

func ParseInstall(inp io.Reader) *InstallInfo {
	s := bufio.NewScanner(inp)

	ii := &InstallInfo{}
	ii.Packages = make(map[int][]string)
	ii.Urls = make([]PackageUrl, 0)

	group := GROUP_NONE
	more := s.Scan()
	for more {
		line := s.Text()
		if len(strings.Trim(line, " ")) == 0 {
			continue
		}

		if strings.Contains(line, "The following extra packages will be installed:") {
			group = GROUP_EXTRA

		} else if strings.Contains(line, "Suggested packages:") {
			group = GROUP_SUGGESTED

		} else if strings.Contains(line, "Recommended packages:") {
			group = GROUP_RECOMMENDED

		} else if strings.Contains(line, "The following NEW packages will be installed:") {
			group = GROUP_INSTALL

		} else if strings.Contains(line, "The following packages will be upgraded:") {
			group = GROUP_UPGRADE

		} else if strings.Contains(line, "After this operation") {
			more = s.Scan()
			for more {
				line = s.Text()
				if len(strings.Trim(line, " ")) == 0 {
					more = s.Scan()
					break
				}

				p := strings.Split(line, " ")
				pu := PackageUrl{
					Url:  strings.Trim(p[0], "'"),
					Name: p[1],
					Size: p[2],
					Hash: p[3],
				}
				ii.Urls = append(ii.Urls, pu)

				more = s.Scan()
			}

		} else {
			more = s.Scan()
			group = GROUP_NONE
		}

		if group != GROUP_NONE {
			names := make([]string, 0)
			for {
				more = s.Scan()
				line = s.Text()
				if !strings.HasPrefix(line, " ") {
					break
				}

				p := strings.Split(strings.Trim(line, " "), " ")
				names = append(names, p...)
			}
			ii.Packages[group] = names
			group = GROUP_NONE
		}

	}

	return ii
}
