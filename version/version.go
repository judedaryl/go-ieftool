package version

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
)

// rawVersion is the current version as a string, as read from the VERSION
// file. This must be a valid semantic version.
//
//go:embed VERSION
var rawVersion string

// dev determines whether the -dev prerelease marker will
// be included in version info. It is expected to be set to "no" using
// linker flags when building binaries for release.
var dev string = "yes"

// Version The main version number that is being run at the moment, populated from the raw version.
var Version string

// Prerelease is a prerelease marker for the version, populated using a combination of the raw version
// and the dev flag.
var Prerelease string

// SemVer is an instance of version.Version representing the main version
// without any prerelease information.
var SemVer *version.Version

func init() {
	semVerFull := version.Must(version.NewVersion(strings.TrimSpace(rawVersion)))
	SemVer = semVerFull.Core()
	Version = SemVer.String()

	if dev == "no" {
		Prerelease = semVerFull.Prerelease()
	} else {
		Prerelease = "dev"
	}
}

// String returns the complete version string, including prerelease
func String() string {
	if Prerelease != "" {
		return fmt.Sprintf("%s-%s", Version, Prerelease)
	}
	return Version
}
