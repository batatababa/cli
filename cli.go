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
	ToHelpString func(c Command, pathToCom []string) string
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

	fullCom, pathToCom, err := tree.FindCommand(appArgs)

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
			helpStr := ""
			if tree.ToHelpString == nil {
				helpStr = ToHelpString(fullCom, pathToCom)
			} else {
				helpStr = tree.ToHelpString(fullCom, pathToCom)
			}

			if helpStr != "" {
				fmt.Println(helpStr)
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
func (tree CommandTree) FindCommand(appArgs []string) (fullCom Command, pathToCom []string, err error) {
	numArgs := len(appArgs)

	if numArgs == 0 {
		err = errors.New("cli: No arguments provided")
		return fullCom, pathToCom, err
	}

	curCommand := &tree.Root
	curArg := appArgs[0]

	if curCommand.Name != curArg {
		errStr := fmt.Sprintf("cli: Command %s not found", curCommand.Name)
		err = errors.New(errStr)
		return fullCom, pathToCom, err
	}

	pathToCom = append(pathToCom, curArg)

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
				pathToCom = append(pathToCom, curArg)
				break
			}
		}

		if argFound == false {
			break
		}
	}
	if pathToCom != nil {
		pathToCom = pathToCom[:len(pathToCom)-1]
	}

	fullCom = *curCommand

	return fullCom, pathToCom, err
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

type Node struct {
	Command
	level int
}

// For the case of mul
func PrintTree(c *Command) {
	slice := CommandToNodeSlice(c)

	for _, node := range slice {
		for j := 0; j < node.level; j++ {
			fmt.Printf("  ")
		}
		fmt.Printf("%d: %s\n", node.level, node.Name)
	}
}

// For the case of mul
func PrintTreeHelp(c *Command) {
	slice := CommandToNodeSlice(c)

	for _, node := range slice {
		fmt.Printf("--------------------------------------------\n")
		fmt.Printf("\"%s\"\n", ToHelpString(node.Command, nil))
	}
}

func addChildrenToSlice(n *Node, slice *[]Node) {
	subCount := len(n.SubCommands)
	for i := 0; i < subCount; i++ {
		child := Node{n.SubCommands[i], n.level + 1} // pointer to a command

		*slice = append(*slice, child)
		addChildrenToSlice(&child, slice)
	}
}

func CommandToNodeSlice(c *Command) (slice []Node) {
	slice = append(slice, Node{*c, 0})
	addChildrenToSlice(&slice[0], &slice)

	return slice
}

// type Iterable interface {
// 	Next() (int, interface{})
// 	Prev() (int, interface{})
// 	First() (int, interface{})
// }

// func GetIterator(iable Iterable) Iterable {
// 	slice := s.ToSlice()
// 	index := -1
// 	size := len(slice)

// 	return Iterable{
// 		Next: func() (int, interface{}) {
// 			index++
// 			if index == size {
// 				return 0, nil
// 			} else {
// 				return index, slice[index]
// 			}
// 		},
// 		Prev: func() (int, interface{}) {
// 			if index <= 0 {
// 				return 0, nil
// 			} else {
// 				index--
// 				return index, slice[index]
// 			}
// 		},
// 		First: func() (int, interface{}) {
// 			index = 0
// 			return index, slice[index]
// 		},
// 	}
// }
