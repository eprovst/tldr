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

	"go.etcd.io/bbolt"
)

// Show shows help for a command
func Show(database *bbolt.DB, commands []string) {
	// Get the page
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			bucket := tx.Bucket(pagesBucket)

			if bucket == nil {
				emptyDatabase()
				return nil
			}

			// Print all the given commands
			for _, command := range commands {
				page := bucket.Get([]byte(command))

				if page == nil {
					pageUnavailable(command)

				} else {
					prettyPrint(page)
				}
			}

			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}
