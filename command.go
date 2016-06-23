package cli

import (
	"fmt"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Flags       []Flag
	Args        []Argument
	ArgSets     []ArgumentSet
	Opts        []Option
	SubCommands []Command
	HideHelp    bool
	Action      func(com Command) error
}

func SubCommandToString(sub *Command) string {
	return fmt.Sprintf("Sub: %-10s, %-s", sub.Name, sub.Description)
}

func SubComArrayToStringArray(subArr []Command) (strArr []string) {
	for _, sub := range subArr {
		strArr = append(strArr, SubCommandToString(&sub))
	}

	return strArr
}

func (c Command) Print() {
	fmt.Println(c.String())
}

func (c Command) String() string {
	sep := "\n.."
	sa := []string{
		fmt.Sprintf("..Name: %s", c.Name),
		fmt.Sprintf("Description: %s", c.Description),
		fmt.Sprintf("Usage: %s", c.Usage),
		strings.Join(SubComArrayToStringArray(c.SubCommands), sep),
		strings.Join(ArgArrayToStringArray(c.Args), sep),
		strings.Join(ArgSetArrayToStringArray(c.ArgSets), sep),
		strings.Join(FlagArrayToStringArray(c.Flags), sep),
		strings.Join(OptArrayToStringArray(c.Opts), sep),
	}

	return strings.Join(sa, sep)
}

func (c Command) hasFlag(flagStr string) (found bool) {
	for _, flag := range c.Flags {
		if flagStr == flag.ShortName || flagStr == flag.LongName {
			return true
		}
	}
	return false
}

func (c Command) hasOption(optStr string) (found bool) {
	for _, opt := range c.Opts {
		if optStr == opt.ShortName || optStr == opt.LongName {
			return true
		}
	}
	return false
}

func (c Command) hasArg(argStr string) (found bool) {
	for _, arg := range c.Args {
		if argStr == arg.Value {
			return true
		}
	}
	return false
}

func (c Command) hasSubCommand(comStr string) (found bool) {
	for _, sub := range c.SubCommands {
		if comStr == sub.Name {
			return true
		}
	}
	return false
}
