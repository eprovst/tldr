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
	"errors"

	"github.com/elecprog/tldr/targets"
	"go.etcd.io/bbolt"
)

// pagesSource is the location from where we will download the pages.
const pagesSource = "https://tldr.sh/assets/tldr.zip"

// commonBucket is the name of the bucket containing the common pages
var commonBucket = []byte("common")

// defaultBucket is the name of the bucket containing english pages
var defaultBucket = []byte("pages")

// Get the correct buckets depending on the current configuration, note any returned value may be nil
func getBuckets(tx *bbolt.Tx) (englishCommon, englishPlatform, langCommon, langPlatform *bbolt.Bucket, err error) {
	// Open the english bucket
	english := tx.Bucket(defaultBucket)

	// Empty database?
	if english == nil {
		return nil, nil, nil, nil, nil
	}

	englishCommon = english.Bucket(commonBucket)

	if targets.OsDir != "common" {
		englishPlatform = english.Bucket([]byte(targets.OsDir))

		if englishPlatform == nil {
			return nil, nil, nil, nil, errors.New("unsupported platform '" + targets.OsDir + "'")
		}
	}

	if targets.CurrentLanguage != "en" {
		lang := tx.Bucket([]byte(targets.CurrentLanguage))

		if lang != nil {
			langCommon = lang.Bucket(commonBucket)

			if targets.OsDir != "common" {
				langPlatform = lang.Bucket([]byte(targets.OsDir))
			}
		}
	}

	return englishCommon, englishPlatform, langCommon, langPlatform, nil
}
