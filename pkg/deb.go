// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

type deb packageSystem

func (p deb) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/apt-get", "update"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	arg := []string{"install", "-y"}
	arg = append(arg, name...)
	return run("/usr/bin/apt-get", arg...)
}

func (deb) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"remove", "-y"}
	arg = append(arg, name...)
	if err := run("/usr/bin/apt-get", arg...); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/apt-get", "autoremove", "-y")
	}
	return nil
}

func (deb) Purge(isMetapackage bool, name ...string) error {
	arg := []string{"purge", "-y"}
	arg = append(arg, name...)
	if err := run("/usr/bin/apt-get", arg...); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/apt-get", "autoremove", "--purge", "-y")
	}
	return nil
}

func (deb) Clean() error {
	return run("/usr/bin/apt-get", "clean")
}

func (deb) Upgrade() error {
	if err := run("/usr/bin/apt-get", "update"); err != nil {
		return err
	}
	return run("/usr/bin/apt-get", "upgrade")
}
