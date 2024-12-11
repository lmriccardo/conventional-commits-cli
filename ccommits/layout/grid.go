package layout

import "github.com/lmriccardo/conventional-commits-cli/ccommits/objects"

type GridElement struct {
	rows []int          // The rows occupied by the current element
	cols []int          // The columns occupied by the current element
	perc float32        // The percentage of screen occupied by the element
	obj  objects.Object // The object for the current element
}

type GridLayout struct {
	nofrows     int            // The total number of rows
	nofcols     int            // The total number of columns
	startpos    objects.Vec2   // The X starting point into the screen
	gridsize    objects.Vec2   // The W and H size of the grid
	elements    []*GridElement // The list of all grid elements
	nofelements int            // Total number of elements added to the grid
	fixed       bool           // If the grid is flexible or not
	currperc    float32        // The current total percentange of remaining space
}

// Create a simple grid layout with the given number of rows, columns,
// starting x and y positions and the width and height.
func GridLayout_new(nr, nc, start_x, start_y, size_w, size_h int, resizable bool) *GridLayout {
	// Create the layout and set all the fields
	layout := new(GridLayout)
	layout.nofrows = nr
	layout.nofcols = nc
	layout.startpos = objects.Vec2{X: start_x, Y: start_y}
	layout.fixed = !resizable
	layout.elements = make([]*GridElement, 0, nr*nc)
	layout.nofelements = 0
	layout.currperc = 100.0
	layout.gridsize = objects.Vec2{X: size_w, Y: size_h}

	return layout
}

// Returns a boolean value indicating if the grid is resizable or not
func (gl *GridLayout) IsResizable() bool {
	return !gl.fixed
}

// Returns the width and the height of the Grid
func (gl *GridLayout) GetGridSize() objects.Vec2 {
	return gl.gridsize
}

// Returns the default single element size in width and height
func (gl *GridLayout) GetDefaultElementSize() objects.Vec2 {

}

// Add an element to the grid that occupied given rows and columns
func (gl *GridLayout) AddGridElement(obj objects.Object, rows, cols []int) {

}
