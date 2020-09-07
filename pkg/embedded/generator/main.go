//go:generate go run main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/vfsgen"
)

var Assets http.FileSystem = http.Dir("../../../bin")

func main() {

	// Keep only "/target/dir" and its contents.
	fs := filter.Keep(Assets, func(path string, fi os.FileInfo) bool {
		return path == "/psh_host.dll" || path == "/"
	})

	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "embedded",
		BuildTags:    "",
		VariableName: "Binaries",
		Filename:     "../binaries.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
