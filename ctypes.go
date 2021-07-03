// Doost!

package segque

// Container types
type CType int

func (c CType) String() string {
	return ctypes[c]
}

const (
	_           CType = iota
	BA                // basic array with direct addressing
	Co2_I_C           // single array choice of 2 using container sequence number
	Co2_I_R           // single array choice of 2 using record sequence number
	Co2_II_C          // double array choice of 2 using container sequence number
	Co2_II_R          // double array choice of 2 using record sequence number
	Co2_II_Rand       // double array choice of 2 with random choice
)

var ctypes = map[CType]string{
	BA:          "BA",
	Co2_I_C:     "Co2_I_C",
	Co2_I_R:     "Co2_I_R",
	Co2_II_C:    "Co2_II_C",
	Co2_II_R:    "Co2_II_R",
	Co2_II_Rand: "Co2_II_Rand",
}
