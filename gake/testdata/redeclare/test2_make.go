// +build gake

package main

import "github.com/kless/osutil/gake/making"

// MakeHello says something.
func MakeHello(*making.M) {
	m.Log("Hello!")
}
