// Doost

package segque

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Params struct {
	// to be provided ------------------------------------
	CLSize  int    // cacheline size
	Degree  int    // size is 1 << degree
	Slots   int    // CLC-n type
	Seqbits int    // default should be minimally log(size) + 4.
	Ctype   CType  // container type
	Htype   HType  // hashfunc type
	Path    string // path to data files
	Verbose bool   // verbose emits
	Trace   bool   // trace all actions

	// computed from provided ----------------------------
	Mask     uint64 // bucket mask is size - 1
	Seqmask  uint64 // seqmask is (1<<seqbits) - 1
	Size     int    // array size is 2^degree
	Capacity int    // total capacity is slots * size
	Memsize  int    // in Kb - array Size * CLSize / 1024 B/Kb
	Warmup   uint64 // computed from capacity
	Runlen   uint64 // fixed at large number for all types
	Filename string
}

func ParseParams() *Params {
	// define defaults here
	var p = Params{
		CLSize:  64,
		Degree:  10,
		Seqbits: 17,
		Slots:   7,
		Ctype:   Co2_II_R,
		Htype:   GomapHash,
		Path:    "data",
		Verbose: false,
		Trace:   false,
	}

	var ctype = int(p.Ctype)
	var htype = int(p.Htype)
	flag.IntVar(&p.CLSize, "cl", p.CLSize, "cahceline size - does not affect result - only for memsize calcs")
	flag.IntVar(&p.Degree, "d", p.Degree, "array degree - size is 2^degree")
	flag.IntVar(&p.Slots, "n", p.Slots, "clc slot count")
	flag.IntVar(&p.Seqbits, "sb", p.Seqbits, "sequence counter bits")
	flag.IntVar(&ctype, "ct", ctype, "container type: 0:BA 1:C2-IC 2:C2-IR 3:C2-IIC 4:C2-IIR 5:C4_IV_C 6:C4_IV_R")
	flag.IntVar(&htype, "ht", htype, "hash type: 1:Blake2b 2:GoMaphash")
	flag.BoolVar(&p.Verbose, "verbose", p.Verbose, "flag - verbose emits")
	flag.BoolVar(&p.Trace, "trace", p.Trace, "flag - trace run")

	flag.Parse()

	p.Ctype = CType(ctype)
	p.Htype = HType(htype)

	p.Initialize()
	p.DebugPrint()

	return &p
}

// fully initialize Params based on partially defined (from CL) struct
func (p *Params) Initialize() {
	p.Size = 1 << p.Degree
	p.Mask = uint64(p.Size - 1)
	p.Seqmask = uint64((1 << p.Seqbits) - 1)
	p.Capacity = p.Size * p.Slots
	p.Memsize = p.Size * p.CLSize / 1024
	p.Warmup = uint64(p.Capacity << 2)
	p.Runlen = uint64(0x100000) + p.Warmup
	if p.Path == "" {
		p.Path = "."
	}
	p.Filename = p.Fname()
}

func (p *Params) DebugPrint() {
	p.Fprint(os.Stdout)
}

// suggested canonical file name based on distinguishing params
func (p *Params) Fname() string {
	return fmt.Sprintf("%s/%s_Clc%d_d%d_s%d_%s.dat", p.Path, p.Ctype, p.Slots, p.Degree, p.Seqbits, p.Htype)
}

func (p *Params) Fprint(w io.Writer) {
	// print p for result output reference
	Emit(p, "--- test parameters --------------------------------\n")
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
	Emit(p, "warmup:         %d\n", p.Warmup)
	Emit(p, "runlen:         %d\n", p.Runlen)
	//	Emit(p, "ref-sizes:      %v\n", p.refsizes)
	//	Emit(p, "ref-caps:       %v\n", p.refcaps)
	Emit(p, "verbose-flag:   %t\n", p.Verbose)
	Emit(p, "trace-flag:     %t\n", p.Trace)
	Emit(p, "----------------------------------------------------\n")
}
