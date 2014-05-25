// +build gake

package main

import (
	"fmt"

	"github.com/kless/osutil/gake/making"
)

// MakeHello says something.
func MakeHello(m *making.M) {
	fmt.Println("Hello!")
	m.Log(`Testing "Hello" function`)
}
