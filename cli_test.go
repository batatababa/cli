package cli

import (
	"strings"
	"testing"
)

func TestPrintCommand(t *testing.T) {
	// comTree.Root.Print()
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

	if foundCom.String() != targetCom.String() {
		t.Errorf("FindCommand:Looking for \"%s\" but found \"%s\"", appArgs, foundCom.Name)
	}
}

func TestParsingArguments(t *testing.T) {
	var expected1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Args: []Argument{{
			Value: "fox",
		}},
	}
	parseHelper(t, "the quick brown fox", *QuickBrown, expected1)

	var expected2 Command = Command{
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
	parseHelper(t, "the quick brown fox dog", *QuickBrown, expected2)
}

func TestParsingFlags(t *testing.T) {
	var expected1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Flags: []Flag{{
			ShortName: "b",
		}},
	}
	parseHelper(t, "the quick brown -b", *QuickBrown, expected1)

	var expected2 Command = Command{
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
	parseHelper(t, "the quick brown --LongD -c", *QuickBrown, expected2)
}

func TestParsingOptions(t *testing.T) {
	var expected1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			LongName: "LongF",
			Value:    "val",
		}},
	}
	parseHelper(t, "the quick brown --LongF val", *QuickBrown, expected1)

	var expected2 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "f",
			Value:     "val",
		}},
	}
	parseHelper(t, "the quick brown -f val", *QuickBrown, expected2)

	var expected3 Command = Command{
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
	parseHelper(t, "the quick brown -g val --LongF val2", *QuickBrown, expected3)
}

func TestParsingAll(t *testing.T) {
	var expected1 Command = Command{
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
	parseHelper(t, "the quick brown fox -g val --LongF val2 --LongD -c", *QuickBrown, expected1)
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

	var expected1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Flags: []Flag{{
			ShortName: "b",
		}},
	}
	parseHelper(t, "the quick brown -bb", *QuickBrown, expected1)
	parseHelper(t, "the quick brown -bbb -bbb", *QuickBrown, expected1)
}

func TestOptionDuplication(t *testing.T) {
	var expected1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "f",
			Value:     "val",
		}},
	}
	parseHelper(t, "the quick brown -f val -f val", *QuickBrown, expected1)
}

func TestOptionShortForm(t *testing.T) {

	var expected1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "f",
			Value:     "val",
		}},
	}
	parseHelper(t, "the quick brown -f=val", *QuickBrown, expected1)

	var expected2 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Opts: []Option{{
			ShortName: "LongF",
			Value:     "val",
		}},
	}
	parseHelperNeg(t, "the quick brown -TooLong=val  ", *QuickBrown, expected2)

}

func TestArgSets(t *testing.T) {
	var testCom1 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Args:        ThreeArg,
		ArgSets:     []ArgumentSet{TwoArgSet},
	}
	assertParsePasses(t, "the quick brown one two", testCom1)
	assertParsePasses(t, "the quick brown one two three", testCom1)

	var testCom2 Command = Command{
		Name:        "brown",
		Description: "the quick brown",
		Usage:       "use a brown?",
		Args:        TwoArg,
		ArgSets:     []ArgumentSet{ThreeArgSet},
	}
	assertParsePasses(t, "the quick brown one two", testCom2)
	assertParsePasses(t, "the quick brown one two three", testCom2)
}

func assertParsePasses(t *testing.T, appArgs string, fullCom Command) {
	argArray := strings.Split(appArgs, " ")
	_, err := ParseArgs(argArray, fullCom)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func assertParseFails(t *testing.T, appArgs string, fullCom Command) {
	argArray := strings.Split(appArgs, " ")
	_, err := ParseArgs(argArray, fullCom)

	if err == nil {
		t.Errorf("ParseArgs Passed when it should have failed")
	}
}

func parseHelper(t *testing.T, appArgs string, fullCom Command, expectedCom Command) {
	argArray := strings.Split(appArgs, " ")
	userCom, err := ParseArgs(argArray, fullCom)

	if err != nil {
		t.Errorf(err.Error())
	}

	if userCom.String() != expectedCom.String() {
		t.Errorf("\n  Parsed Com For: \"%s\" \n%s  \n  NOT EQUAL Expected Com: \n%s", appArgs, userCom.String(), expectedCom.String())
	}
}

func parseHelperNeg(t *testing.T, appArgs string, fullCom Command, expectedCom Command) {
	argArray := strings.Split(appArgs, " ")
	_, err := ParseArgs(argArray, fullCom)

	if err == nil {
		t.Errorf("")
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
