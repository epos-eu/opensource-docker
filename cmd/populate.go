/*
    EPOS Open Source - Local installation with Docker
    Copyright (C) 2022  EPOS ERIC

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
//file: ./cmd/populate.go
package cmd

import (
    _ "embed"
    "github.com/spf13/cobra"
    "path/filepath"
    "os"
    "log"
    "fmt"
)


var populateCmd = & cobra.Command {
    Use: "deploy",
    Short: "Deploy an environment on docker",
    Long: `Deploy an enviroment with .env set up on docker`,
    Run: func(cmd * cobra.Command, args[] string) {

        path, _ := cmd.Flags().GetString("file")
        fileInfo, err := os.Stat(path)
        if err != nil {
            log.Fatal(err)
        }

        setupIPs()

        if fileInfo.IsDir() {
            filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    log.Fatalf(err.Error())
                }
                fmt.Printf("File Name: %s\n", info.Name())
                return nil
            })
        } else {
            // is not a directory
        }
        print_urls()
    },
}

func init() {
    deployCmd.Flags().String("file", "", "--file 1")
}