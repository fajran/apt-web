package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ParseDetail(inp io.Reader) map[string]string {
	var key string

	data := make(map[string]string)

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
