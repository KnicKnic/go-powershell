
## Copy binaries
copy /y ..\native-powershell\host.h .
copy /y ..\native-powershell\x64\Release\psh_host.dll .



cp .\native-powershell\host.h .\native-powershell\native-powershell-bin\
cp .\native-powershell\x64\Release\psh_host.dll .\native-powershell\native-powershell-bin\
cp .\native-powershell\x64\Release\psh_host.pdb .\native-powershell\native-powershell-bin\
cp .\native-powershell\x64\Release\psh_host.lib .\native-powershell\native-powershell-bin\


copy /y native-powershell\host.h .\pkg\powershell\
copy /y native-powershell\x64\Release\psh_host.dll .\pkg\powershell


copy /y native-powershell\x64\Release\psh_host.dll .
copy /y native-powershell\x64\Release\psh_host.dll .\tests\
copy /y ..\psh_host\x64\Debug\psh_host.dll .

go-powershell.exe -logtostderr a
go-powershell.exe -command .\\tests\\t1.ps1 -command .\\tests\\t2.ps1

## forcefull rebuild
go build -a -o go-powershell.exe .\test_app\cmd
go build -a -o go-powershell.exe .\test_app\simple
