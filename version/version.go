package version

import "fmt"

var (
	// Name of application
	Name = "kolide"

	// Description of application
	Description = "osquery command and control"

	// VersionMajor is for an API incompatible changes
	VersionMajor = 0

	// VersionMinor is for functionality
	// in a backwards-compatible manner
	VersionMinor = 1

	// VersionPatch is for backwards-compatible bug fixes
	VersionPatch = 5

	// VersionBuild holds the git commit that was compiled.
	// This will be filled in by the compiler.
	VersionBuild = ""
)

// Version is the specification version that the
// package types support.
var Version = fmt.Sprintf("%d.%d.%d",
	VersionMajor, VersionMinor, VersionPatch)
