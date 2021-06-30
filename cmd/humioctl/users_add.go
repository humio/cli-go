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
	"os"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newUsersAddCmd() *cobra.Command {
	var rootFlag boolPtrFlag
	var firstNameFlag, lastNameFlag, companyFlag, emailFlag, countryCodeFlag stringPtrFlag
	var pictureFlag urlPtrFlag

	cmd := cobra.Command{
		Use:   "add [flags] <username>",
		Short: "Adds a user. [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			username := args[0]

			client := NewApiClient(cmd)

			user, err := client.Users().Add(username, api.UserChangeSet{
				IsRoot:      rootFlag.value,
				FirstName:   firstNameFlag.value,
				LastName:    lastNameFlag.value,
				Company:     companyFlag.value,
				CountryCode: countryCodeFlag.value,
				Email:       emailFlag.value,
				Picture:     pictureFlag.value,
			})

			if err != nil {
				cmd.Println(fmt.Errorf("Error creating the user: %s", err))
				os.Exit(1)
			}

			printUserTable(cmd, user)
		},
	}

	cmd.Flags().Var(&rootFlag, "root", "If true grants root access to the user.")
	cmd.Flags().Var(&firstNameFlag, "first-name", "The first name of the user.")
	cmd.Flags().Var(&lastNameFlag, "last-name", "The last name of the user.")
	cmd.Flags().Var(&countryCodeFlag, "country-code", "A two letter country code.")
	cmd.Flags().Var(&companyFlag, "company", "The company where the user works.")
	cmd.Flags().Var(&pictureFlag, "picture", "A URL to an avatar for user.")
	cmd.Flags().Var(&emailFlag, "email", "The user's email. Does not change the username if the email is used.")

	return &cmd
}
