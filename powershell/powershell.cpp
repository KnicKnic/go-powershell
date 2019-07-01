#include <stddef.h>
#include "powershell.h"
#include "host.h"

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
    InitLibrary(MallocWrapper, free);
}
GenericPowershellObject * MallocCopyGenericPowershellObject(GenericPowershellObject* input, unsigned long long inputCount){
    
    GenericPowershellObject* dest = (GenericPowershellObject*)MallocWrapper(inputCount * sizeof(input[0]));
    memcpy(dest, (GenericPowershellObject *)input, inputCount*sizeof(input[0]));
    return dest;
}

const wchar_t* MallocCopy(const wchar_t* str)
{
    if (str == NULL)
        return NULL;

    size_t s = 0;
    for (; str[s] != '\0'; ++s) {
    }
    ++s;
    wchar_t* dest = (wchar_t*)MallocWrapper(s * sizeof(str[0]));
    memcpy(dest, (wchar_t *)str, s*sizeof(str[0]));
    return (const wchar_t*)dest;
}

    void Logger(void *context, const wchar_t* s)
    {
        logWchart((unsigned long long )context, (wchar_t *)s);
        //printf("My Member Logger: %ws\n", s);
    }
    void Command(void *context, const wchar_t* s, PowerShellObject* input, unsigned long long inputCount, JsonReturnValues * returnValues)
    {        
        commandWchart((unsigned long long) context, (wchar_t *)s, input, inputCount, returnValues);
    }

RunspaceHandle CreateRunspaceHelper(unsigned long long context){
    return CreateRunspace((void*)context, Command, Logger);
    // return CreateRunspace(nullptr, Command, Logger);
}


void SetGenericPowershellString(GenericPowershellObject* object, wchar_t *value, char autoRelease){
    object->type = PowershellObjectTypeString;
    object->instance.string = value;
    object->releaseObject = autoRelease;
}

void SetGenericPowershellHandle(GenericPowershellObject* object, PowerShellObject value,char autoRelease){
    object->type = PowershellObjectHandle;
    object->instance.psObject = value;
    object->releaseObject = autoRelease;
}