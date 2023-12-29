// file.go will hold all the operations that have to do with
// file management.

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func SyncFiles() {
	PrintHeader("Syncing files ...")

	// load config
	c, err := NewConfig(PathDotConfig)
	if err != nil {
		message := fmt.Sprintf("not able to load config file. Make sure the " +
			".dotconfig file is present and points to the correct location")
		PrintBodyError(message)
		return
	}

	// when we have files the sync them
	if len(c.Files) > 0 {
		// for every file track it
		copyAll := false
		for name, path := range c.Files {
			// get full path
			fullPath := fmt.Sprintf("%s%s", HomeDir(), path)
			copyAll = TrackFile(name, fullPath, false, copyAll)
		}
	} else {
		PrintBodyError("there aren't any files being tracked. Begin doing " +
			"so with: `dot add [name] [path]`")
	}
}

// TrackFile will track an individual file, meaning, it will move the original
// file to either the files or backup folder. It will the create a symlink of
// the file in the original location. `name` will be used as the name of the
// folder and key in the config file. `fullPath` has to be the absolute path
// to the file to be tracked.
//
// TrackFile can be called from two contexes:
//
//  1. From SyncFiles, it will read all the tracked files from the config and
//     make track files if necessary.
//
//  2. From CommandAdd, this will add a new file for tracking
//
// TrackFile will make a distinction between a new file and a file that is
// already been tracked:
//
//  1. TrackFile can't find the symlink, but the file is present in the
//     dot_path (the folder that holds all the original files). Then we need to
//     relink it, thus creating a symlink at the correct location. This happens
//     we you run dot on a new 'additional machine'.
//
//  2. TrackFile can't find the symlink, and the file is also not present in
//     the dot_path folder. This will mean that it is a new file were are going
//     to track. So we copy the file to the files folder, create a symlink, and
//     add an entry to the config file.
func TrackFile(name string, fullPath string, push bool, copyAll bool) bool {
	// load config
	c, err := NewConfig(PathDotConfig)
	if err != nil {
		PrintBodyError("not able to find .dotconfig")
		return copyAll
	}

	// Base
	base := path.Base(fullPath)

	// get relative path
	relPath, err := GetRelativePath(fullPath)
	if err != nil {
		PrintBodyError(err.Error())
		return copyAll
	}

	// check if path is present
	_, err = os.Stat(fullPath)
	if err != nil {
		message := fmt.Sprintf("file not present on system: %s", fullPath)
		PrintBodyError(message)

		if copyAll {
			src := fmt.Sprintf("%s%s/files/%s/%s", HomeDir(), c.DotPath, name, base)
			MakeAndCopyToDir(src, fullPath)
			return TrackFile(name, fullPath, push, copyAll)
		}

		PrintBody("Copy file(s) to its destination? [All/Y/N]")
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			log.Fatal(err)
		}

		switch input {
		case "All":
			copyAll = true
			src := fmt.Sprintf("%s%s/files/%s/%s", HomeDir(), c.DotPath, name, base)
			MakeAndCopyToDir(src, fullPath)
			return TrackFile(name, fullPath, push, copyAll)
		case "Y":
			src := fmt.Sprintf("%s%s/files/%s/%s", HomeDir(), c.DotPath, name, base)
			MakeAndCopyToDir(src, fullPath)
			return TrackFile(name, fullPath, push, copyAll)
		default:
			msg := fmt.Sprintf("Ignoring %s", name)
			PrintBodyError(msg)
			return copyAll
		}
	}

	// check if path is already symlinked
	s, err := os.Lstat(fullPath)
	if err != nil {
		return copyAll
	}

	if s.Mode()&os.ModeSymlink == os.ModeSymlink {
		message := fmt.Sprintf("%s is already symlinked", name)
		PrintBody(message)
		return copyAll
	}

	repoPath := fmt.Sprintf("%s%s/files/%s/", HomeDir(), c.DotPath, name)
	if _, err := os.Stat(repoPath); err == nil {
		// no symlink found, already in repo => additional machine
		message := fmt.Sprintf("Symlinking: %s", name)
		PrintBody(message)

		// put in backup folder, set named folder based on `name`, e.g.:
		// `/home/jpbruinsslot/dotfiles/backup/[name]/[base]`
		dst := fmt.Sprintf("%s%s/backup/%s/%s", HomeDir(), c.DotPath, name, base)
		err = MakeAndMoveToDir(fullPath, dst)
		if err != nil {
			msg := fmt.Sprintf("not able to move files to %s (%s)", dst, err)
			PrintBodyError(msg)

			prompt := fmt.Sprintf("Remove %s ? [Y/N]", dst)
			PrintBody(prompt)
			var input string
			_, err := fmt.Scan(&input)
			if err != nil {
				log.Fatal(err)
			}

			if input == "Y" {
				err := os.RemoveAll(dst)
				if err != nil {
					log.Fatal(err)
				}

				err = MakeAndMoveToDir(fullPath, dst)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				msg := fmt.Sprintf("ignoring %s", name)
				PrintBodyError(msg)
				return copyAll
			}
		}

		// trim potential trailing slash for symlink
		fullPath = strings.TrimRight(fullPath, "/")

		// create symlink (os.Symlink(oldname, newname))
		dst = fmt.Sprintf("%s%s/files/%s/%s", HomeDir(), c.DotPath, name, base)
		err = os.Symlink(dst, fullPath)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		// no symlink found, not in repo => new entry
		message := fmt.Sprintf("Symlinking: %s", name)
		PrintBody(message)

		// put in files folder, set named folder based on `name`, e.g.:
		// `/home/jpbruinsslot/dotfiles/files/[name]/[base]`
		dst := fmt.Sprintf("%s%s/files/%s/%s", HomeDir(), c.DotPath, name, base)
		err = MakeAndMoveToDir(fullPath, dst)
		if err != nil {
			log.Fatal(err)
		}

		// trim potential trailing slash for symlink
		fullPath = strings.TrimRight(fullPath, "/")

		// create symlink (os.Symlink(oldname, newname))
		err = os.Symlink(dst, fullPath)
		if err != nil {
			log.Fatal(err)
		}

		// create entry in .dotconfig file
		c.Files[name] = relPath
		c.Save()

		// push changes to repository
		if push {
			GitCommitPush(name, "add")
		}
	}

	return copyAll
}

