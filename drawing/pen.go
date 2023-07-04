package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/goforms/framework/scope"
	"runtime"
)

type Pen struct {
	p *gdip.Pen
}

func newPen(s *Scope, p *gdip.Pen) *Pen {
	pen := &Pen{p: p}
	if s != nil {
		s.Add(pen)
	}
	runtime.SetFinalizer(pen, (*Pen).Dispose)
	return pen
}

func NewPen(s *scope.Scope, color Color) *Pen {
	return NewPenWithWidth(s, color, 1)
}

func NewPenWithWidth(s *scope.Scope, color Color, width float32) *Pen {
	var p *gdip.Pen
	status := gdip.CreatePen1(color.Argb(), width, gdip.UnitWorld, &p)
	checkStatus(status)
	pen := newPen(s, p)
	return pen
}

func NewPenFromBrush(s *scope.Scope, brush *Brush) *Pen {
	return NewPenFromBrushWithWidth(s, brush, 1)
}

func NewPenFromBrushWithWidth(s *scope.Scope, brush *Brush, width float32) *Pen {
	var p *gdip.Pen
	status := gdip.CreatePen2(brush.p, width, gdip.UnitWorld, &p)
	checkStatus(status)
	return newPen(s, p)
}

func (this *Pen) Clone(s *Scope) *Pen {
	var p2 *gdip.Pen
	status := gdip.ClonePen(this.p, &p2)
	checkStatus(status)
	return newPen(s, p2)
}

func (this *Pen) Dispose() {
	if this.p == nil {
		return
	}
	status := gdip.DeletePen(this.p)
	checkStatus(status)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Pen) GetWidth() float32 {
	var width float32
	status := gdip.GetPenWidth(this.p, &width)
	checkStatus(status)
	return width
}

func (this *Pen) SetWidth(width float32) {
	status := gdip.SetPenWidth(this.p, width)
	checkStatus(status)
}

func (this *Pen) SetLineCap(startCap, endCap, dashCap gdip.LineCap) {
	status := gdip.SetPenLineCap197819(this.p, startCap, endCap, gdip.DashCap(dashCap))
	checkStatus(status)
}

func (this *Pen) GetStartCap() gdip.LineCap {
	var cap gdip.LineCap
	status := gdip.GetPenStartCap(this.p, &cap)
	checkStatus(status)
	return gdip.LineCap(cap)
}

func (this *Pen) SetStartCap(cap gdip.LineCap) {
	status := gdip.SetPenStartCap(this.p, gdip.LineCap(cap))
	checkStatus(status)
}

func (this *Pen) GetEndCap() gdip.LineCap {
	var cap gdip.LineCap
	status := gdip.GetPenEndCap(this.p, &cap)
	checkStatus(status)
	return gdip.LineCap(cap)
}

func (this *Pen) SetEndCap(cap gdip.LineCap) {
	status := gdip.SetPenEndCap(this.p, gdip.LineCap(cap))
	checkStatus(status)
}

func (this *Pen) GetDashCap() gdip.LineCap {
	var cap gdip.DashCap
	status := gdip.GetPenDashCap197819(this.p, &cap)
	checkStatus(status)
	return gdip.LineCap(cap)
}

func (this *Pen) SetDashCap(cap gdip.DashCap) {
	status := gdip.SetPenDashCap197819(this.p, cap)
	checkStatus(status)
}

func (this *Pen) GetLineJoin() gdip.LineJoin {
	var join gdip.LineJoin
	status := gdip.GetPenLineJoin(this.p, &join)
	checkStatus(status)
	return gdip.LineJoin(join)
}

func (this *Pen) SetLineJoin(join gdip.LineJoin) {
	status := gdip.SetPenLineJoin(this.p, gdip.LineJoin(join))
	checkStatus(status)
}

func (this *Pen) GetCustomStartCap() *gdip.CustomLineCap {
	var cap *gdip.CustomLineCap
	status := gdip.GetPenCustomStartCap(this.p, &cap)
	checkStatus(status)
	return (*gdip.CustomLineCap)(cap)
}

func (this *Pen) SetCustomLineCap(cap *gdip.CustomLineCap) {
	status := gdip.SetPenCustomStartCap(this.p, (*gdip.CustomLineCap)(cap))
	checkStatus(status)
}

func (this *Pen) GetCustomEndCap() *gdip.CustomLineCap {
	var cap *gdip.CustomLineCap
	status := gdip.GetPenCustomEndCap(this.p, &cap)
	checkStatus(status)
	return (*gdip.CustomLineCap)(cap)
}

