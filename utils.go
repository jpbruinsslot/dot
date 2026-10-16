package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// PrintHeader will print out a colourful header given a string
func PrintHeader(text string) {
	fmt.Printf("==> %s\n", text)
}

// PrintBody will print out a colourful body given a string
func PrintBody(text string) {
	fmt.Printf("... %s\n", text)
}

// PrintBodyError will print out a colourful error given a string
func PrintBodyError(text string) {
	fmt.Printf("... ERROR: %s\n", text)
}

// HomeDir return the home directory of the logged in user.
func HomeDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	return dir
}

// HomePath returns the absolute path of `relPath` inside the home folder.
// `relPath` may be stored with or without a leading slash (older .dotconfig
// files use `/.config/nvim`, newer ones `.config/nvim`).
func HomePath(relPath string) string {
	return filepath.Join(HomeDir(), relPath)
}

// GetRelativePath will remove the home folder from the argument `fullPath`:
//
// `/home/jpbruinsslot/.config/nvim` will become `.config/nvim`
func GetRelativePath(fullPath string) (string, error) {
	relPath, err := filepath.Rel(HomeDir(), filepath.Clean(fullPath))
	if err != nil {
		return "", err
	}

	if relPath == ".." || strings.HasPrefix(relPath, "../") {
		return "", fmt.Errorf("%s is not inside the home folder", fullPath)
	}

	return relPath, nil
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

	PrintHeader("Committing changes to repository ...")

	// setup git command
	cmd := "git"

	// execute git add
	addArgs := []string{"add", "-A"}
	gitAddCmd := exec.Command(cmd, addArgs...)
	gitAddCmd.Dir = c.Dir()
	cmdGitAddOutput, err := gitAddCmd.CombinedOutput()
	if err != nil {
		PrintBodyError("something went wrong with adding the changes " +
			"to the repository, see the output below:")
		log.Println(string(cmdGitAddOutput))
		log.Fatalln(err)
	}

	// setup commit arguments
	commitArgs := []string{"commit", "-a", "-m"}

	// set commit message
	var commitMessage string
	switch action {
	case "add":
		commitMessage = fmt.Sprintf("%s: added %s for tracking", name, name)
	case "rm":
		commitMessage = fmt.Sprintf("%s: removed %s from tracking", name, name)
	}

	// add commitMessage to the commitArgs
	commitArgs = append(commitArgs, commitMessage)

	message := fmt.Sprintf("Committing changes for: %s", name)
	PrintBody(message)

	// execute git commit
	gitCommitCmd := exec.Command(cmd, commitArgs...)
	gitCommitCmd.Dir = c.Dir()
	cmdGitCommitOutput, err := gitCommitCmd.CombinedOutput()
	if err != nil {
		PrintBodyError("something went wrong with commiting the changes " +
			"to the repository, see the output below:")
		log.Println(string(cmdGitCommitOutput))
		log.Fatalln(err)
	}

	// execute git push
	PrintBody("Pushing changes to repository")
	pushArgs := []string{"push", "origin"}
	gitPushCmd := exec.Command(cmd, pushArgs...)
	gitPushCmd.Dir = c.Dir()
	cmdGitPushOutput, err := gitPushCmd.CombinedOutput()
	if err != nil {
		PrintBodyError("something went wrong with pushing the changes " +
			"to the repository, see the output below:")
		log.Println(string(cmdGitPushOutput))
		log.Fatalln(err)
	}
}
