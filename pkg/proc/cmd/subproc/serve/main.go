package main

import (
	"fmt"
	"time"

	"github.com/cybriq/p9/pkg/proc"
	"github.com/cybriq/p9/pkg/qu"
)

func main() {
	p := proc.Serve(
		qu.T(), func(b []byte) (e error) {
			fmt.Print("from parent: ", string(b))
			return
		},
	)
	for {
		_, e := p.Write([]byte("ping"))
		if e != nil {
			fmt.Println("err:", e)
		}
		time.Sleep(time.Second)
	}
}
