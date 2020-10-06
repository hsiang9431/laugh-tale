package cli_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"laugh-tale/pkg/common/cli"
)

// Run with command: go test --count=1 -run TestPrint
// to see the result from stdout
func TestPrint(t *testing.T) {

	cli.PrintHelp(&appCmd, os.Stdout, nil)
	fmt.Println("------------------------------")

	cli.PrintUsage(&appCmd, os.Stdout, nil)
	fmt.Println("------------------------------")

	cli.PrintUsage(&appCmd, os.Stdout, &printSettings)
	fmt.Println("")
	cli.PrintUsage(&all, os.Stdout, &printSettings)
	fmt.Println("")
	cli.PrintUsage(&task1, os.Stdout, &printSettings)
	fmt.Println("")
	cli.PrintUsage(&task2, os.Stdout, &printSettings)
	fmt.Println("")
	cli.PrintUsage(&task3, os.Stdout, &printSettings)

	fmt.Println("------------------------------")
	err := errors.New("test error")
	cli.PrintMisuse(&setup, os.Stdout, err, nil)
	fmt.Println("")
	cli.PrintMisuse(&all, os.Stdout, err, nil)
	fmt.Println("")
	cli.PrintMisuse(&task1, os.Stdout, err, nil)
	fmt.Println("")
	cli.PrintMisuse(&task2, os.Stdout, err, nil)
	fmt.Println("")
	cli.PrintMisuse(&task3, os.Stdout, err, nil)
}

type testParseCmdAns struct {
	cmdName  string
	errIsNil bool
}

func testParseCmdHelper(t *testing.T, testSet [][]string, ansSet []testParseCmdAns, testName string) {
	for i := 0; i < len(testSet); i++ {
		cmd, _, err := cli.ParseCmd(&appCmd, testSet[i])
		if err != nil {
			t.Log("err encountered: ", err.Error())
		}
		// t.Log(cmd)
		// t.Log(err)
		ans := ansSet[i]
		if (cmd != nil && ans.cmdName != cmd.Name) ||
			((err == nil) != ansSet[i].errIsNil) {
			t.Errorf("TestParseCmd failed on #%d %s case: %v", i+1, testName, testSet[i])
		}
	}
}

func TestParseCmd(t *testing.T) {
	// good
	argsGood := [][]string{
		{"app", "help"},
		{"app", "-h"},
		{"app", "--help"},
		{"app", "start"},
		{"app", "status", "--arg11=11", "--arg12=12"},
		{"app", "setup", "task2", "--arg11=11", "--arg12=12"},
	}
	ansGood := []testParseCmdAns{
		{helpShow.Name, true},
		{help1.Name, true},
		{help2.Name, true},
		{start.Name, true},
		{status.Name, true},
		{task2.Name, true},
	}
	testParseCmdHelper(t, argsGood, ansGood, "good")
	// bad
	argsBad := [][]string{
		{"app", "bad"},
		{"app", "setup", "bad"},
	}
	ansBad := []testParseCmdAns{
		{appCmd.Name, false},
		{setup.Name, false},
	}
	testParseCmdHelper(t, argsBad, ansBad, "bad")
}

func TestFlagRequired(t *testing.T) {
	args := [][]string{
		{"app", "setup", "task1", "--arg11=123", "--arg12", "100"},
		{"app", "setup", "task1", "--arg13=123", "--arg12", "100"},
		{"app", "setup", "task2", "--arg21=123", "--arg22", "100"},
		{"app", "setup", "task2", "--arg22=123", "--arg22", "100"},
		{"app", "setup", "task3", "--arg=123", "--arg", "100"},
		{"app", "setup", "task1"},
	}
	noErr := []bool{true, false, true, false, true, false}
	for i, a := range args {
		cmd, ctx, _ := cli.ParseCmd(&appCmd, a)
		ctx, err := cli.ParseCliFlags(cmd, ctx)
		if err != nil {
			t.Log("err encountered: ", err.Error())
		}
		if (err == nil) != noErr[i] {
			t.Errorf("TestValidateCommand failed on #%d case: %v", i+1, args[i])
		}
	}
}

