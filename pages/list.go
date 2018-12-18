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

	"github.com/elecprog/tldr/targets"
	"go.etcd.io/bbolt"
)

// List shows all commands for the current platform
func List(database *bbolt.DB) {
	// TODO: Sorting of output

	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			root := tx.Bucket(rootBucket)

			if root == nil {
				emptyDatabase()
				return nil
			}

			// Print all the common pages
			root.Bucket(commonBucket).ForEach(printPageName)

			// Print all the platform specific pages
			return root.Bucket([]byte(targets.OsDir)).ForEach(printPageName)
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}

// ListAll shows all commands for all platforms
func ListAll(database *bbolt.DB) {
	// TODO: Sorting of output

	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			root := tx.Bucket(rootBucket)

			if root == nil {
				emptyDatabase()
				return nil
			}

			// Print all the common pages
			root.Bucket(commonBucket).ForEach(printPageName)

			// Print all the platform specific pages
			for target := range targets.AllTargets {
				root.Bucket([]byte(target)).ForEach(printPageName)
			}

			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}

func printPageName(page, _ []byte) error {
	fmt.Println(string(page))
	return nil
}
