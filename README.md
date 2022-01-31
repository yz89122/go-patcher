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
