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

package info

import (
	"log"
	"os"
)

// PagesSource is the location from where we will download the pages.
const PagesSource = "https://tldr.sh/assets/tldr.zip"

// GetDatabasePath returns the path to the database or panics if the system
// is unsupported.
func GetDatabasePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	return dir + "/tldr/tldr.bbolt"
}

// Exists checks if a file exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}