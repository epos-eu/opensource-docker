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
//file: ./cmd/delete.go
package cmd

import (
    _ "embed"
    "github.com/spf13/cobra"
    "os/exec"
    "os"
	"github.com/joho/godotenv"
    "log"
)

var deleteCmd = & cobra.Command {
    Use: "delete",
    Short: "Delete an environment on docker",
    Long: `Delete an enviroment with .env set up on docker`,
    Run: func(cmd * cobra.Command, args[] string) {

        env, _ := cmd.Flags().GetString("env")
        dockercomposefile, _ := cmd.Flags().GetString("dockercompose")

        if env == "" {
            env = generateTempFile(configurations)
        }
        if dockercomposefile == "" {
            dockercomposefile = generateTempFile(dockercompose)
        }
        if err := godotenv.Load(env);
        err != nil {
            log.Fatal("Error loading env variables from "+env+"\n")
            log.Fatal(err)
        }

        setupIPs()

        cmd.Printf(">> Delete environment\n   >> Env file: %s \n   >> Docker compose file: %s \n   >> LocalAddress ip %s\n", env, dockercomposefile, os.Getenv("LOCAL_IP"))

        command := exec.Command("docker-compose",
            "-f",
            dockercomposefile,
            "down",
            "-v")

        cmd.Printf(command.String())
        command.Stdout = os.Stdout
        command.Stderr = os.Stderr
        if err := command.Run();
        err != nil {
            print("hello")
        }

    },
}

func init() {
    deleteCmd.Flags().String("env", "", "--env 1")
    deleteCmd.Flags().String("dockercompose", "", "--dockercompose 1")
}