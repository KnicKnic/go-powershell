package powershell

import "github.com/golang/glog"

type LoggerFull interface {
	Warning(args ...interface{})
	Information(args ...interface{})
	Verbose(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Write(args ...interface{})
	Warningln(args ...interface{})
	Informationln(args ...interface{})
	Verboseln(args ...interface{})
	Debugln(args ...interface{})
	Errorln(args ...interface{})
	Writeln(args ...interface{})
}

type LoggerSimple interface {
	Write(args ...interface{})
}

type LogHolder struct {
	Log LoggerFull
}

type SimpleToFull struct {
	simple LoggerSimple
}

func MakeLogHolderFull(logger LoggerFull) LogHolder {
	return LogHolder{logger}
}

func MakeLogHolder(logger LoggerSimple) LogHolder {
	if p, ok := logger.(LoggerFull); ok {
		return MakeLogHolderFull(p)
	}
	return MakeLogHolderFull(SimpleToFull{logger})
}

func Make2ArgInterface(arg1 interface{}, args ...interface{}) []interface{} {

	a := append([]interface{}{arg1}, args...)
	return a
}

func AddInterfaceFront(argEnd interface{}, args ...interface{}) []interface{} {

	a := append(args, argEnd)
	return a
}

func (log SimpleToFull) Warning(args ...interface{}) {

	log.Write(Make2ArgInterface("Warning", args...))
}
func (log SimpleToFull) Information(args ...interface{}) {
	log.Write(Make2ArgInterface("Information: ", args...))
}
func (log SimpleToFull) Verbose(args ...interface{}) {
	log.Write(Make2ArgInterface("Verbose: ", args...))
}
func (log SimpleToFull) Debug(args ...interface{}) {
	log.Write(Make2ArgInterface("Debug: ", args...))
}
func (log SimpleToFull) Error(args ...interface{}) {
	log.Write(Make2ArgInterface("Error: ", args...))
}
func (log SimpleToFull) Write(args ...interface{}) {
	log.Write(args...)
}

func ArgsWithNewLine(level interface{}, args ...interface{}) []interface{} {
	line := Make2ArgInterface(level, args...)
	return AddInterfaceFront('\n', line...)
}
func (log SimpleToFull) Warningln(args ...interface{}) {
	log.Write(ArgsWithNewLine("Warning: ", args...))
}
func (log SimpleToFull) Informationln(args ...interface{}) {
	log.Write(ArgsWithNewLine("Information: ", args...))
}
func (log SimpleToFull) Verboseln(args ...interface{}) {
	log.Write(ArgsWithNewLine("Verbose: ", args...))
}
func (log SimpleToFull) Debugln(args ...interface{}) {
	log.Write(ArgsWithNewLine("Debug: ", args...))
}
func (log SimpleToFull) Errorln(args ...interface{}) {
	log.Write(ArgsWithNewLine("Error: ", args...))
}
func (log SimpleToFull) Writeln(args ...interface{}) {
	log.Write(AddInterfaceFront('\n', args...))
}

type GLogInfoLogger struct {
}

func (logger GLogInfoLogger) Write(args ...interface{}) {
	glog.Info(args...)
}
