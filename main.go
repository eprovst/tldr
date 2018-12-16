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

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/elecprog/tldr/pages"
	"github.com/elecprog/tldr/targets"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

var (
	update          = false
	list            = false
	printCompletion = false
)

var cmd = &cobra.Command{
	Use:     "tldr [-h] [-u] [-l [pattern]] command [command ...]",
	Version: "v0.2.0",
	Short:   "Go command line client for tldr",

	DisableFlagsInUseLine: true,

	Args: func(cmd *cobra.Command, args []string) error {
		// If the bash completion has to be printed we can't have
		// any other arguments
		if printCompletion && (len(args) > 0 || update || list) {
			return errors.New("can't have any other flags or arguments if --completion is set")
		}

		// If we do not have to update the database, list commands
		// nor print the bash completion:
		// we need at least one argument
		if !update && !list && !printCompletion && len(args) == 0 {
			return errors.New("missing argument: command")
		}

		// If we need to list, at most one pattern can be provided
		if list && len(args) > 1 {
			return errors.New("too many arguments: expected at most one pattern")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// We only have to print the bash completion, do so
		if printCompletion {
			fmt.Println(pages.BashCompletion)
			return
		}

		// Get the path where the database is/should be stored
		dbPath := pages.GetDatabasePath()

		// See if it's a first run
		if !pages.PathExists(dbPath) {
			// Create the folder where the database will reside
			err := os.MkdirAll(filepath.Dir(dbPath), 0777)

			if err != nil {
				fmt.Println("error: ", err)
				os.Exit(1)
			}

			// We'll build the database
			update = true
		}

		// Open or create the database, with a timeout of 1 second
		// and read only if we do not have to update it.
		db, err := bbolt.Open(dbPath, 0600,
			&bbolt.Options{
				Timeout:  1 * time.Second,
				ReadOnly: !update,
			})

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		defer db.Close()

		// Update the database if needed
		if update {
			pages.Update(db)
		}

		// Now do the lookups
		if list {
			if len(args) == 1 {
				pages.List(db, args[0])

			} else {
				pages.List(db, "")
			}

		} else {
			pages.Show(db, args)
		}
	},
}

func main() {
	// Add flags
	cmd.Flags().BoolVarP(&update, "update", "u", false, "redownload pages")
	cmd.Flags().BoolVarP(&list, "list", "l", false, "list all pages (which contain the pattern)")

	// Add hidden flags
	cmd.Flags().BoolVar(&printCompletion, "completion", false, "show the bash autocompletion for tldr")
	cmd.Flags().MarkHidden("completion")

	// Change version string
	cmd.SetVersionTemplate("tldr {{.Version}} on " + targets.OsName + "\n")

	// Execute the command
	if err := cmd.Execute(); err != nil {
		// If something went wrong, exit with 1
		os.Exit(1)
	}
}
