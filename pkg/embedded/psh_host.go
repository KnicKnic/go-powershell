package embedded

import "net/http"

// OpenPshHostDll opens a /psh_host.dll
func OpenPshHostDll() http.File {
	file, err := Binaries.Open("/psh_host.dll")
	if err != nil {
		panic(err)
	}
	return file
}
