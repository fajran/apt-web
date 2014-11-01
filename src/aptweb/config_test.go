package aptweb

import (
	"strings"
	"testing"
)

func assert(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected: %v, but received: %v", expected, actual)
	}
}

func TestNewConfigFromJson(t *testing.T) {
	data := `
{
  "apt-get": "/usr/bin/apt-get",
  "apt-cache": "/usr/bin/apt-cache",

  "dist-dir": "virtuals/",

  "dists": [
    {"name": "Ubuntu 14.04 \"Trusty Tahr\" Desktop amd64",
     "path": "ubuntu-14.04-desktop-amd64",
     "arch": "amd64"},
    {"name": "Ubuntu 14.10 \"Utopic Unicorn\" Desktop i386",
     "path": "ubuntu-14.10-desktop-i386",
     "arch": "i386"}
  ]
}
`

	r := strings.NewReader(data)
	config, err := NewConfigFromJson(r)

	if err != nil {
		t.Errorf("Error reading config data: %v", err)
		t.FailNow()
	}

	assert(t, "/usr/bin/apt-get", config.AptGetPath)
	assert(t, "/usr/bin/apt-cache", config.AptCachePath)
	assert(t, "virtuals/", config.DistDir)
	assert(t, 2, len(config.DistList))

	dist := config.DistList[0]
	assert(t, `Ubuntu 14.04 "Trusty Tahr" Desktop amd64`, dist.Name)
	assert(t, `ubuntu-14.04-desktop-amd64`, dist.Path)
	assert(t, "amd64", dist.Arch)

	dist = config.DistList[1]
	assert(t, `Ubuntu 14.10 "Utopic Unicorn" Desktop i386`, dist.Name)
	assert(t, `ubuntu-14.10-desktop-i386`, dist.Path)
	assert(t, "i386", dist.Arch)
}

func TestNewConfigFromJson_Incomplete(t *testing.T) {
	data := `
{
  "apt-get": "/usr/bin/apt-get",

  "dist-dir": "virtuals/",

  "dists": [
    {"name": "Ubuntu 14.04 \"Trusty Tahr\" Desktop amd64",
     "path": "ubuntu-14.04-desktop-amd64",
     "arch": "amd64"},
    {"name": "Ubuntu 14.10 \"Utopic Unicorn\" Desktop i386",
     "path": "ubuntu-14.10-desktop-i386",
     "arch": "i386"}
  ]
}
`

	r := strings.NewReader(data)
	config, err := NewConfigFromJson(r)
	if err == nil || config != nil {
		t.Errorf("Configuration parser should fail")
	}
}
