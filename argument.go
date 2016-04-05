package cli

type Argument struct {
	Name        string
	Value       string
	Description string
}

type ArgumentSlice []Argument

func (slice ArgumentSlice) toColumnArray() (colArr [][]string) {
	if slice != nil {
		for _, arg := range slice {
			colArr = append(colArr, []string{arg.Name, arg.Description, arg.Value})
		}
	}

	return
}
