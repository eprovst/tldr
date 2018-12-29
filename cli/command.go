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

package cli

import (
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
	update        = false
	list          = false
	listAll       = false
	listPlatforms = false
	purge         = false
	platform      = ""
	search        = ""
	render        = ""

	printCompletion = false
)

var tldr = &cobra.Command{
	Use:     "tldr [flags] command [command ...]",
	Version: "v0.3.1",
	Short:   "Go command line client for tldr",
	Args:    validateArguments,

	DisableFlagsInUseLine: true,

	Run: func(cmd *cobra.Command, args []string) {
		// If we only have to print the bash completion, do so
		if printCompletion {
			fmt.Println(bashCompletion)
			return
		}

		// Are we asked to render a page?
		if cmd.Flag("render").Changed {
			pages.Render(render)
			return
		}

		// Get the path where the database is/should be stored
		dbPath := getDatabasePath()

		// See if it's a first run
		if !pathExists(dbPath) {
			// Create the folder where the database will reside
			err := os.MkdirAll(filepath.Dir(dbPath), 0777)

			if err != nil {
				fmt.Fprintln(os.Stderr, "error: ", err)
				os.Exit(1)
			}

			// We'll build the database
			update = true
		}

		// Purge the database if requested
		if purge {
			err := os.Remove(dbPath)

			if err != nil {
				fmt.Fprintln(os.Stderr, "error: ", err)
				os.Exit(1)
			}

			return
		}

		// We open/create the databse with a timeout of one second
		// to not keep on attempting if there is something wrong.
		// The database is opened as read only if we do not have
		// to update it.
		// A pages size of 128 seems to result in the smallest
		// database size for our purposes, setting it so low
		// would result in a major slowdown in large databases
		// but ours is small.
		db, err := bbolt.Open(dbPath, 0600,
			&bbolt.Options{
				Timeout:  1 * time.Second,
				ReadOnly: !update,
				PageSize: 128,
			})

		if err != nil {
			fmt.Fprintln(os.Stderr, "error: ", err)
			os.Exit(1)
		}

		defer db.Close()

		// Overide the operating system
		if cmd.Flag("platform").Changed {
			if platform == "common" {
				fmt.Fprintln(os.Stderr, "error: common is not a platform")
				os.Exit(1)
			}

			targets.OsDir = platform
		}

		// Update the database if needed
		if update {
			pages.Update(db)
			// We might want to do other stuff though
		}

		// Other actions
		// Do we have to list commands?
		if list {
			pages.List(db)
			return
		}

		if listAll {
			pages.ListAll(db)
			return
		}

		// List platforms
		if listPlatforms {
			pages.ListPlatforms(db)
			return
		}

		// Search?
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

func init() {
	// Add flags
	tldr.Flags().StringVarP(&platform, "platform", "p", "", "overide default platform")
	tldr.Flags().BoolVar(&purge, "purge", false, "remove database from disk")
	tldr.Flags().StringVar(&render, "render", "", "render local `page`")
	tldr.Flags().BoolVarP(&update, "update", "u", false, "redownload pages")
	tldr.Flags().StringVarP(&search, "search", "s", "", "list pages matching `regex`")

	// Here for compatibility sake
	tldr.Flags().BoolVarP(&purge, "clear-cache", "c", false, "purge database")
	tldr.Flags().MarkDeprecated("clear-cache", "use --purge instead")

	// Add hidden scripting flags
	tldr.Flags().BoolVar(&printCompletion, "bash-completion", false, "show the bash autocompletion for tldr")
	tldr.Flags().MarkHidden("bash-completion")
	tldr.Flags().BoolVar(&list, "list", false, "list all pages for the current platform")
	tldr.Flags().MarkHidden("list")
	tldr.Flags().BoolVar(&listAll, "list-all", false, "list all available pages")
	tldr.Flags().MarkHidden("list-all")
	tldr.Flags().BoolVar(&listPlatforms, "list-platforms", false, "list all supported platforms")
	tldr.Flags().MarkHidden("list-platforms")

	// Change version string
	tldr.SetVersionTemplate("tldr {{.Version}} on " + targets.OsName + "\n")
}

// Run runs the tldr command
func Run() {
	// Execute the command
	if err := tldr.Execute(); err != nil {
		// If something went wrong, exit with 1
		os.Exit(1)
	}
}
