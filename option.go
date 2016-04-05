package cli

type Option struct {
	ShortName   string
	LongName    string
	Value       string
	Description string
}

type OptionSlice []Option

func (slice OptionSlice) toColumnArray() (colArr [][]string) {
	if slice != nil {
		for _, opt := range slice {
			colArr = append(colArr, []string{opt.ShortName, opt.LongName, opt.Description, opt.Value})
		}
	}

	return
}
