package drawing

import (
	"fmt"
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/scope"
	"math"
)

type Scope = scope.Scope

type Point = gdip.Point
type Rect gdip.Rect
type Size = gdip.Size

type PointF = gdip.PointF
type RectF gdip.RectF
type SizeF = gdip.SizeF

type GdipError gdip.Status

func (me GdipError) Error() string {
	s := fmt.Sprintf("gdip error: %d", me)
	return s
}

func Rc(x, y, w, h int32) Rect {
	return Rect{x, y, w, h}
}

func Pt(x int32, y int32) Point {
	return Point{x, y}
}

func (this *Rect) Location() Point {
	return Point{this.X, this.Y}
}

func (this *Rect) Right() int32 {
	return this.X + this.Width
}

func (this *Rect) Bottom() int32 {
	return this.Y + this.Height
}

func (this *Rect) Contains(x, y int32) bool {
	return this.X <= x &&
		x < this.X+this.Width &&
		this.Y <= y &&
		y < this.Y+this.Height
}

func (this *Rect) ContainsPt(point Point) bool {
	return this.Contains(point.X, point.Y)
}

func (this *Rect) ContainsRect(rect Rect) bool {
	return (this.X <= rect.X) &&
		((rect.X + rect.Width) <= (this.X + this.Width)) &&
		(this.Y <= rect.Y) &&
		((rect.Y + rect.Height) <= (this.Y + this.Height))
}

func (this *Rect) Inflate(dx, dy int32) {
	this.X -= dx
	this.Y -= dy
	this.Width += 2 * dx
	this.Height += 2 * dy
}

func (this *Rect) InflateSize(size Size) {
	this.Inflate(size.Width, size.Height)
}

func (this *Rect) Inflated(dx, dy int32) Rect {
	rect := *this
	rect.Inflate(dx, dy)
	return rect
}

func (this *Rect) Intersect(rect Rect) {
	*this = this.Intersected(rect)
}

func (this *Rect) Intersected(rect Rect) Rect {
	x1 := max(this.X, rect.X)
	x2 := min(this.X+this.Width, rect.X+rect.Width)
	y1 := max(this.Y, rect.Y)
	y2 := min(this.Y+this.Height, rect.Y+rect.Height)

	if x2 >= x1 && y2 >= y1 {
		return Rect{x1, y1, x2 - x1, y2 - y1}
	}
	return Rect{}
}

func (this *Rect) IntersectsWith(rect Rect) bool {
	return (rect.X < this.X+this.Width) &&
		(this.X < (rect.X + rect.Width)) &&
		(rect.Y < this.Y+this.Height) &&
		(this.Y < rect.Y+rect.Height)
}

func (this *Rect) Union(rect Rect) {
	x1 := min(this.X, rect.X)
	x2 := max(this.X+this.Width, rect.X+rect.Width)
	y1 := min(this.Y, rect.Y)
	y2 := max(this.Y+this.Height, rect.Y+rect.Height)

	this.X, this.Y, this.Width, this.Height =
		x1, y1, x2-x1, y2-y1
}

func (this *Rect) Offset(dx, dy int32) {
	this.X += dx
	this.Y += dy
}

func (this *Rect) OffsetPoint(point Point) {
	this.Offset(point.X, point.Y)
}

func (this *Rect) Win32Rect() win32.RECT {
	return win32.RECT{this.X, this.Y, this.X + this.Width, this.Y + this.Height}
}

func (this *RectF) Location() PointF {
	return PointF{this.X, this.Y}
}

func (this *RectF) Contains(x, y float32) bool {
	return this.X <= x &&
		x < this.X+this.Width &&
		this.Y <= y &&
		y < this.Y+this.Height
}

func (this *RectF) ContainsPt(point PointF) bool {
	return this.Contains(point.X, point.Y)
}

func (this *RectF) ContainsRect(rect RectF) bool {
	return (this.X <= rect.X) &&
		((rect.X + rect.Width) <= (this.X + this.Width)) &&
		(this.Y <= rect.Y) &&
		((rect.Y + rect.Height) <= (this.Y + this.Height))
}

func (this *RectF) Inflate(dx, dy float32) {
	this.X -= dx
	this.Y -= dy
	this.Width += 2 * dx
	this.Height += 2 * dy
}

func (this *RectF) InflateSize(size SizeF) {
	this.Inflate(size.Width, size.Height)
}

func (this *RectF) Inflated(dx, dy float32) RectF {
	rect := *this
	rect.Inflate(dx, dy)
	return rect
}

func (this *RectF) Intersect(rect RectF) {
	*this = this.Intersected(rect)
}

func (this *RectF) Intersected(rect RectF) RectF {
	x1 := max(this.X, rect.X)
	x2 := max(this.X+this.Width, rect.X+rect.Width)
	y1 := max(this.Y, rect.Y)
	y2 := max(this.Y+this.Height, rect.Y+rect.Height)

	if x2 >= x1 && y2 >= y1 {
		return RectF{x1, y1, x2 - x1, y2 - y1}
	}
	return RectF{}
}

func (this *RectF) IntersectsWith(rect RectF) bool {
	return (rect.X < this.X+this.Width) &&
		(this.X < (rect.X + rect.Width)) &&
		(rect.Y < this.Y+this.Height) &&
		(this.Y < rect.Y+rect.Height)
}

func (this *RectF) Union(rect RectF) {
	x1 := min(this.X, rect.X)
	x2 := max(this.X+this.Width, rect.X+rect.Width)
	y1 := min(this.Y, rect.Y)
	y2 := max(this.Y+this.Height, rect.Y+rect.Height)

	this.X, this.Y, this.Width, this.Height =
		x1, y1, x2-x1, y2-y1
}

func (this *RectF) Offset(dx, dy float32) {
	this.X += dx
	this.Y += dy
}

func (this *RectF) OffsetPoint(point PointF) {
	this.Offset(point.X, point.Y)
}

func (this *RectF) Round() Rect {
	return Rect{
		int32(math.Round(float64(this.X))),
		int32(math.Round(float64(this.Y))),
		int32(math.Round(float64(this.Width))),
		int32(math.Round(float64(this.Height))),
	}
}

func (this *RectF) Ceiling() Rect {
	return Rect{
		int32(math.Ceil(float64(this.X))),
		int32(math.Ceil(float64(this.Y))),
		int32(math.Ceil(float64(this.Width))),
		int32(math.Ceil(float64(this.Height))),
	}
}

func (this *RectF) Truncate() Rect {
	return Rect{
		int32(this.X),
		int32(this.Y),
		int32(this.Width),
		int32(this.Height),
	}
}
