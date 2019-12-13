package powershell

import (
	"unsafe"
)

func makeUint64FromPtr(v uintptr) uint64 {
	return *((*uint64)(unsafe.Pointer(&v)))
}
func makeUintptrFromUint64(v uint64) uintptr {
	return *((*uintptr)(unsafe.Pointer(&v)))
}

func allocWrapper(size uint64) (uintptr, error) {
	return nativePowerShell_DefaultAlloc(size)
}
func freeWrapper(v uintptr) {
	nativePowerShell_DefaultFree(v)
}

func mallocCopyLogStringHolder(input nativePowerShell_LogString_Holder) uintptr {

	size := uint64(unsafe.Sizeof(input))

	data, err := allocWrapper(size)
	if err != nil {
		panic("Couldn't allocate memory")
	}

	_ = memcpyLogStringHolder(data, input)

	return data
}

func mallocCopyGenericPowerShellObject(input *nativePowerShell_GenericPowerShellObject, inputCount uint64) uintptr {

	size := inputCount * uint64(unsafe.Sizeof(*input))

	data, err := allocWrapper(size)
	if err != nil {
		panic("Couldn't allocate memory")
	}

	_ = memcpyGenericPowerShellObject(data, input, size)

	return data
}

func mallocCopyStr(str string) uintptr {

	size := 2 * uint64((len(str) + 1))
	data, err := allocWrapper(size)
	if err != nil {
		panic("Couldn't allocate memory")
	}
	// safe usage due to data being c pointer
	_ = memcpyStr(data, str)

	return data
}

func wsclen(str uintptr) uint64 {
	var charCode uint16 = 1
	var i uint64 = 0
	for ; charCode != 0; i++ {
		charCode = *((*uint16)(unsafe.Pointer(str + (makeUintptrFromUint64(i) * unsafe.Sizeof(charCode)))))
	}
	return i
}

func uintptrMakeString(ptr uintptr) string {
	return cstrToStr(ptr)
}

func makeCStringUintptr(str string) uintptr {
	return mallocCopyStr(str)
}
