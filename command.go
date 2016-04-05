package cli

import (
	"bytes"
	"fmt"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Flags       FlagSlice
	Args        ArgumentSlice
	Opts        OptionSlice
	SubCommands CommandSlice
	HideHelp    bool
	Action      func(com Command) error
}

type CommandSlice []Command

func (slice CommandSlice) toColumnArray() (colArr [][]string) {
	if slice != nil {
		for _, sub := range slice {
			colArr = append(colArr, []string{sub.Name, sub.Description})
		}
	}
	return
}

func colArrayToString(colArray [][]string, colNums []int) (str string) {
	var buffer bytes.Buffer
	colCount := len(colNums)
	for _, strArr := range colArray {
		buffer.WriteString("{")
		for i, col := range colNums {
			str := strArr[col]
			if str != "" {
				buffer.WriteString(fmt.Sprintf("%s", str))
				if i+1 < colCount {
					buffer.WriteString(",")
				}
			}
		}
		buffer.WriteString("} ")
	}

	str = buffer.String()
	return
}

func (c Command) Print() {
	fmt.Println(c.ToString())
}

func (c Command) ToString() string {

	sa := []string{
		fmt.Sprintf("Name: %s", c.Name),
		fmt.Sprintf("Description: %s", c.Description),
		fmt.Sprintf("Usage: %s", c.Usage),
		fmt.Sprintf("SubCommands: %s", colArrayToString(c.SubCommands.toColumnArray(), []int{0, 1})),
		fmt.Sprintf("Arguments: %s", colArrayToString(c.Args.toColumnArray(), []int{0, 1, 2})),
		fmt.Sprintf("Flags: %s", colArrayToString(c.Flags.toColumnArray(), []int{0, 1})),
		fmt.Sprintf("Options: %s", colArrayToString(c.Opts.toColumnArray(), []int{0, 1, 3})),
	}

	return strings.Join(sa, "\n    ")
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
