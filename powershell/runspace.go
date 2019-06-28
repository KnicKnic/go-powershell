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

type Runspace struct {
	handle C.RunspaceHandle
}

// CreateRunspace think of this kinda like a shell
func CreateRunspace() Runspace {
	return Runspace{C.CreateRunspaceHelper()}
}

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	C.DeleteRunspace(runspace.handle)
}
