package shell

var helpShort string = `ls - list directory contents
cd - change current directory
pwd - print current directory
whoami - print current user
rev - reverse text or file contents
mv - move or rename files and directories
help - show help for commands
exit - exit emulator`

var helpDetails = map[string]string{
	"ls": `NAME
    ls - list directory contents

SYNOPSIS
    ls [path...]

DESCRIPTION
    Lists files and directories in the virtual file system.

    If no path is specified, ls lists the contents of the current directory.

    If one or more paths are specified, each path is processed separately.
    A file path prints the file name.
    A directory path prints the directory contents.

    Multiple directory arguments are displayed with directory headers.

EXAMPLES
    ls
    ls /home/user
    ls /home/user/notes.txt
    ls /home/user /tmp`,

	"cd": `NAME
    cd - change the current directory

SYNOPSIS
    cd [path]
    cd -

DESCRIPTION
    Changes the current working directory.

    If path is omitted, cd changes to /home/user.

    The following path forms are supported:
        .       current directory
        ..      parent directory
        /       root directory
        ~       home directory
        ~/path  path relative to home

    "cd -" changes to the previous working directory.

ERRORS
    An error is returned if the path does not exist, if the target is not
    a directory, or if too many arguments are provided.

EXAMPLES
    cd /home/user/docs
    cd ..
    cd ~
    cd -
    cd`,

	"pwd": `NAME
    pwd - print the current working directory

SYNOPSIS
    pwd

DESCRIPTION
    Prints the path of the current working directory in the virtual file system.

EXAMPLE
    pwd`,

	"whoami": `NAME
    whoami - print the current user name

SYNOPSIS
    whoami

DESCRIPTION
    Prints the user name from the emulator prompt.

    The user name is taken from the part of the prompt before the '@' character.

    This command does not accept arguments.

ERRORS
    If any arguments are provided, the command prints:
        usage: whoami

EXAMPLE
    whoami`,

	"rev": `NAME
    rev - reverse characters in text or files

SYNOPSIS
    rev
    rev file...

DESCRIPTION
    Reverses the characters of each input line.

    If one or more file paths are specified, rev reads each file from the
    virtual file system and prints every line with its characters reversed.

    If no files are specified, rev enters interactive mode. In this mode,
    each entered line is reversed immediately.

    Press Ctrl+D to leave interactive rev mode and return to the normal shell.

ERRORS
    An error is returned if a specified path does not exist or refers to
    a directory.

EXAMPLES
    rev /home/user/notes.txt
    rev file1.txt file2.txt
    rev`,

	"mv": `NAME
    mv - move or rename files and directories

SYNOPSIS
    mv source destination

DESCRIPTION
    Moves or renames a file or directory inside the virtual file system.

    If destination is an existing directory, source is moved into that
    directory and keeps its current name.

    If destination does not exist, the final component of destination is used
    as the new name. The source may therefore be renamed, moved, or both.

    If destination is an existing file and source is also a file, the
    destination file is replaced by source.

    A directory cannot overwrite a file.

    The root directory cannot be moved.

    A directory cannot be moved into itself or into one of its descendants.

    All changes are made only in memory and do not modify the physical VFS
    source file on disk.

ERRORS
    An error is returned if:
        source does not exist;
        the destination parent does not exist;
        the destination parent is not a directory;
        the root directory is used as source;
        a directory is moved into itself;
        a directory is moved into one of its descendants;
        a directory attempts to overwrite a file;
        the number of arguments is invalid.

EXAMPLES
    mv notes.txt renamed.txt
    mv notes.txt /home/user/docs
    mv notes.txt /home/user/docs/renamed.txt
    mv old.txt new.txt`,

	"help": `NAME
    help - display information about shell commands

SYNOPSIS
    help
    help command

DESCRIPTION
    Without arguments, prints the list of available shell commands together
    with a short description of each command.

    If a command name is specified, prints detailed help for that command,
    including its usage, behavior, errors, and examples.

EXAMPLES
    help
    help mv
    help rev`,

	"exit": `NAME
    exit - exit the emulator

SYNOPSIS
    exit

DESCRIPTION
    Closes the emulator shell.

    The command does not accept arguments.

ERRORS
    If any arguments are provided, the command prints:
        usage: exit

EXAMPLE
    exit`,
}

func (s *Shell) helpCmd(args []string) string {
	if len(args) > 1 {
		return "usage: help [command]"
	}

	if len(args) == 1 {
		text, ok := helpDetails[args[0]]
		if !ok {
			return "help: no help topics match"
		}

		return text
	}

	return helpShort
}
