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
//file: ./cmd/methods/exportfunctions.go
package methods

import (
	_ "embed"
)

func ExportVariablesEnvironment(file string, output string) error {

	switch file {
	case "env":
		if err := GenerateFile(GetConfigurationsEmbed(), output+"/configurations.env"); err != nil {
			PrintError("Error on generating file ENV, cause: " + err.Error())
			return err
		}
	case "compose":
		if err := GenerateFile(GetDockerComposeEmbed(), output+"/docker-compose.yaml"); err != nil {
			PrintError("Error on generating file DOCKER-COMPOSE, cause: " + err.Error())
			return err
		}
	default:
		PrintError("Invalid option, available options: [env, compose]")
	}
	return nil
}
