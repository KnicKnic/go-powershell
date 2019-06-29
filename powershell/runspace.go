package powershell

// "bitbucket.org/creachadair/shell"

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import "unsafe"

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
	handle     C.RunspaceHandle
	context    *Context
	contextRef unsafe.Pointer
}

// CreateRunspace think of this kinda like a shell
func CreateRunspace() Runspace {
	context := &Context{MakeLogHolder(GLogInfoLogger{}), callbackTest{}}
	unsafeContext := unsafe.Pointer(context)
	runspace := C.CreateRunspaceHelper(unsafeContext)
	return Runspace{runspace, context, unsafeContext}
}

// CreateRunspace think of this kinda like a shell
// func CreateRunspace2(simpleOrFull LoggerSimple) Runspace {
// 	context := Context{MakeLogHolder(LoggerSimple), CallbackHolder{}}
// 	return Runspace{C.CreateRunspaceHelper(), context}
// }

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	C.DeleteRunspace(runspace.handle)
}
