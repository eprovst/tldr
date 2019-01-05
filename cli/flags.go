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

package cli

import (
	flag "github.com/spf13/pflag"
)

var (
	// Add flags
	update   = flag.BoolP("update", "u", false, "redownload pages")
	help     = flag.BoolP("help", "h", false, "help for tldr")
	platform = flag.StringP("platform", "p", "", "overide default platform")
	language = flag.StringP("language", "l", "", "overide default language")
	search   = flag.StringP("search", "s", "", "list pages matching `regex`")
	purge    = flag.Bool("purge", false, "remove database from disk")
	render   = flag.String("render", "", "render local `page`")
	version  = flag.Bool("version", false, "version for tldr")

	// Add hidden scripting flags
	printBashCompletion = flag.Bool("bash-completion", false, "show the bash autocompletion for tldr")
	list                = flag.Bool("list", false, "list all pages for the current platform")
	listAll             = flag.Bool("list-all", false, "list all available pages")
	listPlatforms       = flag.Bool("list-platforms", false, "list all supported platforms")
	listLanguages       = flag.Bool("list-languages", false, "list all supported languages")
)

func init() {
	// Mark hidden scripting flags
	flag.CommandLine.MarkHidden("bash-completion")
	flag.CommandLine.MarkHidden("list")
	flag.CommandLine.MarkHidden("list-all")
	flag.CommandLine.MarkHidden("list-platforms")
	flag.CommandLine.MarkHidden("list-languages")

	// Here for compatibility sake
	flag.BoolVarP(purge, "clear-cache", "c", false, "purge database")
	flag.CommandLine.MarkDeprecated("clear-cache", "use --purge instead")
}
