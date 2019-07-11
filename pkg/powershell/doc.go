/*
Package powershell allows hosting powershell sessions inside golang enabling bi directional communication between powershell and golang through .Net objects.

Overview

This package can call commandlets, scripts, keep session state between multiple invocations, allows for custom callback routines, and trapping host output redirecting to a custom logger, calling back into golang from powershell.


Lifetimes

All objects that you create or are returned to you must be Closed by you. There are Close routines on all objects that need to be closed. The one exception to this rule is if you are inside a callback and want to return an object that you generated to powershell, and have it closed after powershell has processed it.


Multithreading

A single runspace is not multithreaded, do not try this! Objects generated in one runspace can be used in another assuming there is proper .Net concurrent access to those objects.


Reentrancy

Reentrancy is supported, you can call powershell have it call back to golang, have that code call into powershell as long as this chain has a finite depth that is supported by powershell. This is accomplished by creating nested powershell commands anytime they are created in a callback handler or log routine.


Runspaces

A runspace is where you execute your powershell commands and statements in. It is also a boundary for variables such as "$global". If you specify custom log routines or callback handlers they are also bound to the runspace. This is to enable you to bind context to the log routine or callback handler.

Please see the runspace section for more information on creating a runspace and executing scripts and commands https://godoc.org/github.com/KnicKnic/go-powershell/pkg/powershell#Runspace .


Scripts vs Commands

ExecScript - Use when you want to save variables or create functions or execute a set of lines like how you would normally write in a .ps1 file. Do not use this to execute a .ps1 file (which is known as a command file). See the examples in the Runspace.ExecScript section

ExecCommand - Use when you simply want to call an existing commandlet or function or a .ps1 file. The results are returned. See the examples in the Runspace.ExecCommand section


Scopes

Powershell uses dynamic scoping you can read more about it here https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_scopes?view=powershell-5.1

ExecScript and ExecCommand both have a parameter useLocalScope, what this means is do you wish to create a new child scope for this execution. If true then any variables you save inside powershell will not be accessible. If false then your variables will be persisted at global scope.

If you are still unsure, use true.

*/
package powershell
