package version

//go:generate go run ./update/.

import (
	"fmt"
)

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/p9"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/main"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "c2b540a2739caaf464d9097aa622f73a4c30a0c9"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-05-02T22:19:45+02:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v1.9.25+"
	// PathBase is the path base returned from runtime caller
	PathBase = "/home/loki/src/github.com/p9c/p9/"
	// Major is the major number from the tag
	Major = 1
	// Minor is the minor number from the tag
	Minor = 9
	// Patch is the patch version number from the tag
	Patch = 25
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
