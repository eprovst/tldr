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
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"go.etcd.io/bbolt"
)

// Update fetches all pages and stores them in the database
func Update(database *bbolt.DB) {
	// Download the ZIP file
	zipReader, err := downloadZip(pagesSource)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// Remove all buckets from the old database
	err = database.Update(
		func(tx *bbolt.Tx) error {
			return tx.ForEach(func(bucketName []byte, _ *bbolt.Bucket) error {
				return tx.DeleteBucket(bucketName)
			})
		})

	// Did something go wrong?
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}

	// Now add the files relevant to this platform to the database
	err = database.Update(
		func(tx *bbolt.Tx) error {
			// Create a new default pages bucket
			deflt, _ := tx.CreateBucket(defaultBucket)

			if deflt == nil {
				return errors.New("failed to remove old database")
			}

			// Keep mapping from language to bucket
			langBucket := make(map[string]*bbolt.Bucket)
			langBucket["default"] = deflt

			// Add list of target buckets per language
			pageBucket := make(map[string]map[string]*bbolt.Bucket)
			pageBucket["default"] = make(map[string]*bbolt.Bucket)
			pageBucket["default"]["common"], _ = deflt.CreateBucket(commonBucket)

			// Read in all pages
			for _, file := range zipReader.File {
				// Is it a markdown file?
				if !strings.HasSuffix(path.Base(file.Name), ".md") {
					// No
					continue
				}

				// It's a Markdown file, let us put it in the right spot
				command := strings.TrimSuffix(path.Base(file.Name), ".md")
				dir := path.Dir(file.Name)

				var language, target string

				// Detect language
				// Is it the default: English?
				if strings.HasPrefix(dir, "pages/") {
					language = "default"
					target = strings.TrimPrefix(dir, "pages/")

				} else if strings.HasPrefix(dir, "pages.") {
					// Else extract the language code and target
					split := strings.SplitN(dir, "/", 2)
					language = strings.TrimPrefix(split[0], "pages.")
					target = split[1]

					// See if the language aleady exists?
					if langBucket[language] == nil {
						// No we need to create it
						nw, _ := tx.CreateBucket([]byte(language))
						langBucket[language] = nw

						// Also add a common bucket
						pageBucket[language] = make(map[string]*bbolt.Bucket)
						pageBucket[language]["common"], _ = nw.CreateBucket(commonBucket)
					}

				} else {
					// Not a page to add
					continue
				}

				// Read the page
				contents, err := file.Open()

				if err != nil {
					fmt.Println("warning:", err)
					continue
				}

				out, err := ioutil.ReadAll(contents)
				contents.Close()

				if err != nil {
					fmt.Println("warning:", err)
					continue
				}

				// Do we have to create a new bucket?
				tgtBucket, ok := pageBucket[language][target]

				if !ok {
					tgtBucket, _ = langBucket[language].CreateBucket([]byte(target))
					pageBucket[language][target] = tgtBucket
				}

				// Compress the page and write it to the bucket
				tgtBucket.Put([]byte(command), out)
			}

			// Done!
			return nil
		})

	// Has something gone wrong?
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
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
