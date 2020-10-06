package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

var defaultPS = PrintSettings{
	HelpStr:          "Usage",
	UnknownCmdErrStr: "Invalid usage of command: ",

	DefaultIndent:    4,
	DefaultTabWidth:  8,
	PrintRequiredStr: "<Required>",
}

// PrintUsage prints sub commands and flags of this command
func PrintUsage(cmd *Command, w io.Writer, ps *PrintSettings) {
	if cmd == nil || w == nil {
		return
	}
	ps = validatePrintSettings(ps)
	fmt.Fprintln(w, ps.HelpStr)
	fmt.Fprintln(w, indentString(ps.DefaultIndent)+cmd.DispStr)

	if cmd.Flags != nil {
		cmd.printFlags(w, indentString(ps.DefaultIndent*2), ps)
	} else if cmd.SubCmd != nil {
		cmd.printSubCmds(w, ps)
	}
}

// PrintHelp prints sub commands and flags of this command,
// then print details of all subcommands with BFS
func PrintHelp(cmd *Command, w io.Writer, ps *PrintSettings) {
	if cmd == nil || w == nil {
		return
	}
	ps = validatePrintSettings(ps)
	fmt.Fprintln(w, ps.HelpStr)
	fmt.Fprintln(w, indentString(ps.DefaultIndent)+cmd.DispStr)

	var curCmd Command
	if cmd.SubCmd != nil {
		bfsQueue := append(make([]Command, 0), *cmd)
		for len(bfsQueue) > 0 {
			curCmd, bfsQueue = bfsQueue[0], bfsQueue[1:]
			if curCmd.SubCmd != nil {
				curCmd.printSubCmds(w, ps)
				bfsQueue = append(bfsQueue, curCmd.SubCmd...)
			}
		}
	}
}

// PrintMisuse prints an error message regarding the invalid use of command
// then calls cmd.PrintUsage() with the same io.Writer
func PrintMisuse(cmd *Command, w io.Writer, err error, ps *PrintSettings) {
	if cmd == nil || w == nil {
		return
	}
	ps = validatePrintSettings(ps)
	fmt.Fprintln(w, ps.UnknownCmdErrStr+err.Error())
	PrintUsage(cmd, w, ps)
}

// ----------------------------------------------------
// Unexported functions
// ----------------------------------------------------
func validatePrintSettings(ps *PrintSettings) *PrintSettings {
	if ps == nil {
		return &defaultPS
	}
	if ps.DefaultIndent < 1 ||
		ps.DefaultTabWidth < 1 {
		return &defaultPS
	}
	return ps
}

func indentString(leadingSpaces int) string {
	return strings.Repeat(" ", leadingSpaces)
}

func (cmd *Command) printSubCmds(w io.Writer, ps *PrintSettings) {
	if cmd.SubCmd == nil ||
		cmd.SubCmdDesc == "" {
		return
	}
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, cmd.SubCmdDesc)

	tabW := new(tabwriter.Writer)
	defer tabW.Flush()
	tabW.Init(w, ps.DefaultTabWidth, ps.DefaultTabWidth, 2, '\t', 0)

	indent := indentString(ps.DefaultIndent)
	twoIndent := indentString(ps.DefaultIndent * 2)
	for i := 0; i < len(cmd.SubCmd); i++ {
		curCmd := cmd.SubCmd[i]
		if curCmd.DispStr != "" {
			if curCmd.Flags == nil {
				fmt.Fprintln(tabW, indent+curCmd.DispStr+"\t"+curCmd.Description)
			} else {
				fmt.Fprintln(w, indent+curCmd.DispStr+": "+curCmd.Description)
				curCmd.printFlags(w, twoIndent, ps)
			}
		}
	}
}

func (cmd *Command) printFlags(w io.Writer, indent string, ps *PrintSettings) {
	if cmd.Flags == nil {
		return
	}
	tabW := new(tabwriter.Writer)
	defer tabW.Flush()
	tabW.Init(w, ps.DefaultTabWidth, ps.DefaultTabWidth, 2, '\t', 0)

	for i := 0; i < len(cmd.Flags); i++ {
		curFlag := cmd.Flags[i]
		curFlagDesc := curFlag.Description
		if curFlagDesc != "" {
			if curFlag.Required {
				curFlagDesc += " " + ps.PrintRequiredStr
			}
			fmt.Fprintln(tabW, indent+curFlag.DispStr+"\t"+curFlagDesc)
		}
	}
}
