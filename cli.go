package cli

import (
	"errors"
	"fmt"
	"strings"
)

type SharedParameters struct {
	Flags      []Flag
	Args       []Argument
	ArgSets    []ArgumentSet
	Opts       []Option
	PreAction  func(com Command) error
	PostAction func(com Command) error
}

type CommandTree struct {
	Root         Command
	Shared       SharedParameters
	Author       string
	Copyright    string
	Email        string
	Version      string
	AutoHelp     bool
	ToHelpString func(c Command) string
}

func NewCommandTree() (tree CommandTree) {
	tree.AutoHelp = true
	tree.ToHelpString = ToHelpString
	return tree
}

var autoHelpFlag = Flag{
	ShortName:   "h",
	LongName:    "help",
	Description: "Show help",
}
var autoHelpArg = Argument{
	Name:        "?",
	Value:       "?",
	Description: "Show help",
}

func Run(appArgs []string, tree *CommandTree) (err error) {
	if tree.AutoHelp {
		tree.Shared.Args = append(tree.Shared.Args, autoHelpArg)
		tree.Shared.Flags = append(tree.Shared.Flags, autoHelpFlag)
	}

	fullCom, err := tree.FindCommand(appArgs)

	if err != nil {
		return err
	}

	fullCom.Flags = append(fullCom.Flags, tree.Shared.Flags...)
	fullCom.Args = append(fullCom.Args, tree.Shared.Args...)
	fullCom.ArgSets = append(fullCom.ArgSets, tree.Shared.ArgSets...)
	fullCom.Opts = append(fullCom.Opts, tree.Shared.Opts...)

	userCom, err := ParseArgs(appArgs, fullCom)

	if err != nil {
		return err
	}

	if tree.AutoHelp && !userCom.HideHelp {
		if userCom.hasFlag(autoHelpFlag.ShortName) || userCom.hasFlag(autoHelpFlag.LongName) || userCom.hasArg(autoHelpArg.Value) {
			if tree.ToHelpString == nil {
				fmt.Println(ToHelpString(fullCom))
			} else {
				fmt.Println(tree.ToHelpString(fullCom))
			}
			return nil
		}
	}
	if tree.Shared.PreAction != nil {
		tree.Shared.PreAction(userCom)
		if err != nil {
			return err
		}
	}
	if fullCom.Action != nil {
		fullCom.Action(userCom)
		if err != nil {
			return err
		}
	}
	if tree.Shared.PostAction != nil {
		tree.Shared.PostAction(userCom)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
* predicate refers to the second half of the command, the piece containing the flags,
* options, and arguments of the command string.
 */
func (tree CommandTree) FindCommand(appArgs []string) (fullCom Command, err error) {
	numArgs := len(appArgs)

	if numArgs == 0 {
		return fullCom, errors.New("cli: No arguments provided")
	}

	curCommand := &tree.Root
	curArg := appArgs[0]

	if curCommand.Name != curArg {
		errStr := fmt.Sprintf("cli: Command %s not found", curCommand.Name)
		return fullCom, errors.New(errStr)
	}

	predBegin := 0

	for i := 1; i < numArgs; i++ {
		predBegin++
		curArg = appArgs[i]
		argFound := false

		if strings.HasPrefix(curArg, "-") {
			break
		}

		subCount := len(curCommand.SubCommands)
		for j := 0; j < subCount; j++ {
			curSub := &curCommand.SubCommands[j]

			if curArg == curSub.Name {
				curCommand = curSub
				argFound = true
				break
			}
		}

		if argFound == false {
			break
		}
	}

	fullCom = *curCommand

	return
}

// Big ugly function that does the grunt work of the program. It could be split into functions, but as it is
// they would require a bunch or parameters some of them being pointers and would be just as ugly.
func ParseArgs(appArgs []string, c Command) (userCom Command, err error) {
	predicateStart := 0
	for i, arg := range appArgs {
		if arg == c.Name {
			userCom.Name = c.Name
			userCom.Usage = c.Usage
			userCom.Description = c.Description
			predicateStart = i + 1
			break
		}
	}
	if predicateStart == 0 {
		errStr := fmt.Sprintf("cli: Command %s not found in appArgs %s", c.Name, strings.Join(appArgs, " "))
		return userCom, errors.New(errStr)
	}
	predicate := appArgs[predicateStart:]
	predLen := len(predicate)
	for i := 0; i < predLen; i++ {
		argStr := predicate[i]

		if argStr == " " || argStr == "" {
			continue
		}

		// Must be a 1 piece Option. We split it into a 2 part form and process it as if the user
		// had typed --key value rather than --key=value
		if strings.Contains(argStr, "=") {
			keyValue := strings.Split(argStr, "=")
			predicate[i] = keyValue[1]

			//insert value at i
			predicate = append(predicate, "")
			copy(predicate[i+1:], predicate[i:])
			predicate[i] = keyValue[0]

			// must redo controlling variables after doind an insert on predicate
			argStr = predicate[i]
			predLen = len(predicate)
		}

		if strings.HasPrefix(argStr, "--") {
			i, err = parseLongForm(predicate, i, c, &userCom)

		} else if strings.HasPrefix(argStr, "-") {
			i, err = parseShortForm(predicate, i, c, &userCom)

			// Must be an Argument
		} else {
			a := Argument{Value: argStr}
			userCom.Args = append(userCom.Args, a)
		}

		if err != nil {
			return userCom, err
		}
	}
	return userCom, nil // nil error
}

func parseLongForm(predicate []string, pos int, c Command, userCom *Command) (newPos int, err error) {
	// strip off "--"
	argStr := predicate[pos][2:]
	predLen := len(predicate)

	if c.hasFlag(argStr) {
		if !userCom.hasFlag(argStr) {
			f := Flag{LongName: argStr}
			userCom.Flags = append(userCom.Flags, f)
		}
	} else if c.hasOption(argStr) {
		if !userCom.hasOption(argStr) {
			if pos+1 >= predLen {
				return // fix error later option specified with no value
			}
			o := Option{LongName: argStr, Value: predicate[pos+1]}
			userCom.Opts = append(userCom.Opts, o)
			//If the arg was an option we need to move the iterator an extra position
		}
		pos++
	} else {
		errStr := fmt.Sprintf("cli: Long form input --%s not found", argStr)
		return newPos, errors.New(errStr)
	}
	newPos = pos
	return newPos, nil
}

func parseShortForm(predicate []string, pos int, c Command, userCom *Command) (newPos int, err error) {
	// strip off "-"
	argStr := predicate[pos][1:]
	predLen := len(predicate)

	// For the case of multiple flags like "tar -xvf"
	if len(argStr) > 1 {
		for _, charUtf := range argStr {
			char := string(charUtf)
			if c.hasFlag(char) {
				if !userCom.hasFlag(char) {
					f := Flag{ShortName: char}
					userCom.Flags = append(userCom.Flags, f)
				}
			} else {
				errStr := fmt.Sprintf("cli: Short form input -%s is too Long", argStr)
				err = errors.New(errStr)
				return newPos, err
			}
		}
	} else {
		if c.hasFlag(argStr) {
			if !userCom.hasFlag(argStr) {
				f := Flag{ShortName: argStr}
				userCom.Flags = append(userCom.Flags, f)
			}
		} else if c.hasOption(argStr) {

			if !userCom.hasOption(argStr) {
				if pos+1 >= predLen {
					errStr := fmt.Sprintf("cli: No value provided for option -%s", argStr)
					err = errors.New(errStr)
					return newPos, err
				}
				o := Option{ShortName: argStr, Value: predicate[pos+1]}
				userCom.Opts = append(userCom.Opts, o)
				//If the arg was an option we need to move the iterator an extra position
			}
			pos++
		} else {
			errStr := fmt.Sprintf("cli: Short form input -%s not found", argStr)
			err = errors.New(errStr)
			return newPos, err
		}
	}
	newPos = pos
	return newPos, err
}

// For the case of mul
func PrintTree(c *Command) {
	fmt.Printf("%d: %s\n", 0, c.Name)
	printChildren(c, 0)
}

func printChildren(c *Command, level int) {
	level++
	subCount := len(c.SubCommands)
	for j := 0; j < subCount; j++ {
		curSub := &c.SubCommands[j] // pointer to a command

		for j := 0; j < level; j++ {
			fmt.Printf("  ")
		}
		fmt.Printf("%d: %s\n", level, curSub.Name)
		printChildren(curSub, level)
	}
}
