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
}

// fully initialize Params based on partially defined (from CL) struct
func (p *Params) Initialize() {
	p.Size = 1 << p.Degree
	p.Mask = uint64(p.Size - 1)
	p.Seqmask = uint64((1 << p.Seqbits) - 1)
	p.Capacity = p.Size * p.Slots
	p.Memsize = p.Size * p.CLSize / 1024
}
func (p *Params) Print() {
	p.Fprint(os.Stdout)
}
func (p *Params) Fprint(w io.Writer) {
	// print p for result output reference
	fmt.Println(w, "--- test p -----------------------------------------")
	fmt.Fprintf(w, "cacheline-size  %d\n", p.CLSize)
	fmt.Fprintf(w, "degree:         %d\n", p.Degree)
	fmt.Fprintf(w, "seqnumbits:     %d\n", p.Seqbits)
	fmt.Fprintf(w, "buckets:        %d (degree:%d)\n", p.Size, p.Degree)
	fmt.Fprintf(w, "slots/bucket:   %d\n", p.Slots)
	fmt.Fprintf(w, "capacity:       %d\n", p.Capacity)
	fmt.Fprintf(w, "mem-size:       %d\n", p.Memsize)
	//	fmt.Fprintf(w, "warmup length:  %d\n", p.wup)
	//	fmt.Fprintf(w, "stream length:  %d\n", p.cnt)
	fmt.Fprintf(w, "hashfunc-type:  %s\n", p.Htype)
	fmt.Fprintf(w, "container-type: %s\n", p.Ctype)
	//	fmt.Fprintf(w, "ref-sizes:      %v\n", p.refsizes)
	//	fmt.Fprintf(w, "ref-caps:       %v\n", p.refcaps)
	fmt.Fprintf(w, "verbose-flag:   %t\n", p.Verbose)
	fmt.Fprintf(w, "trace-flag:     %t\n", p.Trace)
	fmt.Fprintln(w, "---------------------------------------------------------")
}
