package powershell

// "bitbucket.org/creachadair/shell"

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import (
	logger "github.com/knicknic/go-powershell/powershell/logger"
)

func init() {
	C.InitLibraryHelper()
}

type runspaceContext struct {
	Log      logger.LoggerFull
	Callback CallbackHolder
}

// Context handle for a runspace, use .Delete() to free
type Runspace struct {
	handle        C.RunspaceHandle
	context       runspaceContext
	contextLookup uint64
}

// CreateRunspace think of this kinda like a shell
//
// You must call Delete when done with this object
func CreateRunspaceSimple() Runspace {
	return CreateRunspace(nil, nil)
}

// CreateRunspace think of this kinda like a shell
//
// You must call Delete when done with this object
func CreateRunspace(loggerCallback logger.LoggerSimple, callback CallbackHolder) Runspace {
	context := runspaceContext{logger.MakeLoggerFull(loggerCallback), callback}
	contextLookup := storeRunspaceContext(context)

	var useLogger C.char = 1
	if loggerCallback == nil {
		useLogger = 0
	}
	var useCommand C.char = 1
	if callback == nil {
		useCommand = 0
	}
	runspace := C.CreateRunspaceHelper(C.ulonglong(contextLookup), useLogger, useCommand)
	return Runspace{runspace, context, contextLookup}
}

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	deleteRunspaceContextLookup(runspace.contextLookup)
	C.DeleteRunspace(runspace.handle)
}