func (this *Pen) SetCustomEndCap(cap *gdip.CustomLineCap) {
	status := gdip.SetPenCustomEndCap(this.p, (*gdip.CustomLineCap)(cap))
	checkStatus(status)
}

func (this *Pen) GetMiterLimit() float32 {
	var limit float32
	status := gdip.GetPenMiterLimit(this.p, &limit)
	checkStatus(status)
	return limit
}

func (this *Pen) SetMiterLimit(limit float32) {
	status := gdip.SetPenMiterLimit(this.p, limit)
	checkStatus(status)
}

func (this *Pen) GetAlignment() gdip.PenAlignment {
	var align gdip.PenAlignment
	status := gdip.GetPenMode(this.p, &align)
	checkStatus(status)
	return gdip.PenAlignment(align)
}

func (this *Pen) SetAlignment(align gdip.PenAlignment) {
	status := gdip.SetPenMode(this.p, gdip.PenAlignment(align))
	checkStatus(status)
}

func (this *Pen) GetTransform(s *Scope) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.GetPenTransform(this.p, pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func (this *Pen) SetTransform(matrix *Matrix) {
	status := gdip.SetPenTransform(this.p, matrix.p)
	checkStatus(status)
}

func (this *Pen) ResetTransform() {
	status := gdip.ResetPenTransform(this.p)
	checkStatus(status)
}

func (this *Pen) MultiplyTransform(matrix *Matrix, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.MultiplyPenTransform(this.p, matrix.p, order)
	checkStatus(status)
}

func (this *Pen) TranslateTransform(dx, dy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.TranslatePenTransform(this.p, dx, dy, order)
	checkStatus(status)
}

func (this *Pen) ScaleTransform(sx, sy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.ScalePenTransform(this.p, sx, sy, order)
	checkStatus(status)
}

func (this *Pen) RotateTransform(angle float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.RotatePenTransform(this.p, angle, order)
	checkStatus(status)
}

func (this *Pen) GetPenType() gdip.PenType {
	var penType gdip.PenType
	status := gdip.GetPenFillType(this.p, &penType)
	checkStatus(status)
	return penType
}

func (this *Pen) GetColor() Color {
	var argb gdip.ARGB
	status := gdip.GetPenColor(this.p, &argb)
	checkStatus(status)
	return ColorOf(argb)
}

func (this *Pen) SetColor(color Color) {
	status := gdip.SetPenColor(this.p, color.Argb())
	checkStatus(status)
}

func (this *Pen) GetBrush(s *Scope) *Brush {
	var pBrush *gdip.Brush
	status := gdip.GetPenBrushFill(this.p, &pBrush)
	checkStatus(status)
	return newBrush(s, pBrush)
}

func (this *Pen) SetBrush(brush *Brush) {
	status := gdip.SetPenBrushFill(this.p, brush.p)
	checkStatus(status)
}

func (this *Pen) GetDashStyle() gdip.DashStyle {
	var style gdip.DashStyle
	status := gdip.GetPenDashStyle(this.p, &style)
	checkStatus(status)
	return gdip.DashStyle(style)
}

func (this *Pen) SetDashStyle(style gdip.DashStyle) {
	status := gdip.SetPenDashStyle(this.p, gdip.DashStyle(style))
	checkStatus(status)
}

func (this *Pen) GetDashOffset() float32 {
	var offset float32
	status := gdip.GetPenDashOffset(this.p, &offset)
	checkStatus(status)
	return offset
}

func (this *Pen) SetDashOffset(offset float32) {
	status := gdip.SetPenDashOffset(this.p, offset)
	checkStatus(status)
}

func (this *Pen) GetDashPattern() []float32 {
	var count32 int32
	gdip.GetPenDashCount(this.p, &count32)
	dashes := make([]float32, int(count32))
	status := gdip.GetPenDashArray(this.p, &dashes[0], count32)
	checkStatus(status)
	return dashes
}

func (this *Pen) SetDashPattern(dashes []float32) {
	status := gdip.SetPenDashArray(this.p, &dashes[0], int32(len(dashes)))
	checkStatus(status)
}

func (this *Pen) GetCompoundArray() []float32 {
	var count32 int32
	gdip.GetPenCompoundCount(this.p, &count32)
	arr := make([]float32, int(count32))
	status := gdip.GetPenCompoundArray(this.p, &arr[0], count32)
	checkStatus(status)
	return arr
}

func (this *Pen) SetCompoundArray(arr []float32) {
	status := gdip.SetPenCompoundArray(this.p, &arr[0], int32(len(arr)))
	checkStatus(status)
}
