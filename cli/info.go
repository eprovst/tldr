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

package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/elecprog/tldr/targets"
	flag "github.com/spf13/pflag"
)

// Version info
const thisVersion = "v0.4.0"

func showHelp() {
	fmt.Fprintln(os.Stderr, "Go command line client for tldr")
	fmt.Fprintln(os.Stderr, "\nUsage:")
	fmt.Fprintln(os.Stderr, "  tldr [flags] [command]")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
}

func showVersion() {
	fmt.Println("tldr", thisVersion, "on", targets.OsName)
}

// TODO: Once Go 1.11 is more commonly available remove this in favour of
// os.UserCacheDir(), ie. when Ubuntu 19.04 is released.
// From Go standard library: os/file.go#L346
func userCacheDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LocalAppData")

		if dir == "" {
			return "", errors.New("%LocalAppData% is not defined")
		}

	case "darwin":
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir += "/Library/Caches"

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}
		dir += "/lib/cache"

	default: // Unix
		dir = os.Getenv("XDG_CACHE_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_CACHE_HOME nor $HOME are defined")
			}
			dir += "/.cache"
		}
	}

	return dir, nil
}

// getDatabasePath returns the path to the database or panics if the system
// does not have a cache directory.
func getDatabasePath() string {
	dir, err := userCacheDir()
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
