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
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an environment on docker",
	Long:  `Deploy an enviroment with .env set up on docker`,
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		externalip, _ := cmd.Flags().GetString("externalip")
		dockercomposefile, _ := cmd.Flags().GetString("dockercompose")
		autoupdate, _ := cmd.Flags().GetString("autoupdate")
		isDefaultEnv := false
		if env == "" {
			env = generateTempFile(configurations)
			isDefaultEnv = true
		}
		if dockercomposefile == "" {
			dockercomposefile = generateTempFile(dockercompose)
		}
		if err := godotenv.Load(env); err != nil {
			printError("Loading env variables from " + env + " cause: " + err.Error())
			os.Exit(0)
		}
		if isDefaultEnv {
			checkImagesUpdate()
		}

		if autoupdate == "true" {
			checkImagesUpdate()
		}
		if externalip == "" {
			setupIPs()
		} else {
			setupProvidedIPs(externalip)
		}
		printSetup(env, dockercomposefile)
		command := exec.Command("docker-compose",
			"-f",
			dockercomposefile,
			"--env-file="+env,
			"up",
			"-d",
			"--build",
			"rabbitmq")

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			printError("Creation of rabbitmq container failed, cause: " + err.Error())
			os.Exit(0)
		}
		printTask("Installing rabbitmq container on the machine")
		time.Sleep(15 * time.Second)
		printTask("Installing all remaining containers on the machine")
		command = exec.Command("docker-compose",
			"-f",
			dockercomposefile,
			"--env-file="+env,
			"up",
			"-d",
			"--build")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			printError("Creation of container failed, cause: " + err.Error())
			os.Exit(0)
		}
		time.Sleep(30 * time.Second)
		printTask("Restarting gateway")
		command = exec.Command("docker-compose",
			"-f",
			dockercomposefile,
			"--env-file="+env,
			"restart",
			"gateway")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			printError("Creation of container failed, cause: " + err.Error())
			os.Exit(0)
		}
		time.Sleep(5 * time.Second)
		print_urls()
	},
}

func init() {
	deployCmd.Flags().String("env", "", "Environment variable file, use default if not provided")
	deployCmd.Flags().String("externalip", "", "IP address used to expose the services, use automatically generated if not provided")
	deployCmd.Flags().String("dockercompose", "", "Docker compose file, use default if not provided")
	deployCmd.Flags().String("autoupdate", "", "Auto update the images versions (true|false)")
}
