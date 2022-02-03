package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

func main() {
	var err = errors.New("test err")
	if replacedErr, ok := err.(interface{ Stack() []uintptr }); ok {
		fmt.Println("successfully replaced")
		for _, pc := range replacedErr.Stack() {
			var f = runtime.FuncForPC(pc)
			var filename, line = f.FileLine(pc)
			fmt.Printf("%s\n\t%s:%d\n", f.Name(), filename, line)
		}
	} else {
		fmt.Fprintln(os.Stderr, `failed to replace package "errors"`)
		os.Exit(1)
	}
}
