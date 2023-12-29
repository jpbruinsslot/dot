Dot - simple dotfile manager
============================

`dot` provides a cli application which will streamline the process of adding
and removing files to and from your archive. Using `dot` in combination with a
hosted repository you'll be able to backup, restore and sync your dotfiles for
several machines. Why would you want your dotfiles hosted in a repository? For
the answer and much more click [here](https://dotfiles.github.io). 

`dot` works by placing a selected dotfile into the archive and symlinking it
to its original location. The archive serves as a repository and can be
persisted using a hosted git repository.

Installation
------------

#### Binary installation

[Download](https://github.com/jpbruinsslot/dot/releases) a
compatible binary for your system. For convenience, place `dot` in a
directory where you can access it from the command line. Usually this is
`/usr/local/bin`.

```bash
$ mv dot /usr/local/bin
```

#### Via Go

If you want, you can also get `dot` via Go:

```bash
$ go get -u github.com/jpbruinsslot/dot
$ cd $GOPATH/src/github.com/jpbruinsslot/dot
$ go install .
```

Setting up
----------

To begin using `dot` we need to create a folder in your home folder where we
will track the dotfiles. And initializing a git repository.

```bash
$ mkdir ~/.dotfiles
$ git init
```

Usage
-----

To start using `dot`, go to the newly created directory and issue the
following command:

```bash
$ dot sync
```

This will initialize the necessary folder and create a `.dotconfig`
configuration file which will automatically be tracked in the archive.

#### Tracking files or folders

You can use the following command to start tracking files or folders:

```bash
# dot add -name [name] -path [path/to/file]
$ dot add nvimrc /home/jpbruinsslot/.nvimrc
```

To remove a file or folder for tracking, use the following command:

```bash
# dot rm -name [name]
$ dot rm -name nvimrc
```

In order to automatically create a git commit message and push to the
repository, pass in the `-push` flag. You can use this for both the
`add` and `rm` command.

```bash
$ dot add -name nvimrc -path /home/jpbruinsslot/.nvimrc -push
$ dot rm -name nvimrc -push
```

#### Additional machines

So you've started tracking your files on one machine but now you want to use
your archive on another machine. Clone your repository on the additional
machine and use the `dot sync` command to start synchronizing your files on the
new machine.
