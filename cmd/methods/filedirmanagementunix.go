//go:build linux
// +build linux

package methods

import (
	"os"
)

func GenerateDirectoryName() string {
	dname := os.TempDir() + "/" + os.Getenv("PREFIX")
	return dname
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
	if _, err := os.Stat(os.TempDir() + "/" + os.Getenv("PREFIX")); os.IsNotExist(err) {
		err := os.Mkdir(os.TempDir()+"/"+os.Getenv("PREFIX"), 0777)
		if err != nil {
			PrintError("Could not create temporary folder, cause " + err.Error() + " error: " + err.Error())
		}
		PrintTask("Directory" + dir + " created successfully")
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
