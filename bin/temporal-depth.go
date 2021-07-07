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

	file, w := segque.CreateDataFile(p)
	defer file.Close()

	for seqnum := uint64(1); seqnum < p.Runlen; seqnum++ {
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

		if seqnum > p.Warmup && evicted > 0 {
			residency := int(seqnum - evicted)
			if residency < 0 {
				panic("bug - negative residency value - warmup too short!")
			}
			segque.WriteInt(w, residency)
		}
	}
	w.Flush()

	segque.Trace(p, "run() -- EXIT\n")
}
