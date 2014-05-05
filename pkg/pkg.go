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
	"os/exec"
)

type Packager interface {
	// Install installs a paquete.
	Install(...string) error

	// Remove removes a paquete.
	Remove(bool, ...string) error

	// Purge removes a paquete and its config files.
	Purge(bool, ...string) error

	// Clean erases downloaded archive files.
	Clean() error

	// Upgrade upgrades all the packages on the system.
	Upgrade() error
}

// * * *

// PackageType represents a package management system.
type PackageType int

const (
	Deb PackageType = iota + 1
	RPM
	Pacman
	Ebuild
	ZYpp
)

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

// * * *

type packagerInfo struct {
	typ PackageType
	pkg Packager
}

// execPackagers is a list of executables of package managers.
var execPackagers = map[string]packagerInfo{
	"apt-get": packagerInfo{Deb, new(deb)},
	"yum":     packagerInfo{RPM, new(rpm)},
	"pacman":  packagerInfo{Pacman, new(pacman)},
	"emerge":  packagerInfo{Ebuild, new(ebuild)},
	"zypper":  packagerInfo{ZYpp, new(zypp)},
}

// Detect tries to get the package manager used in the system, looking for
// executables in directory "/usr/bin".
func Detect() (PackageType, Packager, error) {
	for k, v := range execPackagers {
		_, err := exec.LookPath("/usr/bin/" + k)
		if err == nil {
			return v.typ, v.pkg, nil
		}
	}
	return 0, nil, errors.New("package manager not found in directory /usr/bin")
}

// runc executes a command logging its output if there is not any error.
func run(cmd string, arg ...string) error {
	_, err := exec.Command(cmd, arg...).CombinedOutput()
	if err != nil {
		return err
	}

	// log.Print(string(out)) // DEBUG
	return nil
}

type packageSystem struct {
	isFirstInstall bool
}
