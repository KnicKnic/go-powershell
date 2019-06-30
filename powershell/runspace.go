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

type Context struct {
	Log      LogHolder
	Callback CallbackHolder
}

type Runspace struct {
	handle        C.RunspaceHandle
	context       Context
	contextLookup uint64
}

// CreateRunspace think of this kinda like a shell
func CreateRunspace(logger LoggerSimple, callback CallbackHolder) Runspace {
	context := Context{MakeLogHolder(logger), callback}
	contextLookup := StoreRunspaceContext(context)

	runspace := C.CreateRunspaceHelper(C.ulonglong(contextLookup))
	return Runspace{runspace, context, contextLookup}
}

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	DeleteRunspaceContextLookup(runspace.contextLookup)
	C.DeleteRunspace(runspace.handle)
}
