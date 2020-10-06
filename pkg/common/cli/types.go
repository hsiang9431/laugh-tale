package cli

type Command struct {
	Name        string
	DispStr     string
	Description string

	SubCmd     []Command
	SubCmdDesc string
	Flags      []Flag

	RunFunc func(*Context) error
}

type Flag struct {
	Name        string
	DispStr     string
	Description string

	Required bool
}

type Context struct {
	FlagsParsed bool
	Parents     []string
	rawInput    []string
	mappedFlags map[string][]string
}

type PrintSettings struct {
	HelpStr          string
	UnknownCmdErrStr string

	DefaultIndent    int
	DefaultTabWidth  int
	PrintRequiredStr string
}
