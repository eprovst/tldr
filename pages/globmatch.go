// Copyright © 2018 Evert Provoost
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package pages

import "bytes"

// TODO: Work on rune level instead of bytes? Because '?' currently
// does not support characters longer than one byte.

// We use two special ASCII codes in our pattern building
const (
	asciiGS = byte(0x1D)
	asciiUS = byte(0x1F)
)

// GlobMatcher is the matcher built from a GlobPattern
type GlobMatcher struct {
	useBytes       bool
	pattern        []byte
	dfa            [][]int
	acceptingState int
}

// Match matches the text with the pattern
func (gm GlobMatcher) Match(text []byte) bool {
	// If we didn't have a glob pattern we can symply use the
	// bytes function.
	if gm.useBytes {
		return bytes.Equal(text, gm.pattern)
	}

	// No use this DFA to match a pattern
	s := 0
	for i := range text {
		// Go to the next state
		s = gm.dfa[text[i]][s]

		// Did we come to a point where a match is impossible?
		if s == -1 {
			return false
		}
	}

	// Are we in the accepting state?
	if s == gm.acceptingState {
		return true
	}

	// Else: we are in a nonaccepting state.
	return false
}

// NewGlobMatcher uses a KMP style DFA to match strings which may
// or may not contain globs, this mechod also allows the escaping
// of glob symbols.
// NOTE TO SELF: If this is new, maybe I should write about it...
func NewGlobMatcher(pattern []byte) *GlobMatcher {
	if !bytes.ContainsAny(pattern, "*?") {
		// If no globs are present we shoud use the standard library
		return &GlobMatcher{
			useBytes: true,
			pattern:  pattern,
		}
	}

	// Else use discrete finite automata :)
	// Parse the pattern, replacing * by GS, \* by *, ? by US, \? by ? and \\ by \
	escd := false
	pat := bytes.NewBuffer([]byte{})

	for _, c := range pattern {
		switch c {
		case '\\':
			if escd {
				pat.WriteByte('\\')
				escd = false

			} else {
				escd = true
			}

		case '*':
			if escd {
				pat.WriteByte('*')
				escd = false

			} else {
				pat.WriteByte(asciiGS)
			}

		case '?':
			if escd {
				pat.WriteByte('?')
				escd = false

			} else {
				pat.WriteByte(asciiUS)
			}

		default:
			if escd {
				pat.WriteByte('\\')
			}

			pat.WriteByte(c)
		}
	}

	pattern = pat.Bytes()

	// Initialise the DFA
	// note that Go will write zero's everywhere
	dfa := make([][]int, 256)
	for i := range dfa {
		dfa[i] = make([]int, len(pattern)+1)
	}

	// The aceptingState will be the lenght of the pattern
	acceptingState := len(pattern)

	// Now build the DFA
	// We use a technique similar to KMP however our initial state points to
	// the no match state -1 instead of itself.
	x := 0
	for c := 0; c < 256; c++ {
		dfa[c][0] = -1
	}

	for i := range pattern {
		switch pattern[i] {
		case asciiGS:
			x = i
			fallthrough

		case asciiUS:
			for c := 0; c < 256; c++ {
				dfa[c][i] = i + 1
			}

		default:
			for c := 0; c < 256; c++ {
				dfa[c][i] = dfa[c][x]
			}

			dfa[pattern[i]][i] = i + 1
		}
	}

	// Finally set the accepting state
	for c := 0; c < 256; c++ {
		dfa[c][acceptingState] = dfa[c][x]
	}

	return &GlobMatcher{
		useBytes:       false,
		pattern:        pattern,
		dfa:            dfa,
		acceptingState: acceptingState,
	}
}
