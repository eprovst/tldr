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
	"errors"

	"github.com/spf13/cobra"
)

// validateArguments validates if the combination of flags and arguments is valid
func validateArguments(cmd *cobra.Command, args []string) error {
	// If we don't have to do anything special, we need at least one command
	if cmd.Flags().NFlag() == 0 && len(args) == 0 {
		return errors.New("missing argument: command")
	}

	// See if we have too many flags
	numFlags := cmd.Flags().NFlag()

	// Update and platform don't realy count
	if update {
		numFlags--
	}

	if cmd.Flag("platform").Changed {
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
}
