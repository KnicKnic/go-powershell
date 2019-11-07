package powershell

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func makeUint64FromPtr(v uintptr) uint64 {
	return *((*uint64)(unsafe.Pointer(&v)))
}
func makeUintptrFromUint64(v uint64) uintptr {
	return *((*uintptr)(unsafe.Pointer(&v)))
}

func allocWrapper(size uint64) uintptr {
	return nativePowerShell_DefaultAlloc(size)
}
func freeWrapper(v uintptr) {
	nativePowerShell_DefaultFree(v)
}

func mallocCopy(input uintptr, size uintptr) uintptr {

	u64Size := makeUint64FromPtr(size)
	data := allocWrapper(u64Size)
	_ = memcpy(data, uintptr(unsafe.Pointer(input)), u64Size)

	return data
}

func wsclen(str uintptr) uint64 {
	var charCode uint16 = 1
	var i uint64 = 0
	for ; charCode != 0; i++ {
		charCode = *(*uint16)(unsafe.Pointer(str + (makeUintptrFromUint64(i) * unsafe.Sizeof(charCode))))
	}
	return i
}

func makeString(str uintptr) string {
	count := wsclen(str) + 1
	arr := make([]uint16, count)
	ptrwchar := unsafe.Pointer(&arr[0])

	memcpy(uintptr(ptrwchar), str, count*2)

	s := windows.UTF16ToString(arr)
	return s
}

func uintptrMakeString(ptr uintptr) string {
	return makeString(ptr)
}

func makeCStringUintptr(str string) uintptr {
	cs, _ := windows.UTF16PtrFromString(str)
	ptrwchar := unsafe.Pointer(cs)
	size := 2 * (wsclen(uintptr(ptrwchar)) + 1)

	return mallocCopy(uintptr(ptrwchar), makeUintptrFromUint64(size))
}
