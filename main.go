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

// Flags
var (
	update  = false
	listAll = false
	list    = false
	clear   = false
	search  = ""
	render  = ""

	printCompletion = false
)

var cmd = &cobra.Command{
	Use:     "tldr [flags] command [command ...]",
	Version: "v0.2.0",
	Short:   "Go command line client for tldr",

	DisableFlagsInUseLine: true,

	Args: func(cmd *cobra.Command, args []string) error {
		// If we don't have to do anything special, we need at least one command
		if cmd.Flags().NFlag() == 0 && len(args) == 0 {
			return errors.New("missing argument: command")
		}

		// See if we have too many flags
		numFlags := cmd.Flags().NFlag()

		// Update doesn't realy count
		if update {
			numFlags--
		}

		// We can't have arguments and flags
		if numFlags > 0 && len(args) > 0 {
			return errors.New("too many arguments: expected none")
		}

		// We can't have multiple flags set
		if numFlags > 1 {
			return errors.New("at most one flag can be set")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// If we only have to print the bash completion, do so
		if printCompletion {
			fmt.Println(pages.BashCompletion)
			return
		}

		// Are we asked to render a page?
		if render != "" {
			pages.Render(render)
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
				ReadOnly: !update && !clear,
			})

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		defer db.Close()

		// Clear the database if requested
		if clear {
			pages.Clear(db)
			return
		}

		// Update the database if needed
		if update {
			pages.Update(db)
			// We might want to do other stuff though
		}

		// Do we have to list commands?
		if list || listAll {
			pages.List(db)
			return
		}

		// Search.
		if cmd.Flag("search").Changed {
			pages.Search(db, search)
			return
		}

		// No, we simply want to see tldr pages :)
		if len(args) > 0 {
			pages.Show(db, args)
		}
	},
}

func main() {
	// Add flags
	cmd.Flags().BoolVarP(&update, "update", "u", false, "redownload pages")
	cmd.Flags().BoolVar(&list, "list", false, "list all pages for the current platform")
	cmd.Flags().BoolVar(&clear, "clear-cache", false, "clear database")
	cmd.Flags().StringVarP(&search, "search", "s", "", "show all commands containing `pattern`")
	cmd.Flags().StringVar(&render, "render", "", "render local `page`")

	// Flag is currently useless
	cmd.Flags().BoolVar(&listAll, "list-all", false, "list all available pages")
	cmd.Flags().MarkHidden("list-all")

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
