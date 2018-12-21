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

import (
	"fmt"
	"strings"
)

// We use two special ASCII codes in our pattern building
const (
	asciiGS = rune(0x1D)
	asciiUS = rune(0x1F)
)

// GlobMatcher is the matcher built from a GlobPattern
type GlobMatcher struct {
	simplePat      bool
	pattern        string
	dfa            []state
	acceptingState int
}

// Match matches the string with the pattern
func (gm *GlobMatcher) Match(str string) bool {
	// If we didn't have wildcards in the pattern we can simply use the
	// equality operator.
	if gm.simplePat {
		return str == gm.pattern
	}

	// Else use the DFA to match a pattern
	s := 0
	for _, c := range str {
		// Go to the next state
		s = gm.dfa[s].getNext(c)

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

// NewGlobMatcher builds a DFA capable of matching the provided pattern.
// NOTE TO SELF: If this is new, maybe I should write about it...
func NewGlobMatcher(pattern string) *GlobMatcher {
	if !strings.ContainsAny(pattern, "*?") {
		// If no wildcards are present we shoud use the standard library
		return &GlobMatcher{
			simplePat: true,
			pattern:   pattern,
		}
	}

	// Else use discrete finite automata :)
	// Parse the pattern, replacing * by GS, \* by *, ? by US, \? by ? and \\ by \
	escd := false
	builder := strings.Builder{}
	numberOfStars := 0

	for _, c := range pattern {
		switch c {
		case '\\':
			if escd {
				builder.WriteRune('\\')
				escd = false

			} else {
				escd = true
			}

		case '*':
			if escd {
				builder.WriteByte('*')
				escd = false

			} else {
				builder.WriteRune(asciiGS)
				numberOfStars++
			}

		case '?':
			if escd {
				builder.WriteByte('?')
				escd = false

			} else {
				builder.WriteRune(asciiUS)
			}

		default:
			if escd {
				builder.WriteByte('\\')
			}

			builder.WriteRune(c)
		}
	}

	pattern = builder.String()

	// The number of states is the lenght of the pattern,
	// minus the number of stars plus one (the accepting state)
	numberOfStates := len(pattern) - numberOfStars + 1

	// Initialise the DFA
	dfa := make([]state, numberOfStates)

	for s := range dfa {
		dfa[s].match = make(map[rune]int)
	}

	// Now build the DFA
	x, s := -1, 0
	for _, c := range pattern {
		switch c {
		case asciiGS:
			x = s

		case asciiUS:
			if x != -1 {
				fmt.Println("error: '?' after '*' is not supported at this time")
			}

			dfa[s].defaultNext = s + 1
			s++

		default:
			if x == -1 {
				dfa[s].defaultNext = -1

			} else {
				for r, n := range dfa[x].match {
					dfa[s].match[r] = n
				}

				x = dfa[x].getNext(c)
			}

			dfa[s].match[c] = s + 1
			s++
		}
	}

	// Finally set the accepting state
	acceptingState := numberOfStates - 1
	dfa[acceptingState].defaultNext = x

	return &GlobMatcher{
		simplePat:      false,
		pattern:        pattern,
		dfa:            dfa,
		acceptingState: acceptingState,
	}
}

// state represents a state in a DFA
type state struct {
	defaultNext int
	match       map[rune]int
}

func (s state) getNext(c rune) int {
	if n, e := s.match[c]; e {
		return n
	}

	return s.defaultNext
}
