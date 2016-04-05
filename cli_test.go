package cli

import (
	"strings"
	"testing"
)

func TestPrintCommand(t *testing.T) {
	// fmt.Println(comTree.Root.ToString())
}

func TestPrintTree(t *testing.T) {
	// PrintTree(&comTree.Root)
}

func TestDefaultHelpTemplate(t *testing.T) {
	ToHelpString(comTree.Root)
}

func TestFind(t *testing.T) {
	findHelper(t, "the", &comTree.Root)
	findHelper(t, "the quick", Quick)
	findHelper(t, "the quick brown", QuickBrown)
	findHelper(t, "the quick brown fox", QuickBrownFox)
	findHelper(t, "the quick brown bear", QuickBrownBear)
	findHelper(t, "the quick brown cow", QuickBrownCow)
	findHelper(t, "the quick red", QuickRed)
}

func findHelper(t *testing.T, appArgs string, targetCom *Command) {
	argArray := strings.Split(appArgs, " ")
	foundCom, _ := comTree.FindCommand(argArray)

	if foundCom.ToString() != targetCom.ToString() {
		t.Errorf("FindCommand:Looking for \"%s\" but found \"%s\"", appArgs, foundCom.Name)
	}
}

func TestParsingArguments(t *testing.T) {
	var parsed1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Args: []Argument{{
			Value: "fox",
		}},
	}
	parseHelper(t, "the quick brown fox", *QuickBrown, parsed1)

	var parsed2 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Args: []Argument{
			{
				Value: "fox",
			},
			{
				Value: "dog",
			},
		},
	}
	parseHelper(t, "the quick brown fox dog", *QuickBrown, parsed2)
}

func TestParsingFlags(t *testing.T) {
	var parsed1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Flags: []Flag{{
			ShortName: "b",
		}},
	}
	parseHelper(t, "the quick brown -b", *QuickBrown, parsed1)

	var parsed2 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Flags: []Flag{
			{
				LongName: "LongD",
			},
			{
				ShortName: "c",
			},
		},
	}
	parseHelper(t, "the quick brown --LongD -c", *QuickBrown, parsed2)
}

func TestParsingOptions(t *testing.T) {
	var parsed1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			LongName: "LongF",
			Value:    "val",
		}},
	}
	parseHelper(t, "the quick brown --LongF val", *QuickBrown, parsed1)

	var parsed2 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "LongF",
			Value:     "val",
		}},
	}
	parseHelper(t, "the quick brown --LongF val", *QuickBrown, parsed2)

	var parsed3 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{
			{
				ShortName: "g",
				Value:     "val",
			},
			{
				LongName: "LongF",
				Value:    "val2",
			},
		},
	}
	parseHelper(t, "the quick brown -g val --LongF val2", *QuickBrown, parsed3)
}

func TestParsingAll(t *testing.T) {
	var parsed1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Args: []Argument{{
			Value: "fox",
		}},
		Opts: []Option{
			{
				ShortName: "g",
				Value:     "val",
			},
			{
				LongName: "LongF",
				Value:    "val2",
			},
		},
		Flags: []Flag{
			{
				LongName: "LongD",
			},
			{
				ShortName: "c",
			},
		},
	}
	parseHelper(t, "the quick brown fox -g val --LongF val2 --LongD -c", *QuickBrown, parsed1)
}

func TestRunActions(t *testing.T) {
	ResetActionTesters()
	assertActions(t, false)
	Run(strings.Split("the quick brown fox -i optVal --LongC arg1 arg2", " "), &comTree)
	assertActions(t, true)
}

func TestSharedParameters(t *testing.T) {

	shared := SharedParameters{
		Opts: []Option{
			{
				ShortName: "x",
			},
		},
		Flags: []Flag{
			{
				ShortName: "v",
			},
		},
	}

	tempTree := NewCommandTree()
	tempTree.Root = *Quick
	tempTree.Shared = shared

	Run(strings.Split("quick brown fox -x val -k", " "), &tempTree)
}

func assertActions(t *testing.T, assertion bool) {
	if !(ActionOccured == PreActionOccured == PostActionOccured == assertion) {
		t.Errorf("Actions were not run correctly")
	}
}

func TestAutoHelp(t *testing.T) {
	comTree.AutoHelp = true
	ResetActionTesters()
	assertActions(t, false)
	// This does not fully test auto help, but if auto help is working the actions will
	// never be run
	Run(strings.Split("the quick brown fox -h", " "), &comTree)
	Run(strings.Split("the quick brown fox --help", " "), &comTree)
	Run(strings.Split("the quick brown fox ?", " "), &comTree)
	assertActions(t, false)
}

func TestFlagDuplication(t *testing.T) {

	var parsed1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Flags: []Flag{{
			ShortName: "b",
		}},
	}
	parseHelper(t, "the quick brown -bb", *QuickBrown, parsed1)
	parseHelper(t, "the quick brown -bbb -bbb", *QuickBrown, parsed1)
}

func TestOptionDuplication(t *testing.T) {
	var parsed1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "f",
			Value:     "val",
		}},
	}
	parseHelper(t, "the quick brown -f val -f val", *QuickBrown, parsed1)
}

func TestOptionShortForm(t *testing.T) {
	var parsed1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "f",
			Value:     "val",
		}},
	}
	parseHelper(t, "the quick brown -f=val", *QuickBrown, parsed1)

	var parsed2 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "LongF",
			Value:     "val",
		}},
	}
	parseHelper(t, "the quick brown --LongF=val  ", *QuickBrown, parsed2)

}

func parseHelper(t *testing.T, appArgs string, fullCom Command, parsedCom Command) {
	argArray := strings.Split(appArgs, " ")
	userCom, _ := ParseArgs(argArray, fullCom)

	if userCom.ToString() != parsedCom.ToString() {
		t.Errorf("Bad Parsing \n %s  \nNOT EQUAL TO \n %s", userCom.ToString(), parsedCom.ToString())
	}
}

// func TestReadMeCodeExamples(t *testing.T) {
// 	var BrownCow *Command = &Command{
// 		Name:     "brown",
// 		HideHelp: true,
// 	}
// 	var BrownFox *Command = &Command{
// 		Name: "fox",
// 	}
// 	var Brown *Command = &Command{
// 		Name:        "brown",
// 		Description: "silly example command",
// 		Usage:       "use a brown what?",
// 		SubCommands: []Command{*BrownCow, *BrownFox},
// 		Flags: []Flag{
// 			{
// 				ShortName:   "a",
// 				LongName:    "LongA",
// 				Description: "Description for Flag1",
// 			},
// 		},
// 		Opts: []Option{
// 			{
// 				ShortName:   "a",
// 				LongName:    "LongA",
// 				Description: "Description for Opt A",
// 			},
// 		},
// 		Args: []Argument{
// 			{
// 				Name:        "Arg1",
// 				Description: "Description for Arg1",
// 			},
// 		},
// 	}

// 	tree := NewCommandTree()
// 	tree.Root = *Brown
// 	tree.Version = "1.0.0"

// 	tree.Shared = SharedParameters{
// 		PreAction: func(c Command) (err error) {
// 			fmt.Println(c.Name)
// 			return
// 		},
// 		Flags: []Flag{
// 			{
// 				ShortName:   "v",
// 				LongName:    "verbose",
// 				Description: " Maybe all commands need a verbose flag.te",
// 			},
// 		},
// 	}
// 	Run([]string{"brown", "fox"}, &tree)
// }
