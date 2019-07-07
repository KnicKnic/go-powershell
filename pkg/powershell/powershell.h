
#include "./../../native-powershell/native-powershell-bin/host.h"
#include <stddef.h>
#ifdef __cplusplus
extern "C" {
#endif


wchar_t GetChar(wchar_t *t, int offset);

    unsigned char* MallocWrapper(unsigned long long size);
    void FreeWrapper(void *);

void InitLibraryHelper();

GenericPowershellObject * MallocCopyGenericPowershellObject(GenericPowershellObject* input, unsigned long long inputCount);
const wchar_t* MallocCopy(const wchar_t* str);

    void Logger(const wchar_t* s);

RunspaceHandle CreateRunspaceHelper(unsigned long long, char useLogger, char useCommand);

void SetGenericPowershellString(GenericPowershellObject* object, wchar_t *value,char autoRelease);
void SetGenericPowershellHandle(GenericPowershellObject* object, PowerShellObject handle,char autoRelease);

#ifdef __cplusplus
}
#endif