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
	"bytes"
	"fmt"

	"go.etcd.io/bbolt"
)

// List shows all commands (that contain the pattern, if given)
func List(database *bbolt.DB, pattern string) {
	// We get page names as byte arrays
	patternB := []byte(pattern)

	// Iterate through keys, printing all that match the pattern,
	// as keys are byte ordered, they are also in alphabethic order.
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket, creating it if it doesn't yet exist
			tx.CreateBucketIfNotExists(pagesBucket)
			bucket := tx.Bucket(pagesBucket)

			// Print all the pages that contain the pattern
			return bucket.ForEach(
				func(page, _ []byte) error {
					if bytes.Contains(page, patternB) {
						fmt.Println(string(page))
					}

					return nil
				})
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}
