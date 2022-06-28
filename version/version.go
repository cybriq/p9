package version

//go:generate go run ./update/.

import (
	"fmt"
)

var (

	// URL is the git URL for the repository
	URL = "github.com/cybriq/p9"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/main"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "11aab033ed63b042b610170f679b6232543a0236"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2022-06-28T10:53:45+03:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v0.0.20+"
	// PathBase is the path base returned from runtime caller
	PathBase = "/home/davidvennik/src/github.com/cybriq/p9/"
	// Major is the major number from the tag
	Major = 0
	// Minor is the minor number from the tag
	Minor = 0
	// Patch is the patch version number from the tag
	Patch = 20
	// Meta is the extra arbitrary string field from Semver spec
	Meta = ""
)

// Get returns a pretty printed version information string
func Get() string {
	return fmt.Sprint(
		"\nRepository Information\n"+
		"\tGit repository: "+URL+"\n",
		"\tBranch: "+GitRef+"\n"+
		"\tCommit: "+GitCommit+"\n"+
		"\tBuilt: "+BuildTime+"\n"+
		"\tTag: "+Tag+"\n",
		"\tMajor:", Major, "\n",
		"\tMinor:", Minor, "\n",
		"\tPatch:", Patch, "\n",
		"\tMeta: ", Meta, "\n",
	)
}
