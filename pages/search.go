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
	"fmt"
	"os"
	"regexp"

	"go.etcd.io/bbolt"
)

// Search shows all pages that matches the regex
func Search(database *bbolt.DB, regex string) {
	// Iterate through keys, printing all that match the pattern,
	// as keys are byte ordered, they are also in alphabethic order.
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

			// Create a matcher from the regex
			matcher, err := regexp.CompilePOSIX(regex)

			if err != nil {
				return err
			}

			// TODO strip doubles
			// Search in all relevant buckets
			searchInBucket(englishPlatform, matcher)
			return searchInBucket(englishCommon, matcher)
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func searchInBucket(bucket *bbolt.Bucket, matcher *regexp.Regexp) error {
	// Handle nil buckets
	if bucket == nil {
		return nil
	}

	// Print all the pages that match the pattern
	return bucket.ForEach(
		func(page, _ []byte) error {
			if matcher.Match(page) {
				fmt.Println(string(page))
			}

			return nil
		})
}
