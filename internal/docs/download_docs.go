// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Example HTTP URL
	scopes := []string{"core", "workspace"}
	for _, scope := range scopes {
		url := fmt.Sprintf("https://xata.io/api/openapi?scope=%s", scope)

		// Send a GET request
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Check if the request was successful (status code 200)
		if resp.StatusCode != http.StatusOK {
			log.Printf("HTTP request failed with status code %d\n", resp.StatusCode)
			return
		}

		filename := fmt.Sprintf("%s-openapi.json", scope)

		// Create or overwrite the file for writing
		file, err := os.Create(fmt.Sprintf("%s-openapi.json", scope))
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}

		// Copy the response body into the file
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			fmt.Printf("Error copying response body to file: %v\n", err)
			return
		}

		fmt.Printf("Downloaded %s to %s\n", url, filename)

		err = resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
