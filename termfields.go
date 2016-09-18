//Package termfields creates updateable form fields at specified locations in the console.
package termfields

import (
	"fmt"

	tb "github.com/nsf/termbox-go"
)

var boxRunesMap map[boxStyle][]rune

type (
	boxStyle uint16
	shiftDir uint16
)

// Flags to style a box border around a field
const (
	boxStyleClear boxStyle = iota
	BoxStyleNone
	BoxStyleASCII
	BoxStyleUnicode
)

// Flags to Shift a field in a specified direction
const (
	FieldShiftLeft shiftDir = iota
	FieldShiftRight
	FieldShiftUp
	FieldShiftDown
)

// Field is the identifier for a specific form field on the screen.
type Field struct {
	field
}

type field struct {
	x, y   int
	len    int
	border boxStyle
	text   string
}

func init() {
	boxRunesMap = map[boxStyle][]rune{
		boxStyleClear:   {' ', ' ', ' ', ' ', ' ', ' '},
		BoxStyleNone:    {0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		BoxStyleASCII:   {'+', '+', '+', '+', '-', '|'},
		BoxStyleUnicode: {0x250c, 0x2510, 0x2514, 0x2518, 0x2500, 0x2502},
	}
}

// Init Initializes termfields library. This function should be called before any other functions.
// After successful initialization, the library must be finalized using 'Close' function.
//
// Example usage:
//      err := termfields.Init()
//      if err != nil {
//              panic(err)
//      }
//      defer termfields.Close()
func Init() error {
	return tb.Init()
}

// Close Finalizes termbox library, should be called after successful initialization
// when termbox's functionality isn't required anymore.
func Close() {
	tb.SetCursor(0, 0)
	tb.Close()
}

// Row returns the row of a field.
func (f *field) Row() int {
	return f.y
}

// Column returns the column of a field.
func (f *field) Column() int {
	return f.x
}

// Loc defines a new location for a field.
func (f *field) Loc(y, x int) {
	border := f.border
	f.Update(fmt.Sprintf("%*s", f.len, " "))
	f.DrawBox(boxStyleClear)
	f.x = x
	f.y = y
	f.DrawBox(border)
	f.Update(f.text)
}

// Shift shifts a field a direction based on the value of a moveDir
func (f *field) Shift(dir shiftDir) {
	border := f.border
	f.DrawBox(boxStyleClear)
	switch {
	case dir == FieldShiftLeft:
		f.x--
	case dir == FieldShiftRight:
		f.x++
	case dir == FieldShiftUp:
		f.y--
	case dir == FieldShiftDown:
		f.y++
	}
	f.DrawBox(border)
	f.Update(f.text)
}

// NewField creates a new field at location y,x of lenth len with contents text.
func NewField(y, x, len int, text string) (*Field, error) {
	f := field{
		x:   x,
		y:   y,
		len: len,
	}
	err := f.Update(text)
	if err != nil {
		return nil, err
	}
	return &Field{f}, nil
}

func (f *field) DrawBox(boxType boxStyle) error {
	if !tb.IsInit {
		return fmt.Errorf("Term not Initialized")
	}
	if _, ok := boxRunesMap[boxType]; !ok {
		return fmt.Errorf("Unknown Box Style")
	}

	//Draw Corners
	tb.SetCell(f.x-1, f.y-1, boxRunesMap[boxType][0], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x+f.len+1, f.y-1, boxRunesMap[boxType][1], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x-1, f.y+1, boxRunesMap[boxType][2], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x+f.len+1, f.y+1, boxRunesMap[boxType][3], tb.ColorDefault, tb.ColorDefault)
	//Draw Sides
	tb.SetCell(f.x-1, f.y, boxRunesMap[boxType][5], tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(f.x+f.len+1, f.y, boxRunesMap[boxType][5], tb.ColorDefault, tb.ColorDefault)
	//Draw Top
	for i := 0; i < f.len+1; i++ {
		tb.SetCell(f.x+i, f.y-1, boxRunesMap[boxType][4], tb.ColorDefault, tb.ColorDefault)
		tb.SetCell(f.x+i, f.y+1, boxRunesMap[boxType][4], tb.ColorDefault, tb.ColorDefault)
	}
	tb.Flush()
	f.border = boxType
	return nil
}

func (f *field) Update(s string) error {
	if !tb.IsInit {
		return fmt.Errorf("Term not Initialized")
	}
	for i, c := range s {
		tb.SetCell(f.x+i, f.y, c, tb.ColorDefault, tb.ColorDefault)
	}
	tb.Flush()
	f.text = s
	return nil
}
