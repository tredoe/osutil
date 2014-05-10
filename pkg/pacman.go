// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "github.com/kless/osutil"

type pacman packageSystem

func (p pacman) Install(name ...string) error {
	args := []string{"-S", "--needed", "--noprogressbar"}

	return osutil.Exec("/usr/bin/pacman", append(args, name...)...)
}

func (p pacman) Remove(name ...string) error {
	args := []string{"-R"}

	return osutil.Exec("/usr/bin/pacman", append(args, name...)...)
}

func (p pacman) RemoveMeta(name ...string) error {
	args := []string{"-Rs"}

	if err := osutil.Exec("/usr/bin/pacman", append(args, name...)...);err != nil {
		return err
	}
	return p.Remove(name...)
}

func (p pacman) Purge(name ...string) error {
	args := []string{"-Rn"}

	return osutil.Exec("/usr/bin/pacman", append(args, name...)...)
}

func (p pacman) PurgeMeta(name ...string) error {
	args := []string{"-Rsn"}

	if err := osutil.Exec("/usr/bin/pacman", append(args, name...)...); err != nil {
		return err
	}
	return p.Purge(name...)
}

func (p pacman) Update() error {
	return osutil.Exec("/usr/bin/pacman", "-Syu", "--needed", "--noprogressbar")
}

func (p pacman) Upgrade() error {
	return osutil.Exec("/usr/bin/pacman", "-Syu")
}

func (p pacman) Clean() error {
	return nil
}
