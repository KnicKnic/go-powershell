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
	"sync"
	"sync/atomic"
)

func init() {
	C.InitLibraryHelper()
}

type CallbackHolder interface {
	Callback(str string) string
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

//
// the goal around context cache is to get around golang's problem of not being able to
// marshal to cgo a golang pointer to a golang pointer.
//
var contextCache sync.Map
var contextLookupKey uint64 = 0

// CreateRunspace think of this kinda like a shell
func CreateRunspace() Runspace {
	contextLookup := atomic.AddUint64(&contextLookupKey, 1)
	context := Context{MakeLogHolder(GLogInfoLogger{}), callbackTest{}}
	contextCache.Store(contextLookup, context)

	runspace := C.CreateRunspaceHelper(C.ulonglong(contextLookup))
	return Runspace{runspace, context, contextLookup}
}

// CreateRunspace think of this kinda like a shell
// func CreateRunspace2(simpleOrFull LoggerSimple) Runspace {
// 	context := Context{MakeLogHolder(LoggerSimple), CallbackHolder{}}
// 	return Runspace{C.CreateRunspaceHelper(), context}
// }

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	contextCache.Delete(runspace.contextLookup)
	C.DeleteRunspace(runspace.handle)
}
