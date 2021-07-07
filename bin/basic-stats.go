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
	file, r, size := segque.OpenDataFile(p)
	defer file.Close()

	fmt.Printf("size: %d\n", size)
	var n int = int(size) / 8
	var stats = segque.NewStats(p, n)
	for i := 0; i < n; i++ {
		residency, eof := segque.ReadInt(r)
		if eof != nil {
			panic("bug - got EOF!")
		}
		stats.Update(residency)
		fmt.Printf("%d\n", residency)
	}
	fmt.Printf("%d %d\n", n, size/8)
	p.DebugPrint()
}
