package main

import (
	"github.com/cybriq/p9/pkg/log"
	"github.com/cybriq/p9/pkg/proc"

	"os"
	"time"

	"github.com/cybriq/p9/pkg/qu"
)

func main() {
	// var e error
	log.SetLogLevel("trace")
	// command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := qu.T()
	// splitted := strings.Split(command, " ")
	splitted := os.Args[1:]
	w := proc.LogConsume(
		quit, proc.SimpleLog(splitted[len(splitted)-1]),
		proc.FilterNone, splitted...,
	)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> starting")
	proc.Start(w)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> started")
	time.Sleep(time.Second * 4)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopping")
	proc.Kill(w)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopped")
	// time.Sleep(time.Second * 5)
	// D.Ln(interrupt.GoroutineDump())
	// if e = w.Wait(); E.Chk(e) {
	// }
	// time.Sleep(time.Second * 3)
}
