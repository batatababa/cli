package cli

import (
	"fmt"
	"strings"
)

type Argument struct {
	Name        string
	Value       string
	Description string
}

func (arg *Argument) String() string {
	name := fmt.Sprintf("<%s>", arg.Name)
	return fmt.Sprintf("Arg: %s, %-s %-s", name, arg.Description, arg.Value)
}

func ArgArrayToStringArray(argArr []Argument) (strArr []string) {
	for _, arg := range argArr {
		strArr = append(strArr, arg.String())
	}

	return strArr
}

type ArgumentSet struct {
	Set []Argument
}

func (argSet *ArgumentSet) String() string {
	return fmt.Sprintf("ArgSet: {%s}", strings.Join(ArgArrayToStringArray(argSet.Set), ",\n           "))
}

func ArgSetArrayToStringArray(argSetArr []ArgumentSet) (strArr []string) {
	for _, argSet := range argSetArr {
		strArr = append(strArr, argSet.String())
	}

	return strArr
}
