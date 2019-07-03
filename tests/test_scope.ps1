

write-host  'ab', $local:localVar "ba"
$local:localVar=1
write-host  'ab', $script:scriptVar, "ba"
$script:scriptVar=2
write-host  'ab', $global:globalVar, "ba"
$global:globalVar=3
write-host  'ab', $noscope, "ba"
$noscope =4