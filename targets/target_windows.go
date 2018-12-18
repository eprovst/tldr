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

package targets

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// OsName is the name of the current platform
const OsName = "Windows"

// OsDir is the directory in the tldr pages containing
// the pages for this platform
var OsDir = "windows"

// Windows by default ignores ASCII escape codes,
// however we can change this using this.
// Why is this not the default? No idea...
func init() {
	// Get a handle to the console
	stdOutHandle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)

	if err != nil {
		fmt.Println("Failed to get a handle for standard input, please open an issue, this should work...")
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the current console settings
	var consoleMode uint32 = 0
	err = windows.GetConsoleMode(stdOutHandle, &consoleMode)

	if err != nil {
		fmt.Println("Failed to get current terminal mode, please open an issue, this should work...")
		fmt.Println(err)
		os.Exit(1)
	}

	// Add support for escape codes to those settings
	consoleMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err = windows.SetConsoleMode(stdOutHandle, consoleMode)

	if err != nil {
		fmt.Println("Failed to enable ASCII escape sequences, please open an issue, this should work...")
		fmt.Println(err)
		os.Exit(1)
	}
}
