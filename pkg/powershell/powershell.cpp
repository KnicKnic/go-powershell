#include <stddef.h>
#include "powershell.h"
#include "./../../native-powershell/native-powershell-bin/host.h"

#include <stdio.h>
#include <stdlib.h>
#include "sys/types.h"
#include "_cgo_export.h"

#include <string>


    unsigned char *MallocWrapper(unsigned long long size)
    {
        auto ptr = (unsigned char *)malloc(size);
        if (ptr == nullptr)
        {
            throw "memory alloc returned nullptr";
        }
        return ptr;
    }
void FreeWrapper(void *ptr)
{
    return free(ptr);
    }

void InitLibraryHelper(){
    NativePowerShell_InitLibrary(MallocWrapper, free);
}

void * MallocCopyGeneric(const void * input, unsigned long long byteCount ){
    if(input == nullptr){
        return nullptr;
    }
    void* dest = MallocWrapper(byteCount);
    memcpy(dest, input, byteCount);
    return dest;
}

NativePowerShell_GenericPowerShellObject * MallocCopyGenericPowerShellObject(NativePowerShell_GenericPowerShellObject* input, unsigned long long inputCount){
    return (NativePowerShell_GenericPowerShellObject *)MallocCopyGeneric(input, inputCount * sizeof(input[0]));
}

const wchar_t* MallocCopy(const wchar_t* str)
{
    if (str == NULL)
        return NULL;

    size_t s = wcslen(str) + 1;
    return (const wchar_t *)MallocCopyGeneric(str, s * sizeof(str[0]));
}

    void Logger(void *context, const wchar_t* s)
    {
        logWchart((unsigned long long )context, (wchar_t *)s);
        //printf("My Member Logger: %ws\n", s);
    }
    void Command(void *context, const wchar_t* s, NativePowerShell_PowerShellObject* input, unsigned long long inputCount, NativePowerShell_JsonReturnValues * returnValues)
    {        
        commandWchart((unsigned long long) context, (wchar_t *)s, input, inputCount, returnValues);
    }

NativePowerShell_RunspaceHandle CreateRunspaceHelper(unsigned long long context, char useLogger, char useCommand ){
    NativePowerShell_LogString loggerPtr = Logger;
    NativePowerShell_ReceiveJsonCommand commandPtr = Command;
    if(useLogger == 0){
        loggerPtr = nullptr;
    }
    if(useCommand == 0){
        commandPtr = nullptr;
    }
    return NativePowerShell_CreateRunspace((void*)context, commandPtr, loggerPtr);
}

NativePowerShell_RunspaceHandle CreateRemoteRunspaceHelper(unsigned long long context, char useLogger, const wchar_t * remoteMachine, const wchar_t * userName, const wchar_t * password  ){
    NativePowerShell_LogString loggerPtr = Logger;
    if(useLogger == 0){
        loggerPtr = nullptr;
    }
    return NativePowerShell_CreateRemoteRunspace((void*)context,  loggerPtr, remoteMachine, userName, password);
}


void SetGenericPowershellString(NativePowerShell_GenericPowerShellObject* object, wchar_t *value, char autoRelease){
    object->type = NativePowerShell_PowerShellObjectTypeString;
    object->instance.string = value;
    object->releaseObject = autoRelease;
}

void SetGenericPowerShellHandle(NativePowerShell_GenericPowerShellObject* object, NativePowerShell_PowerShellObject value,char autoRelease){
    object->type = NativePowerShell_PowerShellObjectHandle;
    object->instance.psObject = value;
    object->releaseObject = autoRelease;
}