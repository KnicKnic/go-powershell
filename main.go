package main

import "fmt"

/*

#cgo CFLAGS: -IC:/code/go-net
#cgo LDFLAGS: C:/code/go-net/mscoree.dll
#include <stdio.h>
#include <stdlib.h>
#include <metahost.h>
#pragma comment(lib, "mscoree.lib")

ICLRMetaHost       *pMetaHost       = NULL;
ICLRMetaHostPolicy *pMetaHostPolicy = NULL;
ICLRDebugging      *pCLRDebugging   = NULL;

void myprint(char* s) {
HRESULT hr = CLRCreateInstance(&CLSID_CLRMetaHost, &IID_ICLRMetaHost,
                    (LPVOID*)&pMetaHost);
	printf("%s\n", s);


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
	fmt.Println("hello world")
}
