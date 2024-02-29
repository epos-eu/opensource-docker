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
//file: ./cmd/populate.go
package cmd

import (
	_ "embed"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "Populate the existing environment with metadata information",
	Long:  `Populate the existing environment with metadata information in a specific folder`,
	Run: func(cmd *cobra.Command, args []string) {

		path, _ := cmd.Flags().GetString("folder")
		env, _ := cmd.Flags().GetString("env")

		envname, _ := cmd.Flags().GetString("envname")
		envtag, _ := cmd.Flags().GetString("envtag")

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

		RemoveContents(dname)
		createDirectory(dname)

		fileInfo, err := os.Stat(path)
		if err != nil {
			printError("Loading file folder, cause: " + err.Error())
		}
		if env == "" {
			env = generateTempFile(dname, "configurations", configurations)
		}
		if err := godotenv.Overload(env); err != nil {
			printError("Loading env variables from " + env + " cause: " + err.Error())

		}
		freePortOk := false
		free_port, err := GetFreePort()
		for freePortOk {
			if free_port != 0 {
				freePortOk = true
			} else {
				printError("Free port is not available, cause" + err.Error())
				free_port, err = GetFreePort()
			}
		}
		free_port_string := strconv.Itoa(free_port)
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
				strings.Trim(path, " ")+":/usr/share/nginx/html",
				"nginx")
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Run(); err != nil {
				printError("Creating metadata-cache container, cause " + err.Error())
			}

			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatalf(err.Error())
				}
				if strings.HasSuffix(info.Name(), ".ttl") {
					printTask("Ingestion file into database: " + info.Name())
					posturl := "http://" + os.Getenv("LOCAL_IP") + ":" + os.Getenv("API_PORT") + os.Getenv("DEPLOY_PATH") + os.Getenv("API_PATH") + "/ingestor"
					r, err := http.NewRequest("POST", posturl, nil)
					if err != nil {
						printError("Ingesting file into database, cause " + err.Error())
					}
					r.Header.Add("accept", "*/*")
					r.Header.Add("path", "http://"+os.Getenv("LOCAL_IP")+":"+free_port_string+"/"+info.Name())
					r.Header.Add("securityCode", "changeme")
					r.Header.Add("type", "single")
					r.Header.Add("model", "EPOS-DCAT-AP-V1")

					client := &http.Client{}
					res, err := client.Do(r)
					if err != nil {
						printError("Ingestion failed, cause " + err.Error())

					}
					defer res.Body.Close()
				}
				return nil
			})
			command = exec.Command("docker",
				"rm",
				"-f",
				"tmc")
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Run(); err != nil {
				printError("Deleting metadata-cache container, cause " + err.Error())

			}
		} else {
			printError("You need to define a folder!")

		}
		command := exec.Command("docker",
			"exec",
			"redis-server",
			"redis-cli",
			"FLUSHALL")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			printError("Flushing redis container, cause " + err.Error())

		}
		print_urls()
	},
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		printError("Resolve TCPAddr, cause " + err.Error())

	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		printError("Listening TCPAddr, cause " + err.Error())

	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func init() {
	populateCmd.Flags().String("folder", "", "Fullpath folder where ttl files are located")
	populateCmd.MarkFlagRequired("folder")
	populateCmd.Flags().String("env", "", "Environment variable file")
	populateCmd.Flags().String("envname", "", "Set name of the environment")
	populateCmd.Flags().String("envtag", "", "Set tag of the environment")
}
