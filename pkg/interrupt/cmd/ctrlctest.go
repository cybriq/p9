package main

import (
	"fmt"

	"github.com/cybriq/p9/pkg/interrupt"
)

func main() {
	interrupt.AddHandler(func() {
		fmt.Println("IT'S THE END OF THE WORLD!")
	},
	)
	<-interrupt.HandlersDone
}
