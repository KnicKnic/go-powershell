package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ${SRCDIR}/../../native-powershell/native-powershell-bin/psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import (
	"unsafe"

	"golang.org/x/sys/windows"
	"github.com/KnicKnic/go-powershell/pkg/logger"
)

func init() {
	C.InitLibraryHelper()
}

type runspaceContext struct {
	log      logger.Full
	callback CallbackHolder
	invoking []psCommand // in order list of psCommands we are currently invoking

	// runspaceContext should contain all the datamembers to reconstrut runspace
	handle        C.NativePowerShell_RunspaceHandle
	contextLookup uint64
}

// recreateRunspace will give you a runspace from it's context
func (context *runspaceContext) recreateRunspace() Runspace {
	return Runspace{context}
}

// Runspace  corresponds to a powershell runspace.
//
// Use this object to execute all your powershell commands/scripts, see ExecScript and ExecCommand
//
// use .Close() to free
type Runspace struct {
	*runspaceContext
}

// CreateRunspaceSimple creates a runspace in which to run powershell commands
//
// This function has no callback routines or logging callbacks.
//
// You must call Close when done with this object
func CreateRunspaceSimple() Runspace {
	return CreateRunspace(nil, nil)
}

// CreateRunspace creates a runspace in which to run powershell commands
//
// This function allows you to specify a logging callback as well as a callback routine that processes
// commands from powershell
//
// For more details see logger.Simple and Callback holder types.
//
// You must call Close when done with this object
func CreateRunspace(loggerCallback logger.Simple, callback CallbackHolder) Runspace {
	context := &runspaceContext{log: logger.MakeLoggerFull(loggerCallback),
		callback: callback,
	}
	context.contextLookup = storeRunspaceContext(context)

	var useLogger C.char = 1
	if loggerCallback == nil {
		useLogger = 0
	}
	var useCommand C.char = 1
	if callback == nil {
		useCommand = 0
	}
	context.handle = C.CreateRunspaceHelper(C.ulonglong(context.contextLookup), useLogger, useCommand)
	return context.recreateRunspace()
}


// CreateRemoteRunspace creates a runspace in which to run powershell commands
//
// This function allows you to specify a logging callback
//
// For more details see logger.Simple.
//
// Specify "" for username to not send username and password
//
// You must call Close when done with this object
func CreateRemoteRunspace(loggerCallback logger.Simple, remoteMachine string, username string, password string) Runspace {
	context := &runspaceContext{log: logger.MakeLoggerFull(loggerCallback),
		callback: nil,
	}
	context.contextLookup = storeRunspaceContext(context)

	var useLogger C.char = 1
	if loggerCallback == nil {
		useLogger = 0
	}


	cRemoteMachine, _ := windows.UTF16PtrFromString(remoteMachine)
	ptrRemoteMachine := unsafe.Pointer(cRemoteMachine)

	cUsername, _ := windows.UTF16PtrFromString(username)
	ptrUsername := unsafe.Pointer(cUsername)

	cPassword, _ := windows.UTF16PtrFromString(password)
	ptrPassword := unsafe.Pointer(cPassword)

	context.handle = C.CreateRemoteRunspaceHelper(C.ulonglong(context.contextLookup), useLogger, (*C.wchar_t)(ptrRemoteMachine), (*C.wchar_t)(ptrUsername), (*C.wchar_t)(ptrPassword))
	return context.recreateRunspace()
}

// Close and free a Runspace
func (runspace Runspace) Close() {
	deleteRunspaceContextLookup(runspace.contextLookup)
	C.NativePowerShell_DeleteRunspace(runspace.handle)
}
