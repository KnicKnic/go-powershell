package powershell

import (
	"syscall"
	"unsafe"
)

func loggerCallback(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Write(s)
	}
	return uintptr(0)
}
func loggerCallbackln(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Writeln(s)
	}
	return uintptr(0)
}
func loggerCallbackDebug(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Debug(s)
	}
	return uintptr(0)
}
func loggerCallbackDebugln(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Debugln(s)
	}
	return uintptr(0)
}
func loggerCallbackVerbose(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Verbose(s)
	}
	return uintptr(0)
}
func loggerCallbackVerboseln(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Verboseln(s)
	}
	return uintptr(0)
}
func loggerCallbackInformation(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Information(s)
	}
	return uintptr(0)
}
func loggerCallbackInformationln(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Informationln(s)
	}
	return uintptr(0)
}
func loggerCallbackWarning(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Warning(s)
	}
	return uintptr(0)
}
func loggerCallbackWarningln(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Warningln(s)
	}
	return uintptr(0)
}
func loggerCallbackError(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Error(s)
	}
	return uintptr(0)
}
func loggerCallbackErrorln(context uint64, str uintptr) uintptr {
	if context != 0 {
		s := uintptrMakeString(str)
		contextInterface := getRunspaceContext(context)
		contextInterface.log.Errorln(s)
	}
	return uintptr(0)
}

func commandCallback(context uint64, cMessage uintptr, input uintptr, inputCount uint64, ret uintptr) uintptr {

	var resultsWriter callbackResultsWriter
	if context != 0 {
		contextInterface := getRunspaceContext(context)
		inputArr := make([]Object, inputCount)
		for i := uint32(0); uint64(i) < inputCount; i++ {
			inputArr[i] = makePowerShellObjectIndexed(input, i)
		}
		message := uintptrMakeString(cMessage)
		contextInterface.callback.Callback(contextInterface.recreateRunspace(), message, inputArr, &resultsWriter)
	}
	resultsWriter.filloutResults(ret)
	return uintptr(0)
}

var (
	loggerCallbackHolder nativePowerShell_LogString_Holder = nativePowerShell_LogString_Holder{
		Log:                syscall.NewCallbackCDecl(loggerCallback),
		LogError:           syscall.NewCallbackCDecl(loggerCallbackError),
		LogWarning:         syscall.NewCallbackCDecl(loggerCallbackWarning),
		LogInformation:     syscall.NewCallbackCDecl(loggerCallbackInformation),
		LogVerbose:         syscall.NewCallbackCDecl(loggerCallbackVerbose),
		LogDebug:           syscall.NewCallbackCDecl(loggerCallbackDebug),
		LogLine:            syscall.NewCallbackCDecl(loggerCallbackln),
		LogErrorLine:       syscall.NewCallbackCDecl(loggerCallbackErrorln),
		LogWarningLine:     syscall.NewCallbackCDecl(loggerCallbackWarningln),
		LogInformationLine: syscall.NewCallbackCDecl(loggerCallbackInformationln),
		LogVerboseLine:     syscall.NewCallbackCDecl(loggerCallbackVerboseln),
		LogDebugLine:       syscall.NewCallbackCDecl(loggerCallbackDebugln),
	}
	loggerCallbackPointer  unsafe.Pointer = unsafe.Pointer(&loggerCallbackHolder)
	commandCallbackPointer uintptr        = syscall.NewCallbackCDecl(commandCallback)
)

func createRunspaceHelper(context uint64, useLogger byte, useCommand byte) nativePowerShell_RunspaceHandle {
	var commandPtr uintptr = 0
	var loggerPtr uintptr = 0
	if useLogger != 0 {
		loggerPtr = uintptr(loggerCallbackPointer)
	}
	if useCommand != 0 {
		commandPtr = commandCallbackPointer
	}
	return nativePowerShell_CreateRunspace(makeUintptrFromUint64(context), commandPtr, loggerPtr)
}

func createRemoteRunspaceHelper(context uint64, useLogger byte, remoteMachine *uint16, userName *uint16, password *uint16) nativePowerShell_RunspaceHandle {
	var loggerPtr uintptr = 0
	if useLogger != 0 {
		loggerPtr = uintptr(loggerCallbackPointer)
	}
	return nativePowerShell_CreateRemoteRunspace(makeUintptrFromUint64(context), loggerPtr, remoteMachine, userName, password)
}
