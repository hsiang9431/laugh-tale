package cli

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCmd    = errors.New("Invalid command, pointer to command can not be nil")
	ErrInvalidCtx    = errors.New("Invalid context, pointer to context can not be nil")
	ErrCmdParse      = errors.New("Failed to parse command")
	ErrFlagNotParsed = errors.New("Unparsed flags")
	ErrNotAbleToRun  = errors.New("Not function assigned to command")
)

// ParseCmd checks if the commands are valid and return corresponding
// cli.Command struct. It returns the last parsed command and a
// cli.Context structure which is supposed to contain only command line flags.
// This function should be called on the root of Cmd tree of an application
// The parse is succeeded if the command returned has no sub command.
func ParseCmd(app *Command, args []string) (*Command, *Context, error) {
	if app == nil {
		return nil, nil, ErrInvalidCmd
	}
	if len(args) < 2 {
		if app.SubCmd == nil {
			return app, &Context{rawInput: args[1:]}, nil
		}
		return app, &Context{rawInput: args[1:]}, ErrCmdParse
	}
	curNode := app
	argIdx := 1
	subCmdIdx := 0
	var parents []string
	for subCmdIdx < len(curNode.SubCmd) &&
		argIdx < len(args) {
		if curNode.SubCmd[subCmdIdx].Name == args[argIdx] {
			parents = append(parents, curNode.Name)
			curNode = &curNode.SubCmd[subCmdIdx]
			subCmdIdx = 0
			argIdx++
		} else {
			subCmdIdx++
		}
	}
	if len(curNode.SubCmd) != 0 {
		return curNode, &Context{Parents: parents, rawInput: args[argIdx:]}, ErrCmdParse
	}
	return curNode, &Context{Parents: parents, rawInput: args[argIdx:]}, nil
}

// ParseCliFlags parses command line flags for the Cmd from which
// its called and return the result in a parsed Context struct,
// which can be access with receiver functions of Flag structures.
func ParseCliFlags(cmd *Command, ctx *Context) (*Context, error) {
	if cmd == nil {
		return ctx, ErrInvalidCmd
	}
	if ctx == nil {
		return ctx, ErrInvalidCtx
	}
	argsMap := make(map[string][]string)
	var lastOpt string
	for i := 0; i < len(ctx.rawInput); i++ {
		curStr := ctx.rawInput[i]
		if lastOpt != "" {
			addFlag(argsMap, lastOpt, curStr)
			lastOpt = ""
		}
		if strings.HasPrefix(curStr, "--") {
			equalSign := strings.Index(curStr, "=")
			if equalSign > -1 {
				addFlag(argsMap, curStr[2:equalSign], curStr[equalSign+1:])
			} else {
				lastOpt = curStr[2:]
			}
		}
	}
	if lastOpt != "" {
		if _, ok := argsMap[lastOpt]; !ok {
			argsMap[lastOpt] = []string{""}
		}
	}
	var retErr error
	missingFlags := ""
	for _, f := range cmd.Flags {
		if f.Required {
			vals, ok := argsMap[f.Name]
			if !ok || len(vals) == 0 {
				missingFlags += ", " + f.Name
			}
		}
	}
	if missingFlags != "" {
		retErr = errors.New("Missing required flags: " + missingFlags[2:])
	}
	ctx.FlagsParsed = true
	ctx.mappedFlags = argsMap
	return ctx, retErr
}

func Call(cmd *Command, ctx *Context) error {
	if cmd == nil {
		return ErrInvalidCmd
	}
	if !ctx.FlagsParsed {
		return ErrFlagNotParsed
	}
	if cmd.RunFunc == nil {
		return ErrNotAbleToRun
	}
	return cmd.RunFunc(ctx)
}

func addFlag(m map[string][]string, key, val string) {
	if _, ok := m[key]; ok {
		m[key] = append(m[key], val)
	} else {
		m[key] = []string{val}
	}
}
