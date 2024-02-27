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
	"os"
	"os/exec"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment on docker",
	Long:  `Delete an enviroment with .env set up on docker`,
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		dockercomposefile, _ := cmd.Flags().GetString("dockercompose")
		envname, _ := cmd.Flags().GetString("envname")
		envtag, _ := cmd.Flags().GetString("envtag")
		isDefaultEnv := false

		envtagname := ""

		if envname != "" {
			envtagname += envname
		}
		if envtag != "" {
			envtagname += envtag
		}
		if envtagname != "" {
			envtagname += "-"
		}
		envtagname = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(envtagname, "-")

		os.Setenv("PREFIX", envtagname)

		dname := os.TempDir() + os.Getenv("PREFIX")

		RemoveContents(dname)
		createDirectory(dname)

		if env == "" {
			env = generateTempFile(dname, "configurations", configurations)
			isDefaultEnv = true
		}
		if dockercomposefile == "" {
			dockercomposefile = generateTempFile(dname, "dockercompose", dockercompose)
		}
		if err := godotenv.Load(env); err != nil {
			printError("Error loading env variables from " + env + " cause: " + err.Error())

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

		}

	},
}

func init() {
	deleteCmd.Flags().String("env", "", "Environment variable file, use default if not provided")
	deleteCmd.Flags().String("dockercompose", "", "Docker compose file, use default if not provided")
	deleteCmd.Flags().String("envname", "", "Set name of the environment")
	deleteCmd.Flags().String("envtag", "", "Set tag of the environment")
}
