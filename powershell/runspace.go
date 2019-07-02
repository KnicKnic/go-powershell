package powershell

// "bitbucket.org/creachadair/shell"

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"

func init() {
	C.InitLibraryHelper()
}

type runspaceContext struct {
	Log      LoggerFull
	Callback CallbackHolder
}

// Context handle for a runspace, use .Delete() to free
type Runspace struct {
	handle        C.RunspaceHandle
	context       runspaceContext
	contextLookup uint64
}

type emptyLogger struct {
}

func (emptyLogger) Write(args ...interface{}) {
}

type emptyCallback struct {
}

func (emptyCallback) Callback(string, []PowershellObject, CallbackResultsWriter) {
}

// CreateRunspace think of this kinda like a shell
//
// You must call Delete when done with this object
func CreateRunspaceSimple() Runspace {
	return CreateRunspace(emptyLogger{}, emptyCallback{})
}

// CreateRunspace think of this kinda like a shell
//
// You must call Delete when done with this object
func CreateRunspace(logger LoggerSimple, callback CallbackHolder) Runspace {
	context := runspaceContext{makeLoggerFull(logger), callback}
	contextLookup := storeRunspaceContext(context)

	runspace := C.CreateRunspaceHelper(C.ulonglong(contextLookup))
	return Runspace{runspace, context, contextLookup}
}

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	deleteRunspaceContextLookup(runspace.contextLookup)
	C.DeleteRunspace(runspace.handle)
}
