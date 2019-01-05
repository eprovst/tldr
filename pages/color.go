// Copyright © 2019 Evert Provoost
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
	"strconv"
)

type color uint16

const (
	colBlack color = iota + 30
	colRed
	colGreen
	colYellow
	colBlue
	colMagenta
	colCyan
	colGray

	colBrightBlack color = (iota - 8) + 90
	colBrightRed
	colBrightGreen
	colBrightYellow
	colBrightBlue
	colBrightMagenta
	colBrightCyan
	colBrightGray

	colDefault color = 39
	colReset   color = 0

	modBold   color = 2 << 8
	modFaint  color = 2 << 9
	modItalic color = 2 << 10

	colMask uint64 = 2<<8 - 1
)

type coloredText struct {
	value interface{}
	color color
}

func colorize(val interface{}, col color) coloredText {
	return coloredText{value: val, color: col}
}

func toEscapeCode(col color) string {
	res := "\033["

	res += strconv.FormatUint(uint64(col)&colMask, 10)

	if col&modBold > 0 {
		res += ";1"
	}

	if col&modFaint > 0 {
		res += ";2"
	}

	if col&modItalic > 0 {
		res += ";3"
	}

	return res + "m"
}

func (ct coloredText) Format(f fmt.State, c rune) {
	if ct.color == colDefault {
		fmt.Fprint(f, ct.value)
	} else {
		fmt.Fprint(f, toEscapeCode(ct.color), ct.value, "\033[0m")
	}
}
