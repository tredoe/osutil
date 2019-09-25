// Copyright 2012, Jeramey Crawford <jeramey@antihe.ro>
// Copyright 2013, Jonas mg
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package common

import (
	"testing"
	"strings"
	"strconv"
)

var _Salt = &Salt{
	MagicPrefix: []byte("$foo$"),
	SaltLenMin:  1,
	SaltLenMax:  8,
	RoundsMin: 1000,
	RoundsMax: 999999999,
	RoundsDefault: 5000,
}

func TestGenerateSalt(t *testing.T) {
	magicPrefixLen := len(_Salt.MagicPrefix)
	
	salt := _Salt.Generate(0)
	if len(salt) != magicPrefixLen+1 {
		t.Errorf("Expected len 1, got len %d", len(salt))
	}

	for i := _Salt.SaltLenMin; i <= _Salt.SaltLenMax; i++ {
		salt = _Salt.Generate(i)
		if len(salt) != magicPrefixLen+i {
			t.Errorf("Expected len %d, got len %d", i, len(salt))
		}
	}

	salt = _Salt.Generate(9)
	if len(salt) != magicPrefixLen+8 {
		t.Errorf("Expected len 8, got len %d", len(salt))
	}
}

func TestGenerateSaltWRounds(t *testing.T) {
	const rounds = 5001
	salt := _Salt.GenerateWRounds(_Salt.SaltLenMax, rounds)
	if salt == nil {
		t.Errorf("salt should not be nil")
	}

	expectedPrefix := string(_Salt.MagicPrefix) + "rounds=" + strconv.Itoa(rounds) + "$"
	if !strings.HasPrefix(string(salt), expectedPrefix) {
		t.Errorf("salt '%s' should start with prefix '%s' but didn't", salt, expectedPrefix)

	}
}
