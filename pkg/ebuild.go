// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

type ebuild packageSystem

func (p ebuild) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/emerge", "--sync"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}
	return run("/usr/bin/emerge", name...)
}

func (ebuild) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"--unmerge"}
	arg = append(arg, name...)
	if err := run("/usr/bin/emerge", arg...); err != nil {
		return err
	}

	if isMetapackage {
		return run("/usr/bin/emerge", "--depclean")
	}
	return nil
}

func (ebuild) Purge(isMetapackage bool, name ...string) error {
	return nil
}

func (ebuild) Clean() error {
	return nil
}

func (ebuild) Upgrade() error {
	if err := run("/usr/bin/emerge", "--sync"); err != nil {
		return err
	}
	return run("/usr/bin/emerge", "--update", "--deep", "--with-bdeps=y", "--newuse world")
}
