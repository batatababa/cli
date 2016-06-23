package cli

import "fmt"

type Option struct {
	ShortName   string
	LongName    string
	Value       string
	Description string
}

func (opt *Option) String() string {
	return fmt.Sprintf("Opt: -%-2s --%-10s, %-s %-s", opt.ShortName, opt.LongName, opt.Description, opt.Value)
}

func OptArrayToStringArray(optArr []Option) (strArr []string) {
	for _, opt := range optArr {
		strArr = append(strArr, opt.String())
	}

	return strArr
}
