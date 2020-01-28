// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Module struct {
	Path    string       // module path
	Version string       // module version
	Main    bool         // is this the main module?
	Dir     string       // directory holding files for this module, if any
	GoMod   string       // path to go.mod file for this module, if any
	Error   *ModuleError // error loading module
}

type ModuleError struct {
	Err string // the error itself
}

func (me *ModuleError) Error() string {
	return me.Err
}

type CachedModule struct {
	Path     string // module path
	Version  string // module version
	Error    string // error loading module
	Info     string // absolute path to cached .info file
	GoMod    string // absolute path to cached .mod file
	Zip      string // absolute path to cached .zip file
	Dir      string // absolute path to cached source root directory
	Sum      string // checksum for path, version (as in go.sum)
	GoModSum string // checksum for go.mod (as in go.sum)
}

func (cm *CachedModule) String() string {
	return cm.Path + "@" + cm.Version
}

// invokeGo returns the stdout of a go command invocation.
// The implementation is based on golang.org/x/tools/go/packages, but greatly
// simplified.
func invokeGo(env []string, verb string, args ...string) ([]byte, error) {
	args = append([]string{verb}, args...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command("go", args...)
	cmd.Env = env
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		// Just return the error, including the stderr output as is.
		// Make sure to also return stdout.Bytes() if there is some data, since
		// it may be important.
		args := strings.Trim(fmt.Sprint(args), "[]")
		var buf []byte
		if stdout.Len() > 0 {
			buf = stdout.Bytes()
		}

		return buf, fmt.Errorf("go %v: %w: %s", args, err, stderr)
	}

	return stdout.Bytes(), nil
}

// GetModule returns the current main module.
func GetModule() (*Module, error) {
	buf, err := invokeGo(nil, "list", "-m", "-e", "-json")
	if err != nil {
		return nil, fmt.Errorf("get module: %w", err)
	}

	mod := new(Module)
	if err := json.Unmarshal(buf, mod); err != nil {
		return nil, fmt.Errorf("get module: JSON unmarshalling: %w", err)
	}
	if mod.Error != nil {
		return nil, fmt.Errorf("get module: %w", mod.Error)
	}
	if mod.GoMod == "" {
		return nil, fmt.Errorf("get module: not inside a module")
	}

	return mod, nil
}

// DownloadModule downloads the specified module and return the cached module
// informations.
func DownloadModule(env []string, module, version string) (*CachedModule, error) {
	modpath := module
	if version == "" {
		// Use @master to force cmd/go to find the remote module.
		modpath = modpath + "@master"
	} else {
		modpath = modpath + "@" + version
	}

	buf, err := invokeGo(env, "mod", "download", "-json", modpath)
	if err != nil && buf == nil {
		// Special case, since go mod download -json returns an exit status 1
		// in case of errors, in addition of setting Module.Error, as with go
		// list -m -e -json.
		return nil, fmt.Errorf("download module %q: %w", module, err)
	}

	mod := new(CachedModule)
	if err := json.Unmarshal(buf, mod); err != nil {
		return nil, fmt.Errorf("download module: JSON unmarshalling: %w", err)
	}
	if mod.Error != "" {
		return nil, fmt.Errorf("download module: %s", mod.Error)
	}

	return mod, nil
}
