package main

import (
	"fmt"

	"github.com/cybriq/p9/pkg/proc"
)

func main() {
	proc.AddHandler(
		func() {
			fmt.Println("IT'S THE END OF THE WORLD!")
		},
	)
	<-proc.HandlersDone
}
