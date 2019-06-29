
## Copy binaries
copy /y ..\native-powershell\host.h .
copy /y ..\native-powershell\x64\Release\psh_host.dll .


copy /y ..\psh_host\host.h .\powershell\
copy /y ..\psh_host\x64\Release\psh_host.dll .
copy /y ..\psh_host\x64\Release\psh_host.dll .\powershell\
copy /y ..\psh_host\x64\Debug\psh_host.dll .

go-powershell.exe -logtostderr a
go-powershell.exe -command c:\\code\\go-net\\t1.ps1 -command c:\\code\\go-net\\t2.ps1 -logtostderr

## forcefull rebuild
go build -a .
