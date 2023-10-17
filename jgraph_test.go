package graphiso

import (
	"bytes"
	"testing"

	"gonum.org/v1/gonum/mat"
)

var testGraph = `[
        {"from":"u1","to":"u2","weight":5.0},
        {"from":"u1","to":"u3","weight":8.0},
        {"from":"u1","to":"u4","weight":6.0},
        {"from":"u2","to":"u3","weight":5.0},
        {"from":"u2","to":"u4","weight":1.0},
        {"from":"u3","to":"u4","weight":2.0}
]`

func TestReadJsonGraph(t *testing.T) {

	buf := bytes.NewBuffer([]byte(testGraph))

	edges, err := ParseJsonGraph(buf)

	if err != nil {
		t.Errorf("ParseJsonGraph error %v", err)
	}

	if len(edges) != 6 {
		t.Errorf("Incorrect number of edges parsed %d entries should be 6", len(edges))
	}
}

func TestIntMapForEdges(t *testing.T) {

	buf := bytes.NewBuffer([]byte(testGraph))

	edges, err := ParseJsonGraph(buf)
	if err != nil {
		t.Errorf("ParseJsonGraph error %v", err)
	}

	m := IntMapForEdges(edges)

	if len(m) != 4 {
		t.Errorf("Incorrect number of entries for graph map %d should be 4", len(m))
	}

	t.Logf("Map %+v", m)
}

func TestSymMatrixFromStream(t *testing.T) {
	buf := bytes.NewBuffer([]byte(testGraph))

	m, err := SymMatrixFromStream(buf)

	if err != nil {
		t.Errorf("SymMatrixFromStream error %v", err)
	}
	t.Logf("M : \n%v\n", mat.Formatted(m))
}
