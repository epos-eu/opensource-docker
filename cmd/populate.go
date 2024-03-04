/*
   EPOS Open Source - Local installation with Docker
   Copyright (C) 2023  EPOS ERIC

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
//file: ./cmd/populate.go
package cmd

import (
	_ "embed"

	"github.com/epos-eu/opensource-docker/cmd/methods"
	"github.com/spf13/cobra"
)

var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Populate the existing environment with metadata information",
	Long:  `Populate the existing environment with metadata information in a specific folder`,
	RunE: func(cmd *cobra.Command, args []string) error {

		path, _ := cmd.Flags().GetString("folder")
		env, _ := cmd.Flags().GetString("env")

		envname, _ := cmd.Flags().GetString("envname")
		envtag, _ := cmd.Flags().GetString("envtag")
		if err := methods.PopulateEnvironment(env, path, envname, envtag); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	populateCmd.Flags().String("folder", "", "Fullpath folder where ttl files are located")
	populateCmd.MarkFlagRequired("folder")
	populateCmd.Flags().String("env", "", "Environment variable file")
	populateCmd.Flags().String("envname", "", "Set name of the environment")
	populateCmd.Flags().String("envtag", "", "Set tag of the environment")
}
