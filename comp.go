package graphiso

import (
	"gonum.org/v1/gonum/mat"
)

func SymFromDense(M *mat.Dense) *mat.SymDense {
	m, n := M.Dims()
	sM := mat.NewSymDense(m, nil)
	for i := 0; i < m; i++ {
		for j := i; j < n; j++ {
			sM.SetSym(i, j, M.At(i, j))
		}
	}
	return sM
}

func abs(v float64) float64 {
	if v < 0.0 {
		return -v
	}
	return v
}

func SliceFn(s []float64, f func(float64) float64) {
	for i, v := range s {
		s[i] = f(v)
	}
}

// create a dense matrix that reverses column order
// with absolute values
func DenseReverseColsWithAbs(m *mat.Dense) {
	_, nc := m.Dims()
	for j, k := 0, nc-1; j < nc/2; j, k = j+1, k-1 {
		jcol := mat.Col(nil, j, m)
		kcol := mat.Col(nil, k, m)
		SliceFn(jcol, abs)
		SliceFn(kcol, abs)
		m.SetCol(j, kcol)
		m.SetCol(k, jcol)
	}
}

func calcU(m *mat.Dense) *mat.Dense {
	s := SymFromDense(m)
	var eigM mat.EigenSym
	eigM.Factorize(s, true)
	var eigenVec mat.Dense
	eigM.VectorsTo(&eigenVec)
	DenseReverseColsWithAbs(&eigenVec)
	return &eigenVec
}

func Abs(m *mat.Dense) {
	_, nc := m.Dims()
	for j := 0; j < nc; j++ {
		jcol := mat.Col(nil, j, m)
		SliceFn(jcol, abs)
		m.SetCol(j, jcol)
	}
}

func calcUSym(m *mat.SymDense) *mat.Dense {
	var eigM mat.EigenSym
	eigM.Factorize(m, true)
	var eigenVec mat.Dense
	eigM.VectorsTo(&eigenVec)
	//DenseReverseColsWithAbs(&eigenVec)
	Abs(&eigenVec)
	return &eigenVec
}

// construct permutation matrix
func PermM(a, b []int64) *mat.Dense {
	n := len(a)
	P := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		P.Set(int(a[i]), int(b[i]), 1)
	}
	return P
}

// evaluate costs
func Costs(m mat.Dense) ([]int64, []int64) {
	nr, nc := m.Dims()
	n := nr * nc
	a := make([]int64, n)
	b := make([]int64, n)
	var costs []float64
	for i := 0; i < nr; i++ {
		costs = append(costs, m.RawRowView(i)...)
	}
	retval := lsa(int64(nr), int64(nc), costs, true, a, b)
	println("retval : ", retval)
	return a, b
}

func ComputePermutationMatrix(Ag, Ah *mat.SymDense) *mat.Dense {
	nr, nc := Ag.Dims()
	if nr != nc {
		return nil
	}
	Ugs := calcUSym(Ag)
	Uhs := calcUSym(Ah)

	var Us mat.Dense

	Us.Mul(Uhs, Ugs.T())
	a, b := Costs(Us)

	return PermM(a[:nr], b[:nr])
}
