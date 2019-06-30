package powershell

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"
import "unsafe"

type PowershellObject struct {
	handle uint64
}

func (obj PowershellObject) toCHandle() C.PowerShellObject {
	return *((*C.PowerShellObject)(unsafe.Pointer(&obj.handle)))
}

func (obj PowershellObject) Close() {
	C.ClosePowerShellObject(obj.toCHandle())
}

func (obj PowershellObject) AddRef() PowershellObject {
	handle := C.AddPSObjectHandle(obj.toCHandle())
	return makePowerShellObject(handle)
}

func (obj PowershellObject) IsNull() bool {
	return C.IsPSObjectNullptr(obj.toCHandle()) == C.char(1)
}

func (obj PowershellObject) Type() string {
	if obj.IsNull() {
		return "nullptr"
	}

	var str *C.wchar_t = C.GetPSObjectType(obj.toCHandle())
	defer C.FreeWrapper(unsafe.Pointer(str))
	return makeString(str)
}

func (obj PowershellObject) ToString() string {
	if obj.IsNull() {
		return "nullptr"
	}

	var str *C.wchar_t = C.GetPSObjectToString(obj.toCHandle())
	defer C.FreeWrapper(unsafe.Pointer(str))
	return makeString(str)
}
