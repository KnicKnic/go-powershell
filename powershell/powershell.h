
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
    void FreeWrapper(void *);

void InitLibraryHelper();

GenericPowershellObject * MallocCopyGenericPowershellObject(GenericPowershellObject* input, unsigned long long inputCount);
const wchar_t* MallocCopy(const wchar_t* str);

    void Logger(const wchar_t* s);

RunspaceHandle CreateRunspaceHelper(unsigned long long);

void SetGenericPowershellString(GenericPowershellObject* object, wchar_t *value,char autoRelease);
void SetGenericPowershellHandle(GenericPowershellObject* object, unsigned long long handle,char autoRelease);

void CClosePowerShellObject(unsigned long long handle);
unsigned long long ClonePowerShellObject(unsigned long long handle);
#ifdef __cplusplus
}
#endif