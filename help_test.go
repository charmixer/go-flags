package flags

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

type helpOptions struct {
	Verbose          []bool       `short:"v" long:"verbose" description:"Show verbose debug information" ini-name:"verbose"`
	Call             func(string) `short:"c" description:"Call phone number" ini-name:"call"`
	PtrSlice         []*string    `long:"ptrslice" description:"A slice of pointers to string"`
	EmptyDescription bool         `long:"empty-description"`

	Default           string            `long:"default" default:"Some\nvalue" description:"Test default value"`
	DefaultArray      []string          `long:"default-array" default:"Some value" default:"Other\tvalue" description:"Test default array value"`
	DefaultMap        map[string]string `long:"default-map" default:"some:value" default:"another:value" description:"Testdefault map value"`
	OptionWithArgName string            `long:"opt-with-arg-name" value-name:"something" description:"Option with named argument"`
	OptionWithChoices string            `long:"opt-with-choices" value-name:"choice" choice:"dog" choice:"cat" description:"Option with choices"`
	Hidden            string            `long:"hidden" description:"Hidden option" hidden:"yes"`

	HiddenOptionWithVeryLongName bool `long:"this-hidden-option-has-a-ridiculously-long-name" hidden:"yes"`

	OnlyIni string `ini-name:"only-ini" description:"Option only available in ini"`

	Other struct {
		StringSlice []string       `short:"s" default:"some" default:"value" description:"A slice of strings"`
		IntMap      map[string]int `long:"intmap" default:"a:1" description:"A map from string to int" ini-name:"int-map"`
	} `group:"Other Options"`

	HiddenGroup struct {
		InsideHiddenGroup string `long:"inside-hidden-group" description:"Inside hidden group"`
		Padder            bool   `long:"this-option-in-a-hidden-group-has-a-ridiculously-long-name"`
	} `group:"Hidden group" hidden:"yes"`

	GroupWithOnlyHiddenOptions struct {
		SecretFlag bool `long:"secret" description:"Hidden flag in a non-hidden group" hidden:"yes"`
	} `group:"Non-hidden group with only hidden options"`

	Group struct {
		Opt                  string `long:"opt" description:"This is a subgroup option"`
		HiddenInsideGroup    string `long:"hidden-inside-group" description:"Hidden inside group" hidden:"yes"`
		NotHiddenInsideGroup string `long:"not-hidden-inside-group" description:"Not hidden inside group" hidden:"false"`

		Group struct {
			Opt string `long:"opt" description:"This is a subsubgroup option"`
		} `group:"Subsubgroup" namespace:"sap"`
	} `group:"Subgroup" namespace:"sip"`

	Bommand struct {
		Hidden bool `long:"hidden" description:"A hidden option" hidden:"yes"`
	} `command:"bommand" description:"A command with only hidden options"`

	Command struct {
		ExtraVerbose []bool `long:"extra-verbose" description:"Use for extra verbosity"`
	} `command:"command" alias:"cm" alias:"cmd" description:"A command"`

	HiddenCommand struct {
		ExtraVerbose []bool `long:"extra-verbose" description:"Use for extra verbosity"`
	} `command:"hidden-command" description:"A hidden command" hidden:"yes"`

	ParentCommand struct {
		Opt        string `long:"opt" description:"This is a parent command option"`
		SubCommand struct {
			Opt string `long:"opt" description:"This is a sub command option"`
		} `command:"sub" description:"A sub command"`
	} `command:"parent" description:"A parent command"`

	Args struct {
		Filename     string  `positional-arg-name:"filename" description:"A filename with a long description to trigger line wrapping"`
		Number       int     `positional-arg-name:"num" description:"A number"`
		HiddenInHelp float32 `positional-arg-name:"hidden-in-help" required:"yes"`
	} `positional-args:"yes"`
}

