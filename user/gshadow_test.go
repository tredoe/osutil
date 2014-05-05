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

func TestGShadowParser(t *testing.T) {
	f, err := os.Open(_GSHADOW_FILE)
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

		if _, err = parseGShadow(string(line)); err != nil {
			t.Error(err)
		}
	}
}

func TestGShadowFull(t *testing.T) {
	entry, err := LookupGShadow("root")
	if err != nil || entry == nil {
		t.Error(err)
	}

	entries, err := LookupInGShadow(GS_PASSWD, "!", -1)
	if err != nil || entries == nil {
		t.Error(err)
	}

	entries, err = LookupInGShadow(GS_ALL, "", -1)
	if err != nil || len(entries) == 0 {
		t.Error(err)
	}
}

func TestGShadowCount(t *testing.T) {
	count := 5
	entries, err := LookupInGShadow(GS_ALL, "", count)
	if err != nil || len(entries) != count {
		t.Error(err)
	}
}

func TestGShadowError(t *testing.T) {
	var err error

	if _, err = LookupGShadow("!!!???"); err != ErrNoFound {
		t.Error("expected to report ErrNoFound")
	}

	if _, err = LookupInGShadow(GS_MEMBER, "", 0); err != ErrSearch {
		t.Error("expected to report ErrSearch")
	}

	gs := &GShadow{}
	if err = gs.Add(nil); err != RequiredError("Name") {
		t.Error("expected to report RequiredError")
	}
}
