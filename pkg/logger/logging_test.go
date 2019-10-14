package logger

import (
	"fmt"
	"testing"
)

type simpleTester struct {
	t        *testing.T
	previous string
}

func (log *simpleTester) ValidateWarning(arg string) {

	if log.previous != fmt.Sprint("Warning: ", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateInformation(arg string) {
	if log.previous != fmt.Sprint("Information: ", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateVerbose(arg string) {
	if log.previous != fmt.Sprint("Verbose: ", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateDebug(arg string) {
	if log.previous != fmt.Sprint("Debug: ", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateError(arg string) {
	if log.previous != fmt.Sprint("Error: ", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateWrite(arg string) {
	if log.previous != arg {
		log.t.Fail()
	}
}

func (log *simpleTester) ValidateWarningln(arg string) {
	if log.previous != fmt.Sprintln("Warning:", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateInformationln(arg string) {
	if log.previous != fmt.Sprintln("Information:", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateVerboseln(arg string) {
	if log.previous != fmt.Sprintln("Verbose:", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateDebugln(arg string) {
	if log.previous != fmt.Sprintln("Debug:", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateErrorln(arg string) {
	if log.previous != fmt.Sprintln("Error:", arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) ValidateWriteln(arg string) {
	if log.previous != fmt.Sprintln(arg) {
		log.t.Fail()
	}
}
func (log *simpleTester) Write(arg string) {
	log.previous = arg
}

func TestMakeLoggerFull(t *testing.T) {
	// test null simple
	if nil != MakeLoggerFull(nil) {
		t.Fail()
	}

	// test simple wrapper
	simple := &simpleTester{t: t}
	full := MakeLoggerFull(simple)

	// test all variations of simple
	testWord := " testWord with spaces "
	full.Warning(testWord)
	simple.ValidateWarning(testWord)
	full.Information(testWord)
	simple.ValidateInformation(testWord)
	full.Verbose(testWord)
	simple.ValidateVerbose(testWord)
	full.Debug(testWord)
	simple.ValidateDebug(testWord)
	full.Error(testWord)
	simple.ValidateError(testWord)
	full.Write(testWord)
	simple.ValidateWrite(testWord)
	full.Warningln(testWord)
	simple.ValidateWarningln(testWord)
	full.Informationln(testWord)
	simple.ValidateInformationln(testWord)
	full.Verboseln(testWord)
	simple.ValidateVerboseln(testWord)
	full.Debugln(testWord)
	simple.ValidateDebugln(testWord)
	full.Errorln(testWord)
	simple.ValidateErrorln(testWord)
	full.Writeln(testWord)
	simple.ValidateWriteln(testWord)

	// test Full
	full2 := MakeLoggerFull(full)
	if full2 != full {
		t.Fail()
	}
}

func TestSimpleFmtPrint(t *testing.T) {
	SimpleFmtPrint{}.Write("hello\n")
}

func TestSimpleFuncPtr(t *testing.T) {
	var lastWrite string
	logger := SimpleFuncPtr{func(line string) { lastWrite = line }}
	logger.Write("test")
	if lastWrite != "test" {
		t.Fail()
	}
}
