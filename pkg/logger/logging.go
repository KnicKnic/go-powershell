package logger

import "fmt"

// Full the full logging interface with all functions
type Full interface {
	Warning(arg string)
	Information(arg string)
	Verbose(arg string)
	Debug(arg string)
	Error(arg string)
	Write(arg string)
	Warningln(arg string)
	Informationln(arg string)
	Verboseln(arg string)
	Debugln(arg string)
	Errorln(arg string)
	Writeln(arg string)
}

// Simple is the simplest logging interface you can have
//
// If this is specified, it will get wrapped into the full interface by prepending the category "Warning: ", "Information: "... and appending a \n for the *ln functions (Writeln,...)
type Simple interface {
	Write(arg string)
}

type logHolder struct {
	Log Full
}

type simpleToFull struct {
	simple Simple
}

// SimpleFmtPrint is a Simple logger that calls fmt.Print
type SimpleFmtPrint struct {
}

func (SimpleFmtPrint) Write(arg string) {
	fmt.Print(arg)
}

// SimpleFuncPtr is a Simple logger that allows you to pass in a function pointer for the Write call
type SimpleFuncPtr struct {
	FuncPtr func(string)
}

func (holder SimpleFuncPtr) Write(arg string) {
	holder.FuncPtr(arg)
}

// func makeLogHolderFull(logger Full) Full {
// 	return logger
// }

// MakeLoggerFull returns a wrapper class that provides Full semantics,
// utilizing a simple Simple.write() function
func MakeLoggerFull(logger Simple) Full {
	if logger == nil {
		return nil
	}
	if p, ok := logger.(Full); ok {
		return p
	}
	return simpleToFull{logger}
}

func (log simpleToFull) Warning(arg string) {

	log.Write(fmt.Sprint("Warning: ", arg))
}
func (log simpleToFull) Information(arg string) {
	log.Write(fmt.Sprint("Information: ", arg))
}
func (log simpleToFull) Verbose(arg string) {
	log.Write(fmt.Sprint("Verbose: ", arg))
}
func (log simpleToFull) Debug(arg string) {
	log.Write(fmt.Sprint("Debug: ", arg))
}
func (log simpleToFull) Error(arg string) {
	log.Write(fmt.Sprint("Error: ", arg))
}
func (log simpleToFull) Write(arg string) {
	log.simple.Write(arg)
}

func (log simpleToFull) Warningln(arg string) {
	log.Write(fmt.Sprintln("Warning: ", arg))
}
func (log simpleToFull) Informationln(arg string) {
	log.Write(fmt.Sprintln("Information: ", arg))
}
func (log simpleToFull) Verboseln(arg string) {
	log.Write(fmt.Sprintln("Verbose: ", arg))
}
func (log simpleToFull) Debugln(arg string) {
	log.Write(fmt.Sprintln("Debug: ", arg))
}
func (log simpleToFull) Errorln(arg string) {
	log.Write(fmt.Sprintln("Error: ", arg))
}
func (log simpleToFull) Writeln(arg string) {
	log.Write(fmt.Sprintln(arg))
}
