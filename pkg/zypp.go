// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

type zypp packageSystem

func (p zypp) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/zypper", "refresh"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	arg := []string{"install", "--auto-agree-with-licenses"}
	arg = append(arg, name...)
	return run("/usr/bin/zypper", arg...)
}

func (zypp) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"remove"}
	arg = append(arg, name...)
	return run("/usr/bin/zypper", arg...)
}

func (zypp) Purge(isMetapackage bool, name ...string) error {
	return nil
}

func (zypp) Clean() error {
	return run("/usr/bin/zypper", "clean")
}

func (zypp) Upgrade() error {
	if err := run("/usr/bin/zypper", "refresh"); err != nil {
		return err
	}
	return run("/usr/bin/zypper", "up", "--auto-agree-with-licenses")
}
