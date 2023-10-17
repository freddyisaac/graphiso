package graphiso

import (
	"encoding/json"
	"io"
	"os"

	"gonum.org/v1/gonum/mat"
)

type Edge struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Weight float64 `json:"weight"`
}

type Edges []Edge

func ParseJsonGraph(r io.Reader) (Edges, error) {
	var edges Edges

	b, err := io.ReadAll(r)
	if err != nil {
		return edges, err
	}

	err = json.Unmarshal(b, &edges)

	return edges, nil
}

func ReadJsonGraph(name string) (Edges, error) {

	var edges Edges

	rc, err := os.Open(name)
	if err != nil {
		return edges, err
	}

	edges, err = ParseJsonGraph(rc)

	return edges, err
}

func IntMapForEdges(edges Edges) map[string]int {
	index := 0
	m := make(map[string]int)
	for _, edge := range edges {
		if _, ok := m[edge.From]; !ok {
			m[edge.From] = index
			index++
		}
		if _, ok := m[edge.To]; !ok {
			m[edge.To] = index
			index++
		}
	}
	return m
}

func EdgeListFromMap(labelMap map[string]int) []string {
	ss := make([]string, len(labelMap))
	for label, index := range labelMap {
		ss[index] = label
	}
	return ss
}

func SymMatrixFromStream(r io.Reader) (*mat.SymDense, error) {
	var symM *mat.SymDense
	edges, err := ParseJsonGraph(r)
	if err != nil {
		return symM, err
	}
	m := IntMapForEdges(edges)

	symM = mat.NewSymDense(len(m), nil)
	for _, edge := range edges {
		fromIndex := m[edge.From]
		toIndex := m[edge.To]
		symM.SetSym(fromIndex, toIndex, edge.Weight)
	}
	return symM, nil
}

func GraphLabelMatch(P *mat.Dense, labelsU, labelsV []string) map[string]string {
	nodeMap := make(map[string]string)
	nr, nc := P.Dims()
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			if P.At(i, j) == 1 {
				nodeMap[labelsU[i]] = labelsV[j]
			}
		}
	}
	return nodeMap
}
