package gotools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type Module struct {
	Path       string       // module path
	Version    string       // module version
	Versions   []string     // available module versions (with -versions)
	Replace    *Module      // replaced by this module
	Time       *time.Time   // time version was created
	Update     *Module      // available update, if any (with -u)
	Main       bool         // is this the main module?
	Indirect   bool         // is this module only an indirect dependency of main module?
	Dir        string       // directory holding files for this module, if any
	GoMod      string       // path to go.mod file for this module, if any
	GoVersion  string       // go version used in module
	Deprecated string       // deprecation message, if any (with -u)
	Error      *ModuleError // error loading module
}

type ModuleError struct {
	Err string // the error itself
}

func GoList(path, dir string) (*Module, error) {
	cmd := exec.Command("go", "list", "-json", "-m", path)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run '%v': %w", cmd, err)
	}

	var resp Module
	if err := json.NewDecoder(&out).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return &resp, nil
}
