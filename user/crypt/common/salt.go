// Copyright 2012, Jeramey Crawford <jeramey@antihe.ro>
// Copyright 2013, Jonas mg
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file.

package common

import (
	"crypto/rand"
	"errors"
	"strconv"
)

var (
	ErrSaltPrefix = errors.New("invalid magic prefix")
	ErrSaltFormat = errors.New("invalid salt format")
	ErrSaltRounds = errors.New("invalid rounds")
)

// Salt represents a salt.
type Salt struct {
	MagicPrefix []byte

	SaltLenMin int
	SaltLenMax int

	RoundsMin     int
	RoundsMax     int
	RoundsDefault int
}

// Generate generates a random salt of a given length.
//
// The length is set thus:
//
//   length > SaltLenMax: length = SaltLenMax
//   length < SaltLenMin: length = SaltLenMin
func (s *Salt) Generate(length int) []byte {
	switch {
	case length > s.SaltLenMax:
		length = s.SaltLenMax
	case length < s.SaltLenMin:
		length = s.SaltLenMin
	default:
	}

	saltLen := (length * 6 / 8)
	if (length*6)%8 != 0 {
		saltLen += 1
	}
	salt := make([]byte, saltLen)
	rand.Read(salt)

	magicPrefixLen := len(s.MagicPrefix)
	out := make([]byte, magicPrefixLen+length)
	copy(out, s.MagicPrefix)
	copy(out[magicPrefixLen:], Base64_24Bit(salt))
	return out
}

// GenerateWRounds creates a random salt with the random bytes being of the
// length provided, and the rounds parameter set as specified.
//
// The parameters are set thus:
//
//   length > SaltLenMax: length = SaltLenMax
//   length < SaltLenMin: length = SaltLenMin
//
//   rounds < 0: rounds = RoundsDefault
//   rounds < RoundsMin: rounds = RoundsMin
//   rounds > RoundsMax: rounds = RoundsMax
//
// If rounds is equal to RoundsDefault, then the "rounds=" part of the salt is
// removed.
func (s *Salt) GenerateWRounds(length, rounds int) []byte {
	switch {
	case length > s.SaltLenMax:
		length = s.SaltLenMax
	case length < s.SaltLenMin:
		length = s.SaltLenMin
	default:
	}
	
	switch {
	case rounds > s.RoundsMax:
		rounds = s.RoundsMax
	case rounds < s.RoundsMin:
		rounds = s.RoundsMin
	default:
		rounds = s.RoundsDefault
	}

	saltLen := (length * 6 / 8)
	if (length*6)%8 != 0 {
		saltLen += 1
	}
	salt := make([]byte, saltLen)
	rand.Read(salt)

	roundsText := ""
	if rounds != s.RoundsDefault {
		roundsText = "rounds=" + strconv.Itoa(rounds)
	}

	magicPrefixLen := len(s.MagicPrefix)
	roundsTextLen := len(roundsText)
	out := make([]byte, magicPrefixLength+roundsTextLen+length)
	copy(out, s.MagicPrefix)
	copy(out[magicPrefixLength:], []byte(roundsText))
	copy(out[magicPrefixLength+roundsTextLen:], Base64_24Bit(salt))
	return out
}
