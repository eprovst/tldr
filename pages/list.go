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
	"errors"
	"fmt"
	"os"

	"github.com/elecprog/tldr/targets"
	"go.etcd.io/bbolt"
)

// List shows all commands for the current platform
func List(database *bbolt.DB) {
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			root := tx.Bucket(rootBucket)

			if root == nil {
				emptyDatabase()
				return nil
			}

			// Get the local en common pages
			common := root.Bucket(commonBucket)
			local := root.Bucket([]byte(targets.OsDir))

			// If the platform is not supported print an error
			if local == nil {
				return errors.New("unsupported platform '" + targets.OsDir + "'")
			}

			// Print all the common pages
			common.ForEach(printPageName)

			// Print all the platform specific pages
			return local.ForEach(printPageName)
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func printPageName(name, _ []byte) error {
	fmt.Println(string(name))
	return nil
}

// ListAll shows all commands for all platforms
func ListAll(database *bbolt.DB) {
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			root := tx.Bucket(rootBucket)

			if root == nil {
				emptyDatabase()
				return nil
			}

			// Print all the pages
			printPages(root)

			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func printPages(bucket *bbolt.Bucket) error {
	return bucket.ForEach(
		func(page, contents []byte) error {
			if contents != nil {
				fmt.Println(string(page))

			} else {
				printPages(bucket.Bucket(page))
			}
			return nil
		})
}

// ListPlatforms shows all available platforms
func ListPlatforms(database *bbolt.DB) {
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			root := tx.Bucket(rootBucket)

			if root == nil {
				emptyDatabase()
				return nil
			}

			// Print all the platforms, which are buckets in the root
			return root.ForEach(
				func(name []byte, value []byte) error {
					if value == nil && !bytes.Equal(name, commonBucket) {
						fmt.Println(string(name))
					}

					return nil
				})
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
