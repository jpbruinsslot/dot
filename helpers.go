package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/fatih/color"
)

// PrintHeader will print out a colourful header given a string
func PrintHeader(text string) {
	color := color.New(color.FgYellow, color.Bold)
	color.Printf("==> %s\n", text)
}

// PrintBody will print out a colourful body given a string
func PrintBody(text string) {
	color := color.New(color.FgCyan)
	color.Printf("    %s\n", text)
}

// PrintBodyBold will print out a colourful body give a string
func PrintBodyBold(text string) {
	color := color.New(color.FgCyan, color.Bold)
	color.Printf("    %s\n", text)
}

// PrintBodyError will print out a colourful error given a string
func PrintBodyError(text string) {
	color := color.New(color.FgRed)
	color.Printf("    Error: %s\n", text)
}

// HomeDir return the home directory of the logged in user.
// FIXME: this can trip up some operation systems
func HomeDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return user.HomeDir
}

// GetRelativePathFromCwd will remove the home folder from the current working
// directory:
// `/home/erroneousboat/dotfiles/` will become `dotfiles/`
func GetRelativePathFromCwd() (string, error) {
	// get current working directory, this returns absolute path
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// remove HomeDir() from the currentWorkingDir to get relative path
	relPath := strings.Split(currentWorkingDir, HomeDir())

	if len(relPath) != 2 {
		err := errors.New("not able to uncover relative path")
		return "", err
	} else {
		return relPath[1], nil
	}
}

// GetRelativePath will remove the home folder from the argument `fullPath`:
// `/home/erroneousboat/.nvimrc` will become `.nvimrc`
func GetRelativePath(fullPath string) (string, error) {
	relPath := strings.Split(fullPath, HomeDir())

	if len(relPath) != 2 {
		err := errors.New("not able to uncover relative path")
		return "", err
	}

	return relPath[1], nil
}

// GitCommitPush will execute the command git commit -a -m [message] and
// git push origin. This function will be called when a user specifies it
// wants to commit the changes made to its repository in the form of the
// `-p` flag used in combination with the `dot add` and `dot rm` commands.
func GitCommitPush(name, action string) {

	// load config
	c, err := NewConfig(PathDotConfig)
	if err != nil {
		message := fmt.Sprintf("not able to load config file. Make sure the " +
			".dotconfig file is present and points to the correct location")
		PrintBodyError(message)
		return
	}

	// change current working directory to DotPath
	os.Chdir(c.DotPath)

	PrintHeader("Committing changes to repository ...")

	// setup git command
	cmd := "git"

	// execute git add
	addArgs := []string{"add", "-A"}
	cmdGitAddOutput, err := exec.Command(cmd, addArgs...).Output()
	if err != nil {
		PrintBodyError("something went wrong with adding the changes " +
			"to the repository, see the output below:")
		log.Println(cmdGitAddOutput)
		log.Fatalln(err)
	}

	// setup commit arguments
	commitArgs := []string{"commit", "-a", "-m"}

	// set commit message
	var commitMessage string
	if action == "add" {
		commitMessage = fmt.Sprintf("%s: added %s for tracking", name, name)
	} else if action == "rm" {
		commitMessage = fmt.Sprintf("%s: removed %s from tracking", name, name)
	}

	// add commitMessage to the commitArgs
	commitArgs = append(commitArgs, commitMessage)

	message := fmt.Sprintf("Committing changes for: %s", name)
	PrintBody(message)

	// execute git commit
	cmdGitCommitOutput, err := exec.Command(cmd, commitArgs...).Output()
	if err != nil {
		PrintBodyError("something went wrong with commiting the changes " +
			"to the repository, see the output below:")
		log.Println(cmdGitCommitOutput)
		log.Fatalln(err)
	}

	// execute git push
	PrintBody("Pushing changes to repository")
	pushArgs := []string{"push", "origin"}
	cmdGitPushOutput, err := exec.Command(cmd, pushArgs...).Output()
	if err != nil {
		PrintBodyError("something went wrong with pushing the changes " +
			"to the repository, see the output below:")
		log.Println(cmdGitPushOutput)
		log.Fatalln(err)
	}

	return
}
