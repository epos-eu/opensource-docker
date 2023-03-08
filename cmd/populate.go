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
	"net/http"
    "os/exec"
    "strings"
    "net"
    "strconv"
	"github.com/joho/godotenv"
)


var populateCmd = & cobra.Command {
    Use: "populate",
    Short: "Populate the existing environment with metadata information",
    Long: `Populate the existing environment with metadata information in a specific folder`,
    Run: func(cmd * cobra.Command, args[] string) {

        path, _ := cmd.Flags().GetString("folder")
        env, _ := cmd.Flags().GetString("env")

        fileInfo, err := os.Stat(path)
        if err != nil {
            log.Fatal(err)
        }

        if env == "" {
            env = generateTempFile(configurations)
        }
        if err := godotenv.Load(env);

        err != nil {
            log.Fatal("Error loading env variables from "+env+"\n")
            log.Fatal(err)
        }

        free_port, err := GetFreePort()
        if free_port != 0 {
            fmt.Println("Free port is: ", free_port)
        } else {
            fmt.Println("Free port is not available. Error is: ", err)
        }

        free_port_string := strconv.Itoa(free_port)
        fmt.Println(free_port_string, free_port)

        setupIPs()

        if fileInfo.IsDir() {
            command := exec.Command("docker",
                    "run",
                    "-idt",
                    "--name",
                    "tmc",
                    "-p",
                    free_port_string+":80",
                    "-v",
                    strings.Trim(path," ")+":/usr/share/nginx/html",
                    "nginx")

            cmd.Println(command.String())

            command.Stdout = os.Stdout
            command.Stderr = os.Stderr
            if err := command.Run();
            err != nil {
                log.Fatal("Error creating metadata-cache container...\n")
                log.Fatal(err)
            }

            filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    log.Fatalf(err.Error())
                }
                if strings.HasSuffix(info.Name(), ".ttl") {
                    fmt.Println("Ingesting File: ", info.Name())

                    posturl := "http://"+os.Getenv("LOCAL_IP")+":"+os.Getenv("API_PORT")+os.Getenv("DEPLOY_PATH")+os.Getenv("API_PATH")+"/ingestor"
                    
                    r, err := http.NewRequest("POST", posturl, nil)
                    if err != nil {
                        panic(err)
                    }

                    r.Header.Add("accept", "*/*")
                    r.Header.Add("path", "http://"+os.Getenv("LOCAL_IP")+":"+free_port_string+"/"+info.Name())
                    r.Header.Add("securityCode", "CodiceDiTest")
                    r.Header.Add("type", "single")

                    client := &http.Client{}
                    res, err := client.Do(r)
                    if err != nil {
                        panic(err)
                    }
                    defer res.Body.Close()

                    fmt.Println("Resulting request File: ", res.Body)
                }
                return nil
            })
            command = exec.Command("docker",
                    "rm",
                    "-f",
                    "tmc")

            cmd.Println(command.String())

            command.Stdout = os.Stdout
            command.Stderr = os.Stderr
            if err := command.Run();
            err != nil {
                log.Fatal("Error deleting metadata-cache container...\n")
                log.Fatal(err)
            }
        } else {
            fmt.Println("You need to define a folder")
        }
        print_urls()
    },
}

func GetFreePort() (int, error) {
    addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
    if err != nil {
            return 0, err
    }

    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
            return 0, err
    }
    defer l.Close()
    return l.Addr().(*net.TCPAddr).Port, nil
}

func init() {
    populateCmd.Flags().String("folder", "", "Folder where ttl files are located")
    populateCmd.MarkFlagRequired("folder")
    populateCmd.Flags().String("env", "", "Environment variable file")
}