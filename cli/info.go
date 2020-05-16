// Copyright © 2020 Evert Provoost
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

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/elecprog/tldr/targets"
	flag "github.com/spf13/pflag"
)

// Version info
const thisVersion = "v0.4.1"
const thisSpec = "1.2 (partially)"

func showHelp() {
	fmt.Fprintln(os.Stderr, "Go command line client for tldr")
	fmt.Fprintln(os.Stderr, "\nUsage:")
	fmt.Fprintln(os.Stderr, "  tldr [flags] [command]")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
}

func showVersion() {
	fmt.Println("Go client for tldr", "("+thisVersion+")", "on", targets.OsName)
	fmt.Println("Implements tldr spec", thisSpec)
}

// getDatabasePath returns the path to the database or panics if the system
// does not have a cache directory.
func getDatabasePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	return filepath.Join(dir, "tldr", "tldr.bbolt")
}

// pathExists checks if a path/file exists
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
