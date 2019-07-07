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
}

// Runspace a context handle for a runspace, use .Close() to free
type Runspace struct {
	handle        C.RunspaceHandle
	context       runspaceContext
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

// Close and free a Runspace
func (runspace Runspace) Close() {
	deleteRunspaceContextLookup(runspace.contextLookup)
	C.DeleteRunspace(runspace.handle)
}
