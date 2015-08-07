package main

import (
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"testing"
)

const payload string = `
{
	"dot_path": "/path/to/dotfiles",
	"files": {
		"test_file_1": "path-to-test-file-1",
		"test_file_2": "path-to-test-file-2"
	}
}
`

// Define here the setup and teardown functions
func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	os.Exit(retCode)
}

// setUp() will create a config file for the other tests so we can test the
// methods associated with the config struct
func setUp() {
	// create a test dotconfig file
	f, err := ioutil.TempFile("", ".dotconfig")
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Unlink(f.Name())

	ioutil.WriteFile(
		f.Name(),
		[]byte(payload),
		0644,
	)

	// try to read the file
	_, err = NewConfig(f.Name())
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewConfig(t *testing.T) {
	// create a test dotconfig file
	f, err := ioutil.TempFile("", ".dotconfig")
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Unlink(f.Name())

	ioutil.WriteFile(
		f.Name(),
		[]byte(payload),
		0644,
	)

	// try to read the file
	c, err := NewConfig(f.Name())
	if err != nil {
		log.Fatal(err)
	}

	// c.DotPath should be the same as payload
	if c.DotPath != "/path/to/dotfiles" {
		t.Error("c.DotPath doesn't match")
	}

	// c.Files should be the same as payload
	if c.Files["test_file_1"] != "path-to-test-file-1" {
		t.Error("c.Files doesn't match")
	}

	// c.Files should be the same as payload
	if c.Files["test_file_2"] != "path-to-test-file-2" {
		t.Error("c.Files doesn't match")
	}
}
