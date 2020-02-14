# mark

Mark files and execute commands on them.

**mark** was created because every time I wanted to mass move or copy files in a terminal, something always went wrong.

Example scenario:

```sh
# I want all the .h
$ mark *.h
# And for some reason the main.c
$ mark main.c
$ cd ../some/other/place
# I want to move the files here
# Lets be sure what I am about to do 
$ mark
/the/original/place/name.h
/the/original/place/other.h
/the/original/place/main.cpp
# Ok lets do it
$ mark -mv
# done
```

There are the shortcuts `-cp`, `-mv`, `-rm` and a general `-exec` that allows you to pass every marked file into a shell command:

```mark
mark -exec 'curl -F file=@{} http:/... && echo Uploaded: {}'
```

Also If some king of selector is available (eg. fzf, dmenu), it can be used by piping them into mark:

```sh
ls | fzf -m | mark -
```

## Usage

**mark** keeps a file with the absolute paths of all the marked files. A selection of items can be passed by a stdin pipe or as arguments, and the operation (add, remove, toggle) to be over the selection can be specified by the argument `-op`.

- `-, -c, -clear, --clear`: remove all marked items.
- `-p, -op, --op` followed by one of `+, add, -, sub, ^, xor`: set the operation.
- `-e, -exec, --exec` followed by a shell command: set the command to execute with every marked item.

In the shell commands that are provided to the `-exec`:
- `{}` is replaced with one of the marked items.
- `{ }` is replace with all the marked items separated by spaces.
- If there is no `{}` or `{ }` and the shell command does not end with a `;`, a ` {}` is implied at the end.
- If only the shell command contains one or more `{ }` and no `{}`, it will be executed a single time.-

The command line arguments `-cp`, `-mv`, `-rm` are shortcuts to `-exec 'cp { } .'`, `-exec 'mv { } .'` and `-exec 'rm { }'`.

For all the command line arguments check the `mark --help`.

## Install

```
git clone https://github.com/MaanooAk/mark.git
cd mark
go install
```

Or

```
go get -u github.com/MaanooAk/mark
```

