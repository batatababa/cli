# **CLI**
**Version  0.1.0**

## Whats New?
All the things!

## Overview
The cli package was created to be a quick, simple, and extensible command line package. Almost everything in it is exported so, that it may be easily modified. The underlying structure for your command line may be defined however you like, but below is a short example of how it is intended to be used.

**Example**
```
    // In my project a define my commands in other files grouped by category
    var BrownCow *cli.Command = &cli.Command{
        Name: "cow",
    }
    var BrownFox *cli.Command = &cli.Command{
        Name: "fox",
    }
    var Brown *cli.Command = &cli.Command{
        Name:        "brown",
        Description: "silly example command",
        Usage:       "use a brown what?",
        SubCommands: []cli.Command{*BrownCow, *BrownFox},
    }
    
    func main() {
        tree := cli.NewCommandTree()
        tree.Root = *Brown
        tree.Version = "1.0.0"
    
        cli.Run(os.Args, &tree)
    }
```

In the example a command tree struture is defined by first defining stand alone commands and then tying them together. In main() the call to cli.Run() finds the command in the tree that corresponds to the user input specified by os.Args, an array of strings. cli.Run() will verify that the user input is valid and then call the Action associated with the command that was found. It is within this Action function that you will need to do the work of your program. The Command object that is passed into the Action function contains all of the Arguments, Flags, and Options that were specfied by the user.

## Expectations
I wrote this package, because I was unable to find a simple cli package that treated all subcommands equally. For example, I came across packages that would only print help for leafs in the command tree. Anyways, because it was written just to be used in a larger project, I most likely won't be adding any new features or accepting many pull requests. At this time, bug fixes are the only expected changes that will be made. That said, if you think something is missing from the package, feel free to submit a pull request or feature requests. If it is in line with the intentions of the package, it may make it in.

## Building
Requires Go 1.6 due to a vendored dependency. The dependency is very small and should go away before cli 1.0.0. Currently CLI uses the go tool to build.

## Command Trees
Command trees define the structure of your application's command line. A command tree is made up primarily of the following:
* **Top Level Properties:** Properties that apply to your application.
* **Root:** The root of the actual tree of commands that make up your program.
* **Shared Parameters:** Contains Arguments, Flags, and Options to be applied to all commands. It provides function hooks, "PreAction & PostAction", to be run before and after all commands.

 
**Example**
```
tree := cli.NewCommandTree()
tree.Root = rootCommandOfTree
tree.Version = "1.0.0"
```
## Commands
A command is made up primarily of the following:
* **Flags, Arguments, Options:** Define the accepted input to the command
* **SubCommands:** Define commands nested under the command
* **Action:** Function executed for this command

For a command, the arguments, flags, and options given by the user are checked against the ones defined on the command object.

## Flags
Flags provide a binary input to the program.
* User provides flags in either short form "-h" and/or long form "--help". 
* There is support for combining the short form flags like in the common usage of the linux tar commnand, "tar -xvf filename".
* There is support for multiple occurences of the same flag as in "-hhh" or "-h -h -h".

**Example**
```
var Brown *cli.Command = &cli.Command{
    Name:        "brown",
    Description: "silly example command",
    Usage:       "use a brown what?",
    Flags:       []cli.Flag{
        {
            ShortName:   "a",
            LongName:    "LongA",
            Description: "Description for Flag A",
        },
    },
}
```

## Options
Options provide a key value pair input to the program.
* User provides an options in either short form "-k value", "-k=value" and/or long form "--key value", "--key=value". 
* There is support for multiple occurences of the same option as in "-k value -k value".

```
var Brown *cli.Command = &cli.Command{
    Name:        "brown",
    Description: "silly example command",
    Usage:       "use a brown what?",
    Opts:        []cli.Option{
        {
            ShortName:   "a",
            LongName:    "LongA",
            Description: "Description for Opt A",
        },
    },
}
```
## Arguments
Arguments provide a value input into a program.
* User provides an argument as a string, "programName argValue"
```
var Brown *cli.Command = &cli.Command{
    Name:        "brown",
    Description: "silly example command",
    Usage:       "use a brown what?",
    Args:        []cli.Argument{
        {
            Name:        "Arg1",
            Description: "Description for Arg1",
        },
    },
}
```
## Sub Commands
Subcmmands are commands nested inside of other commands. It is often useful to use subcommands to break commands into categories.
```
var BrownFox *cli.Command = &cli.Command{}
var BrownCow *cli.Command = &cli.Command{}

var Brown *cli.Command = &cli.Command{
    Name:        "brown",
    Description: "silly example command",
    Usage:       "use a brown what?",
    SubCommands: []cli.Command{*BrownCow, *BrownFox},
}
```
## Auto Help 
By default automatic help generation is turned on. The user can use "-h", "--help", or "?" to show the help for a command. Automatic help can be turned off on the Command Tree.
```
tree := cli.NewCommandTree()
tree.AutoHelp = false
```
## Hiding Help
Even with autohelp turned on, you can turn off help for individual commands.
```
var BrownCow *cli.Command = &cli.Command{
    Name:     "brown",
    HideHelp: true,
}
```

## Shared Inputs
Contains Arguments, Flags, and Options to be applied to all commands. It provides hooks, "PreAction & PostAction", to allow functions to be run before and after all commands.
```
tree := cli.NewCommandTree()

tree.Shared = cli.SharedParameters{
    Flags: []cli.Flag{
        {
            ShortName:   "v",
            LongName:    "verbose",
            Description: " Maybe all commands need a verbose flag.te",
        },
    },
}
```

## Pre-Actions and Post-Actions
Pre-Actions and Post-Actions provide a way to run functions before and after the action of a command. It can be specified using the "shared" property on the command tree.
```
tree := cli.NewCommandTree()
tree.Shared = cli.SharedParameters{
    PreAction: func(c Command) (err error) {
        // do stuff
        return
    },
}
```