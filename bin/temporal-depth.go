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
	segque.Trace(p, "run() -- ENTER\n")
	// create container and hashfunc per params
	container := segque.NewContainer(p.Ctype, p.Size, p.Slots, p.Seqmask)
	hfunc, pfunc := segque.NewHashFunc(p.Htype)
	segque.Emit(p, "%v %v %v\n", container, hfunc, pfunc)
	rndparams := []interface{}{pfunc(0), pfunc(1)}

	// TODO these should be cmdline args via flags
	var warmup = uint64(p.Capacity << 4)
	var runlen = warmup << 1
	segque.Emit(p, "warmup: %d runlen: %d\n", warmup, runlen)

	file, w := segque.CreateDataFile(p)
	defer file.Close()

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
			segque.Trace(p, "%7d  %+f\n", r, nr)

			segque.WriteUint64(w, seqnum)
			segque.WriteInt(w, r)
			segque.WriteFloat64(w, nr)
		}
	}
	w.Flush()

	segque.Trace(p, "run() -- EXIT\n")
}
