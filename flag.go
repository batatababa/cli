package cli

type Flag struct {
	ShortName   string
	LongName    string
	Description string
}

type FlagSlice []Flag

func (slice FlagSlice) toColumnArray() (colArr [][]string) {
	if slice != nil {
		for _, flag := range slice {
			colArr = append(colArr, []string{flag.ShortName, flag.LongName, flag.Description})
		}
	}

	return
}
