
#include "./../../native-powershell/native-powershell-bin/host.h"
#include <stddef.h>
#ifdef __cplusplus
extern "C" {
#endif


wchar_t GetChar(wchar_t *t, int offset);

    unsigned char* MallocWrapper(unsigned long long size);
    void FreeWrapper(void *);

void InitLibraryHelper();

NativePowerShell_GenericPowerShellObject * MallocCopyGenericPowerShellObject(NativePowerShell_GenericPowerShellObject* input, unsigned long long inputCount);
const wchar_t* MallocCopy(const wchar_t* str);

    void Logger(const wchar_t* s);

NativePowerShell_RunspaceHandle CreateRunspaceHelper(unsigned long long, char useLogger, char useCommand);
NativePowerShell_RunspaceHandle CreateRemoteRunspaceHelper(unsigned long long context, char useLogger, const wchar_t * remoteMachine, const wchar_t * userName, const wchar_t * password  );

void SetGenericPowershellString(NativePowerShell_GenericPowerShellObject* object, wchar_t *value,char autoRelease);
void SetGenericPowerShellHandle(NativePowerShell_GenericPowerShellObject* object, NativePowerShell_PowerShellObject handle,char autoRelease);

#ifdef __cplusplus
}
#endif