func TestParseBoolAndString(t *testing.T) {
	args := []string{"app", "setup", "task1", "--arg11=ewret", "--arg12", "rgshgs", "--arg12=rgsffhgs", "--arg13"}
	cmd, ctx, err := cli.ParseCmd(&appCmd, args)
	ctx, err = cli.ParseCliFlags(cmd, ctx)
	if err != nil {
		t.Error("TestParseBoolAndString failed", err.Error())
	}
	t.Log(ctx)

	// string
	strAns, err := task1Fg1.GetString(ctx)
	t.Log(strAns)
	if strAns != "ewret" {
		t.Error("TestParseBoolAndString failed wrong string")
	}
	if err != nil {
		t.Error("TestParseBoolAndString failed", err.Error())
	}
	strSliceAns, err := task1Fg2.GetStringSlice(ctx)
	t.Log(strSliceAns)
	if len(strSliceAns) != 2 {
		t.Error("TestParseBoolAndString failed wrong string slice")
	}
	if err != nil {
		t.Error("TestParseBoolAndString failed", err.Error())
	}

	// bool
	boolAnsT := task1Fg1.GetBool(ctx)
	if !boolAnsT {
		t.Error("TestParseBoolAndString failed wrong bool")
	}
	boolAnsF := task1Fg3.GetBool(ctx)
	if !boolAnsF {
		t.Error("TestParseBoolAndString failed wrong bool")
	}
}

func TestParseFlagsInt(t *testing.T) {

	args := []string{"app", "setup", "task1", "--arg11=123", "--arg12", "100", "--arg11", "122", "--arg11=124", "--arg11", "121", "--arg11=125"}
	cmd, ctx, err := cli.ParseCmd(&appCmd, args)
	ctx, err = cli.ParseCliFlags(cmd, ctx)
	if err != nil {
		t.Error("TestParseFlagsInt failed", err.Error())
	}
	t.Log(ctx)

	// GetInt and GetUint calls to their 64 bit respectively
	ansInt, err := task1Fg1.GetInt(ctx)
	if ansInt != 123 {
		t.Error("TestParseFlagsInt failed wrong int value")
	}
	if err != nil {
		t.Error("TestParseFlagsInt failed", err.Error())
	}
	ansUint, err := task1Fg1.GetUint(ctx)
	if ansUint != 123 {
		t.Error("TestParseFlagsInt failed wrong uint value")
	}
	if err != nil {
		t.Error("TestParseFlagsInt failed", err.Error())
	}
	ansIntSlice, err := task1Fg1.GetIntSlice(ctx)
	t.Log(ansIntSlice)
	if len(ansIntSlice) != 5 {
		t.Error("TestParseFlagsInt failed wrong int slice")
	}
	if err != nil {
		t.Error("TestParseFlagsInt failed ansIntSlice", err.Error())
	}
	ansUintSlice, err := task1Fg1.GetUintSlice(ctx)
	t.Log(ansUintSlice)
	if len(ansIntSlice) != 5 {
		t.Error("TestParseFlagsInt failed wrong uint slice")
	}
	if err != nil {
		t.Error("TestParseFlagsInt failed ansUintSlice", err.Error())
	}
}

func TestParseFlagFloat(t *testing.T) {
	args := []string{"app", "setup", "task1", "--arg11=123e2", "--arg12", "5456.8964", "--arg11=122.123", "--arg11=124.123", "--arg11=121e10", "--arg11=125e-2"}
	cmd, ctx, err := cli.ParseCmd(&appCmd, args)
	ctx, err = cli.ParseCliFlags(cmd, ctx)
	if err != nil {
		t.Error("TestParseFlagFloat failed", err.Error())
	}
	t.Log(ctx)

	// GetInt and GetUint calls to their 64 bit respectively
	ansFloat, err := task1Fg1.GetFloat32(ctx)
	t.Log(ansFloat)
	if err != nil {
		t.Error("TestParseFlagFloat failed", err.Error())
	}
	ansFloatSlice, err := task1Fg1.GetFloat32Slice(ctx)
	t.Log(ansFloatSlice)
	if err != nil {
		t.Error("TestParseFlagFloat failed ansFloatSlice", err.Error())
	}
}

func TestParents(t *testing.T) {
	args := [][]string{
		{"app", "setup", "task1"},
		{"app", "setup", "task2"},
		{"app", "setup", "task3"},
		{"app", "setup"},
	}
	for i, a := range args {
		_, ctx, _ := cli.ParseCmd(&appCmd, a)
		t.Log(ctx.Parents)
		// check if parents match
		for j := 0; j < len(args[i])-1; j++ {
			if ctx.Parents[j] != args[i][j] {
				t.Error("Parent not match", ctx.Parents[j])
			}
		}
	}
}
