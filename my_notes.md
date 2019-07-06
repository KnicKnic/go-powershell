
## Copy binaries
copy /y ..\native-powershell\host.h .
copy /y ..\native-powershell\x64\Release\psh_host.dll .


copy /y native-powershell\host.h .\pkg\powershell\
copy /y native-powershell\x64\Release\psh_host.dll .
copy /y native-powershell\x64\Release\psh_host.dll .\pkg\powershell
copy /y native-powershell\x64\Release\psh_host.dll .\tests\
copy /y ..\psh_host\x64\Debug\psh_host.dll .

go-powershell.exe -logtostderr a
go-powershell.exe -command .\\tests\\t1.ps1 -command .\\tests\\t2.ps1

## forcefull rebuild
go build -a -o go-powershell.exe .\test_app\cmd
go build -a -o go-powershell.exe .\test_app\simple
