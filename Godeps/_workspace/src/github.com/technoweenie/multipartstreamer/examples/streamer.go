package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"teadiller/Godeps/_workspace/src/github.com/technoweenie/multipartstreamer"
)

func main() {
	defaultPath, _ := os.Getwd()
	defaultFile := filepath.Join(defaultPath, "streamer.go")
	fullpath := flag.String("path", defaultFile, "Path to the include in the multipart data.")
	flag.Parse()

	ms := multipartstreamer.New()

	fmt.Println("Adding the file to the multipart writer")
	ms.WriteFile("file", *fullpath)
	reader := ms.GetReader()

	fmt.Println("Writing the multipart data to a file")
	file, _ := os.Create("streamtest")
	io.Copy(file, reader)
}
