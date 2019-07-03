package logger

import "fmt"

// LoggerFull the full logging interface with all functions
type LoggerFull interface {
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

// LoggerSimple is the simplest logging interface you can have
//
// If this is specified, it will get wrapped into the full interface by prepending the category "Warning: ", "Information: "... and appending a \n for the *ln functions (Writeln,...)
type LoggerSimple interface {
	Write(arg string)
}

type logHolder struct {
	Log LoggerFull
}

type simpleToFull struct {
	simple LoggerSimple
}

// func makeLogHolderFull(logger LoggerFull) LoggerFull {
// 	return logger
// }

// MakeLoggerFull returns a wrapper class that provides LoggerFull semantics,
// utilizing a simple LoggerSimple.write() function
func MakeLoggerFull(logger LoggerSimple) LoggerFull {
	if logger == nil {
		return nil
	}
	if p, ok := logger.(LoggerFull); ok {
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
	log.Write(fmt.Sprintln('\n', arg))
}
