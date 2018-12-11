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
	"time"

	"github.com/elecprog/tldr/info"
	"github.com/elecprog/tldr/pages"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

var update = false

var cmd = &cobra.Command{
	Use:     "tldr [-h] [-u] command [command ...]",
	Version: "v0.1.0",
	Short:   "Go command line client for tldr",

	DisableFlagsInUseLine: true,

	Args: func(cmd *cobra.Command, args []string) error {
		// If we do not have to update the database:
		// we need at least one argument
		if !update && len(args) == 0 {
			return errors.New("missing argument: command")
		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// See if it's a fresh database
		update = update || !info.Exists(info.GetDatabasePath())

		// Open or create the database, with a timeout of 1 second
		// and readonly if we do not have to update it.
		db, err := bbolt.Open(info.GetDatabasePath(), 0600,
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
		pages.Show(db, args)
	},
}

func main() {
	cmd.Flags().BoolVarP(&update, "update", "u", false, "redownload pages")

	cmd.SetVersionTemplate("tldr {{.Version}} on " + info.OsName + "\n")

	// Execute the command
	if err := cmd.Execute(); err != nil {
		// If something went wrong, exit with 1
		os.Exit(1)
	}
}
