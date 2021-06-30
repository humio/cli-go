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
	"github.com/spf13/cobra"
)

func newParsersRemoveCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "remove [flags] <repo> <parser>",
		Short: "Remove (uninstall) a parser from a repository.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			parserId := args[1]

			client := NewApiClient(cmd)

			apiError := client.Parsers().Remove(repo, parserId)
			exitOnError(cmd, apiError, "Error removing parser")
		},
	}

	return &cmd
}
