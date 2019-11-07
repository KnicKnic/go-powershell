#generate

to generate zpsh_host.go
`go generate`

You will have to change `modpsh_host = windows.NewLazySystemDLL("psh_host.dll")` to `modpsh_host = windows.NewLazyDLL("psh_host.dll")`
