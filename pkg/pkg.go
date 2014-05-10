// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package pkg handles basic operations in the management of packages in
// operating systems.
//
// Important
//
// If you are going to use a package manager different to Deb, then you should
// check the options since I cann't test all.
//
// TODO
//
// Add managers of BSD systems.
//
// Use flag to do not show questions.
package pkg

import (
	"errors"
	"log"
	"os/exec"
)

// Packager is the common interface to handle different package systems.
type Packager interface {
	// Install installs a package.
	Install(name ...string) error

	// Remove removes a package.
	Remove(name ...string) error

	// RemoveMeta removes a meta-package.
	RemoveMeta(name ...string) error

	// Purge removes a package and its configuration files.
	Purge(name ...string) error

	// PurgeMeta removes a meta-package and theirs configuration files.
	PurgeMeta(name ...string) error

	// Update resynchronizes the package index files from their sources.
	Update() error

	// Upgrade upgrades all the packages on the system.
	Upgrade() error

	// Clean erases packages downloaded.
	Clean() error
}

// PackageType represents a package management system.
type PackageType int

const (
	// Linux
	Deb PackageType = iota + 1
	RPM
	Pacman
	Ebuild
	ZYpp
)

func (pkg PackageType) String() string {
	switch pkg {
	case Deb:
		return "Deb"
	case RPM:
		return "RPM"
	case Pacman:
		return "Pacman"
	case Ebuild:
		return "Ebuild"
	case ZYpp:
		return "ZYpp"
	}
	panic("unreachable")
}

// New returns the interface to handle the package manager.
func New(p PackageType) Packager {
	switch p {
	case Deb:
		return new(deb)
	case RPM:
		return new(rpm)
	case Pacman:
		return new(pacman)
	case Ebuild:
		return new(ebuild)
	case ZYpp:
		return new(zypp)
	}
	panic("unreachable")
}

// execPackagers is a list of executables of package managers.
var execPackagers = [...]string{
	Deb:    "apt-get",
	RPM:    "yum",
	Pacman: "pacman",
	Ebuild: "emerge",
	ZYpp:   "zypper",
}

// Detect tries to get the package manager used in the system, looking for
// executables in directory "/usr/bin".
func Detect() (Packager, error) {
	for k, v := range execPackagers {
		_, err := exec.LookPath("/usr/bin/" + v)
		if err == nil {
			return New(PackageType(k)), nil
		}
	}
	return nil, errors.New("package manager not found in directory '/usr/bin'")
}

// * * *

type packageSystem byte

func init() {
	log.SetFlags(0)
	log.SetPrefix("[pkg] ")
}