// UntrackFile will remove a file from tracking. `name` will be the key
// in the config file that points to the initial location of the file
func UntrackFile(name string, push bool) {
	// open config file
	c, err := NewConfig(HomeDir() + "/" + ConfigFileName)
	if err != nil {
		PrintBodyError("not able to find .dotconfig")
		return
	}

	// check if `name` is present in c.Files
	path := c.Files[name]
	if path == "" {
		message := fmt.Sprintf("'%s' is not being tracked. Get the list of "+
			"tracked files with `dot list`", name)
		PrintBodyError(message)
		return
	}

	// check if path (the symlink) is present
	pathSymlink := fmt.Sprintf("%s%s", HomeDir(), path)
	f, err := os.Lstat(pathSymlink)
	if err != nil {
		message := fmt.Sprintf("not able to find: %s", path)
		PrintBodyError(message)
		return
	}

	// check if path is symlink
	if f.Mode()&os.ModeSymlink != os.ModeSymlink {
		message := fmt.Sprintf("%s is not a symlink", path)
		PrintBodyError(message)
		return
	}

	// check if src is present
	src := fmt.Sprintf("%s%s/files/%s%s", HomeDir(), c.DotPath, name, path)
	if _, err = os.Stat(src); err != nil {
		message := fmt.Sprintf("not able to find %s", src)
		PrintBodyError(message)
		return
	}

	// remove symlink
	err = os.Remove(pathSymlink)
	if err != nil {
		message := fmt.Sprintf("not able to remove %s", pathSymlink)
		PrintBodyError(message)
		return
	}

	// move the file or directory
	dst := fmt.Sprintf("%s%s", HomeDir(), path)

	message := fmt.Sprintf("Moving %s back to %s", name, dst)
	PrintBody(message)

	err = MakeAndCopyToDir(src, dst)
	if err != nil {
		log.Fatal(err)
	}

	// remove tracked files from repo dir
	entry := fmt.Sprintf("%s%s/files/%s", HomeDir(), c.DotPath, name)
	err = os.RemoveAll(entry)
	if err != nil {
		log.Fatal(err)
	}

	// remove entry from config and save config
	delete(c.Files, name)
	c.Save()

	// push changes to repository
	if push {
		GitCommitPush(name, "rm")
	}
}

// MakeAndMoveToDir will move the source file/folder `src` to the destination
// `dst` (`dst` will be absolute path to the destination).
func MakeAndMoveToDir(src string, dst string) error {
	err := MakeAndCopyToDir(src, dst)
	if err != nil {
		return err
	}

	err = os.RemoveAll(src)
	if err != nil {
		return err
	}

	return nil
}

func MakeAndCopyToDir(src string, dst string) error {
	// folder or file
	f, err := os.Stat(src)
	if err != nil {
		return err
	}

	if f.IsDir() {
		err = CopyDir(src, dst)
		if err != nil {
			return err
		}
	} else {
		// get directory
		dir, _ := filepath.Split(dst)

		// create destination dir
		// err = os.MkdirAll(dir, f.Mode())
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		// rename the file
		err = CopyFile(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	// Check src
	f, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !f.IsDir() {
		return errors.New("source is not a directory")
	}

	// Check dst
	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return nil
	}
	if err == nil {
		return errors.New("dst already exist")
	}

	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	os.MkdirAll(dst, f.Mode())
	if err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		if file.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			if file.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(src string, dst string) error {
	inFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		return err
	}

	err = outFile.Sync()
	if err != nil {
		return err
	}

	f, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, f.Mode())
	if err != nil {
		return err
	}

	return nil
}
