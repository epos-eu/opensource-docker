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
    "os"
    _ "embed"

    "github.com/spf13/cobra"
)


var (

    //go:embed "docker-compose/docker-compose.yaml"
    dockercompose []byte
    
    //go:embed "configurations/env.env"
    configurations []byte
    
    )

var rootCmd = & cobra.Command {
    Use: "epos-docker-cli",
    Short: "EPOS Open Source CLI installer",
    Version: "1.0.0",
    Long: `EPOS Open Source CLI installer to deploy the EPOS System using docker-compose`,
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(deployCmd)
    rootCmd.AddCommand(deleteCmd)
    rootCmd.AddCommand(populateCmd)
    rootCmd.AddCommand(exportCmd)
}
