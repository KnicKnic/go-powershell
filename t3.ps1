
write-error "asf"

write-host $ErrorActionPreference

write-host $localScope
$localScope = "local"

write-host $script:scriptScope
$script:scriptScope = "script"

write-host $global:globalScope
$global:globalScope = "global"


