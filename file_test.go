package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSyncFiles(t *testing.T) {
}

func TestTrackFile(t *testing.T) {
}

func TestUntrackFile(t *testing.T) {
}

// Test if MakeAndMoveToDir will be able to move a directory, we wil also
// test it when there is a file inside the directory
func TestMakeAndMoveToDirDirectory(t *testing.T) {
	// create, and change to current working
	// directory to temporary directory
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	os.Chdir(tempDir)

	// create temporary directory with a temporary file, also create a
	// temporary directory in which the first one is going to be placed,
	// we will end up with:
	//
	// tempDir
	// |_ dirOne
	//    |_ tempFile
	// |_ dirTwo
	dirOne, err := ioutil.TempDir(tempDir, "dirOne")
	dirTwo, err := ioutil.TempDir(tempDir, "dirTwo")

	// create a temporary file in dirOne
	tempFile, err := ioutil.TempFile(dirOne, "tempfile")
	if err != nil {
		t.Error(err)
	}

	// we want the base of dirOne
	baseDirOne := filepath.Base(dirOne)

	// we want the file name of the tempFile
	_, fileNameTempFile := filepath.Split(tempFile.Name())

	// set destination (/dirTwo/dirOne/)
	dst := fmt.Sprintf("%s/%s", dirTwo, baseDirOne)

	// move dirOne inside dirTwo
	err = MakeAndMoveToDir(dirOne, dst)
	if err != nil {
		t.Error(err)
	}

	// we should have:
	// tempDir
	// |_ dirTwo
	//    |_ dirOne
	//       |_ tempFile

	// check if dirOne is present in dirTwo
	resultDir := fmt.Sprintf("%s/%s", dirTwo, baseDirOne)
	if _, err = os.Stat(resultDir); err != nil {
		t.Error(err)
	}

	// check if file is also present in the correct dir
	resultFile := fmt.Sprintf("%s/%s/%s", dirTwo, baseDirOne, fileNameTempFile)
	if _, err = os.Stat(resultFile); err != nil {
		t.Error(err)
	}

	// check if directory is gone from source
	if _, err = os.Stat(dirOne); err == nil {
		t.Error(err)
	}
}

// Test if MakeAndMoveToDir will be able to move a file
func TestMakeAndMoveToDirFile(t *testing.T) {
	// create, and change to current working
	// directory to temporary directory
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	os.Chdir(tempDir)

	// create temporary directory with a temporary file, also create a
	// temporary directory in which the first one is going to be placed,
	// we will end up with:
	//
	// tempDir
	// |_ dirOne
	// |  |_ tempFile
	// |_ dirTwo
	dirOne, err := ioutil.TempDir(tempDir, "dirOne")
	dirTwo, err := ioutil.TempDir(tempDir, "dirTwo")

	// create a temporary file which we will move
	tempFile, err := ioutil.TempFile(dirOne, "tempfile")
	if err != nil {
		t.Error(err)
	}

	// we want the base of dirOne
	baseDirOne := filepath.Base(dirOne)

	// we want the file name of the tempFile
	_, fileNameTempFile := filepath.Split(tempFile.Name())

	// set destination (/dirTwo/dirOne/tempFile)
	dst := fmt.Sprintf("%s/%s/%s", dirTwo, baseDirOne, fileNameTempFile)

	// move the file
	err = MakeAndMoveToDir(tempFile.Name(), dst)
	if err != nil {
		t.Error(err)
	}

	// we should end up with the following:
	// tempDir
	// |_ dirTwo
	//    |_ dirOne
	//       |_ tempFile

	// check if file is at destination
	result := fmt.Sprintf("%s/%s/%s", dirTwo, baseDirOne, fileNameTempFile)
	if _, err = os.Stat(result); err != nil {
		t.Error(err)
	}

	// check if file is gone from source
	if _, err = os.Stat(tempFile.Name()); err == nil {
		t.Error(err)
	}
}
