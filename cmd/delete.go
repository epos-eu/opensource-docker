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
//file: ./cmd/delete.go
package cmd

import (
	_ "embed"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment on docker",
	Long:  `Delete an enviroment with .env set up on docker`,
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		dockercomposefile, _ := cmd.Flags().GetString("dockercompose")
		isDefaultEnv := false
		if env == "" {
			env = generateTempFile(configurations)
			isDefaultEnv = true
		}
		if dockercomposefile == "" {
			dockercomposefile = generateTempFile(dockercompose)
		}
		if err := godotenv.Load(env); err != nil {
			printError("Error loading env variables from " + env + " cause: " + err.Error())
			os.Exit(0)
		}
		if isDefaultEnv {
			checkImagesUpdate()
		}
		setupIPs()
		printSetup(env, dockercomposefile)
		command := exec.Command("docker-compose",
			"-f",
			dockercomposefile,
			"down",
			"-v")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			printError("Error deleting environment, cause: " + err.Error())
			os.Exit(0)
		}

	},
}

func init() {
	deleteCmd.Flags().String("env", "", "Environment variable file, use default if not provided")
	deleteCmd.Flags().String("dockercompose", "", "Docker compose file, use default if not provided")
}
