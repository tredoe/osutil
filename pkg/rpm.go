// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

type rpm packageSystem

func (p rpm) Install(name ...string) error {
	if p.isFirstInstall {
		if err := run("/usr/bin/yum", "update"); err != nil {
			return err
		}
		p.isFirstInstall = false
	}

	arg := []string{"install"}
	arg = append(arg, name...)
	return run("/usr/bin/yum", arg...)
}

func (rpm) Remove(isMetapackage bool, name ...string) error {
	arg := []string{"remove"}
	arg = append(arg, name...)
	return run("/usr/bin/yum", arg...)
}

func (rpm) Purge(isMetapackage bool, name ...string) error {
	return nil
}

func (rpm) Clean() error {
	return run("/usr/bin/yum", "clean", "packages")
}

func (rpm) Upgrade() error {
	return run("/usr/bin/yum", "update")
}
