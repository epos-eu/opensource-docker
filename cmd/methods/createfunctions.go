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
//file: ./cmd/methods/createfunctions.go
package methods

import (
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

func CreateEnvironment(env string, dockercomposefile string, externalip string, envname string, envtag string, update string, autoupdate string) error {

	if update != "true" && update != "false" {
		update = "false"
	}

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
		PrintError("Error on removing the content from directory " + err.Error())
		return err
	}
	if err := CreateDirectory(dname); err != nil {
		PrintError("Error on creating the directory " + err.Error())
		return err
	}

	isDefaultEnv := false
	if env == "" {
		ret_env, err := GenerateTempFile(dname, "configurations", GetConfigurationsEmbed())
		if err != nil {
			return err
		}
		env = ret_env
		isDefaultEnv = true
	}

	if dockercomposefile == "" {

		ret_dockercomposefile, err := GenerateTempFile(dname, "dockercompose", GetDockerComposeEmbed())
		if err != nil {
			return err
		}
		dockercomposefile = ret_dockercomposefile
	}
	if err := godotenv.Overload(env); err != nil {
		PrintError("Loading env variables from " + env + " cause: " + err.Error())
		return err
	}

	if autoupdate == "true" || isDefaultEnv {
		if err := CheckImagesUpdate(); err != nil {
			PrintError("Error on updating the docker container images " + err.Error())
			return err
		}
	}
	if err := OverridePorts(update); err != nil {
		PrintError("Error during overriding ports if update=true " + err.Error())
		return err
	}

	if externalip == "" {
		if err := SetupIPs(); err != nil {
			PrintError("Error on setting the IPs " + err.Error())
			return err
		}
	} else {
		if err := SetupProvidedIPs(externalip); err != nil {
			PrintError("Error on setting the IPs using the provided IP " + err.Error())
			return err
		}
	}

	PrintSetup(env, dockercomposefile)

	PrintTask("Installing rabbitmq container on the machine")

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
		PrintError("Creation of rabbitmq container failed, cause: " + err.Error())
		return err
	}
	time.Sleep(15 * time.Second)
	PrintTask("Installing all remaining containers on the machine")
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
		PrintError("Creation of container failed, cause: " + err.Error())
		return err
	}
	PrintWait("Waiting for the containers to be up and running...")
	time.Sleep(40 * time.Second)
	PrintTask("Restarting gateway")
	command = exec.Command("docker-compose",
		"-f",
		dockercomposefile,
		"--env-file="+env,
		"restart",
		"gateway")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		PrintError("Creation of container failed, cause: " + err.Error())
		return err
	}
	time.Sleep(5 * time.Second)
	PrintUrls()
	return nil
}
