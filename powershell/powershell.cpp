#include <stddef.h>
#include "powershell.h"
#include "host.h"

#include <stdio.h>
#include <stdlib.h>
#include "sys/types.h"
#include "_cgo_export.h"
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

void MemoryCopy(void * dest, wchar_t * src, int size){
    for (int i=0; i<size/2;++i) {
        ((wchar_t*)dest)[i] = ((wchar_t*)src)[i];
    }
    // memcpy(dest, src,size);
}

wchar_t MakeNullTerminator(){
    return L'\0';
}

wchar_t GetChar(wchar_t *t, int offset){
    return t[offset];
}

    unsigned char* MallocWrapper(unsigned long long size) {
        return (unsigned char*)malloc(size);
    }
    void FreeWrapper(void *ptr){
        return free(ptr);
    }

void InitLibraryHelper(){
    InitLibrary(MallocWrapper, free);
}

const wchar_t* MallocCopy(const wchar_t* str)
{
    if (str == NULL)
        return NULL;

    size_t s = 0;
    for (; str[s] != '\0'; ++s) {
    }
    ++s;
    wchar_t* dest = (wchar_t*)malloc(s * sizeof(str[0]));
     MemoryCopy(dest, (wchar_t *)str, s*2);
    return (const wchar_t*)dest;
}

    void Logger(void *, const wchar_t* s)
    {
        logWchart((wchar_t *)s);
        //printf("My Member Logger: %ws\n", s);
    }
    const wchar_t* Command(void *, const wchar_t* s)
    {
        printf("My Member Logger: %ws\n", s);
        return MallocCopy(s);
    }

RunspaceHandle CreateRunspaceHelper(){
    return CreateRunspace(nullptr, Command, Logger);
}
