# Status
This project is currently not really suitable for consumption. It does work(can call scripts, communicate from golang to powershell, and powershell back to golang). However it is under production. Come back in a few weeks. You can always email me if you are interested.

# Goal
The goal of this project is to enable you to quickly write golang code and interact with windows via powershell. Because powershell is a powerfull scripting language you will sometimes want to call back into golang. This is also permitted. Also due to sometimes wanting to host .net and powershell giving you an easy way to wrap .net modules and functions and objects, this project also enables that.

## Dependencies
This project has a dependency on [native-powershell](https://github.com/KnicKnic/native-powershell). This is a c++/cli project that enables interacting with powershell through a C DLL interface.

### Using native-powershell
1. copy host.h into the powershell folder
1. Copy the compiled psh_host.dll into
    1. powershell folder
    1. any folder that uses the powershell package
    1. the same folder whereever you distribute the golang binary

### Getting cgo (so you can compile)
Windows - install dependencies - Use choco (easiest way to install gcc)

1. `choco install mingw -y`


# Docs
https://grokbase.com/t/gg/golang-nuts/154m672a6t/go-nuts-linking-cgo-with-visual-studio-x64-release-libraries-on-windows
