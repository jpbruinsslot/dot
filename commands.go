package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

// CommandSync will do several things depending configuration:
//  1. Sync files when there is a .dotconfig present in the correct location
//  2. Create a .dotconfig in the correct location when it isn't
//  3. Create a new setup of dot, including a .dotconfig, files and backup
//     folders
func CommandSync() {
	// get current working directory
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// path to .dotconfig in current working dir
	pathDotConfigCwd := fmt.Sprintf(
		"%s/files/dotconfig/%s", currentWorkingDir, ConfigFileName)

	// here we try to uncover 3 possibilities:
	// 1. .dotconfig (symlink) is already on machine in correct location
	// 2. .dotconfig (regular file) is already on machine
	// 3. .dotconfig not in home folder but in current working directory
	// 4. .dotconfig not in home folder and not in current working directory
	if _, err := os.Lstat(PathDotConfig); err == nil {

		// .dotconfig (symlink) found in home dir => SyncFiles
		PrintBody("The .dotconfig file is present")

		// relink everything
		SyncFiles()

	} else if _, err := os.Stat(HomeDir() + "/" + ConfigFileName); err == nil {
		// .dotconfig (regular file, not symlinked) found in home dir =>
		// symlink .dotconfig
		PrintBody("Found .dotconfig file in home folder")

		// make sure .dotconfig is present in DotPath
		if _, err := os.Stat(pathDotConfigCwd); err != nil {
			PrintBodyError("couldn't find .dotconfig in your archive, make sure it is present")
			return
		}

		// remove found .dotconfig
		err = os.Remove(HomeDir() + "/" + ConfigFileName)
		if err != nil {
			log.Fatal(err)
		}

		// make symlink for .dotconfig
		dotconfigOld := fmt.Sprintf("%s/files/dotconfig/%s",
			currentWorkingDir, ConfigFileName)

		dotconfigNew := fmt.Sprintf("%s/%s", HomeDir(), ConfigFileName)

		err = os.Symlink(dotconfigOld, dotconfigNew)
		if err != nil {
			log.Fatal(err)
		}

		// relink everything
		SyncFiles()
	} else if _, err := os.Stat(pathDotConfigCwd); err == nil {

		// .dotconfig not found in home dir,
		// .dotconfig found in current working dir => symlink .dotconfig
		PrintBody("Found .dotconfig file in repository folder")

		// make a symlink for .dotconfig file
		dotconfigOld := fmt.Sprintf("%s/files/dotconfig/%s",
			currentWorkingDir, ConfigFileName)

		dotconfigNew := fmt.Sprintf("%s/%s", HomeDir(), ConfigFileName)

		err = os.Symlink(dotconfigOld, dotconfigNew)
		if err != nil {
			log.Fatal(err)
		}

		// relink everything
		SyncFiles()
	} else {

		// .dotconfig not found in home dir,
		// .dotconfig not found in current working dir => new setup
		PrintBody("Couldn't find the .dotconfig file, do you want to create a new one? [Y/N]")

		// get input
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			log.Fatal(err)
		}

		if input == "y" || input == "Y" {
			// setup initial machine
			// create new .dotconfig file
			SetupInitialMachine(PathDotConfig)

			PrintBody("You're now ready to use dot! Type 'dot -help' for help")
		} else {
			return
		}
	}
}

// CommandAdd will add a file or folder for tracking.
func CommandAdd(name, path string, push, force bool) {
	PrintHeader("Adding new entry for tracking ...")
	_ = TrackFile(name, path, push, false)
}

// CommandRemove will remove a file from tracking.
func CommandRemove(name string, push bool) {
	PrintHeader("Removing entry from tracking ...")
	UntrackFile(name, push)
}

// CommandList will output the list of files that are being tracked by dot.
func CommandList() {
	PrintHeader("Following files are being tracked by dot ...")

	// open config file
	config, err := NewConfig(HomeDir() + "/" + ConfigFileName)
	if err != nil {
		PrintBodyError("not able to find .dotconfig")
		return
	}

	// check if there is anything to display
	if len(config.Files) == 0 {
		PrintBodyError(
			"there are no files being tracked. Begin doing so, with `dot add -name [name] -path [path]`",
		)
		return
	}

	// print out the tracked files
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "name\tpath")
	for name, path := range config.Files {
		line := fmt.Sprintf("%s\t%s%s", name, HomeDir(), path)
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
