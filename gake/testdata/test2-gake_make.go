// +build gake

package main

import (
	"fmt"

	"github.com/kless/osutil/gake/making"
)

// MakeBye says something.
func MakeBye(m *making.M) {
	fmt.Println("Bye!")
	//m.Log(`Testing "Bye" function`)
}
