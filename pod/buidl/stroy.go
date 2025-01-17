//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/cybriq/p9/pkg/fork"

	"github.com/cybriq/p9/pkg/appdata"
	"github.com/cybriq/p9/pkg/apputil"
)

var (
	URL       string
	GitRef    string
	GitCommit string
	BuildTime string
	Tag       string
)

type command struct {
	name string
	args []string
}

var ldFlags []string

func main() {
	fmt.Fprintln(os.Stderr, os.Args)
	var e error
	var ok bool
	var home string
	var list []string
	if home, ok = os.LookupEnv("HOME"); !ok {
		panic(e)
	}
	if len(os.Args) > 1 {
		folderName := "test0"
		var datadir string
		if len(os.Args) > 2 {
			datadir = os.Args[2]
		} else {
			datadir = filepath.Join(home, folderName)
		}
		if list, ok = commands[os.Args[1]]; ok {
			for i := range list {
				// inject the data directory
				var split []string
				out := strings.ReplaceAll(list[i], "%datadir", datadir)
				split = strings.Split(out, " ")
				fmt.Fprintf(
					os.Stderr,
					"executing item %d of list '%v' '%v' '%v'\n",
					i, os.Args[1], split[0], split[1:],
				)
				var cmd *exec.Cmd
				scriptPath := filepath.Join(
					appdata.Dir("buidl", false),
					"buidl.sh",
				)
				apputil.EnsureDir(scriptPath)
				if len(os.Args) > 2 {
					split = append(split, os.Args[2:]...)
				}
				if e = ioutil.WriteFile(
					scriptPath,
					[]byte(strings.Join(split, " ")),
					0700,
				); e != nil {
				} else {
					cmd = exec.Command("sh", scriptPath)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
				}
				if cmd == nil {
					panic("cmd is nil")
				}
				var e error
				if e = cmd.Start(); e != nil {
					fmt.Fprintln(os.Stderr, e)
					os.Exit(1)
				}
				if e := cmd.Wait(); e != nil {
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("command", os.Args[1], "not found")
		}
	} else {
		fmt.Println("no command requested, available:")
		for i := range commands {
			fmt.Println(i)
			for j := range commands[i] {
				fmt.Println("\t" + commands[i][j])
			}
		}
		fmt.Println()
		fmt.Println(
			"add a path after the command name to change the home/root for" +
				" the build",
		)
	}
}
