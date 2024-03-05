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
//file: ./cmd/methods/deletefunctions.go
package methods

import (
	_ "embed"
	"os"
	"os/exec"
	"regexp"

	"github.com/joho/godotenv"
)

func DeleteEnvironment(env string, dockercomposefile string, envname string, envtag string) error {
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

	if err := RemoveContents(dname); err != nil {
		PrintError("Error on removing the content from directory " + err.Error())
		return err
	}
	if err := CreateDirectory(dname); err != nil {
		PrintError("Error on creating the directory " + err.Error())
		return err
	}

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
		PrintError("Error loading env variables from " + env + " cause: " + err.Error())
		return err
	}
	if isDefaultEnv {
		CheckImagesUpdate()
	}

	SetupIPs()

	PrintSetup(env, dockercomposefile)

	if err := ExecuteCommand(exec.Command("docker-compose",
		"-f",
		dockercomposefile,
		"down",
		"-v")); err != nil {
		PrintError("Creation of container failed, cause: " + err.Error())
		return err
	}
	return nil
}
