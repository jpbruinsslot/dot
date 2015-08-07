dot - simple dotfiles tracking
------------------------------

`dot` provides a cli application which will streamline the process of adding
and removing files to and from your archive. Using `dot` in combination with
an hosted repository you'll be able to backup, restore and sync your dotfiles
for several machines. Why would you want your dotfiles hosted in a
repository? For the answer and much more click
[here](https://dotfiles.github.io). 

`dot` works by placing a selected dotfile into the archive and symlinking it
to its original location. The archive serves as a repository and can be
persisted using a hosted git repository.

# Beta
Take note that the current version is still in beta. Bugs are to be expected,
and if you uncover one please add an
[issue](https://github.com/erroneousboat/dot/issues).

# Installation
For this version make sure you have [Go](golang.org) installed. And run the
following commands:

```bash
$ go get github.com/erroneousboat/dot
$ go install github.com/erroneousboat/dot
$ dot --help
```

# Setting up
To begin using `dot` we need to create a folder in your home folder where we
will track the dotfiles. And initializing a git repository.

```bash
$ mkdir ~/dotfiles/
$ git init
```

# Usage
To start using `dot`, go to the newly created directory and issue the
following command:

```bash
$ dot up
```

This will initialize the necessary folder and create a `.dotconfig`
configuration file which will automatically be tracked in the archive.

## Tracking files or folders
You can use the following command to start tracking files or folders:

```bash
# dot add [name] [path/to/file]
$ dot add nvimrc /home/erroneousboat/.nvimrc
```

To remove a file or folder for tracking, use the following command:
```bash
# dot rm [name]
$ dot rm nvimrc
```

In order to automatically create a git commit message and push to the
repository, pass in the `-p` or `--push` flag. You can use this for both the
`add` and `remove` command.

```bash
$ dot add -p nvimrc /home/erroneousboat/.nvimrc
```

## Additional machines
So you've started tracking your files on one machine but now you want to use
your archive on another machine. Clone your repository on the additional
machine and use the `dot up` command to start synchronizing your files on the
new machine.
