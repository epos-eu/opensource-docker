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
//file: ./cmd/root.go
package cmd

import (
	_ "embed"

	"github.com/epos-eu/opensource-docker/cmd/methods"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "opensource-docker",
	Short:   "EPOS Open Source CLI installer",
	Version: methods.GetVersion(),
	Long:    `EPOS Open Source CLI installer to deploy the EPOS System using docker-compose`,
}

func ExecuteStandAlone() error {
	methods.GetLastTag()
	err := rootCmd.Execute()
	if err != nil {
		methods.PrintError("Error on executing rootCMD, cause: " + err.Error())
		return err
	}
	return nil
}

func Execute(args []string) error {
	methods.GetLastTag()
	if args != nil {
		rootCmd.SetArgs(args)
	}
	err := rootCmd.Execute()
	if err != nil {
		methods.PrintError("Error on executing rootCMD, cause: " + err.Error())
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(populateCmd)
	rootCmd.AddCommand(exportCmd)
}
