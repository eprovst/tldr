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
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/elecprog/tldr/targets"
	"go.etcd.io/bbolt"
)

// Update fetches all pages and stores them in the database
func Update(database *bbolt.DB) {
	// Download the ZIP file
	zipReader, err := downloadZip(pagesSource)

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// Purge the old database
	Purge(database)

	// Now add the files relevant to this platform to the database
	err = database.Update(
		func(tx *bbolt.Tx) error {
			// Create a new pages bucket
			root, _ := tx.CreateBucket(rootBucket)

			if root == nil {
				return errors.New("failed to remove old database")
			}

			// Add buckets for all supported platforms and for common pages
			targetBuckets := make(map[string]*bbolt.Bucket)
			targetBuckets["common"], _ = root.CreateBucket(commonBucket)

			for target := range targets.AllTargets {
				targetBuckets[target], _ = root.CreateBucket([]byte(target))
			}

			// Read in all pages
			for _, file := range zipReader.File {
				// Only add english pages
				if strings.HasPrefix(file.Name, "pages/") {
					target := strings.TrimPrefix(path.Dir(file.Name), "pages/")
					command := strings.TrimSuffix(path.Base(file.Name), ".md")

					// Read the page
					contents, err := file.Open()

					if err != nil {
						fmt.Println("warning:", err)
					}

					out, err := ioutil.ReadAll(contents)

					if err != nil {
						fmt.Println("warning:", err)
					}

					// Write the page to the correct bucket
					if targetBuckets[target] != nil {
						targetBuckets[target].Put([]byte(command), out)
					}
				}
			}

			// Done!
			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Println("warning:", err)
	}
}

func downloadZip(url string) (*zip.Reader, error) {
	// Download the ZIP file
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	// Read the entire body into a byte array
	zipFile, err := ioutil.ReadAll(resp.Body)

	// Close the body
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	// Turn this array into a zip reader
	zipReader, err := zip.NewReader(
		bytes.NewReader(zipFile),
		int64(len(zipFile)),
	)

	if err != nil {
		return nil, err
	}

	return zipReader, nil
}
