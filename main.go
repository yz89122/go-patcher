package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defaultPatchesDir = "patches"
	patchesFilename   = "patches.json"
)

type Patch struct {
	ImportPath string
	Filename   string
}

type Patches []*Patch

type OverlayJSON struct {
	Replace map[string]string
}

func main() {
	var patches Patches = getPatches()

	var overlayJSON = patchesToOverlayJSON(patches)

	var tmpFilePath = writeOverlayJSONToTmpFile(overlayJSON)

	fmt.Fprint(os.Stdout, tmpFilePath)
}

func printlnAndExit(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func printfAndExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func getPatches() Patches {
	var patchesDir = defaultPatchesDir
	if len(os.Args) > 1 {
		patchesDir = os.Args[1]
	}

	var patches Patches
	if err := fs.WalkDir(os.DirFS(patchesDir), ".", func(path string, f fs.DirEntry, err error) error {
		if f.IsDir() {
			return nil
		}

		var filename = filepath.Base(path)
		if strings.HasSuffix(filename, ".go") {
			patches.Add(filepath.Dir(path), filename)
		}

		// TODO: add diff type

		return nil
	}); err != nil {
		printlnAndExit("failed to list patches, err:", err)
	}
	return patches
}

func (p *Patches) Add(importPath, filename string) {
	*p = append(*p, &Patch{ImportPath: importPath, Filename: filename})
}

func (p *Patches) ImportPathes() []string {
	var importPathes []string
	var importPathesMap = make(map[string]struct{})
	for _, p := range *p {
		if _, seen := importPathesMap[p.ImportPath]; seen {
			continue
		}
		importPathesMap[p.ImportPath] = struct{}{}
		importPathes = append(importPathes, p.ImportPath)
	}
	return importPathes
}

func patchesToOverlayJSON(patches Patches) OverlayJSON {
	var listCommand = exec.Command("go", append([]string{"list", "-json"}, patches.ImportPathes()...)...)

	var stdout io.Reader
	{
		var err error
		if stdout, err = listCommand.StdoutPipe(); err != nil {
			printlnAndExit("failed to get stdout pipe:", err)
		}
		if err = listCommand.Start(); err != nil {
			printlnAndExit("failed to start go list command:", err)
		}
	}

	type goListPackage struct {
		Dir        string
		ImportPath string
	}

	var packages []*goListPackage
	{
		var decoder = json.NewDecoder(stdout)
		for decoder.More() {
			var pkg goListPackage
			if err := decoder.Decode(&pkg); err != nil {
				printlnAndExit("failed to decode JSON from go list command:", err)
			}
			packages = append(packages, &pkg)
		}
	}

	var packagesMap = make(map[string]*goListPackage) // key = import path
	for _, pkg := range packages {
		packagesMap[pkg.ImportPath] = pkg
	}

	var overlayJSON OverlayJSON
	{
		overlayJSON.Replace = make(map[string]string)

		var patchesDir = defaultPatchesDir
		if len(os.Args) > 1 {
			patchesDir = os.Args[1]
		}

		for _, patch := range patches {
			var pkg = packagesMap[patch.ImportPath]
			if pkg == nil {
				printfAndExit("package %s not found", patch.ImportPath)
			}
			overlayJSON.Replace[filepath.Join(pkg.Dir, patch.Filename)] =
				filepath.Join(patchesDir, patch.ImportPath, patch.Filename)
		}
	}

	if err := listCommand.Wait(); err != nil {
		printlnAndExit("go list command failed:", err)
	}

	return overlayJSON
}

func writeOverlayJSONToTmpFile(overlayJSON OverlayJSON) string {
	var jsonBytes []byte
	{
		var err error
		if jsonBytes, err = json.Marshal(overlayJSON); err != nil {
			printlnAndExit("failed to marshal overlay JSON:", err)
		}
	}

	var tmpDir string
	{
		var err error
		if tmpDir, err = os.MkdirTemp(os.TempDir(), "go-patcher"); err != nil {
			printlnAndExit("failed to create temporary directory:", err)
		}
	}

	var tmpFilePath string = filepath.Join(tmpDir, patchesFilename)

	if err := os.WriteFile(tmpFilePath, jsonBytes, 0644); err != nil {
		printfAndExit("failed to write to temporary file [%s]: %v\n", tmpFilePath, err)
	}

	return tmpFilePath
}
