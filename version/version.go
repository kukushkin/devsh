package version

import (
	_ "embed"
	"strings"
)

var (
	//go:embed version.txt
	_version string

	Version string = strings.TrimSpace(_version)
)
