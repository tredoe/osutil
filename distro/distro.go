// Copyright 2013 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package distro detects the Linux distribution.
package distro

import (
	"os"

	"github.com/tredoe/osutil/config/shconf"
)

// Distro represents a distribution of Linux system.
type Distro int

// Most used Linux distributions.
const (
	DistroUnknown Distro = iota
	Arch
	CentOS
	Debian
	Fedora
	Manjaro
	OpenSUSE
	Ubuntu
)

var distroNames = [...]string{
	DistroUnknown: "unknown distribution",
	Arch:          "Arch",
	CentOS:        "CentOS",
	Debian:        "Debian",
	Fedora:        "Fedora",
	Manjaro:       "Manjaro",
	OpenSUSE:      "openSUSE",
	Ubuntu:        "Ubuntu",
}

func (s Distro) String() string { return distroNames[s] }

var idToDistro = map[string]Distro{
	"arch":                Arch,
	"manjaro":             Manjaro, // based on Arch
	"centos":              CentOS,
	"debian":              Debian,
	"fedora":              Fedora,
	"opensuse-leap":       OpenSUSE,
	"opensuse-tumbleweed": OpenSUSE,
	"ubuntu":              Ubuntu,
}

// Detect returns the Linux distribution.
func Detect() (Distro, error) {
	var id string
	var err error

	if _, err = os.Stat("/etc/os-release"); !os.IsNotExist(err) {
		cfg, err := shconf.ParseFile("/etc/os-release")
		if err != nil {
			return 0, err
		}
		if id, err = cfg.Get("ID"); err != nil {
			return 0, err
		}

		if v, found := idToDistro[id]; found {
			return v, nil
		}
	}

	return DistroUnknown, nil
}
