# Go Patcher

Generate overlay patches JSON for go build `-overlay` option.

## Usage

```bash
# build this tool
go install 'github.com/yz89122/go-patcher@latest'
# use with go build to build your binary
go build -overlay "$(go-patcher the/patches/dir)" -o main ./main.go
```

## Arguments

The second argument (which is `the/pathes/dir` in above example) is the path to the patches directory, it's optional and default to `$PWD/patches`.

## Structure of the Patches directory

In the patches directory, the structure of the directories will be taken as import pathes. For example, `errors/errors.go` will replace the file `errors.go` in the  standard package `errors`.

## Example Output

```
go-patcher example
/tmp/go-patcher574244752/patches.json
```

## Patcher Example

In this example, we replaced the standard `errors.New()` to return a new error with stack trace, and printed the stack trace in the main function.

```go
// main.go
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
```

### Run

```bash
# clone this repository
git clone 'git@github.com:yz89122/go-patcher.git'
# to cloned repository
cd 'go-patcher'
# to example directory
cd 'run-example'
# run example
./run-example.sh
```

#### Example Output

```
successfully replaced
runtime.main
        /usr/local/go/src/runtime/proc.go:264
runtime.goexit
        /usr/local/go/src/runtime/asm_amd64.s:1582
```
