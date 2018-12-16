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
func Show(database *bbolt.DB, commands []string) {
	// Get the page
	err := database.View(
		func(tx *bbolt.Tx) error {
			// Open the pages bucket, creating it if it doesn't yet exist
			tx.CreateBucketIfNotExists(pagesBucket)
			buck := tx.Bucket(pagesBucket)

			// Print all the given commands
			for _, command := range commands {
				page := buck.Get([]byte(command))
				prettyPrint(command, page)
			}

			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}

func prettyPrint(command string, page []byte) {
	// The page is not in the database
	if page == nil {
		fmt.Println()
		fmt.Println(aurora.Bold(command), "documentation is not available.")
		fmt.Println("Consider making a Pull Request to https://github.com/tldr-pages/tldr")
		fmt.Println()
		return
	}

	// Add an blank line in front of the page
	fmt.Println()

	// Pretty print the lines in the page
	for _, line := range strings.Split(string(page), "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			line = strings.TrimPrefix(line, "#")

			for _, part := range processLine(line) {
				fmt.Print(aurora.Bold(part))
			}
			fmt.Println()

		} else if strings.HasPrefix(line, ">") {
			line = strings.TrimPrefix(line, ">")
			line = strings.TrimSpace(line)

			for _, part := range processLine(line) {
				fmt.Print(part)
			}
			fmt.Println()

		} else if strings.HasPrefix(line, "-") {
			line = strings.TrimPrefix(line, "-")
			line = strings.TrimSpace(line)

			fmt.Print("\n- ")

			for _, part := range processLine(line) {
				switch part := part.(type) {
				case aurora.Value:
					fmt.Print(part)

				default:
					fmt.Print(aurora.Green(part).Bold())
				}
			}
			fmt.Println()

		} else if line == "" {
			// Skip empty lines

		} else {
			// Print the parsed line
			fmt.Print("  ")

			for _, part := range processLine(line) {
				fmt.Print(part)
			}
			fmt.Println()
		}
	}

	// Add an extra blank line at the end of the page
	fmt.Println()
}

func processLine(line string) []interface{} {
	// Remove unneeded spaces
	line = strings.TrimSpace(line)

	// Ugly parsing number one
	parts := []interface{}{}

	// Our parsing method would fail on ``, but as
	// these a no-ops we can safely remove them.
	line = strings.Replace(line, "``", "", -1)

	inVerbatim := line[0] == '`'
	for _, part := range strings.Split(line, "`") {
		// Skip empty strings
		if len(part) == 0 {
			continue
		}

		if inVerbatim {
			// Verbatim
			parts = append(parts, processVerbatim(part)...)

		} else {
			// Normal text
			parts = append(parts, part)
		}

		inVerbatim = !inVerbatim
	}

	// As you might have noticed, we never check if the backticks are balanced.
	// But that check is not regular, and pages should be valid,
	// so in theory we never have a case where the backticks aren't balanced.

	return parts
}

func processVerbatim(verbatim string) []interface{} {
	// Ugly parsing number two
	parts := []interface{}{}

	// Our parsing method would fail on {{}} or }}{{, but as
	// these a no-ops we can safely remove them.
	verbatim = strings.Replace(verbatim, "{{}}", "", -1)
	verbatim = strings.Replace(verbatim, "}}{{", "", -1)

	inOptional := strings.HasPrefix(verbatim, "{{")
	for _, segment := range strings.Split(verbatim, "{{") {
		for _, part := range strings.Split(segment, "}}") {
			if inOptional {
				// Optional
				parts = append(parts, part)

			} else {
				// Verbatim
				parts = append(parts, aurora.Red(part).Bold())
			}

			inOptional = !inOptional
		}
	}

	// As you might have noticed, we never check if the braces are balanced.
	// But that check is not regular, and pages should be valid,
	// so in theory we never have a case where the braces aren't balanced.

	// Return the parsed line
	return parts
}
