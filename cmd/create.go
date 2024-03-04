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
//file: ./cmd/create.go
package cmd

import (
	"github.com/epos-eu/opensource-docker/cmd/methods"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an environment on docker",
	Long:  `Deploy an enviroment with .env set up on docker`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, _ := cmd.Flags().GetString("env")
		externalip, _ := cmd.Flags().GetString("externalip")
		dockercomposefile, _ := cmd.Flags().GetString("dockercompose")
		autoupdate, _ := cmd.Flags().GetString("autoupdate")
		update, _ := cmd.Flags().GetString("update")
		envname, _ := cmd.Flags().GetString("envname")
		envtag, _ := cmd.Flags().GetString("envtag")

		if err := methods.CreateEnvironment(env, dockercomposefile, externalip, envname, envtag, update, autoupdate); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	deployCmd.Flags().String("env", "", "Environment variable file, use default if not provided")
	deployCmd.Flags().String("externalip", "", "IP address used to expose the services, use automatically generated if not provided")
	deployCmd.Flags().String("dockercompose", "", "Docker compose file, use default if not provided")
	deployCmd.Flags().String("envname", "", "Set name of the environment")
	deployCmd.Flags().String("envtag", "", "Set tag of the environment")
	deployCmd.Flags().String("autoupdate", "", "Auto update the images versions (true|false)")
	deployCmd.Flags().String("update", "", "Update of an existing deployment (true|false), default false")
}
