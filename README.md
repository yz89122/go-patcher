# Go Patcher

Generate overlay patches JSON for go build `-overlay` option.

## Usage

```bash
# build this tool
go install 'github.com/yz89122/go-patcher@latest'
# use with go build to build your binary
go build -overlay "$(go-patcher the/patches/dir)" -o main ./main.go
```
