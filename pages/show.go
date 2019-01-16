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

	"go.etcd.io/bbolt"
)

// Show shows help for a command
func Show(database *bbolt.DB, commands []string) {
	// Get the page
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages buckets
			englishCommon, englishPlatform, langCommon, langPlatform, err := getBuckets(tx)

			if err != nil {
				return err
			}

			// This one should exist
			if englishCommon == nil {
				emptyDatabase()
				return nil
			}

			// Join the subcommands with a `-`
			command := commands[0]
			for _, subcommand := range commands[1:] {
				command += "-" + subcommand
			}

			// Print the given command
			var page []byte

			if langPlatform != nil {
				page = langPlatform.Get([]byte(command))
			}

			if page == nil && langCommon != nil {
				page = langCommon.Get([]byte(command))
			}

			if page == nil && englishPlatform != nil {
				page = englishPlatform.Get([]byte(command))
			}

			if page == nil {
				page = englishCommon.Get([]byte(command))
			}

			if page == nil {
				pageUnavailable(command)

			} else {
				prettyPrint(page)
			}

			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
