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
//file: ./cmd/functions.go
package cmd
import (
    "log"
	"os"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "net"
)

func setupIPs() {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().( * net.UDPAddr)

    os.Setenv("API_HOST", "http://" + localAddr.IP.String() + ":" + os.Getenv("API_PORT")+os.Getenv("DEPLOY_PATH")+"/api")
    os.Setenv("EXECUTE_HOST", "http://" + localAddr.IP.String() + ":" + os.Getenv("API_PORT"))
    os.Setenv("HOST", "http://" + localAddr.IP.String() + ":" + os.Getenv("GUI_PORT"))
    os.Setenv("LOCAL_IP", localAddr.IP.String())
}

func print_urls() {

    fmt.Print(`Open Source Docker deploy

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
    Copyright (C) 2022  EPOS ERIC`);

    fmt.Print("++++++++++++++++++ EPOS ACCESS POINTS ++++++++++++++++++\n")
    fmt.Print("--------------------------------------------------------\n")
    fmt.Print("EPOS Data Portal: http://"+os.Getenv("LOCAL_IP")+":"+os.Getenv("GUI_PORT")+os.Getenv("DEPLOY_PATH")+"\n")
    fmt.Print("--------------------------------------------------------\n")
    fmt.Print("EPOS Backoffice: http://"+os.Getenv("LOCAL_IP")+":"+os.Getenv("BACKOFFICE_GUI_PORT")+os.Getenv("DEPLOY_PATH")+"\n")
    fmt.Print("--------------------------------------------------------\n")
    fmt.Print("EPOS API Gateway: http://"+os.Getenv("LOCAL_IP")+":"+os.Getenv("API_PORT")+os.Getenv("DEPLOY_PATH")+os.Getenv("API_PATH")+"/ui/\n")
    fmt.Print("--------------------------------------------------------\n")
    fmt.Print("++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
}


func generateTempFile( text []byte ) string {
    tmpFile, err := ioutil.TempFile("", fmt.Sprintf("%s-", filepath.Base(os.Args[0])))
    if err != nil {
        log.Fatal("Could not create temporary file", err)
    }
    defer tmpFile.Close()

    fmt.Println("Created temp file: ", tmpFile.Name())
    name := tmpFile.Name()

    fmt.Println("Writing some data to the temp file")
    if _, err = tmpFile.Write(text); err != nil {
        log.Fatal("Unable to write to temporary file", err)
    } else {
        fmt.Println("Data should have been written")
    }

    return name
}

func generateFile( text []byte, filePath string ) {
    err := ioutil.WriteFile(filePath, text, 0777)
    if err != nil {
        log.Fatal("Could not create file", err)
    }

    fmt.Println("Created file ", filePath)
}