func TestHelp(t *testing.T) {
	var opts helpOptions
	p := NewNamedParser("TestHelp", HelpFlag)
	p.AddGroup("Application Options", "The application options", &opts)

	_, err := p.ParseArgs([]string{"--help"})

	if err == nil {
		t.Fatalf("Expected help error")
	}

	if e, ok := err.(*Error); !ok {
		t.Fatalf("Expected flags.Error, but got %T", err)
	} else {
		if e.Type != ErrHelp {
			t.Errorf("Expected flags.ErrHelp type, but got %s", e.Type)
		}

		var expected string

		if runtime.GOOS == "windows" {
			expected = `Usage:
  TestHelp [OPTIONS] [filename] [num] hidden-in-help <bommand | command | parent>

Application Options:
  /v, /verbose                              Show verbose debug information
  /c:                                       Call phone number
      /ptrslice:                            A slice of pointers to string
      /empty-description
      /default:                             Test default value (default:
                                            "Some\nvalue")
      /default-array:                       Test default array value (default:
                                            Some value, "Other\tvalue")
      /default-map:                         Testdefault map value (default:
                                            some:value, another:value)
      /opt-with-arg-name:something          Option with named argument
      /opt-with-choices:choice[dog|cat]     Option with choices

Other Options:
  /s:                                       A slice of strings (default: some,
                                            value)
      /intmap:                              A map from string to int (default:
                                            a:1)

Subgroup:
      /sip.opt:                             This is a subgroup option
      /sip.not-hidden-inside-group:         Not hidden inside group

Subsubgroup:
      /sip.sap.opt:                         This is a subsubgroup option

Help Options:
  /?                                        Show this help message
  /h, /help                                 Show this help message

Arguments:
  filename:                                 A filename with a long description
                                            to trigger line wrapping
  num:                                      A number

Available commands:
  bommand  A command with only hidden options
  command  A command (aliases: cm, cmd)
  parent   A command with a sub command
`
		} else {
			expected = `Usage:
  TestHelp [OPTIONS] [filename] [num] hidden-in-help <bommand | command | parent>

Application Options:
  -v, --verbose                             Show verbose debug information
  -c=                                       Call phone number
      --ptrslice=                           A slice of pointers to string
      --empty-description
      --default=                            Test default value (default:
                                            "Some\nvalue")
      --default-array=                      Test default array value (default:
                                            Some value, "Other\tvalue")
      --default-map=                        Testdefault map value (default:
                                            some:value, another:value)
      --opt-with-arg-name=something         Option with named argument
      --opt-with-choices=choice[dog|cat]    Option with choices

Other Options:
  -s=                                       A slice of strings (default: some,
                                            value)
      --intmap=                             A map from string to int (default:
                                            a:1)

Subgroup:
      --sip.opt=                            This is a subgroup option
      --sip.not-hidden-inside-group=        Not hidden inside group

Subsubgroup:
      --sip.sap.opt=                        This is a subsubgroup option

Help Options:
  -h, --help                                Show this help message

Arguments:
  filename:                                 A filename with a long description
                                            to trigger line wrapping
  num:                                      A number

Available commands:
  bommand  A command with only hidden options
  command  A command (aliases: cm, cmd)
  parent   A parent command
`
		}

		assertDiff(t, e.Message, expected, "help message")
	}
}

