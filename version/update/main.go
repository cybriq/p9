package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var (
	URL                 string
	GitRef              string
	GitCommit           string
	BuildTime           string
	Tag                 string
	Major, Minor, Patch int
	Meta                string
	PathBase            string
)

func main() {

	// D.Ln(version.Get())
	BuildTime = time.Now().Format(time.RFC3339)
	var cwd string
	var e error
	if cwd, e = os.Getwd(); E.Chk(e) {
		return
	}
	var bumptag bool
	if len(os.Args) > 1 {

		switch {

		case os.Args[1] == "patch":
			bumptag = true

		case os.Args[1] == "minor":
			bumptag = true

		case os.Args[1] == "major":
			bumptag = true
		}
	}
	if !bumptag {
		cwd = filepath.Dir(cwd)
	}
	D.Ln(cwd)
	var repo *git.Repository
	if repo, e = git.PlainOpen(cwd); E.Chk(e) {
		return
	}
	var rr []*git.Remote
	if rr, e = repo.Remotes(); E.Chk(e) {
		return
	}
	for i := range rr {
		rs := rr[i].String()
		if strings.HasPrefix(rs, "origin") {
			rss := strings.Split(rs, "git@")
			if len(rss) > 1 {
				rsss := strings.Split(rss[1], ".git")
				URL = strings.ReplaceAll(rsss[0], ":", "/")
				break
			}
			rss = strings.Split(rs, "https://")
			if len(rss) > 1 {
				rsss := strings.Split(rss[1], ".git")
				URL = rsss[0]
				break
			}

		}
	}
	var tr *git.Worktree
	if tr, e = repo.Worktree(); E.Chk(e) {
	}
	var rh *plumbing.Reference
	if rh, e = repo.Head(); E.Chk(e) {
		return
	}
	rhs := rh.Strings()
	GitRef = rhs[0]
	GitCommit = rhs[1]
	var rt storer.ReferenceIter
	if rt, e = repo.Tags(); E.Chk(e) {
		return
	}
	var maxVersion int
	var maxString string
	var maxIs bool
	if e = rt.ForEach(
		func(pr *plumbing.Reference) (e error) {
			s := strings.Split(pr.String(), "/")
			prs := s[2]
			if strings.HasPrefix(prs, "v") {
				var va [3]int
				var meta string
				_, _ = fmt.Sscanf(
					prs, "v%d.%d.%d%s", &va[0], &va[1], &va[2],
					&meta,
				)
				vn := va[0]*1000000 + va[1]*1000 + va[2]
				if maxVersion < vn {
					maxVersion = vn
					maxString = prs
					Major = va[0]
					Minor = va[1]
					Patch = va[2]
					Meta = meta
				}
				if pr.Hash() == rh.Hash() {
					maxIs = true
					return
				}
			}
			return
		},
	); E.Chk(e) {
		return
	}
	if !maxIs {
		maxString += "+"
	}
	if bumptag {

		switch {

		case os.Args[1] == "patch":
			Patch++

		case os.Args[1] == "minor":
			Minor++

		case os.Args[1] == "major":
			Major++
		}
	}

	Tag = maxString
	PathBase = tr.Filesystem.Root() + "/"
	// I.Ln(PathBase)
	if bumptag {
		Tag = fmt.Sprintf("v%d.%d.%d", Major, Minor, Patch)
	}

	versionFile := `package version

` + `//go:generate go run ./update/.

import (
	"fmt"
)

var (

	// URL is the git URL for the repository
	URL = "%s"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "%s"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "%s"
	// BuildTime stores the time when the current binary was built
	BuildTime = "%s"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "%s"
	// PathBase is the path base returned from runtime caller
	PathBase = "%s"
	// Major is the major number from the tag
	Major = %d
	// Minor is the minor number from the tag
	Minor = %d
	// Patch is the patch version number from the tag
	Patch = %d
	// Meta is the extra arbitrary string field from Semver spec
	Meta = "%s"
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
`
	versionFileOut := fmt.Sprintf(
		versionFile,
		URL,
		GitRef,
		GitCommit,
		BuildTime,
		Tag,
		PathBase,
		Major,
		Minor,
		Patch,
		Meta,
	)
	path := filepath.Join(filepath.Join(PathBase, "version"), "version.go")
	if e = ioutil.WriteFile(path, []byte(versionFileOut), 0666); E.Chk(e) {
	}
	D.Ln("updated version.go written")
	if bumptag {
		bumpTag(Major, Minor, Patch)
	}
	return
}

func bumpTag(major, minor, patch int) {

	I.Ln(os.Args)
	getwd, _ := os.Getwd()
	I.Ln("cwd", getwd)
	tag := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	I.Ln("tag", tag)
	cmd := exec.Command("git", "tag", tag)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	commitMsg := "bumping version"
	I.Ln(os.Args)
	if len(os.Args) > 2 {
		commitMsg = strings.Join(os.Args[2:], " ")
	}

	cmd = exec.Command("git", "commit", fmt.Sprintf("-m%s", commitMsg))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("git", "push", "origin", "main")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("git", "push", "origin", tag)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()

	I.Ln("pushed tag", tag, "and updated repo")
}
