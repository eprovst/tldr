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

	"github.com/elecprog/tldr/targets"
	"go.etcd.io/bbolt"
)

// Search shows all pages that contain the pattern
func Search(database *bbolt.DB, pattern string) {
	// TODO: Sorting of output

	// Iterate through keys, printing all that match the pattern,
	// as keys are byte ordered, they are also in alphabethic order.
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			root := tx.Bucket(rootBucket)

			if root == nil {
				emptyDatabase()
				return nil
			}

			searchInBucket(root.Bucket(commonBucket), []byte(pattern))
			return searchInBucket(root.Bucket([]byte(targets.OsDir)), []byte(pattern))
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}

func searchInBucket(bucket *bbolt.Bucket, pattern []byte) error {
	// Print all the pages that start with the pattern
	bucket.ForEach(
		func(page, _ []byte) error {
			if bytes.HasPrefix(page, pattern) {
				fmt.Println(string(page))
			}

			return nil
		})

	// Print all the pages that contain the pattern but don't start with it
	return bucket.ForEach(
		func(page, _ []byte) error {
			if !bytes.HasPrefix(page, pattern) && bytes.Contains(page, pattern) {
				fmt.Println(string(page))
			}

			return nil
		})
}
