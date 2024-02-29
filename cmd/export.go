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
//file: ./cmd/export.go
package cmd

import (
	_ "embed"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export configuration files in output folder, options: [env, compose]",
	Long:  `Export configuration files for customization in output folder, options: [env, compose]`,
	RunE: func(cmd *cobra.Command, args []string) error {

		file, _ := cmd.Flags().GetString("file")
		output, _ := cmd.Flags().GetString("output")

		switch file {
		case "env":
			if err := generateFile(configurations, output+"/configurations.env"); err != nil {
				printError("Error on generating file ENV, cause: " + err.Error())
				return err
			}
		case "compose":
			if err := generateFile(dockercompose, output+"/docker-compose.yaml"); err != nil {
				printError("Error on generating file DOCKER-COMPOSE, cause: " + err.Error())
				return err
			}
		default:
			printError("Invalid option, available options: [env, compose]")
		}
		return nil
	},
}

func init() {
	exportCmd.Flags().String("file", "", "File to export, available options: [env, compose]")
	exportCmd.MarkFlagRequired("file")
	exportCmd.Flags().String("output", "", "Full path utput folder")
	exportCmd.MarkFlagRequired("output")
}
