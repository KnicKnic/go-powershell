# install
Windows - install dependencies - Use choco

1. `choco install mingw -y`

## hacks
1. `copy "C:\Program Files (x86)\Windows Kits\NETFXSDK\4.7.2\Include\um\metahost.h" .`
1. `copy "c:\windows\system32\mscoree.dll" .`

## Correct
copy /y ..\native-powershell\host.h .
copy /y ..\native-powershell\x64\Release\psh_host.dll .


copy /y ..\psh_host\host.h .
copy /y ..\psh_host\x64\Release\psh_host.dll .
copy /y ..\psh_host\x64\Debug\psh_host.dll .

.\go-net.exe -logtostderr a
.\go-net.exe -command c:\\code\\go-net\\t1.ps1 -command c:\\code\\go-net\\t2.ps1 -logtostderr

# Docs
https://docs.microsoft.com/en-us/dotnet/framework/unmanaged-api/hosting/clrcreateinstance-function

https://grokbase.com/t/gg/golang-nuts/154m672a6t/go-nuts-linking-cgo-with-visual-studio-x64-release-libraries-on-windows

https://docs.microsoft.com/en-us/dotnet/core/tutorials/netcore-hosting

hosting .net / .netcore
https://github.com/dotnet/samples/blob/master/core/hosting/HostWithMscoree/host.cpp

