
#include "host.h"
#include <stddef.h>
#ifdef __cplusplus
extern "C" {
#endif
void myprint(void* unknown) ;

wchar_t * MakeWchar(void *unknown);

void MemoryCopy(void * dest, wchar_t * src, int size);

wchar_t MakeNullTerminator();

wchar_t GetChar(wchar_t *t, int offset);

    unsigned char* MallocWrapper(unsigned long long size);

void InitLibraryHelper();

const wchar_t* MallocCopy(const wchar_t* str);

    void Logger(const wchar_t* s);

RunspaceHandle CreateRunspaceHelper();

#ifdef __cplusplus
}
#endif