package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -IC:/code/psh_host
#cgo LDFLAGS: c:/code/psh_host/x64/Release/psh_host.dll
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
    long ret = startpowershell(s);
	printf("\n%ws, old_main %d\n", s,ret);


}



*/
import "C"

func Example() {
	// cs := C.CString("Hello from stdio\n")
	cs, _ := windows.UTF16PtrFromString("c:\\fuzzy2")

	ptrwchar := unsafe.Pointer(cs) // that is what you will get from Windows
	// cs := C CString("Hello from stdio\n")
	C.myprint(ptrwchar)
	// C.free(unsafe.Pointer(cs))
}
func main() {
	Example()
}
