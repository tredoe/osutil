// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

type pacman packageSystem

func (p pacman) Install(name ...string) error {
	if p.isFirstInstall {
		p.isFirstInstall = false
		arg := []string{"-Syu", "--needed", "--noprogressbar"}
		arg = append(arg, name...)
		return run("/usr/bin/pacman", arg...)
	}

	arg := []string{"-S", "--needed", "--noprogressbar"}
	arg = append(arg, name...)
	return run("/usr/bin/pacman", arg...)
}

func (pacman) Remove(isMetapackage bool, name ...string) error {
	if isMetapackage {
		arg := []string{"-Rs"}
		arg = append(arg, name...)
		return run("/usr/bin/pacman", arg...)
	}

	arg := []string{"-R"}
	arg = append(arg, name...)
	return run("/usr/bin/pacman", arg...)
}

func (pacman) Purge(isMetapackage bool, name ...string) error {
	if isMetapackage {
		arg := []string{"-Rsn"}
		arg = append(arg, name...)
		return run("/usr/bin/pacman", arg...)
	}

	arg := []string{"-Rn"}
	arg = append(arg, name...)
	return run("/usr/bin/pacman", arg...)
}

func (pacman) Clean() error {
	return nil
}

func (pacman) Upgrade() error {
	return run("/usr/bin/pacman", "-Syu")
}
