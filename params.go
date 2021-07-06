// Doost

package segque

import (
	"fmt"
	"io"
	"os"
)

type Params struct {
	// to be provided ------------------------------------
	CLSize  int   // cacheline size
	Degree  int   // size is 1 << degree
	Slots   int   // CLC-n type
	Seqbits int   // default should be minimally log(size) + 4.
	Ctype   CType // container type
	Htype   HType // hashfunc type
	Verbose bool  // verbose emits
	Trace   bool  // trace all actions

	// computed from provided ----------------------------
	Mask     uint64 // bucket mask is size - 1
	Seqmask  uint64 // seqmask is (1<<seqbits) - 1
	Size     int    // array size is 2^degree
	Capacity int    // total capacity is slots * size
	Memsize  int    // in Kb - array Size * CLSize / 1024 B/Kb
	Filename string
}

// fully initialize Params based on partially defined (from CL) struct
func (p *Params) Initialize() {
	p.Size = 1 << p.Degree
	p.Mask = uint64(p.Size - 1)
	p.Seqmask = uint64((1 << p.Seqbits) - 1)
	p.Capacity = p.Size * p.Slots
	p.Memsize = p.Size * p.CLSize / 1024
	p.Filename = p.Fname()
}

func (p *Params) DebugPrint() {
	p.Fprint(os.Stdout)
}

// suggested canonical file name based on distinguishing params
func (p *Params) Fname() string {
	return fmt.Sprintf("%s_Clc%d_d%d_s%d_%s", p.Ctype, p.Slots, p.Degree, p.Seqbits, p.Htype)
}

func (p *Params) Fprint(w io.Writer) {
	// print p for result output reference
	Emit(p, "--- test p -----------------------------------------\n")
	Emit(p, "cacheline-size  %d\n", p.CLSize)
	Emit(p, "degree:         %d\n", p.Degree)
	Emit(p, "seqnumbits:     %d\n", p.Seqbits)
	Emit(p, "buckets:        %d (degree:%d)\n", p.Size, p.Degree)
	Emit(p, "slots/bucket:   %d\n", p.Slots)
	Emit(p, "capacity:       %d\n", p.Capacity)
	Emit(p, "mem-size:       %d Kb\n", p.Memsize)
	Emit(p, "filename:       %s\n", p.Filename)
	//	Emit(p, "warmup length:  %d\n", p.wup)
	//	Emit(p, "stream length:  %d\n", p.cnt)
	Emit(p, "hashfunc-type:  %s\n", p.Htype)
	Emit(p, "container-type: %s\n", p.Ctype)
	//	Emit(p, "ref-sizes:      %v\n", p.refsizes)
	//	Emit(p, "ref-caps:       %v\n", p.refcaps)
	Emit(p, "verbose-flag:   %t\n", p.Verbose)
	Emit(p, "trace-flag:     %t\n", p.Trace)
	Emit(p, "----------------------------------------------------\n")
}
