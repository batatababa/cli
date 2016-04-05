package cli

import "fmt"

var PreActionOccured bool = false
var ActionOccured bool = false
var PostActionOccured bool = false

var Quick *Command = &Command{
	Name:        "quick",
	Description: "the quick",
	Usage:       "use a quick?",
	Flags:       MultiFlag,
	Args:        MultiArg,
	Opts:        MultiOpt,
	HideHelp:    false,
	Action:      ActionPrintCommand,
	SubCommands: []Command{*QuickBrown, *QuickRed},
}
var QuickRed *Command = &Command{
	Name:        "red",
	Description: "the quick red",
	Usage:       "use a red?",
	Flags:       MultiFlag,
	Args:        MultiArg,
	Opts:        MultiOpt,
	HideHelp:    false,
	Action:      ActionPrintCommand,
}

var QuickBrown *Command = &Command{
	Name:        "brown",
	Description: "the quick brown",
	Usage:       "use a brown?",
	Flags:       MultiFlag,
	Args:        MultiArg,
	Opts:        MultiOpt,
	HideHelp:    false,
	Action:      ActionPrintCommand,
	SubCommands: []Command{*QuickBrownFox, *QuickBrownBear, *QuickBrownCow},
}

var QuickBrownFox *Command = &Command{
	Name:        "fox",
	Description: "the quick brown fox",
	Usage:       "use a fox?",
	Flags:       MultiFlag,
	Args:        MultiArg,
	Opts:        MultiOpt,
	HideHelp:    false,
	Action:      ActionTester,
}
var QuickBrownBear *Command = &Command{
	Name:        "bear",
	Description: "the quick brown bear",
	Usage:       "use a bear?",
	Flags:       MultiFlag,
	Args:        MultiArg,
	Opts:        MultiOpt,
	HideHelp:    true,
	Action:      ActionPrintCommand,
}

var QuickBrownCow *Command = &Command{
	Name:        "cow",
	Description: "the quick brown cow",
	Usage:       "use a cow?",
	Flags:       MultiFlag,
	Args:        MultiArg,
	Opts:        MultiOpt,
	HideHelp:    true,
	Action:      ActionPrintCommand,
}

var comTree CommandTree = CommandTree{
	Author:       "Jeff Williams",
	Version:      "0.1",
	AutoHelp:     true,
	ToHelpString: ToHelpString,
	Root: Command{
		Name:        "the",
		Description: "The What?",
		Usage:       "How should I be doing this?",
		Args:        MultiArg,
		Flags:       MultiFlag,
		Opts:        MultiOpt,
		SubCommands: []Command{*Quick},
	},
	Shared: SharedParameters{
		PreAction:  PreActionTester,
		PostAction: PostActionTester,
	},
}

var MultiArg []Argument = []Argument{
	{
		Name:        "Arg1",
		Description: "Description for Arg1",
	},
	{
		Name:        "Arg2",
		Description: "Description for Arg2",
	},

	{
		Name:        "Arg2",
		Description: "Description for Arg2",
	},
}

var SingleArg []Argument = []Argument{{
	Name:        "Arg",
	Description: "Description for Arg",
}}

var MultiFlag []Flag = []Flag{
	{
		LongName:    "LongA",
		Description: "Description for Flag1",
	},
	{
		ShortName:   "b",
		Description: "Description for Flag2",
	},

	{
		LongName:  "LongC",
		ShortName: "c",
	},
	{
		LongName:    "LongD",
		ShortName:   "d",
		Description: "Description for Flag",
	},
}

var SingleFlag []Flag = []Flag{{
	LongName:    "LongE",
	ShortName:   "e",
	Description: "Description for Flag LongB",
}}

var MultiOpt []Option = []Option{
	{
		LongName:    "LongF",
		ShortName:   "f",
		Description: "Description for Opt LongC",
	},
	{
		LongName:    "LongG",
		ShortName:   "g",
		Description: "Description for Opt LongD",
	},

	{
		LongName:    "LongI",
		ShortName:   "i",
		Description: "Description for Opt LongE",
	},
}

var SingleOpt []Option = []Option{{
	LongName:    "LongJ",
	ShortName:   "j",
	Description: "Description for Option LongF",
}}

func ActionPrintCommand(c Command) (err error) {
	fmt.Println(c.ToString())
	return
}

func PreActionTester(c Command) (err error) {
	PreActionOccured = true
	return
}

func ActionTester(c Command) (err error) {
	ActionOccured = true
	return
}

func PostActionTester(c Command) (err error) {
	PostActionOccured = true
	return
}

func ResetActionTesters() {
	PostActionOccured = false
	ActionOccured = false
	PreActionOccured = false
}

// var commandTree CommandTree = CommandTree{
// 	Author:  "Jeff Williams",
// 	Version: "0.1",
// 	Root: Command{
// 		Description: "description for a test app",
// 		Name:        "the",
// 		Usage:       "how to use the program",
// 		Args:        MultiArg,
// 		Flags:       MultiFlag,
// 		Opts:        MultiOpt,
// 		SubCommands: []Command{
// 			{
// 				Name:        "quick",
// 				Description: "the quick",
// 				Flags:       MultiFlag,
// 				Args:        MultiArg,
// 				Opts:        MultiOpt,
// 				Action:      ActionPrintCommand,
// 				SubCommands: []Command{
// 					{
// 						Name:        "brown",
// 						Description: "the quick brown",
// 						Flags:       SingleFlag,
// 						Args:        SingleArg,
// 						Opts:        SingleOpt,
// 						Action:      ActionPrintCommand,
// 					},
// 					{
// 						Name:        "red",
// 						Description: "the quick red",
// 					},
// 				},
// 			},
// 			{
// 				Name:        "slow",
// 				Description: "the slow",
// 				SubCommands: []Command{
// 					{
// 						Name:        "brown",
// 						Description: "the slow brown",
// 					},
// 					{
// 						Name:        "red",
// 						Description: "the slow red",
// 					},
// 				},
// 			},
// 		},
// 	},
// }
