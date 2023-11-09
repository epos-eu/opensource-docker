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
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v52/github"
	"github.com/hashicorp/go-version"
	"github.com/jedib0t/go-pretty/v6/table"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func setupIPs() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		printError("Dial udp 8.8.8.8:80")
		os.Exit(0)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	val, present := os.LookupEnv("API_HOST_ENV")
	if present {
		fmt.Print(val)
	} else {
		os.Setenv("API_HOST_ENV", localAddr.IP.String())
	}

	os.Setenv("API_HOST", "http://"+os.Getenv("API_HOST_ENV")+":"+os.Getenv("API_PORT")+os.Getenv("DEPLOY_PATH")+"/api")
	os.Setenv("EXECUTE_HOST", "http://"+os.Getenv("API_HOST_ENV")+":"+os.Getenv("API_PORT"))
	os.Setenv("HOST", "http://"+os.Getenv("API_HOST_ENV")+":"+os.Getenv("GUI_PORT"))
	os.Setenv("LOCAL_IP", os.Getenv("API_HOST_ENV"))
}

func print_urls() {

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
	t.AppendRow(table.Row{"EPOS API Gateway", "http://" + os.Getenv("API_HOST_ENV") + ":" + os.Getenv("API_PORT") + os.Getenv("DEPLOY_PATH") + os.Getenv("API_PATH") + "/ui/"})
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	fmt.Println(t.Render())
}

func printError(message string) {
	fmt.Println(string(colorRed), "[ERROR] "+message, string(colorReset))
}
func printTask(message string) {
	fmt.Println(string(colorGreen), "[TASK] "+message, string(colorReset))
}

func generateTempFile(text []byte) string {
	tmpFile, err := ioutil.TempFile("", fmt.Sprintf("%s-", filepath.Base(os.Args[0])))
	if err != nil {
		printError("Could not create temporary file, cause " + err.Error())
		os.Exit(0)
	}
	defer tmpFile.Close()
	name := tmpFile.Name()
	if _, err = tmpFile.Write(text); err != nil {
		printError("Unable to write to temporary file, cause " + err.Error())
		os.Exit(0)
	}
	return name
}

func generateFile(text []byte, filePath string) {
	err := ioutil.WriteFile(filePath, text, 0777)
	if err != nil {
		printError("Could not create file, cause " + err.Error())
		os.Exit(0)
	}
}

func getLastTag() {
	client := github.NewClient(nil)
	tags, _, err := client.Repositories.ListTags(context.Background(), "epos-eu", "opensource-docker", nil)
	if err != nil {
		printError("Could not retrieve tags of the repository, cause " + err.Error())
		os.Exit(0)
	}
	if len(tags) > 0 {
		latestTag := tags[0]
		currentVersion := getVersion()
		v1, _ := version.NewVersion(currentVersion)
		v2, _ := version.NewVersion(latestTag.GetName())
		if v1.LessThan(v2) {
			fmt.Println(string(colorPurple), "New version available "+v1.String()+" ---> "+v2.String(), string(colorReset))
		}
	}
}

func printSetup(env string, dockercomposefile string) {

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

func checkImagesUpdate() {
	t := table.NewWriter()
	t.SetTitle("Docker Images updated")
	t.AppendHeader(table.Row{"Default Image", "New Image"})
	for _, env := range os.Environ() {
		splitted := strings.Split(env, "=")
		if strings.HasSuffix(splitted[0], "_IMAGE") && splitted[0] != "MESSAGE_BUS_IMAGE" && splitted[0] != "REDIS_IMAGE" {
			imageRepositoryName := strings.Split(splitted[1], ":")
			latestImageTag := getLastDockerImageTag(imageRepositoryName[0])
			if latestImageTag != imageRepositoryName[1] {
				os.Setenv(splitted[0], imageRepositoryName[0]+":"+latestImageTag)
				t.AppendRow(table.Row{splitted[1], imageRepositoryName[0] + ":" + latestImageTag})
			}
		}
	}
	t.SetStyle(table.StyleColoredRedWhiteOnBlack)
	fmt.Println(t.Render())
}

func getLastDockerImageTag(repo string) string {
	response := Response{}
	namespace := "epos"
	resp, err := http.Get("https://hub.docker.com/v2/repositories/" + namespace + "/" + repo + "/tags?page_size=2")
	if err != nil {
		printError("Can't retrieve tags from dockerhub")
		os.Exit(0)
	}
	json.NewDecoder(resp.Body).Decode(&response)
	defer resp.Body.Close()
	return response.Results[1].Name
}

func getVersion() string {
	return "0.3.4"
}
