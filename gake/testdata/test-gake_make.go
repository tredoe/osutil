// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build task foo

package main

import "github.com/kless/osutil/gake/making"

// MakeHello says something.
func MakeHello(*making.M) {
	m.Log("Hello!")
}
