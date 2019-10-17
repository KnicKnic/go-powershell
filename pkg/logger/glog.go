package logger

import (
	"github.com/golang/glog"
)

// Glog a type that can be used to log messages to glog
type Glog struct {
	ErrorLevel       glog.Level
	WarningLevel     glog.Level
	InformationLevel glog.Level
	VerboseLevel     glog.Level
	DebugLevel       glog.Level
}

func (log Glog) Warning(arg string) {
	if glog.V(log.WarningLevel) {
		glog.Warning(arg)
	}
}
func (log Glog) Information(arg string) {
	if glog.V(log.InformationLevel) {
		glog.Info(arg)
	}
}
func (log Glog) Verbose(arg string) {
	if glog.V(log.VerboseLevel) {
		glog.Info(arg)
	}
}
func (log Glog) Debug(arg string) {
	if glog.V(log.DebugLevel) {
		glog.Info(arg)
	}
}
func (log Glog) Error(arg string) {
	if glog.V(log.ErrorLevel) {
		glog.Error(arg)
	}
}
func (log Glog) Write(arg string) {
	if glog.V(log.InformationLevel) {
		glog.Info(arg)
	}
}

func (log Glog) Warningln(arg string) {
	if glog.V(log.WarningLevel) {
		glog.Warningln(arg)
	}
}
func (log Glog) Informationln(arg string) {
	if glog.V(log.InformationLevel) {
		glog.Infoln(arg)
	}
}
func (log Glog) Verboseln(arg string) {
	if glog.V(log.VerboseLevel) {
		glog.Infoln(arg)
	}
}
func (log Glog) Debugln(arg string) {
	if glog.V(log.DebugLevel) {
		glog.Infoln(arg)
	}
}
func (log Glog) Errorln(arg string) {
	if glog.V(log.ErrorLevel) {
		glog.Errorln(arg)
	}
}
func (log Glog) Writeln(arg string) {
	if glog.V(log.InformationLevel) {
		glog.Infoln(arg)
	}
}
