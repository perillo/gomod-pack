// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
)

func joinList(pathlist []string) string {
	return strings.Join(pathlist, string(os.PathListSeparator))
}

// UpdateEnviron updates the environment with a custom GOPRIVATE, GOPROXY and
// PATH environment variables.
func UpdateEnviron(modpath, path string) {
	// Set modpath private to avoid the use of GOPROXY and GOSUMDB and use a
	// direct connection to force the use of git.
	os.Setenv("GOPRIVATE", modpath)
	os.Setenv("GOPROXY", "direct")

	// Set the path with the custom git script at the front.
	oldpath := os.Getenv("PATH")
	if oldpath == "" {
		os.Setenv("PATH", path)
	} else {
		os.Setenv("PATH", joinList([]string{path, oldpath}))
	}
}
