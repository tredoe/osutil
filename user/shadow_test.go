// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestShadowParser(t *testing.T) {
	f, err := os.Open(_SHADOW_FILE)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			continue
		}

		if _, err = parseShadow(string(line)); err != nil {
			t.Error(err)
		}
	}
}

func TestShadowFull(t *testing.T) {
	entry, err := LookupShadow("root")
	if err != nil || entry == nil {
		t.Error(err)
	}

	entries, err := LookupInShadow(S_PASSWD, "!", -1)
	if err != nil || entries == nil {
		t.Error(err)
	}

	entries, err = LookupInShadow(S_ALL, nil, -1)
	if err != nil || len(entries) == 0 {
		t.Error(err)
	}
}

func TestShadowCount(t *testing.T) {
	count := 2
	entries, err := LookupInShadow(S_MIN, 0, count)
	if err != nil || len(entries) != count {
		t.Error(err)
	}

	count = 5
	entries, err = LookupInShadow(S_ALL, nil, count)
	if err != nil || len(entries) != count {
		t.Error(err)
	}
}

func TestShadowError(t *testing.T) {
	var err error

	if _, err = LookupShadow("!!!???"); err != ErrNoFound {
		t.Error("expected to report ErrNoFound")
	}

	if _, err = LookupInShadow(S_MIN, 0, 0); err != ErrSearch {
		t.Error("expected to report ErrSearch")
	}

	s := &Shadow{}
	if err = s.Add(nil); err != RequiredError("Name") {
		t.Error("expected to report RequiredError")
	}
}
