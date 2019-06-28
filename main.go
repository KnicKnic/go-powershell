package main

// "bitbucket.org/creachadair/shell"
import (
	"flag"
	"strings"
	"unsafe"

	"github.com/golang/glog"

	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


#include <stddef.h>
#include "powershell.h"

*/
import "C"

func MakeString(str *C.wchar_t) string {
	var count C.int = 0
	var zero C.wchar_t = C.MakeNullTerminator()
	for ; C.GetChar(str, count) != zero; count++ {
	}
	count++
	arr := make([]uint16, count)
	arrPtr := &arr[0]
	ptrwchar := unsafe.Pointer(arrPtr)

	C.MemoryCopy(ptrwchar, str, count*2)

	s := windows.UTF16ToString(arr)
	return s
}

//export LogWchart
func LogWchart(str *C.wchar_t) {
	s := MakeString(str)
	glog.Info(s)
}

type Runspace struct {
	handle C.RunspaceHandle
}

type PowershellCommand struct {
	handle C.PowershellHandle
}

type PowershellObject struct {
	handle uint64
}

type InvokeResults struct {
	objects   []PowershellObject
	count     uint32
	exception PowershellObject
}

func MakePowerShellObjectIndexed(objects *C.PowerShellObject, index uint32) PowershellObject {
	// I don't get why I have to use unsafe.Pointer on C memory
	ptr := uintptr(unsafe.Pointer(objects))
	handle := ptr + (uintptr(index) * unsafe.Sizeof(*objects))
	var obj uint64 = *(*uint64)(unsafe.Pointer(handle))
	return PowershellObject{obj}
}

func MakePowerShellObject(object C.PowerShellObject) PowershellObject {
	var obj uint64 = *(*uint64)(unsafe.Pointer(&object))
	return PowershellObject{obj}
}

func MakeInvokeResults(objects *C.PowerShellObject, count C.uint, exception C.PowerShellObject) (results InvokeResults) {
	results.count = uint32(count)
	results.objects = make([]PowershellObject, count)
	for i := uint32(0); i < results.count; i++ {
		results.objects[i] = MakePowerShellObjectIndexed(objects, i)
	}
	results.exception = MakePowerShellObject(exception)
	return
}

// CreateRunspace think of this kinda like a shell
func CreateRunspace() Runspace {
	return Runspace{C.CreateRunspaceHelper()}
}

// Delete and free a Runspace
func (runspace Runspace) Delete() {
	C.DeleteRunspace(runspace.handle)
}

// CreatePowershellCommand using a runspace, still need to create a command in the powershell command
func (runspace Runspace) CreatePowershellCommand() PowershellCommand {
	return PowershellCommand{C.CreatePowershell(runspace.handle)}
}

// Delete and free a PowershellCommand
func (command PowershellCommand) Delete() {
	C.DeletePowershell(command.handle)
}

// AddCommand to an existing powershell command
func (command PowershellCommand) AddCommand(commandlet string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(commandlet)

	ptrwchar := unsafe.Pointer(cs)

	if useLocalScope {
		_ = C.AddCommandSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 1)

	} else {
		_ = C.AddCommandSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 0)

	}
}

// AddCommand to an existing powershell command
func (command PowershellCommand) AddScript(script string, useLocalScope bool) {
	cs, _ := windows.UTF16PtrFromString(script)

	ptrwchar := unsafe.Pointer(cs)

	if useLocalScope {
		_ = C.AddScriptSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 1)

	} else {
		_ = C.AddScriptSpecifyScope(command.handle, C.MakeWchar(ptrwchar), 0)

	}
}

// AddArgument to an existing powershell command
func (command PowershellCommand) AddArgument(argument string) {
	cs, _ := windows.UTF16PtrFromString(argument)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddArgument(command.handle, C.MakeWchar(ptrwchar))
}

// Invoke the powershell command, do not reuse afterwards
func (command PowershellCommand) Invoke() {

	var objects *C.PowerShellObject
	var count C.uint
	exception := C.InvokeCommand(command.handle, &objects, &count)
	MakeInvokeResults(objects, count, exception)
}

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string) {
	command := runspace.CreatePowershellCommand()
	defer command.Delete()

	// fields, ok := shell.Split(commandStr)
	// if !ok {
	// 	panic("command was invalid {" + commandStr + "}")
	// }

	if strings.HasSuffix(commandStr, ".ps1") {
		command.AddCommand(commandStr, *useLocalScope)
	} else {
		command.AddScript(commandStr, *useLocalScope)
	}
	// for i := 1; i < len(fields); i++ {
	// 	command.AddArgument(fields[i])
	// }
	command.Invoke()
}

// Example on how to use powershell wrappers
func Example() {
	C.InitLibraryHelper()
	runspace := CreateRunspace()
	defer runspace.Delete()

	for i := 0; i < len(commandFlags); i++ {
		commandFlags[i] = strings.ReplaceAll(commandFlags[i], "\\", "\\\\")
		runspace.ExecStr(commandFlags[i])
	}
}

type arrayCommandFlags []string

func (i *arrayCommandFlags) String() string {
	return "my string representation"
}

func (i *arrayCommandFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var commandFlags arrayCommandFlags
var useLocalScope = flag.Bool("useLocalScope", false, "True if should execute scripts in the local scope")

func main() {
	flag.Var(&commandFlags, "command", "Command to run in powershell")
	flag.Parse()
	glog.Info(*useLocalScope, commandFlags)
	Example()
}
