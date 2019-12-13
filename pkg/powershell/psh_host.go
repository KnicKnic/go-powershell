package powershell

//go:generate go run $GOROOT/src/syscall/mksyscall_windows.go -output zpsh_host.go psh_host.go

//https://stackoverflow.com/questions/35154074/how-do-i-manage-windows-user-accounts-in-go
// copy ~\go\src\golang.org\x\sys\windows\mksyscall.go .

const (
	nativePowerShell_InvalidHandleValue = 0
	nativePowerShell_InvalidPointer     = uintptr(0)
)

type (
	nativePowerShell_StringPtr        = *uint16
	nativePowerShell_LogString        func(context uintptr, messages nativePowerShell_StringPtr) uintptr
	nativePowerShell_LogString_Holder struct {
		Log                uintptr
		LogError           uintptr
		LogWarning         uintptr
		LogInformation     uintptr
		LogVerbose         uintptr
		LogDebug           uintptr
		LogLine            uintptr
		LogErrorLine       uintptr
		LogWarningLine     uintptr
		LogInformationLine uintptr
		LogVerboseLine     uintptr
		LogDebugLine       uintptr
	}
	nativePowerShell_RunspaceHandle       = uint64
	nativePowerShell_PowerShellHandle     = uint64
	nativePowerShell_PowerShellObject     = uint64
	nativePowerShell_PowerShellObjectType = uint32
)

const (
	nativePowerShell_PowerShellObjectTypeString = 0
	nativePowerShell_PowerShellObjectHandle     = 1
)

type (
	nativePowerShell_GenericPowerShellObject struct {
		typeEnum      nativePowerShell_PowerShellObjectType
		object        nativePowerShell_PowerShellObject // can also be union with string
		releaseObject byte
	}

	nativePowerShell_JsonReturnValues struct {
		objects uintptr // *nativePowerShell_GenericPowerShellObject
		count   uint32
	}
	nativePowerShell_ReceiveJsonCommand func(context uintptr, command nativePowerShell_StringPtr, inputrs *nativePowerShell_PowerShellObject, inputCount uint64, returnValues *nativePowerShell_JsonReturnValues) uintptr
)

//sys	nativePowerShell_CreatePowerShell(handle nativePowerShell_RunspaceHandle) (status nativePowerShell_PowerShellHandle)= psh_host.NativePowerShell_CreatePowerShell
//sys   nativePowerShell_CreatePowerShellNested(handle nativePowerShell_PowerShellHandle) (status nativePowerShell_PowerShellHandle)= psh_host.NativePowerShell_CreatePowerShellNested

//sys	nativePowerShell_DeletePowershell(handle nativePowerShell_PowerShellHandle) = psh_host.NativePowerShell_DeletePowershell

// to fix command was of type *nativePowerShell_ReceiveJsonCommand
// to fix holder was of type *nativePowerShell_LogString_Holder
//sys	nativePowerShell_CreateRunspace(context uintptr, command uintptr, holder uintptr) (status nativePowerShell_RunspaceHandle) = psh_host.NativePowerShell_CreateRunspace
//sys   nativePowerShell_CreateRemoteRunspace(context uintptr, holder uintptr, computerName *uint16, username *uint16, password *uint16) (status nativePowerShell_RunspaceHandle) = psh_host.NativePowerShell_CreateRemoteRunspace

//sys	nativePowerShell_DeleteRunspace(handle nativePowerShell_RunspaceHandle) = psh_host.NativePowerShell_DeleteRunspace

//sys   nativePowerShell_AddCommand(handle nativePowerShell_PowerShellHandle, command *uint16) (status uint) = psh_host.NativePowerShell_AddCommand
//sys   nativePowerShell_AddCommandSpecifyScope( handle nativePowerShell_PowerShellHandle, command *uint16, useLocalScope byte) (status uint) = psh_host.NativePowerShell_AddCommandSpecifyScope
//sys   nativePowerShell_AddParameterString( handle nativePowerShell_PowerShellHandle, name *uint16, value *uint16) (status uint) = psh_host.NativePowerShell_AddParameterString
//sys   nativePowerShell_AddParameterObject( handle nativePowerShell_PowerShellHandle, name *uint16, object nativePowerShell_PowerShellObject) (status uint) = psh_host.NativePowerShell_AddParameterObject
//sys	nativePowerShell_AddArgument( handle nativePowerShell_PowerShellHandle, argument *uint16) (status uint) = psh_host.NativePowerShell_AddArgument
//sys	nativePowerShell_AddPSObjectArgument( handle nativePowerShell_PowerShellHandle, object nativePowerShell_PowerShellObject) (status uint) = psh_host.NativePowerShell_AddPSObjectArgument
//sys	nativePowerShell_AddPSObjectArguments( handle nativePowerShell_PowerShellHandle,  objects *nativePowerShell_PowerShellObject, count uint) (status uint) = psh_host.NativePowerShell_AddPSObjectArguments

// caller is responsible for calling ClosePowerShellObject on all returned objects, as well as
// calling the appropriate free routine on objects assuming it is not nullptr

// to fix objects was really of type nativePowerShell_PowerShellObject **

//sys	nativePowerShell_InvokeCommand(handle nativePowerShell_PowerShellHandle,  objects *uintptr, objectCount *uint) (status nativePowerShell_PowerShellObject) = psh_host.NativePowerShell_InvokeCommand
//sys   nativePowerShell_AddScript(handle nativePowerShell_PowerShellHandle , path *uint16) (status int) = psh_host.NativePowerShell_AddScript
//sys   nativePowerShell_AddScriptSpecifyScope(handle nativePowerShell_PowerShellHandle, path *uint16, useLocalScope byte) (status int)= psh_host.NativePowerShell_AddScriptSpecifyScope
//sys   nativePowerShell_ClosePowerShellObject(psobject nativePowerShell_PowerShellObject) = psh_host.NativePowerShell_ClosePowerShellObject

// to fix status is really nativePowerShell_StringPtr
//sys   nativePowerShell_GetPSObjectType(handle nativePowerShell_PowerShellObject) (status uintptr)= psh_host.NativePowerShell_GetPSObjectType

// to fix status is really nativePowerShell_StringPtr
//sys   nativePowerShell_GetPSObjectToString(handle nativePowerShell_PowerShellObject) (status uintptr)= psh_host.NativePowerShell_GetPSObjectToString
//sys   nativePowerShell_IsPSObjectNullptr(handle nativePowerShell_PowerShellObject) (status byte) = psh_host.NativePowerShell_IsPSObjectNullptr
//sys   nativePowerShell_AddPSObjectHandle(handle nativePowerShell_PowerShellObject) (status nativePowerShell_PowerShellObject)= psh_host.NativePowerShell_AddPSObjectHandle
//sys 	nativePowerShell_DefaultAlloc(size uint64) (status uintptr) = psh_host.NativePowerShell_DefaultAlloc
//sys 	nativePowerShell_DefaultFree(address uintptr) = psh_host.NativePowerShell_DefaultFree

//sys memcpy(dest uintptr, src uintptr, size uint64) (ptr uintptr) = ntdll.memcpy
