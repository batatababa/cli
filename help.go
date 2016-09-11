package cli

import (
	"bytes"
	"fmt"

	"github.com/ryanuber/columnize"
)

func ToHelpString(c Command) (help string) {
	var buf bytes.Buffer
	config := columnize.DefaultConfig()
	config.Prefix = "  "

	buf.WriteString(fmt.Sprintf("%s\n", c.Description))
	buf.WriteString("Usage: \n")

	uniqueArgs := make(map[*Argument]bool)
	if c.Args != nil {
		var args []string
		for _, a := range c.Args {
			args = append(args, fmt.Sprintf("<%s>", a.Name))
			uniqueArgs[&a] = true
		}
		buf.WriteString(columnize.Format(args, config))
	}

	if c.ArgSets != nil {
		var args []string
		for _, argSet := range c.ArgSets {
			var setString string
			for _, a := range argSet.Set {
				setString += fmt.Sprintf("<%s> ", a.Name)
			}
			args = append(args, setString)
			buf.WriteString(columnize.Format(args, config))
		}
	}

	if c.SubCommands != nil {
		buf.WriteString(" SubCommands:\n")
		var subs []string
		for _, s := range c.SubCommands {
			subs = append(subs, fmt.Sprintf("%s:|%s", s.Name, s.Description))
		}
		buf.WriteString(columnize.Format(subs, config))
		buf.WriteString("\n\n")
	}

	if c.Flags != nil {
		buf.WriteString(" Flags:\n")
		var flags []string
		for _, f := range c.Flags {
			flags = append(flags, toShortLongDescString(f.ShortName, f.LongName, f.Description))
		}
		buf.WriteString(columnize.Format(flags, config))
		buf.WriteString("\n\n")
	}

	if c.Opts != nil {
		buf.WriteString(" Options:\n")
		var opts []string
		for _, o := range c.Opts {
			opts = append(opts, toShortLongDescString(o.ShortName, o.LongName, o.Description))
		}
		buf.WriteString(columnize.Format(opts, config))
		buf.WriteString("\n\n")
	}
	help = buf.String()
	return
}

func toShortLongDescString(short string, long string, description string) (str string) {
	var buf bytes.Buffer
	if short != "" {
		buf.WriteString(fmt.Sprintf("-%s", short))
	}
	buf.WriteString("|")
	if long != "" {
		buf.WriteString(fmt.Sprintf("--%s", long))
	}
	buf.WriteString(fmt.Sprintf(",|%s", description))
	str = buf.String()
	return
}
