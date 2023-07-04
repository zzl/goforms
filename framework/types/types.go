package types

import (
	"fmt"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type SignedInt interface {
	int8 | int16 | int32 | int64 | int
}

type UnsignedInt interface {
	uint8 | uint16 | uint32 | uint64 | uint
}

type Int interface {
	SignedInt | UnsignedInt
}

type Float interface {
	float32 | float64
}

type Number interface {
	SignedInt | UnsignedInt | Float
}

type Bool int32

type Numeric32BitsOrMore interface {
	int32 | uint32 | int64 | uint64 | int | uint | float32 | float64 | Bool
}

type PtrCompatible interface {
	Number | uintptr | unsafe.Pointer // | *int16 | *uint16 | *int8 | *uint8
}

type Initable interface {
	Init()
}

type Disposable interface {
	Dispose()
}

type LifecycleAware interface {
	Initable
	Disposable
}

type Flusher interface {
	Flush() error
}

type ValueText struct {
	Value int
	Text  string
}

type Formatter interface {
	Format(value any) string
}

type FormatterFunc func(value any) string

func (me FormatterFunc) Format(value any) string {
	return me(value)
}

//

type Rect struct {
	Left   int
	Top    int
	Right  int
	Bottom int
}

func (this *Rect) String() string {
	return fmt.Sprintf("(%d,%d-%d,%d)(%dx%d)",
		this.Left, this.Top, this.Right, this.Bottom,
		this.Width(), this.Height())
}

func (this *Rect) ToRECT() win32.RECT {
	return win32.RECT{
		Left:   int32(this.Left),
		Top:    int32(this.Top),
		Right:  int32(this.Right),
		Bottom: int32(this.Bottom),
	}
}

func (this *Rect) IsEmpty() bool {
	return this.Left == this.Right && this.Top == this.Bottom
}

func (this *Rect) GetPos() Point {
	return Point{this.Left, this.Top}
}

func (this *Rect) SetPos(left int, top int) {
	this.Right += left - this.Left
	this.Left = left

	this.Bottom += top - this.Top
	this.Top = top
}

func (this *Rect) GetSize() Size {
	return Size{this.Right - this.Left, this.Bottom - this.Top}
}

func (this *Rect) SetSize(width, height int) {
	this.Right = this.Left + width
	this.Bottom = this.Top + height
}

func (this *Rect) Width() int {
	return this.Right - this.Left
}

func (this *Rect) Height() int {
	return this.Bottom - this.Top
}

func (this *Rect) SetRect(x1, y1, x2, y2 int) {
	this.Left = x1
	this.Top = y1
	this.Right = x2
	this.Bottom = y2
}

func (this *Rect) SetRectBySize(x1, y1, width, height int) {
	this.Left = x1
	this.Top = y1
	this.Right = x1 + width
	this.Bottom = y1 + height
}

func (this *Rect) SetHeight(height int) {
	this.Bottom = this.Top + height
}

func (this *Rect) Normalize() {
	if this.Left > this.Right {
		this.Left, this.Right = this.Right, this.Left
	}
	if this.Top > this.Bottom {
		this.Top, this.Bottom = this.Bottom, this.Top
	}
}

func (this *Rect) Offset(dx int, dy int) {
	this.Left += dx
	this.Right += dx
	this.Top += dy
	this.Bottom += dy
}

type Point struct {
	X, Y int
}

type Size struct {
	Width, Height int
}

type BoundsAware interface {
	SetBounds(left, top, width, height int)
	GetBounds() Rect
}

type DataAware interface {
	SetData(key string, value any)
	GetData(key string) any
}

type NameAware interface {
	SetName(name string)
	GetName() string
}
