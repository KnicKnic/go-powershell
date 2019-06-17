package main

/*

#cgo CFLAGS: -IC:/code/go-net
#cgo LDFLAGS: C:/code/go-net/mscoree.dll
#include <stddef.h>
#include "main.h"
#include <stdio.h>
#include <stdlib.h>
// #include <metahost.h>
// #pragma comment(lib, "mscoree.lib")

// ICLRMetaHost       *pMetaHost       = NULL;
// ICLRMetaHostPolicy *pMetaHostPolicy = NULL;
// ICLRDebugging      *pCLRDebugging   = NULL;

void myprint(char* s) {
// HRESULT hr = CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost,
//                     (LPVOID*)&pMetaHost);

    wchar_t* arr[] = {(wchar_t*)L"",(wchar_t*)L"c:\\someassembly.dll"};
    int ret = oldmain(2,arr);
	printf("\n%s, old_main %d\n", s,ret);


}



*/
import "C"

import "unsafe"

func Example() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}
func main() {
	Example()
}
