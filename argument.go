package cli

import "fmt"

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
