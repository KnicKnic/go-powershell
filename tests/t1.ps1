

write-error "someerror"
write-host  'ab', $script:blah "ba"
$script:blah=5
write-host  'ab', $script:blah1, "ba"
write-host  'ab', $global:blah1, "ba"
$script:blah1=5
$global:blah1=5

write-host $(1..3 | send-hostcommand -message "asdfasdf")

write-host $(send-hostcommand -message "two")
write-host $($null | send-hostcommand -message "three")
write-host $(@($null, $null) | send-hostcommand -message "four")


function global:Hi($string){
write-host $string
}