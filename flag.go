package cli

import "fmt"

type Flag struct {
	ShortName   string
	LongName    string
	Description string
}

func (flag *Flag) String() string {
	return fmt.Sprintf("Flg: -%-2s --%-10s, %-s", flag.ShortName, flag.LongName, flag.Description)
}

func FlagArrayToStringArray(flagArr []Flag) (strArr []string) {
	for _, flag := range flagArr {
		strArr = append(strArr, flag.String())
	}

	return strArr
}
