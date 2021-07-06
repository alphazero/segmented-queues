// Doost!

package main

import (
	"flag"
	"fmt"
	"github.com/alphazero/segque"
)

// define defaults here
var p = segque.Params{
	CLSize:  64,
	Degree:  10,
	Seqbits: 17,
	Slots:   7,
	Ctype:   segque.Co2_II_R,
	Htype:   segque.GomapHash,
	Verbose: false,
	Trace:   false,
}

func init() {
	var ctype = int(p.Ctype)
	var htype = int(p.Htype)
	flag.IntVar(&p.CLSize, "cl", p.CLSize, "cahceline size - does not affect result - only for memsize calcs")
	flag.IntVar(&p.Degree, "d", p.Degree, "array degree - size is 2^degree")
	flag.IntVar(&p.Slots, "n", p.Slots, "clc slot count")
	flag.IntVar(&p.Seqbits, "sb", p.Seqbits, "sequence counter bits")
	flag.IntVar(&ctype, "ct", ctype, "container type: 1:BA 2:C2-IC 3:C2-IR 4:C2-IIC 5:C2-IIR 6:C2-IIRand ")
	flag.IntVar(&htype, "ht", htype, "hash type: 1:Blake2b 2:GoMaphash")
	flag.BoolVar(&p.Verbose, "verbose", p.Verbose, "flag - verbose emits")
	flag.BoolVar(&p.Trace, "trace", p.Trace, "flag - trace run")

	flag.Parse()

	p.Ctype = segque.CType(ctype)
	p.Htype = segque.HType(htype)

	p.Initialize()
	p.DebugPrint()
}

func main() {
	fmt.Printf("Salaam Sultan of Love!\n")

	run(&p)
}

func run(p *segque.Params) {
	segque.Trace(p, "run() -- ENTER\n")
	// create container and hashfunc per params
	container := segque.NewContainer(p.Ctype, p.Size, p.Slots, p.Seqmask)
	hfunc, pfunc := segque.NewHashFunc(p.Htype)
	segque.Emit(p, "%v %v %v\n", container, hfunc, pfunc)
	rndparams := []interface{}{pfunc(0), pfunc(1)}
	// TODO these should be cmdline args via flags
	var warmup = uint64(p.Capacity << 4)
	var runlen = warmup // << 1
	segque.Emit(p, "warmup: %d runlen: %d\n", warmup, runlen)

	for seqnum := uint64(1); seqnum < runlen; seqnum++ {
		var evicted uint64
		switch p.Ctype {
		case segque.BA:
			k0 := segque.Randu64(hfunc, seqnum, rndparams[0])
			nopStub := uint64(0)
			evicted = container.Update(p, seqnum, k0, nopStub)
		default:
			k0 := segque.Randu64(hfunc, seqnum, rndparams[0])
			k1 := segque.Randu64(hfunc, seqnum, rndparams[1])
			evicted = container.Update(p, seqnum, k0, k1)
		}

		if seqnum > warmup && evicted > 0 {
			r := int(seqnum-evicted) - p.Capacity
			nr := float64(r) / float64(p.Capacity)
			fmt.Printf("%7d  %+f\n", r, nr)
		}
	}

	segque.Trace(p, "run() -- EXIT\n")
}
