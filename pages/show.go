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
	"strings"

	"github.com/logrusorgru/aurora"
	"go.etcd.io/bbolt"
)

// Show shows help for a command
func Show(database *bbolt.DB, command string) {
	// Get the page
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket, creating it if it doesn't yet exist
			tx.CreateBucketIfNotExists(pagesBucket)
			buck := tx.Bucket(pagesBucket)

			// Read the page
			page := buck.Get([]byte(command))
			prettyPrint(command, page)

			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}

func prettyPrint(command string, page []byte) {
	// No page in database
	if page == nil {
		println()
		fmt.Println(aurora.Bold(command), "documentation is not available")
		fmt.Println("Consider contributing a Pull Request to https://github.com/tldr-pages/tldr")
		println()
		return
	}

	// Add an blank line in front of the page
	println()

	// Pretty print the lines in the page
	for _, line := range strings.Split(string(page), "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimSpace(line)

			fmt.Println(aurora.Bold(line))

		} else if strings.HasPrefix(line, ">") {
			line = strings.TrimPrefix(line, ">")
			line = strings.TrimSpace(line)

			fmt.Println(line)

		} else if strings.HasPrefix(line, "-") {
			line = strings.TrimPrefix(line, "-")
			line = strings.TrimSpace(line)

			fmt.Println("\n-", aurora.Green(line).Bold())

		} else if strings.HasPrefix(line, "`") && strings.HasSuffix(line, "`") {
			line = strings.TrimPrefix(line, "`")
			line = strings.TrimSuffix(line, "`")

			line = strings.TrimSpace(line)

			// Some ugly parsing...
			parts := []interface{}{}

			// Our parsing method would fail on }}{{, but as
			// it's a no-op we can safely remove it.
			line = strings.Replace(line, "}}{{", "", -1)

			inVerbatim := !strings.HasPrefix(line, "{{")
			for _, segment := range strings.Split(line, "{{") {
				for _, part := range strings.Split(segment, "}}") {
					if inVerbatim {
						parts = append(parts, aurora.Red(part).Bold())

					} else {
						parts = append(parts, part)
					}

					inVerbatim = !inVerbatim
				}
			}

			// Print the parsed line
			print("  ")
			for _, part := range parts {
				fmt.Print(part)
			}
			println()

		} else if line == "" {
			// Skip empty lines

		} else {
			// Simply show invalid lines without formatting
			fmt.Println(line)
		}
	}

	// Add an extra blank line at the end of the page
	println()
}
