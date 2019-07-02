package powershell

// LoggerFull the full logging interface with all functions
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

// LoggerSimple is the simplest logging interface you can have
//
// If this is specified, it will get wrapped into the full interface by prepending the category "Warning: ", "Information: "... and appending a \n for the *ln functions (Writeln,...)
type LoggerSimple interface {
	Write(args ...interface{})
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
func makeLoggerFull(logger LoggerSimple) LoggerFull {
	if logger == nil {
		return nil
	}
	if p, ok := logger.(LoggerFull); ok {
		return p
	}
	return simpleToFull{logger}
}

func make2ArgInterface(arg1 interface{}, args ...interface{}) []interface{} {

	a := append([]interface{}{arg1}, args...)
	return a
}

func addInterfaceFront(argEnd interface{}, args ...interface{}) []interface{} {

	a := append(args, argEnd)
	return a
}

func (log simpleToFull) Warning(args ...interface{}) {

	log.Write(make2ArgInterface("Warning", args...))
}
func (log simpleToFull) Information(args ...interface{}) {
	log.Write(make2ArgInterface("Information: ", args...))
}
func (log simpleToFull) Verbose(args ...interface{}) {
	log.Write(make2ArgInterface("Verbose: ", args...))
}
func (log simpleToFull) Debug(args ...interface{}) {
	log.Write(make2ArgInterface("Debug: ", args...))
}
func (log simpleToFull) Error(args ...interface{}) {
	log.Write(make2ArgInterface("Error: ", args...))
}
func (log simpleToFull) Write(args ...interface{}) {
	log.simple.Write(args...)
}

func argsWithNewLine(level interface{}, args ...interface{}) []interface{} {
	line := make2ArgInterface(level, args...)
	return addInterfaceFront('\n', line...)
}
func (log simpleToFull) Warningln(args ...interface{}) {
	log.Write(argsWithNewLine("Warning: ", args...))
}
func (log simpleToFull) Informationln(args ...interface{}) {
	log.Write(argsWithNewLine("Information: ", args...))
}
func (log simpleToFull) Verboseln(args ...interface{}) {
	log.Write(argsWithNewLine("Verbose: ", args...))
}
func (log simpleToFull) Debugln(args ...interface{}) {
	log.Write(argsWithNewLine("Debug: ", args...))
}
func (log simpleToFull) Errorln(args ...interface{}) {
	log.Write(argsWithNewLine("Error: ", args...))
}
func (log simpleToFull) Writeln(args ...interface{}) {
	log.Write(addInterfaceFront('\n', args...))
}