func TestMan(t *testing.T) {
	var opts helpOptions
	p := NewNamedParser("TestMan", HelpFlag)
	p.ShortDescription = "Test manpage generation"
	p.LongDescription = "This is a somewhat `longer' description of what this does.\nWith multiple lines."
	p.AddGroup("Application Options", "The application options", &opts)

	for _, cmd := range p.Commands() {
		cmd.LongDescription = fmt.Sprintf("Longer `%s' description", cmd.Name)
	}

	var buf bytes.Buffer
	p.WriteManPage(&buf)

	got := buf.String()

	tt := time.Now()

	expected := fmt.Sprintf(`.TH TestMan 1 "%s"
.SH NAME
TestMan \- Test manpage generation
.SH SYNOPSIS
\fBTestMan\fP [OPTIONS]
.SH DESCRIPTION
This is a somewhat \fBlonger\fP description of what this does.
With multiple lines.
.SH OPTIONS
.SS Application Options
The application options
.TP
\fB\fB\-v\fR, \fB\-\-verbose\fR\fP
Show verbose debug information
.TP
\fB\fB\-c\fR\fP
Call phone number
.TP
\fB\fB\-\-ptrslice\fR\fP
A slice of pointers to string
.TP
\fB\fB\-\-empty-description\fR\fP
.TP
\fB\fB\-\-default\fR <default: \fI"Some\\nvalue"\fR>\fP
Test default value
.TP
\fB\fB\-\-default-array\fR <default: \fI"Some value", "Other\\tvalue"\fR>\fP
Test default array value
.TP
\fB\fB\-\-default-map\fR <default: \fI"some:value", "another:value"\fR>\fP
Testdefault map value
.TP
\fB\fB\-\-opt-with-arg-name\fR \fIsomething\fR\fP
Option with named argument
.TP
\fB\fB\-\-opt-with-choices\fR \fIchoice\fR\fP
Option with choices
.SS Other Options
.TP
\fB\fB\-s\fR <default: \fI"some", "value"\fR>\fP
A slice of strings
.TP
\fB\fB\-\-intmap\fR <default: \fI"a:1"\fR>\fP
A map from string to int
.SS Subgroup
.TP
\fB\fB\-\-sip.opt\fR\fP
This is a subgroup option
.TP
\fB\fB\-\-sip.not-hidden-inside-group\fR\fP
Not hidden inside group
.SS Subsubgroup
.TP
\fB\fB\-\-sip.sap.opt\fR\fP
This is a subsubgroup option
.SH COMMANDS
.SS bommand
A command with only hidden options

Longer \fBbommand\fP description
.SS command
A command

Longer \fBcommand\fP description

\fBUsage\fP: TestMan [OPTIONS] command [command-OPTIONS]
.TP

\fBAliases\fP: cm, cmd

.TP
\fB\fB\-\-extra-verbose\fR\fP
Use for extra verbosity
.SS parent
A parent command

Longer \fBparent\fP description

\fBUsage\fP: TestMan [OPTIONS] parent [parent-OPTIONS]
.TP
.TP
\fB\fB\-\-opt\fR\fP
This is a parent command option
.SS parent sub
A sub command

\fBUsage\fP: TestMan [OPTIONS] parent [parent-OPTIONS] sub [sub-OPTIONS]
.TP
.TP
\fB\fB\-\-opt\fR\fP
This is a sub command option
`, tt.Format("2 January 2006"))

	assertDiff(t, got, expected, "man page")
}

type helpCommandNoOptions struct {
	Command struct {
	} `command:"command" description:"A command"`
}

func TestHelpCommand(t *testing.T) {
	var opts helpCommandNoOptions
	p := NewNamedParser("TestHelpCommand", HelpFlag)
	p.AddGroup("Application Options", "The application options", &opts)

	_, err := p.ParseArgs([]string{"command", "--help"})

	if err == nil {
		t.Fatalf("Expected help error")
	}

	if e, ok := err.(*Error); !ok {
		t.Fatalf("Expected flags.Error, but got %T", err)
	} else {
		if e.Type != ErrHelp {
			t.Errorf("Expected flags.ErrHelp type, but got %s", e.Type)
		}

		var expected string

		if runtime.GOOS == "windows" {
			expected = `Usage:
  TestHelpCommand [OPTIONS] command

Help Options:
  /?              Show this help message
  /h, /help       Show this help message
`
		} else {
			expected = `Usage:
  TestHelpCommand [OPTIONS] command

Help Options:
  -h, --help      Show this help message
`
		}

		assertDiff(t, e.Message, expected, "help message")
	}
}

