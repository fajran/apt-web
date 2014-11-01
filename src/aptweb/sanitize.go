package main

import (
	"regexp"
)

func SanitizePackages(names []string) []string {
	res := make([]string, 0)
	for _, name := range names {
		name = SanitizePackage(name)
		if len(name) > 0 {
			res = append(res, name)
		}
	}
	return res
}

func SanitizePackage(name string) string {
	re := regexp.MustCompile("[^a-z0-9.+-]")
	return re.ReplaceAllString(name, "")
}
