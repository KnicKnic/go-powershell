package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"

type PowershellObject struct {
	handle uint64
}

func (obj PowershellObject) Close() {
	C.CClosePowerShellObject(C.ulonglong(obj.handle))
}

func (obj PowershellObject) Clone() PowershellObject {
	var handle C.ulonglong = C.ClonePowerShellObject(C.ulonglong(obj.handle))
	return PowershellObject{uint64(handle)}
}

// StringPtr GetPSObjectType(PowerShellObject handle);
// StringPtr GetPSObjectToString(PowerShellObject handle);
// char IsPSObjectNullptr(PowerShellObject handle);
