package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateDotConfigFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}

	pathDotConfig := fmt.Sprintf("%s/.dotconfig", tempDir)

	// set current working directory to home dir, this will
	// make the GetRelativePathFromCwd work
	os.Chdir(HomeDir())

	// test creating .dotconfig file
	err = CreateDotConfigFile(pathDotConfig)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateDotFolders(t *testing.T) {
	// change current working directory
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	os.Chdir(tempDir)

	// create the folders
	err = CreateDotFolders()
	if err != nil {
		t.Error(err)
	}
}
