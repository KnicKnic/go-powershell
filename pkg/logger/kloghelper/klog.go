package kloghelper

import (
	"k8s.io/klog"
)

// Klog a type that can be used to log messages to klog
type Klog struct {
	ErrorLevel       klog.Level
	WarningLevel     klog.Level
	InformationLevel klog.Level
	VerboseLevel     klog.Level
	DebugLevel       klog.Level
}

func (log Klog) Warning(arg string) {
	if klog.V(log.WarningLevel) {
		klog.Warning(arg)
	}
}
func (log Klog) Information(arg string) {
	if klog.V(log.InformationLevel) {
		klog.Info(arg)
	}
}
func (log Klog) Verbose(arg string) {
	if klog.V(log.VerboseLevel) {
		klog.Info(arg)
	}
}
func (log Klog) Debug(arg string) {
	if klog.V(log.DebugLevel) {
		klog.Info(arg)
	}
}
func (log Klog) Error(arg string) {
	if klog.V(log.ErrorLevel) {
		klog.Error(arg)
	}
}
func (log Klog) Write(arg string) {
	if klog.V(log.InformationLevel) {
		klog.Info(arg)
	}
}

func (log Klog) Warningln(arg string) {
	if klog.V(log.WarningLevel) {
		klog.Warningln(arg)
	}
}
func (log Klog) Informationln(arg string) {
	if klog.V(log.InformationLevel) {
		klog.Infoln(arg)
	}
}
func (log Klog) Verboseln(arg string) {
	if klog.V(log.VerboseLevel) {
		klog.Infoln(arg)
	}
}
func (log Klog) Debugln(arg string) {
	if klog.V(log.DebugLevel) {
		klog.Infoln(arg)
	}
}
func (log Klog) Errorln(arg string) {
	if klog.V(log.ErrorLevel) {
		klog.Errorln(arg)
	}
}
func (log Klog) Writeln(arg string) {
	if klog.V(log.InformationLevel) {
		klog.Infoln(arg)
	}
}
