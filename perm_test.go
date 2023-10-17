package graphiso

import (
	"bytes"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestComputePermutationMatrix(t *testing.T) {
	Ag := AgSym()
	Ah := AhSym()

	P := ComputePermutationMatrix(Ag, Ah)

	pCorrect := mat.NewDense(4, 4, []float64{
		0, 0, 1, 0,
		0, 0, 0, 1,
		1, 0, 0, 0,
		0, 1, 0, 0,
	})

	nr, nc := P.Dims()
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			if P.At(i, j)-pCorrect.At(i, j) != 0 {
				t.Errorf("Permutation matrix failed at indices [%d,%d]\n", i, j)
			}
		}
	}
}

var testGraph1 = `[
        {"from":"u1","to":"u2","weight":5.0},
        {"from":"u1","to":"u3","weight":8.0},
        {"from":"u1","to":"u4","weight":6.0},
        {"from":"u2","to":"u3","weight":5.0},
        {"from":"u2","to":"u4","weight":1.0},
        {"from":"u3","to":"u4","weight":2.0}
]`

var testGraph2 = `[
        {"from":"v1","to":"v2","weight":1.0},
        {"from":"v1","to":"v3","weight":8.0},
        {"from":"v1","to":"v4","weight":4.0},
        {"from":"v2","to":"v3","weight":5.0},
        {"from":"v2","to":"v4","weight":2.0},
        {"from":"v3","to":"v4","weight":5.0}
]`

func TestComputePermutationMatrixFromStream(t *testing.T) {
	buf := bytes.NewBuffer([]byte(testGraph1))
	U, err := SymMatrixFromStream(buf)
	if err != nil {
		t.Errorf("SymMatrixFromStream error %v", err)
	}

	buf = bytes.NewBuffer([]byte(testGraph2))
	V, err := SymMatrixFromStream(buf)
	if err != nil {
		t.Errorf("SymMatrixFromStream error %v", err)
	}

	P := ComputePermutationMatrix(U, V)

	pCorrect := mat.NewDense(4, 4, []float64{
		0, 0, 1, 0,
		0, 0, 0, 1,
		1, 0, 0, 0,
		0, 1, 0, 0,
	})
	nr, nc := P.Dims()
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			if P.At(i, j)-pCorrect.At(i, j) != 0 {
				t.Errorf("Permutation matrix failed at indices [%d,%d]\n", i, j)
			}
		}
	}
}

func TestGraphLabelMatch(t *testing.T) {
	buf := bytes.NewBuffer([]byte(testGraph1))
	U, err := SymMatrixFromStream(buf)
	if err != nil {
		t.Errorf("SymMatrixFromStream error %v", err)
	}

	buf = bytes.NewBuffer([]byte(testGraph2))
	V, err := SymMatrixFromStream(buf)
	if err != nil {
		t.Errorf("SymMatrixFromStream error %v", err)
	}

	P := ComputePermutationMatrix(U, V)

	labelsU := []string{"u1", "u2", "u3", "u4"}
	labelsV := []string{"v1", "v2", "v3", "v4"}

	nodeMap := GraphLabelMatch(P, labelsU, labelsV)

	t.Logf("node map : %v", nodeMap)
}
