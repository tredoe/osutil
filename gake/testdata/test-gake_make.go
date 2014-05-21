// +build gake

package main

import "github.com/kless/osutil/gake/making"

// MakeHello says something.
func MakeHello(m *making.M) {
	m.Log("Hello!")
}
