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
//file: ./cmd/functions.go
package methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/hashicorp/go-version"
	"github.com/jedib0t/go-pretty/v6/table"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"
const colorBlue = "\033[34m"
const colorPurple = "\033[35m"
const colorCyan = "\033[36m"
const colorWhite = "\033[37m"

type Response struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Creator int `json:"creator"`
		ID      int `json:"id"`
		Images  []struct {
			Architecture string    `json:"architecture"`
			Features     string    `json:"features"`
			Variant      any       `json:"variant"`
			Digest       string    `json:"digest"`
			Os           string    `json:"os"`
			OsFeatures   string    `json:"os_features"`
			OsVersion    any       `json:"os_version"`
			Size         int       `json:"size"`
			Status       string    `json:"status"`
			LastPulled   any       `json:"last_pulled"`
			LastPushed   time.Time `json:"last_pushed"`
		} `json:"images"`
		LastUpdated         time.Time `json:"last_updated"`
		LastUpdater         int       `json:"last_updater"`
		LastUpdaterUsername string    `json:"last_updater_username"`
		Name                string    `json:"name"`
		Repository          int       `json:"repository"`
		FullSize            int       `json:"full_size"`
		V2                  bool      `json:"v2"`
		TagStatus           string    `json:"tag_status"`
		TagLastPulled       any       `json:"tag_last_pulled"`
		TagLastPushed       time.Time `json:"tag_last_pushed"`
		MediaType           string    `json:"media_type"`
		ContentType         string    `json:"content_type"`
		Digest              string    `json:"digest"`
	} `json:"results"`
}

func OverridePorts(update string) error {

	if update == "true" {
		PrintNotification("No ports check needed, update=true")
		return nil
	} else {
		ports := [2]string{"API_PORT", "DATA_PORTAL_PORT"}

		for i := 0; i < len(ports); i++ {
			PrintNotification("Checking availability of " + ports[i] + " " + os.Getenv(ports[i]))
			isPortAvailable, err := IsPortAvailable(os.Getenv(ports[i]))
			if err != nil {
				PrintError("Problem on retrieving the availability for the port for " + ports[i] + " error: " + err.Error())
				return err
			}
			if isPortAvailable {
				PrintNotification("Port " + ports[i] + " " + os.Getenv(ports[i]) + " available")
			} else {
				port, err := GetAvailablePort()
				if err != nil {
					PrintError("Problem on assigning a free port for " + ports[i] + " error: " + err.Error())
					return err
				}
				os.Setenv(ports[i], port)
				PrintNotification("Port " + ports[i] + " " + os.Getenv(ports[i]) + " available")
			}
		}
	}
	return nil
}

func SetupIPs() error {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		PrintError("Dial udp 8.8.8.8:80")

	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	val, present := os.LookupEnv("API_HOST_ENV")
	if present {
		os.Setenv("API_HOST_ENV", val)
	} else {
		os.Setenv("API_HOST_ENV", localAddr.IP.String())
	}

	os.Setenv("API_HOST", "http://"+os.Getenv("API_HOST_ENV")+":"+os.Getenv("API_PORT")+os.Getenv("DEPLOY_PATH")+os.Getenv("API_PATH"))
	os.Setenv("EXECUTE_HOST", "http://"+os.Getenv("API_HOST_ENV")+":"+os.Getenv("API_PORT"))
	os.Setenv("HOST", "http://"+os.Getenv("API_HOST_ENV")+":"+os.Getenv("DATA_PORTAL_PORT"))
	os.Setenv("LOCAL_IP", os.Getenv("API_HOST_ENV"))
	return nil
}

func SetupProvidedIPs(externalip string) error {
	os.Setenv("API_HOST", "http://"+externalip+":"+os.Getenv("API_PORT")+os.Getenv("DEPLOY_PATH")+"/api")
	os.Setenv("EXECUTE_HOST", "http://"+externalip+":"+os.Getenv("API_PORT"))
	os.Setenv("HOST", "http://"+externalip+":"+os.Getenv("DATA_PORTAL_PORT"))
	os.Setenv("LOCAL_IP", externalip)
	return nil
}

