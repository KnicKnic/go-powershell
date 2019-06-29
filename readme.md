# install
Windows - install dependencies - Use choco

1. `choco install mingw -y`

## Correct
copy /y ..\native-powershell\host.h .
copy /y ..\native-powershell\x64\Release\psh_host.dll .


copy /y ..\psh_host\host.h .
copy /y ..\psh_host\host.h .\powershell\
copy /y ..\psh_host\x64\Release\psh_host.dll .
copy /y ..\psh_host\x64\Release\psh_host.dll .\powershell\
copy /y ..\psh_host\x64\Debug\psh_host.dll .

.\go-net.exe -logtostderr a
.\go-net.exe -command c:\\code\\go-net\\t1.ps1 -command c:\\code\\go-net\\t2.ps1 -logtostderr

# Docs
https://grokbase.com/t/gg/golang-nuts/154m672a6t/go-nuts-linking-cgo-with-visual-studio-x64-release-libraries-on-windows
