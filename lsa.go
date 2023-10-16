package graphiso

import (
	"fmt"
	"math"
	"sort"
)

var _ = fmt.Printf

const (
	RECTANGULAR_LSAP_INFEASIBLE = -1
	RECTANGULAR_LSAP_INVALID    = -2
)

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func minFloat64Array(a []float64) float64 {
	min := math.MaxFloat64
	for _, v := range a {
		min = minFloat64(min, v)
	}
	return min
}

func maxFloat64(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

func setInt64Array(a []int64, val int64) {
	for i := range a {
		a[i] = val
	}
}

func setFloat64Array(f []float64, val float64) {
	for i := range f {
		f[i] = val
	}
}

func setBoolArray(b []bool, val bool) {
	for i := range b {
		b[i] = val
	}
}

// the C++ code is a template
// but I think is is only used for intptr_t which whould be int64
// this looks right from the c++ but might need to verfiy
func argsortIter(arg []int64) []int64 {
	index := make([]int64, len(arg))
	for i := range index {
		index[i] = int64(i)
	}
	// I think that this is right - not sure
	// sorts index based on order arg
	sort.Slice(index, func(i, j int) bool { return arg[i] < arg[j] })
	return index
}

func augmenting_path(
	nc int64,
	cost []float64,
	u []float64,
	v []float64,
	path []int64,
	row4col []int64,
	shortestPathCosts []float64,
	i int64,
	SR []bool,
	SC []bool,
	pMinVal *float64) int64 {

	minVal := float64(0.0)

	num_remaining := nc
	remaining := make([]int64, nc)
	for i := int64(0); i < nc; i++ {
		remaining[i] = nc - i - 1
	}

	setBoolArray(SR, false)
	setBoolArray(SC, false)
	setFloat64Array(shortestPathCosts, math.MaxFloat64)

	sink := int64(-1)

	for sink == -1 {

		index := int64(-1)
		lowest := math.MaxFloat64
		SR[i] = true

		for it := int64(0); it < num_remaining; it++ {
			j := remaining[it]
			r := minVal + cost[i*nc+j] - u[i] - v[j]
			if r < shortestPathCosts[j] {
				path[j] = i
				shortestPathCosts[j] = r
			}

			// When multiple nodes have the minimum cost, we select one which
			// gives us a new sink node. This is particularly important for
			// integer cost matrices with small co-efficients.
			if shortestPathCosts[j] < lowest ||
				shortestPathCosts[j] == lowest && row4col[j] == -1 {
				lowest = shortestPathCosts[j]
				index = it
			}
		}

		minVal = lowest
		j := remaining[index]
		if minVal == math.MaxFloat64 {
			return -1
		}

		if row4col[j] == -1 {
			sink = j
		} else {
			i = row4col[j]
		}

		SC[j] = true
		num_remaining--
		remaining[index] = remaining[num_remaining]
		remaining = remaining[:num_remaining]
	}

	*pMinVal = minVal
	return sink
}

// square NxN case

func lsa(nr, nc int64, input_cost []float64, maximizing bool, a, b []int64) int {
	if nr == 0 || nc == 0 {
		return 0
	}

	cost := make([]float64, int64(len(input_cost))+(nr*nc))

	var transpose bool = nr < nc
	if transpose {
		for i := int64(0); i < nr; i++ {
			for j := int64(0); j < nc; j++ {
				cost[j*nr+i] = input_cost[i*nc+j]
			}
		}
		nr, nc = nc, nr
	} else {
		cost = input_cost
	}

	// assume we are always maximizing
	if maximizing {
		for i := int64(0); i < nr*nc; i++ {
			cost[i] = -cost[i]
		}
	}

	minval := minFloat64Array(cost)

	// build non-negative cost matrix
	for i := int64(0); i < nr*nc; i++ {
		shift := cost[i] - minval
		cost[i] = shift

		// test for NaN
		// seems hokey
		if math.IsNaN(shift) || shift == -math.MaxFloat64 {
			return RECTANGULAR_LSAP_INVALID
		}
	}

	// create working storage for algorithm
	u := make([]float64, nr)
	v := make([]float64, nc)
	shortestPathCosts := make([]float64, nc)
	path := make([]int64, nc)
	setInt64Array(path, -1)
	col4row := make([]int64, nr)
	setInt64Array(col4row, -1)
	row4col := make([]int64, nc)
	setInt64Array(row4col, -1)
	SR := make([]bool, nr)
	SC := make([]bool, nc)

	// build iterative solution

	for curRow := int64(0); curRow < nr; curRow++ {
		var minVal float64
		sink := augmenting_path(nc, cost, u, v, path, row4col, shortestPathCosts,
			curRow, SR, SC, &minVal)
		if sink < int64(0) {
			return RECTANGULAR_LSAP_INFEASIBLE
		}

		// update dual variables
		u[curRow] += minVal
		for i := int64(0); i < nr; i++ {
			if SR[i] && i != curRow {
				u[i] += minVal - shortestPathCosts[col4row[i]]
			}
		}
		for j := int64(0); j < nc; j++ {
			if SC[j] {
				v[j] -= minVal - shortestPathCosts[j]
			}
		}

		// augment previous solution
		j := sink
		for {
			i := path[j]
			row4col[j] = i
			col4row[i], j = j, col4row[i]
			if i == curRow {
				break
			}
		}
	}

	if transpose {
		// needs to do the argsort_it
		vec := make([]int64, len(col4row))
		copy(vec, col4row)
		for i, val := range vec {
			a[i] = col4row[val]
			b[i] = val
			i++
		}
	} else {
		for i := int64(0); i < nr; i++ {
			a[i] = i
			b[i] = col4row[i]
		}
	}

	return 0
}

/*
func main() {
	fmt.Printf("starting...\n")

	var retVal int
	var nr int64 = 3
	var nc int64 = 3
	maximize := false
	costs := []float64{
            -0.736381, -1.150318, -1.193758,
            -0.788229, -1.193758, -1.255909,
            -0.593773, -0.736381, -0.788229,
	}
	var a []int64 = make([]int64, 9)
	var b []int64 = make([]int64, 9)

	retVal = lsa(nr, nc, costs, maximize, a, b)
	fmt.Printf("r : %v\n", retVal)
	fmt.Printf("a: %+v\n", a)
	fmt.Printf("b: %+v\n", b)

}
*/
