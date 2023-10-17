package graphiso

import (
	"gonum.org/v1/gonum/mat"
)

func Tm() *mat.Dense {
	return mat.NewDense(3, 3, []float64{
		1.0, 2.1, 5.6,
		2.1, 5.0, -1.0,
		5.6, -1.0, 5.1,
	})
}

func Ag() *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		0.0, 5.0, 8.0, 6.0,
		5.0, 0.0, 5.0, 1.0,
		8.0, 5.0, 0.0, 2.0,
		6.0, 1.0, 2.0, 0.0,
	})
}

func AgSym() *mat.SymDense {
	return mat.NewSymDense(4, []float64{
		0.0, 5.0, 8.0, 6.0,
		5.0, 0.0, 5.0, 1.0,
		8.0, 5.0, 0.0, 2.0,
		6.0, 1.0, 2.0, 0.0,
	})
}

func Ah() *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		0.0, 1.0, 8.0, 4.0,
		1.0, 0.0, 5.0, 2.0,
		8.0, 5.0, 0.0, 5.0,
		4.0, 2.0, 5.0, 0.0,
	})
}

func AhSym() *mat.SymDense {
	return mat.NewSymDense(4, []float64{
		0.0, 1.0, 8.0, 4.0,
		1.0, 0.0, 5.0, 2.0,
		8.0, 5.0, 0.0, 5.0,
		4.0, 2.0, 5.0, 0.0,
	})
}
