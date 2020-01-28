// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
)

// pack will pre-fill the Go module cache with the specified module.
// It will execute go mod download with a custom environment, in order to trick
// the go tool to resolve the remote module path to the local module root
// directory.
//
// pack will return the cached module or an error.
func pack(module *Module) (*CachedModule, error) {
	gitpath, err := InstallGit(module.Path, module.Dir)
	if err != nil {
		return nil, err
	}
	gitdir := filepath.Dir(gitpath)
	env := NewEnviron(module.Path, gitdir)

	return DownloadModule(env, module.Path, module.Version)
}

var version = flag.String("version", "master", "module version")

func main() {
	flag.Parse()
	log.SetFlags(0)

	mod, err := GetModule()
	if err != nil {
		log.Fatal(err)
	}
	mod.Version = *version

	cmod, err := pack(mod)
	if err != nil {
		log.Fatal(err)
	}

	// Print the versioned module path so that it can be used in a go.mod
	// require directive.
	fmt.Println(cmod)
}

func debugf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