func PrintUrls() {

	fmt.Println(string(colorCyan), `Open Source Docker deploy

     &&&&&&&&&&&&&&&&&& *&&&&&&&%&&&%               *****************               &&&&&&/         
     &&&&&&&&&&&&&&&&&& *&&&&&&&&&&&&&&&&&       **  **********  *******       &&&&&&&&&&&&&&&&&    
     &&&&&&&&&&&%&&&&&& *&&&&&&&%    &&&&&&&   ,************     *********    &&%&&&&&&&&&&&&&      
     &&&&&&             *&&&&&&        &&&&&( ************   **   ********** &&&&&&#                
     &&&&&&             *&&&&&&(       &&&&& ****** * *****  **  *********** &&&&&&&&#              
     &&&&&&&&&&&&&&&&.  *&&&&&&&&&&&&&&&&&&& *******   *   , *    *********** &&&&&&&&&&&&&&&&      
     &&&%&&&&&&&%&&&&.  *&&&&&&&%&&&&&&&%&   *******                 ,*******    &&&&&&&%&&&&&&&    
     &&&&&&             *&&&&&&               *                   , ********              &&&&&&.   
     &&&&&&             *&&&&&&               .    ******  *,    ******* **    &&         &&&&&&    
     &&&&&&&&&&&&&&&&&& *&&&&&&                 ************** *         *   &&&&&&&&&&&&&&&&&&&    
     &&&&&&&&&&&%&&&&&& *&&&&&&                   ************* ,*******     &&&%&&&&&&&&&&&&       
                                                      **************                             
    Copyright (C) 2023  EPOS ERIC`, string(colorReset))
	t := table.NewWriter()
	t.SetTitle("EPOS ACCESS POINTS")
	t.AppendRow(table.Row{"EPOS Data Portal", "http://" + os.Getenv("API_HOST_ENV") + ":" + os.Getenv("DATA_PORTAL_PORT")})
	t.AppendRow(table.Row{"EPOS API Gateway", "http://" + os.Getenv("API_HOST_ENV") + ":" + os.Getenv("API_PORT") + os.Getenv("DEPLOY_PATH") + os.Getenv("API_PATH") + "/ui/"})
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	fmt.Println(t.Render())
}

func PrintError(message string) {
	fmt.Println(string(colorRed), "[ERROR] "+message, string(colorReset))
}
func PrintTask(message string) {
	fmt.Println(string(colorGreen), "[TASK] "+message, string(colorReset))
}
func PrintWait(message string) {
	fmt.Println(string(colorYellow), "[WAITING] "+message, string(colorReset))
}
func PrintNotification(message string) {
	fmt.Println(string(colorPurple), "[NOTIFICATION] "+message, string(colorReset))
}
func PrintNewVersionAvailable(message string) {
	fmt.Println(string(colorYellow), "[NEW VERSION AVAILABLE] "+message, string(colorReset))
}

func GenerateTempFile(dname string, filetype string, text []byte) (string, error) {

	tmpFile, err := os.CreateTemp(dname, filetype)
	if err != nil {
		PrintError("Could not create temporary file, cause " + err.Error() + " error: " + err.Error())
		return "", err
	}
	defer tmpFile.Close()
	name := tmpFile.Name()
	if _, err = tmpFile.Write(text); err != nil {
		PrintError("Unable to write to temporary file, cause " + err.Error() + " error: " + err.Error())
		return "", err
	}
	PrintNotification("File " + name + " created successfully")

	return name, nil
}

func CreateDirectory(dir string) error {
	if _, err := os.Stat(os.TempDir() + os.Getenv("PREFIX")); os.IsNotExist(err) {
		err := os.Mkdir(os.TempDir()+os.Getenv("PREFIX"), 0777)
		if err != nil {
			PrintError("Could not create temporary folder, cause " + err.Error() + " error: " + err.Error())
		}
		PrintTask("Directory" + dir + " created successfully")
	} else {
		PrintNotification("Directory " + dir + " already exists, using it")
	}
	return nil
}

func RemoveContents(dir string) error {
	if _, err := os.Stat(os.TempDir() + os.Getenv("PREFIX")); err == nil {
		d, err := os.Open(dir)
		if err != nil {
			return err
		}
		defer d.Close()
		names, err := d.Readdirnames(-1)
		if err != nil {
			return err
		}
		for _, name := range names {
			err = os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}
	} else {
		PrintNotification("Directory " + dir + " already exists, using it")
	}
	return nil
}

