// Doost!

package segque

type CType int

func (c CType) String() string {
	return Strategies[c]
}

const (
	_ CType = iota
	BA
	Co2_I_C
	Co2_I_R
	Co2_II_C
	Co2_II_R
	Co2_II_Rand
)

type Strategy struct {
	Ctype CType
	Name  string
}

var Strategies = map[CType]string{
	BA:          "BA",
	Co2_I_C:     "Co2_I_C",
	Co2_I_R:     "Co2_I_R",
	Co2_II_C:    "Co2_II_C",
	Co2_II_R:    "Co2_II_R",
	Co2_II_Rand: "Co2_II_Rand",
}
