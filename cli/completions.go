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

// bashCompletion is a bashcompletion script for tldr
const bashCompletion = `#!/bin/bash
_tldr_completion()
{
	if [[ "${COMP_WORDS[$COMP_CWORD - 1]}" == "-"* ]]; then
		if [[ "${COMP_WORDS[$COMP_CWORD - 1]}" == "-p" || "${COMP_WORDS[$COMP_CWORD - 1]}" == "--platform" ]]; then
			COMPREPLY=($(compgen -W "$(tldr --list-platforms)" -- ${COMP_WORDS[$COMP_CWORD]}))
		else
			COMPREPLY=()
		fi
	else
		if [[ "${COMP_WORDS[$COMP_CWORD]}" == "-"* ]]; then
			COMPREPLY=($(compgen -W "--help --platform --purge --render --search --update --version" -- ${COMP_WORDS[$COMP_CWORD]}))
		else
			COMPREPLY=($(tldr --search "^${COMP_WORDS[$COMP_CWORD]}" 2> /dev/null))
		fi
	fi
}

complete -o default -F _tldr_completion tldr`