func GenerateFile(text []byte, filePath string) error {
	err := os.WriteFile(filePath, text, 0777)
	if err != nil {
		PrintError("Could not create file, cause " + err.Error())
		return err
	}
	return nil
}

func GetLastTag() error {
	client := github.NewClient(nil)
	tags, _, err := client.Repositories.ListTags(context.Background(), "epos-eu", "opensource-docker", nil)
	if err != nil {
		PrintError("Could not retrieve tags of the repository, cause " + err.Error())
		return err
	}
	if len(tags) > 0 {
		latestTag := tags[0]
		currentVersion := GetVersion()
		v1, _ := version.NewVersion(currentVersion)
		v2, _ := version.NewVersion(latestTag.GetName())
		if v1.LessThan(v2) {
			PrintNewVersionAvailable(v1.String() + " ---> " + v2.String())
		}
	}
	return nil
}

func PrintSetup(env string, dockercomposefile string) {

	t := table.NewWriter()
	t.SetTitle("DEPLOY CONFIGURATION")
	t.AppendHeader(table.Row{"Name", "Value"})
	for _, env := range os.Environ() {
		splitted := strings.Split(env, "=")
		if strings.HasSuffix(splitted[0], "_IMAGE") {
			t.AppendRow(table.Row{splitted[0], os.Getenv(splitted[0])})
		}
	}
	t.AppendRow(table.Row{"Env File", env})
	t.AppendRow(table.Row{"Docker Compose File", dockercomposefile})
	t.AppendRow(table.Row{"Local IP", os.Getenv("API_HOST_ENV")})
	t.SetStyle(table.StyleColoredGreenWhiteOnBlack)

	fmt.Println(t.Render())

}

func CheckImagesUpdate() error {
	t := table.NewWriter()
	t.SetTitle("Docker Images updated")
	t.AppendHeader(table.Row{"Default Image", "New Image"})
	for _, env := range os.Environ() {
		splitted := strings.Split(env, "=")
		if strings.HasSuffix(splitted[0], "_IMAGE") && splitted[0] != "MESSAGE_BUS_IMAGE" && splitted[0] != "REDIS_IMAGE" {
			imageRepositoryName := strings.Split(splitted[1], ":")
			latestImageTag, err := GetLastDockerImageTag(imageRepositoryName[0])
			if err != nil {
				return err
			}
			if latestImageTag != imageRepositoryName[1] {
				os.Setenv(splitted[0], imageRepositoryName[0]+":"+latestImageTag)
				t.AppendRow(table.Row{splitted[1], imageRepositoryName[0] + ":" + latestImageTag})
			}
		}
	}
	t.SetStyle(table.StyleColoredRedWhiteOnBlack)
	fmt.Println(t.Render())
	return nil
}

func GetLastDockerImageTag(repo string) (string, error) {
	response := Response{}
	namespace := "epos"
	resp, err := http.Get("https://hub.docker.com/v2/repositories/" + namespace + "/" + repo + "/tags?page_size=2")
	if err != nil {
		PrintError("Can't retrieve tags from dockerhub, error: " + err.Error())
		return "", err
	}
	json.NewDecoder(resp.Body).Decode(&response)
	defer resp.Body.Close()
	return response.Results[1].Name, nil
}

func IsPortAvailable(port string) (bool, error) {
	portInt, err := strconv.Atoi(port)
	if err != nil || portInt < 1 || portInt > 65535 {
		return false, err
	}
	ln, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		return true, nil
	}
	defer ln.Close()
	return false, nil
}

func GetAvailablePort() (string, error) {
	const maxAttempts = 10
	for i := 0; i < maxAttempts; i++ {
		ln, err := net.Listen("tcp", ":0")
		if err != nil {
			return "", err
		}
		defer ln.Close()
		addr := ln.Addr().String()
		parts := strings.Split(addr, ":")
		port := parts[len(parts)-1]
		return port, nil
	}
	return "", fmt.Errorf("could not find an available port")
}

func GetVersion() string {
	return "1.1.2"
}
