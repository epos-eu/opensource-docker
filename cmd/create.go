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
//file: ./cmd/create.go
package cmd

import (
    "github.com/spf13/cobra"
	"github.com/joho/godotenv"
    "os/exec"
    "os"
    "log"
    "time"
)


var deployCmd = & cobra.Command {
    Use: "deploy",
    Short: "Deploy an environment on docker",
    Long: `Deploy an enviroment with .env set up on docker`,
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

        cmd.Printf(">> Deploy environment\n   >> Env file: %s \n   >> Docker compose file: %s \n   >> LocalAddress ip %s\n", env, dockercomposefile, os.Getenv("LOCAL_IP"))
      
        command := exec.Command("docker-compose",
            "-f",
            dockercomposefile,
            "--env-file=" + env,
            "up",
            "-d",
            "--build",
            "rabbitmq")

        cmd.Printf(command.String())

        command.Stdout = os.Stdout
        command.Stderr = os.Stderr
        if err := command.Run();
        err != nil {
            log.Fatal("Error creating rabbitmq container...\n")
            log.Fatal(err)
        }

        cmd.Printf("Installing rabbitmq container on the machine...\n")
        time.Sleep(8 * time.Second)

        cmd.Printf("Installing all remaining containers on the machine...\n")
        command = exec.Command("docker-compose",
            "-f",
            dockercomposefile,
            "--env-file=" + env,
            "up",
            "-d",
            "--build")
        cmd.Printf(command.String())
        command.Stdout = os.Stdout
        command.Stderr = os.Stderr
        if err := command.Run();
        err != nil {
            log.Fatal("Error creating containers...\n")
            log.Fatal(err)
        }

        time.Sleep(20 * time.Second)

        print_urls()

    },
}

func init() {
    deployCmd.Flags().String("env", "", "Environment variable file, use default if not provided")
    deployCmd.Flags().String("dockercompose", "", "Docker compose file, use default if not provided")
}