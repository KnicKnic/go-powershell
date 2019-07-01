

write-error "someerror"
write-host $script:blah
$script:blah=5
write-host $script:blah1
write-host $global:blah1
$script:blah1=5
$global:blah1=5

write-host $(1..3 | send-hostcommand -message "asdfasdf")

write-host $(send-hostcommand -message "two")
write-host $($null | send-hostcommand -message "three")
write-host $(@($null, $null) | send-hostcommand -message "four")


function global:Hi($string){
write-host $string
}