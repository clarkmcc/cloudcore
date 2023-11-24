/*
Copyright (C) 2020 Print Tracker, LLC - All Rights Reserved

Unauthorized copying of this file, via any medium is strictly prohibited
as this source code is proprietary and confidential. Dissemination of this
information or reproduction of this material is strictly forbidden unless
prior written permission is obtained from Print Tracker, LLC.
*/

package version

import (
	"fmt"
	"runtime"
)

var Version = "x.y.z"
var Hash = "aabbccdd"

type Info struct {
	GoVersion string `json:"go"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
	Version   string `json:"version"`
	Hash      string `json:"hash"`
}

// Get wraps the provided application version with some metadata about the
// go version this binary was compiled with, and the platform it was compiled on
func Get() Info {
	return Info{
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Version:   Version,
		Hash:      Hash,
	}
}
