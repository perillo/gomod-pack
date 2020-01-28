// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"text/template"
)

const (
	gitTemplate = `#!/bin/sh
{{ .GitPath }} -c url."ssh://{{ .User }}@localhost:{{ .Path }}".insteadOf=https://{{ .ModPath }}
`
)

type gitContext struct {
	GitPath string
	User    string
	ModPath string
	Path    string
}

// tempExecutable returns a temporary executable file with the given base name,
// with an 0o744 permission.
func tempExecutable(name string) (*os.File, error) {
	// Create the file in a temporary directory, in order to ensure that the
	// base name is set correctly.
	dirpath, err := ioutil.TempDir("", "gomod-pack")
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dirpath, name)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o744)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// InstallGit installs a custom temporary git script that will setup an URL
// rewrite in order to allow modpath to be resolved to the local filesystem
// path via an SSH connection to localhost.
func InstallGit(modpath, path string) (string, error) {
	// Find the current OS user and git command path.
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	gitpath, err := exec.LookPath("git")
	if err != nil {
		return "", err
	}

	// Create the custom git script.
	f, err := tempExecutable("git")
	if err != nil {
		return "", err
	}
	defer f.Close()

	t := template.Must(template.New("git").Parse(gitTemplate))
	ctx := gitContext{
		GitPath: gitpath,
		User:    user.Username,
		ModPath: modpath,
		Path:    path,
	}
	if err := t.Execute(f, ctx); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}

	return f.Name(), nil
}
