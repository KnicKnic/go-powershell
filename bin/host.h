#pragma once

/*

This header file for native-powershell (psh_host.dll)

*/


#ifdef __cplusplus
extern "C" {
#endif
    
    #define NativePowerShell_InvalidPointer (void *)0;
    typedef void (*NativePowerShell_FreePointer)(void*);
    typedef unsigned char* (*NativePowerShell_AllocPointer)(unsigned long long size);
    
    typedef void (*NativePowerShell_LogString)(void * context, const wchar_t* messages);
    
    // If Log is Null none of the logging works
    typedef struct NativePowerShell_LogString_Holder_ {
        NativePowerShell_LogString Log;
        NativePowerShell_LogString LogError;
        NativePowerShell_LogString LogWarning;
        NativePowerShell_LogString LogInformation;
        NativePowerShell_LogString LogVerbose;
        NativePowerShell_LogString LogDebug;
        NativePowerShell_LogString LogLine;
        NativePowerShell_LogString LogErrorLine;
        NativePowerShell_LogString LogWarningLine;
        NativePowerShell_LogString LogInformationLine;
        NativePowerShell_LogString LogVerboseLine;
        NativePowerShell_LogString LogDebugLine;
    }NativePowerShell_LogString_Holder, * PNativePowerShell_LogString_Holder;

    void NativePowerShell_InitLibrary( NativePowerShell_AllocPointer, NativePowerShell_FreePointer);
    
    unsigned char* NativePowerShell_DefaultAlloc(unsigned long long size);
    void NativePowerShell_DefaultFree(void* address);

	//typedef struct NativePowerShell_RunspaceHandle_ {} *NativePowerShell_RunspaceHandle;
 //   typedef struct NativePowerShell_PowerShellHandle_ {} *NativePowerShell_PowerShellHandle;
 //   typedef struct NativePowerShell_PowerShellObject_ {} *NativePowerShell_PowerShellObject;
#define NativePowerShell_InvalidHandleValue 0
    typedef unsigned long long NativePowerShell_RunspaceHandle;
    typedef unsigned long long NativePowerShell_PowerShellHandle;
    typedef unsigned long long NativePowerShell_PowerShellObject;
	//typedef NativePowerShell_RunspaceHandle_d * NativePowerShell_RunspaceHandle;
       
    typedef const wchar_t* NativePowerShell_StringPtr;



    typedef enum NativePowerShell_PowerShellObjectType_ { NativePowerShell_PowerShellObjectTypeString, NativePowerShell_PowerShellObjectHandle }NativePowerShell_PowerShellObjectType;

    typedef struct NativePowerShell_GenericPowerShellObject_ {
        NativePowerShell_PowerShellObjectType type;
        union NativePowerShell_PowerShellObjectInstance {
            NativePowerShell_StringPtr string;
            NativePowerShell_PowerShellObject psObject;
            // continue for other ones such as UInt64...
        } instance;
        char releaseObject; // if true reciever of this object will release the instance
    }NativePowerShell_GenericPowerShellObject, * PNativePowerShell_GenericPowerShellObject;


    typedef struct NativePowerShell_JsonReturnValues_ {
        NativePowerShell_GenericPowerShellObject* objects;
        unsigned long           count;
    }NativePowerShell_JsonReturnValues,*NativePowerShell_PJsonReturnValues;
    typedef void (*NativePowerShell_ReceiveJsonCommand)(void* context, const wchar_t* command, NativePowerShell_PowerShellObject * inputs, unsigned long long inputCount, NativePowerShell_JsonReturnValues* returnValues);

	NativePowerShell_PowerShellHandle NativePowerShell_CreatePowerShell(NativePowerShell_RunspaceHandle handle);
    NativePowerShell_PowerShellHandle NativePowerShell_CreatePowerShellNested(NativePowerShell_PowerShellHandle handle);


	void NativePowerShell_DeletePowershell(NativePowerShell_PowerShellHandle handle);


	NativePowerShell_RunspaceHandle NativePowerShell_CreateRunspace(void* context, NativePowerShell_ReceiveJsonCommand, PNativePowerShell_LogString_Holder);
    NativePowerShell_RunspaceHandle NativePowerShell_CreateRemoteRunspace(void* context, PNativePowerShell_LogString_Holder, const wchar_t* computerName, const wchar_t* username, const wchar_t * password);

	void NativePowerShell_DeleteRunspace(NativePowerShell_RunspaceHandle handle);


    long NativePowerShell_AddCommand(NativePowerShell_PowerShellHandle handle, NativePowerShell_StringPtr command);
    long NativePowerShell_AddCommandSpecifyScope(NativePowerShell_PowerShellHandle handle, NativePowerShell_StringPtr command, char useLocalScope);
    long NativePowerShell_AddParameterString(NativePowerShell_PowerShellHandle handle, NativePowerShell_StringPtr name, NativePowerShell_StringPtr value);
    long NativePowerShell_AddParameterObject(NativePowerShell_PowerShellHandle handle, NativePowerShell_StringPtr name, NativePowerShell_PowerShellObject object);
	long NativePowerShell_AddArgument(NativePowerShell_PowerShellHandle handle, NativePowerShell_StringPtr argument);
	long NativePowerShell_AddPSObjectArgument(NativePowerShell_PowerShellHandle handle, NativePowerShell_PowerShellObject object);
	long NativePowerShell_AddPSObjectArguments(NativePowerShell_PowerShellHandle handle, NativePowerShell_PowerShellObject* objects, unsigned int count);

    // caller is responsible for calling ClosePowerShellObject on all returned objects, as well as
    // calling the appropriate free routine on objects assuming it is not nullptr
    NativePowerShell_PowerShellObject NativePowerShell_InvokeCommand(NativePowerShell_PowerShellHandle handle, NativePowerShell_PowerShellObject** objects, unsigned int* objectCount);
    long NativePowerShell_AddScript(NativePowerShell_PowerShellHandle handle, NativePowerShell_StringPtr path);
    long NativePowerShell_AddScriptSpecifyScope(NativePowerShell_PowerShellHandle handle, NativePowerShell_StringPtr path, char useLocalScope);
    void NativePowerShell_ClosePowerShellObject(NativePowerShell_PowerShellObject psobject);

    NativePowerShell_StringPtr NativePowerShell_GetPSObjectType(NativePowerShell_PowerShellObject handle);
    NativePowerShell_StringPtr NativePowerShell_GetPSObjectToString(NativePowerShell_PowerShellObject handle);
    char NativePowerShell_IsPSObjectNullptr(NativePowerShell_PowerShellObject handle);
    NativePowerShell_PowerShellObject NativePowerShell_AddPSObjectHandle(NativePowerShell_PowerShellObject handle);

#ifdef __cplusplus
}
#endif