package main

import (
	"flag"
	"os"
	"unsafe"

	"github.com/golang/glog"

	"bitbucket.org/creachadair/shell"
	"golang.org/x/sys/windows"
)

/*

#cgo CFLAGS: -I.
#cgo LDFLAGS: ./psh_host.dll


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
func (command PowershellCommand) AddCommand(commandlet string) {
	cs, _ := windows.UTF16PtrFromString(commandlet)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddCommand(command.handle, C.MakeWchar(ptrwchar))
}

// AddArgument to an existing powershell command
func (command PowershellCommand) AddArgument(argument string) {
	cs, _ := windows.UTF16PtrFromString(argument)

	ptrwchar := unsafe.Pointer(cs)

	_ = C.AddArgument(command.handle, C.MakeWchar(ptrwchar))
}

// Invoke the powershell command, do not reuse afterwards
func (command PowershellCommand) Invoke() {

	_ = C.InvokeCommand(command.handle)
}

// ExecStr - executes a commandline in powershell
func (runspace Runspace) ExecStr(commandStr string) {
	command := runspace.CreatePowershellCommand()
	defer command.Delete()

	fields, ok := shell.Split(commandStr)
	if !ok {
		panic("command was invalid {" + commandStr + "}")
	}
	command.AddCommand(fields[0])
	for i := 1; i < len(fields); i++ {
		command.AddArgument(fields[i])
	}
	command.Invoke()
}

// Example on how to use powershell wrappers
func Example() {
	C.InitLibraryHelper()
	runspace := CreateRunspace()
	defer runspace.Delete()

	for i := 1; i < len(os.Args); i++ {
		runspace.ExecStr("mkdir " + os.Args[i])
	}
}
func main() {
	flag.Parse()
	Example()
}
