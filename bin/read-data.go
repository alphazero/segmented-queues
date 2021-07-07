// Doost!

package main

import (
	"fmt"
	"github.com/alphazero/segque"
)

func main() {
	fmt.Printf("Salaam Sultan of Love!\n")
	p := segque.ParseParams()
	run(p)
}

func run(p *segque.Params) {
	file, r := segque.OpenDataFile(p)
	defer file.Close()

	for {
		seqnum, eof := segque.ReadUint64(r)
		if eof != nil {
			break
		}
		res, _ := segque.ReadInt(r)
		nres, _ := segque.ReadFloat64(r)
		fmt.Printf("%d : %d : %f\n", seqnum, res, nres)
	}
}
