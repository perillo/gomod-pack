// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
)

type environ []string

func (env *environ) Set(key, value string) {
	*env = append(*env, key+"="+value)
}

func (env *environ) SetList(key string, values ...string) {
	value := joinList(values...)
	env.Set(key, value)
}

func joinList(paths ...string) string {
	return strings.Join(paths, string(os.PathListSeparator))
}

// NewEnviron returns a minimal environment with a custom GOPRIVATE, GOPROXY
// and PATH environment variables.
func NewEnviron(modpath, path string) environ {
	env := make(environ, 0, 10)
	env.Set("HOME", os.Getenv("HOME"))
	env.Set("GO111MODULE", "on")

	// Set modpath private to avoid the use of GOPROXY and GOSUMDB and use a
	// direct connection to force the use of git.
	env.Set("GOPRIVATE", modpath)
	env.Set("GOPROXY", "direct")

	// Set GOPATH to only have the custom git script path.
	env.Set("PATH", path)

	return env
}
