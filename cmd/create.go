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
	"regexp"
	"time"

	"github.com/joho/godotenv"
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
		envname, _ := cmd.Flags().GetString("envname")
		envtag, _ := cmd.Flags().GetString("envtag")

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

		if err := RemoveContents(dname); err != nil {
			printError("Error on removing the content from directory " + err.Error())
			return err
		}
		if err := createDirectory(dname); err != nil {
			printError("Error on creating the directory " + err.Error())
			return err
		}

		isDefaultEnv := false
		if env == "" {
			ret_env, err := generateTempFile(dname, "configurations", configurations)
			if err != nil {
				return err
			}
			env = ret_env
			isDefaultEnv = true
		}

		if dockercomposefile == "" {
			ret_dockercomposefile, err := generateTempFile(dname, "dockercompose", dockercompose)
			if err != nil {
				return err
			}
			dockercomposefile = ret_dockercomposefile
		}
		if err := godotenv.Overload(env); err != nil {
			printError("Loading env variables from " + env + " cause: " + err.Error())
			return err
		}

		if autoupdate == "true" || isDefaultEnv {
			if err := checkImagesUpdate(); err != nil {
				printError("Error on updating the docker container images " + err.Error())
				return err
			}
		}
		if externalip == "" {
			if err := setupIPs(); err != nil {
				printError("Error on setting the IPs " + err.Error())
				return err
			}
		} else {
			if err := setupProvidedIPs(externalip); err != nil {
				printError("Error on setting the IPs using the provided IP " + err.Error())
				return err
			}
		}

		printSetup(env, dockercomposefile)

		printTask("Installing rabbitmq container on the machine")

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
			return err
		}
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
			return err
		}
		printWait("Waiting for the containers to be up and running...")
		time.Sleep(40 * time.Second)
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
			return err
		}
		time.Sleep(5 * time.Second)
		print_urls()
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
}