func TestHiddenCommandNoBuiltinHelp(t *testing.T) {
	// no auto added help group
	p := NewNamedParser("TestHelpCommand", 0)
	// and no usage information either
	p.Usage = ""

	// add custom help group which is not listed in --help output
	var help struct {
		ShowHelp func() error `short:"h" long:"help"`
	}
	help.ShowHelp = func() error {
		return &Error{Type: ErrHelp}
	}
	hlpgrp, err := p.AddGroup("Help Options", "", &help)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	hlpgrp.Hidden = true
	hlp := p.FindOptionByLongName("help")
	hlp.Description = "Show this help message"
	// make sure the --help option is hidden
	hlp.Hidden = true

	// add a hidden command
	var hiddenCmdOpts struct {
		Foo        bool `short:"f" long:"very-long-foo-option" description:"Very long foo description"`
		Bar        bool `short:"b" description:"Option bar"`
		Positional struct {
			PositionalFoo string `positional-arg-name:"<positional-foo>" description:"positional foo"`
		} `positional-args:"yes"`
	}
	cmdHidden, err := p.Command.AddCommand("hidden", "Hidden command description", "Long hidden command description", &hiddenCmdOpts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// make it hidden
	cmdHidden.Hidden = true
	if len(cmdHidden.Options()) != 2 {
		t.Fatalf("unexpected options count")
	}
	// which help we ask for explicitly
	_, err = p.ParseArgs([]string{"hidden", "--help"})

	if err == nil {
		t.Fatalf("Expected help error")
	}
	if e, ok := err.(*Error); !ok {
		t.Fatalf("Expected flags.Error, but got %T", err)
	} else {
		if e.Type != ErrHelp {
			t.Errorf("Expected flags.ErrHelp type, but got %s", e.Type)
		}

		var expected string

		if runtime.GOOS == "windows" {
			expected = `Usage:
  TestHelpCommand hidden [hidden-OPTIONS] [<positional-foo>]

Long hidden command description

[hidden command arguments]
  <positional-foo>:         positional foo
`
		} else {
			expected = `Usage:
  TestHelpCommand hidden [hidden-OPTIONS] [<positional-foo>]

Long hidden command description

[hidden command arguments]
  <positional-foo>:         positional foo
`
		}
		h := &bytes.Buffer{}
		p.WriteHelp(h)

		assertDiff(t, h.String(), expected, "help message")
	}
}

func TestHelpDefaults(t *testing.T) {
	var expected string

	if runtime.GOOS == "windows" {
		expected = `Usage:
  TestHelpDefaults [OPTIONS]

Application Options:
      /with-default:               With default (default: default-value)
      /without-default:            Without default
      /with-programmatic-default:  With programmatic default (default:
                                   default-value)

Help Options:
  /?                               Show this help message
  /h, /help                        Show this help message
`
	} else {
		expected = `Usage:
  TestHelpDefaults [OPTIONS]

Application Options:
      --with-default=              With default (default: default-value)
      --without-default=           Without default
      --with-programmatic-default= With programmatic default (default:
                                   default-value)

Help Options:
  -h, --help                       Show this help message
`
	}

	tests := []struct {
		Args   []string
		Output string
	}{
		{
			Args:   []string{"-h"},
			Output: expected,
		},
		{
			Args:   []string{"--with-default", "other-value", "--with-programmatic-default", "other-value", "-h"},
			Output: expected,
		},
	}

	for _, test := range tests {
		var opts struct {
			WithDefault             string `long:"with-default" default:"default-value" description:"With default"`
			WithoutDefault          string `long:"without-default" description:"Without default"`
			WithProgrammaticDefault string `long:"with-programmatic-default" description:"With programmatic default"`
		}

		opts.WithProgrammaticDefault = "default-value"

		p := NewNamedParser("TestHelpDefaults", HelpFlag)
		p.AddGroup("Application Options", "The application options", &opts)

		_, err := p.ParseArgs(test.Args)

		if err == nil {
			t.Fatalf("Expected help error")
		}

		if e, ok := err.(*Error); !ok {
			t.Fatalf("Expected flags.Error, but got %T", err)
		} else {
			if e.Type != ErrHelp {
				t.Errorf("Expected flags.ErrHelp type, but got %s", e.Type)
			}

			assertDiff(t, e.Message, test.Output, "help message")
		}
	}
}

func TestHelpRestArgs(t *testing.T) {
	opts := struct {
		Verbose bool `short:"v"`
	}{}

	p := NewNamedParser("TestHelpDefaults", HelpFlag)
	p.AddGroup("Application Options", "The application options", &opts)

	retargs, err := p.ParseArgs([]string{"-h", "-v", "rest"})

	if err == nil {
		t.Fatalf("Expected help error")
	}

	assertStringArray(t, retargs, []string{"-v", "rest"})
}

func TestWrapText(t *testing.T) {
	s := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

	got := wrapText(s, 60, "      ")
	expected := `Lorem ipsum dolor sit amet, consectetur adipisicing elit,
      sed do eiusmod tempor incididunt ut labore et dolore magna
      aliqua. Ut enim ad minim veniam, quis nostrud exercitation
      ullamco laboris nisi ut aliquip ex ea commodo consequat.
      Duis aute irure dolor in reprehenderit in voluptate velit
      esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
      occaecat cupidatat non proident, sunt in culpa qui officia
      deserunt mollit anim id est laborum.`

	assertDiff(t, got, expected, "wrapped text")
}

func TestWrapParagraph(t *testing.T) {
	s := "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\n\n"
	s += "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n\n"
	s += "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.\n\n"
	s += "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n"

	got := wrapText(s, 60, "      ")
	expected := `Lorem ipsum dolor sit amet, consectetur adipisicing elit,
      sed do eiusmod tempor incididunt ut labore et dolore magna
      aliqua.

      Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat.

      Duis aute irure dolor in reprehenderit in voluptate velit
      esse cillum dolore eu fugiat nulla pariatur.

      Excepteur sint occaecat cupidatat non proident, sunt in
      culpa qui officia deserunt mollit anim id est laborum.
`

	assertDiff(t, got, expected, "wrapped paragraph")
}

func TestHelpDefaultMask(t *testing.T) {
	var tests = []struct {
		opts    interface{}
		present string
	}{
		{
			opts: &struct {
				Value string `short:"v" default:"123" description:"V"`
			}{},
			present: "V (default: 123)\n",
		},
		{
			opts: &struct {
				Value string `short:"v" default:"123" default-mask:"abc" description:"V"`
			}{},
			present: "V (default: abc)\n",
		},
		{
			opts: &struct {
				Value string `short:"v" default:"123" default-mask:"-" description:"V"`
			}{},
			present: "V\n",
		},
		{
			opts: &struct {
				Value string `short:"v" description:"V"`
			}{Value: "123"},
			present: "V (default: 123)\n",
		},
		{
			opts: &struct {
				Value string `short:"v" default-mask:"abc" description:"V"`
			}{Value: "123"},
			present: "V (default: abc)\n",
		},
		{
			opts: &struct {
				Value string `short:"v" default-mask:"-" description:"V"`
			}{Value: "123"},
			present: "V\n",
		},
	}

	for _, test := range tests {
		p := NewParser(test.opts, HelpFlag)
		_, err := p.ParseArgs([]string{"-h"})
		if flagsErr, ok := err.(*Error); ok && flagsErr.Type == ErrHelp {
			err = nil
		}
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		h := &bytes.Buffer{}
		w := bufio.NewWriter(h)
		p.writeHelpOption(w, p.FindOptionByShortName('v'), p.getAlignmentInfo())
		w.Flush()
		if strings.Index(h.String(), test.present) < 0 {
			t.Errorf("Not present %q\n%s", test.present, h.String())
		}
	}
}

func TestWroteHelp(t *testing.T) {
	type testInfo struct {
		value  error
		isHelp bool
	}
	tests := map[string]testInfo{
		"No error":    {value: nil, isHelp: false},
		"Plain error": {value: errors.New("an error"), isHelp: false},
		"ErrUnknown":  {value: newError(ErrUnknown, "an error"), isHelp: false},
		"ErrHelp":     {value: newError(ErrHelp, "an error"), isHelp: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := WroteHelp(test.value)
			if test.isHelp != res {
				t.Errorf("Expected %t, got %t", test.isHelp, res)
			}
		})
	}
}
