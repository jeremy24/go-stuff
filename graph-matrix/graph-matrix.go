package GraphMatrix

//import "fmt"

const (
	ON  uint32 = 0x1
	OFF uint32 = 0x0
)

// this is the basic interface that is used for a graph edge graph-matrix
type GraphMatrix interface {
	Has(row int, col int) bool
	//Weight(i int32, j int32) float64
	Connect(rows int, cols int)
	//ConnectMany([]int, []int)
	Dims() (int, int)
	Remove(row int, col int)
	Weight(row int, col int) float32
	AddWeight(row int, col int, weight float32)
}

type bitMatrix struct {
	edges [][]uint32
	weights [][]float32
	rows  int
	cols  int
}


// make sure x <= y
func Order(x int, y int) (uint, uint) {
	if x < y {
		return uint(x), uint(y)
	}
	return uint(y), uint(x)
}

/*
	This will initialize all of the values in a new matrix for us.
	It only stores the values as a triangle matrix in order to save space
 */
func NewMatrix(rows int, cols int) GraphMatrix {
	edges := make([][]uint32, rows) // initialize a slice of dy slices

	for i := 0; i < rows; i++ {
		edges[i] = make([]uint32, cols - i + 1) // initialize a slice of dx unit8 in each of dy slices
	}

	weights := make([][]float32, rows) // initialize a slice of dy slices

	for i := 0; i < rows; i++ {
		weights[i] = make([]float32, cols - i + 1) // initialize a slice of dx unit8 in each of dy slices
	}

	return bitMatrix{edges, weights,rows, cols}
}

func (g bitMatrix) Dims() (int, int) {
	return g.rows, g.cols
}

// Check if an edge exists
func (g bitMatrix) Has(i int, j int) bool {
	row, col := Order(i, j)
	bit := g.GetBit(uint(row), uint(col))
	ret := bit == true
	return ret
}

func (g bitMatrix) Connect(i int, j int) {
	row, col := Order(i, j)
	g.SetBit(uint(row), uint(col), ON)
}

func (g bitMatrix) Weight(i int, j int) float32 {
	row, col := Order(i, j)
	return g.weights[row][col]
}

func (g bitMatrix) AddWeight(i int, j int, weight float32) {
	row, col := Order(i, j)
	g.weights[row][col] = weight
}

func (g bitMatrix) Remove(i int, j int) {
	row, col := Order(i, j)
	g.SetBit(uint(row), uint(col), OFF)
}




// Below here are internal functions




func (g bitMatrix) SetBit(i uint, j uint, flag uint32) {
	row := i
	col := uint32(j) / uint32(32)
	offset := uint32(j) % uint32(32)
	chunk := g.edges[row][col]

	// set the bit to zero
	if flag == OFF {
		g.edges[row][col] = chunk - (chunk & offset)
		return
	}

	// set the bit to one
	g.edges[row][col] = chunk | offset
}

// the return for this will ALWAYS be 0 or 1
func (g bitMatrix) GetBit(i uint, j uint) bool {
	row := i
	col := uint32(j) / uint32(32)
	offset := uint32(j) % uint32(32)
	chunk := g.edges[row][col]

	// get the bit and mask off the rest
	ret := (chunk & offset) == offset
	return ret
}

