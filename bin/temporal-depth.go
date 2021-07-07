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

	container := segque.NewContainer(p.Ctype, p.Size, p.Slots, p.Seqmask)
	hfunc, pfunc := segque.NewHashFunc(p.Htype)
	segque.Emit(p, "%v %v %v\n", container, hfunc, pfunc)

	rndparams := make([]interface{}, container.ArrCnt())
	for i := 0; i < len(rndparams); i++ {
		rndparams[i] = pfunc(i)
	}

	file, w := segque.CreateDataFile(p)
	defer file.Close()

	for seqnum := uint64(1); seqnum < p.Runlen; seqnum++ {
		var evicted uint64
		keys := make([]uint64, 4)
		for i := 0; i < len(rndparams); i++ {
			keys[i] = segque.Randu64(hfunc, seqnum, rndparams[i])
		}
		evicted = container.Update(p, seqnum, keys...)
		/*
			switch p.Ctype {
			case segque.BA:
				keys[0] = segque.Randu64(hfunc, seqnum, rndparams[0])
				evicted = container.Update(p, seqnum, keys...)
			case segque.Co2_I_R, segque.Co2_II_R, segque.Co2_I_C, segque.Co2_II_C:
			default:
				k0 := segque.Randu64(hfunc, seqnum, rndparams[0])
				k1 := segque.Randu64(hfunc, seqnum, rndparams[1])
				evicted = container.Update(p, seqnum, k0, k1)
			}
		*/
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
