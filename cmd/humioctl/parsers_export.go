// Copyright © 2018 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func newParsersExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <repo> <parser>",
		Short: "Export a parser <parser> in <repo> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			parserId := args[1]

			if outputName == "" {
				outputName = parserId
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			yamlData, apiErr := client.Parsers().Export(repo, parserId)
			if apiErr != nil {
				cmd.Println(fmt.Errorf("Error fetching parsers: %s", apiErr))
				os.Exit(1)
			}

			outFilePath := outputName + ".yaml"

			writeErr := ioutil.WriteFile(outFilePath, []byte(yamlData), 0644)
			if writeErr != nil {
				cmd.Println(fmt.Errorf("Error saving the parser file: %s", writeErr))
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the parser should be written. Defaults to ./<parser-name>.yaml")

	return &cmd
}
