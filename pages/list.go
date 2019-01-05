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
	"bytes"
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

// List shows all commands for the current platform
func List(database *bbolt.DB) {
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages buckets, only english concerns us here
			// as other languages will only contain translations
			englishCommon, englishPlatform, _, _, err := getBuckets(tx)

			if err != nil {
				return err
			}

			// This one should exist
			if englishCommon == nil {
				emptyDatabase()
				return nil
			}

			// TODO strip doubles
			// Print all the pages
			printPageNames(englishPlatform)
			return printPageNames(englishCommon)
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func printPageNames(bucket *bbolt.Bucket) error {
	// Does the bucket exist?
	if bucket == nil {
		return nil
	}

	// Print all the names of the pages in the bucket
	return bucket.ForEach(
		func(name, _ []byte) error {
			fmt.Println(string(name))
			return nil
		})
}

// ListAll shows all commands for all platforms
func ListAll(database *bbolt.DB) {
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket
			root := tx.Bucket(defaultBucket)

			if root == nil {
				emptyDatabase()
				return nil
			}

			// Print all the pages
			// TODO remove doubles
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
			// Open the pages bucket, we assume that all pages appear in EN
			root := tx.Bucket(defaultBucket)

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

// ListLanguages shows all available platforms
func ListLanguages(database *bbolt.DB) {
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Not even the default bucket -> empty database
			if tx.Bucket(defaultBucket) == nil {
				emptyDatabase()
				return nil
			}

			// Print all the languages, which are buckets in the default
			return tx.ForEach(
				func(name []byte, _ *bbolt.Bucket) error {
					if !bytes.Equal(name, defaultBucket) {
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
