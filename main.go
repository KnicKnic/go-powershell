package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll
#include <stddef.h>
#include "host.h"
#include <stdio.h>
#include <stdlib.h>
#include "sys/types.h"
// #include <metahost.h>
// #pragma comment(lib, "mscoree.lib")

// ICLRMetaHost       *pMetaHost       = NULL;
// ICLRMetaHostPolicy *pMetaHostPolicy = NULL;
// ICLRDebugging      *pCLRDebugging   = NULL;

void myprint(void* unknown) {
    wchar_t* s = (wchar_t*)unknown;
// HRESULT hr = CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost,
//                     (LPVOID*)&pMetaHost);
	printf("\n%ws, old_main %d\n", s);
}

wchar_t * MakeWchar(void *unknown){
	wchar_t* s = (wchar_t*)unknown;;
	return s;
}

long StartPowershell(RunspaceHandle runspace, void * unknown){
	wchar_t* s = (wchar_t*)unknown;
	return startpowershell(runspace, s);
}



*/
import "C"

func Example() {
	// cs := C.CString("Hello from stdio\n")
	handle := C.CreateRunspace()
	defer C.DeleteRunspace(handle)
	cs, _ := windows.UTF16PtrFromString("c:\\fuzzy3")

	ptrwchar := unsafe.Pointer(cs) // that is what you will get from Windows
	// cs := C CString("Hello from stdio\n")
	// C.myprint(ptrwchar)

	_ = C.startpowershell(handle, C.MakeWchar(ptrwchar))
	// C.free(unsafe.Pointer(cs))
}
func main() {
	Example()
}
