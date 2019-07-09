package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ${SRCDIR}/../../native-powershell/native-powershell-bin/psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import (
	"github.com/KnicKnic/go-powershell/pkg/logger"
)

func init() {
	C.InitLibraryHelper()
}

type runspaceContext struct {
	Log      logger.Full
	Callback CallbackHolder
	invoking []psCommand // in order list of psCommands we are currently invoking

	// runspaceContext should contain all the datamembers to reconstrut runspace
	handle C.RunspaceHandle
	contextLookup uint64
}

// recreateRunspace will give you a runspace from it's context
func (context *runspaceContext) recreateRunspace() Runspace{
	return Runspace{context.handle, context, context.contextLookup}
}

// Runspace a context handle for a runspace, use .Close() to free
type Runspace struct {
	handle        C.RunspaceHandle
	context       *runspaceContext
	contextLookup uint64
}

// CreateRunspaceSimple think of this kinda like creating a shell
//
// You must call Close when done with this object
func CreateRunspaceSimple() Runspace {
	return CreateRunspace(nil, nil)
}

// CreateRunspace think of this kinda like creating a shell
//
// You must call Close when done with this object
func CreateRunspace(loggerCallback logger.Simple, callback CallbackHolder) Runspace {
	context := &runspaceContext{Log: logger.MakeLoggerFull(loggerCallback), 
								Callback: callback,
								invoking: nil,
								handle: C.ulonglong(0),
								contextLookup: 0,
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

// Close and free a Runspace
func (runspace Runspace) Close() {
	deleteRunspaceContextLookup(runspace.contextLookup)
	C.DeleteRunspace(runspace.handle)
}
