package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// SetupInitialMachine will setup a machine to be able to use dot it'll do the
// following:
// 1. Create a .dotconfig file
// 2. Create the files and backup folders
// 3. Add the .dotconfig for tracking
func SetupInitialMachine(pathDotConfig string) {
	// create .dotconfig file
	err := CreateDotConfigFile(pathDotConfig)
	if err != nil {
		PrintBodyError(err.Error())
		return
	}

	// create dot folders
	err = CreateDotFolders()
	if err != nil {
		PrintBodyError(err.Error())
		return
	}

	// add .dotconfig for tracking
	TrackFile("dotconfig", pathDotConfig, false)
}

// CreateDotConfigFile will create a .dotconfig file in the specified path
func CreateDotConfigFile(pathDotConfig string) error {
	// remove HomeDir() from the currentWorkingDir to get relative path
	// we need this relative path to put in the .dotconfig file
	relPath, err := GetRelativePathFromCwd()
	if err != nil {
		return err
	}

	// create .dotconfig file
	payload := fmt.Sprintf("{\"dot_path\": \"%s\", \"files\": {}}", relPath)
	err = ioutil.WriteFile(pathDotConfig, []byte(payload), 0755)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Creating new .dotconfig file: %s", pathDotConfig)
	PrintBody(message)

	return nil
}

// CreateDotFolders() will create the files and backup folders
func CreateDotFolders() error {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	folders := [2]string{
		fmt.Sprintf("%s/files", currentWorkingDir),
		fmt.Sprintf("%s/backup", currentWorkingDir),
	}

	for _, folder := range folders {
		message := fmt.Sprintf("Creating folder: %s", folder)
		PrintBody(message)
		err = os.Mkdir(folder, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
