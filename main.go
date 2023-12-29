package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	version = "0.3.0"
)

var (
	syncCmd = flag.NewFlagSet("sync", flag.ExitOnError)
	addCmd  = flag.NewFlagSet("add", flag.ExitOnError)
	rmCmd   = flag.NewFlagSet("rm", flag.ExitOnError)
	listCmd = flag.NewFlagSet("list", flag.ExitOnError)

	// Flags for 'add' command
	addName = addCmd.String("name", "", "Name for the data")
	addPath = addCmd.String("path", "", "Path to the data")
	addPush = addCmd.Bool("push", false, "Push changes to a git repository")

	// Flags for 'rm' command
	rmName = rmCmd.String("name", "", "Name of the data to remove")
	rmPush = rmCmd.Bool("push", false, "Push changes to a git repository")
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "sync":
		syncCmd.Parse(os.Args[2:])

		if len(syncCmd.Args()) > 0 {
			printUsage()
			os.Exit(1)
		}

		CommandSync()
	case "add":
		addCmd.Parse(os.Args[2:])

		if *addName == "" || *addPath == "" {
			addCmd.PrintDefaults()
			os.Exit(1)
		}

		CommandAdd(*addName, *addPath, *addPush, false)
	case "rm":
		rmCmd.Parse(os.Args[2:])

		if *rmName == "" {
			rmCmd.PrintDefaults()
			os.Exit(1)
		}

		CommandRemove(*rmName, *rmPush)
	case "list":
		listCmd.Parse(os.Args[2:])

		if len(listCmd.Args()) > 0 {
			printUsage()
			os.Exit(1)
		}

		CommandList()
	default:
		printUsage()
		os.Exit(0)
	}
}

func printUsage() {
	usage := fmt.Sprintf(`Dot - simple dotfile manager

Usage:

    dot [command] [arguments]

Version: %s

Commands:

    sync    syncs all files that are being tracked
    add     add a file or folder for tracking
    rm      remove a file from tracking
    list    list all files that are being tracked

Use "dot [command] -help" for more information about a command.
`, version)

	print(usage)
}
