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
//file: ./cmd/methods/populatefunctions.go
package methods

import (
	_ "embed"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func PopulateEnvironment(env string, path string, envname string, envtag string) error {

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

	dname := GenerateDirectoryName()

	RemoveContents(dname)
	CreateDirectory(dname)

	fileInfo, err := os.Stat(path)
	if err != nil {
		PrintError("Loading file folder, cause: " + err.Error())
		return err
	}
	if env == "" {
		ret_env, err := GenerateTempFile(dname, "configurations", GetConfigurationsEmbed())
		if err != nil {
			return err
		}
		env = ret_env
	}

	if err := godotenv.Overload(env); err != nil {
		PrintError("Loading env variables from " + env + " cause: " + err.Error())
		return err
	}
	freePortOk := false
	free_port, err := GetFreePort()
	if err != nil {
		return err
	}
	for freePortOk {
		if free_port != 0 {
			freePortOk = true
		} else {
			PrintError("Free port is not available, cause" + err.Error())
			free_port, err = GetFreePort()
			if err != nil {
				return err
			}
		}
	}
	free_port_string := strconv.Itoa(free_port)
	if err := SetupIPs(); err != nil {
		PrintError("Error on setting the IPs " + err.Error())
		return err
	}
	if fileInfo.IsDir() {

		if err := ExecuteCommand(exec.Command("docker",
			"run",
			"-idt",
			"--name",
			"tmc",
			"-p",
			free_port_string+":80",
			"-v",
			strings.Trim(path, " ")+":/usr/share/nginx/html",
			"nginx")); err != nil {
			PrintError("Creating metadata-cache container, cause " + err.Error())
			return err
		}

		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				PrintError("Filepath error " + err.Error())
				return err
			}
			if strings.HasSuffix(info.Name(), ".ttl") {
				PrintTask("Ingestion file into database: " + info.Name())
				posturl := "http://" + os.Getenv("LOCAL_IP") + ":" + os.Getenv("API_PORT") + os.Getenv("DEPLOY_PATH") + os.Getenv("API_PATH") + "/ingestor"
				r, err := http.NewRequest("POST", posturl, nil)
				if err != nil {
					PrintError("Ingesting file into database, cause " + err.Error())
					return err
				}
				r.Header.Add("accept", "*/*")
				r.Header.Add("path", "http://"+os.Getenv("LOCAL_IP")+":"+free_port_string+"/"+info.Name())
				r.Header.Add("securityCode", "changeme")
				r.Header.Add("type", "single")
				r.Header.Add("model", "EPOS-DCAT-AP-V1")

				client := &http.Client{}
				res, err := client.Do(r)
				if err != nil {
					PrintError("Ingestion failed, cause " + err.Error())
					return err
				}
				defer res.Body.Close()
			}
			return nil
		})

		if err := ExecuteCommand(exec.Command("docker",
			"rm",
			"-f",
			"tmc")); err != nil {
			PrintError("Deleting metadata-cache container, cause " + err.Error())
			return err
		}

		if err := ExecuteCommand(exec.Command("docker",
			"restart",
			envtagname+"converter-service")); err != nil {
			PrintError("Error restarting converter service, cause " + err.Error())
			return err
		}
	} else {
		PrintError("You need to define a folder!")
		return err
	}
	PrintUrls()
	return nil
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		PrintError("Resolve TCPAddr, cause " + err.Error())
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		PrintError("Listening TCPAddr, cause " + err.Error())
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
