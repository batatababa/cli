package cli

import (
	"bytes"
	"fmt"

	"github.com/ryanuber/columnize"
)

func ToHelpString(c Command, pathToCom []string) (help string) {
	var helpBuf bytes.Buffer
	config := columnize.DefaultConfig()
	config.Prefix = "  "

	helpBuf.WriteString(fmt.Sprintf("%s: %s\n", c.Name, c.Description))
	helpBuf.WriteString(fmt.Sprintf("Usage: %s\n\n", c.Usage))

	if c.SubCommands != nil {
		helpBuf.WriteString(" SubCommands:\n")
		var subs []string
		for _, s := range c.SubCommands {
			subs = append(subs, fmt.Sprintf("%s:|%s", s.Name, s.Description))
		}
		helpBuf.WriteString(columnize.Format(subs, config))
		helpBuf.WriteString("\n\n")
	}

	if c.Args != nil {
		helpBuf.WriteString(" Arguments:\n")
		var args []string
		for _, a := range c.Args {
			args = append(args, fmt.Sprintf("<%s>|%s", a.Name, a.Description))
		}
		helpBuf.WriteString(columnize.Format(args, config))
		helpBuf.WriteString("\n\n")
	}

	if c.Flags != nil {
		helpBuf.WriteString(" Flags:\n")
		var flags []string
		for _, f := range c.Flags {
			flags = append(flags, toShortLongDescString(f.ShortName, f.LongName, f.Description))
		}
		helpBuf.WriteString(columnize.Format(flags, config))
		helpBuf.WriteString("\n\n")
	}

	if c.Opts != nil {
		helpBuf.WriteString(" Options:\n")
		var opts []string
		for _, o := range c.Opts {
			opts = append(opts, toShortLongDescString(o.ShortName, o.LongName, o.Description))
		}
		helpBuf.WriteString(columnize.Format(opts, config))
		helpBuf.WriteString("\n\n")
	}
	help = helpBuf.String()
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
