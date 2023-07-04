package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"runtime"
)

type Matrix struct {
	p *gdip.Matrix
}

func newMatrix(s *Scope, pMatrix *gdip.Matrix) *Matrix {
	matrix := &Matrix{p: pMatrix}
	if s != nil {
		s.Add(matrix)
	}
	runtime.SetFinalizer(matrix, (*Matrix).Dispose)
	return matrix
}

func NewMatrix(s *Scope) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.CreateMatrix(&pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func NewMatrixWithValues(s *Scope, m11, m12, m21, m22, dx, dy float32) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.CreateMatrix2(m11, m12, m21, m22, dx, dy, &pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func NewMatrixFromRectPoints(s *Scope, rect Rect, plgpts []Point) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.CreateMatrix3I((*gdip.Rect)(&rect), &plgpts[0], &pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func NewMatrixFromRectPointsF(s *Scope, rect RectF, plgpts []PointF) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.CreateMatrix3((*gdip.RectF)(&rect), &plgpts[0], &pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func (this *Matrix) Dispose() {
	if this.p == nil {
		return
	}
	status := gdip.DeleteMatrix(this.p)
	checkStatus(status)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Matrix) Clone(s *Scope) *Matrix {
	var pMatrix2 *gdip.Matrix
	status := gdip.CloneMatrix(this.p, &pMatrix2)
	checkStatus(status)
	return newMatrix(s, pMatrix2)
}

func (this *Matrix) GetElements() []float32 {
	elems := make([]float32, 6)
	status := gdip.GetMatrixElements(this.p, &elems[0])
	checkStatus(status)
	return elems
}

func (this *Matrix) GetOffsetX() float32 {
	return this.GetElements()[4]
}

func (this *Matrix) GetOffsetY() float32 {
	return this.GetElements()[5]
}

func (this *Matrix) Reset() {
	status := gdip.SetMatrixElements(this.p, 1, 0, 0, 1, 0, 0)
	checkStatus(status)
}

func (this *Matrix) Multiply(matrix *Matrix, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderPrepend
	}
	status := gdip.MultiplyMatrix(this.p, matrix.p, order)
	checkStatus(status)
}

func (this *Matrix) Translate(offsetX, offsetY float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderPrepend
	}
	status := gdip.TranslateMatrix(this.p, offsetX, offsetY, order)
	checkStatus(status)
}

func (this *Matrix) Scale(scaleX, scaleY float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderPrepend
	}
	status := gdip.ScaleMatrix(this.p, scaleX, scaleY, order)
	checkStatus(status)
}

func (this *Matrix) Rotate(angle float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderPrepend
	}
	status := gdip.RotateMatrix(this.p, angle, order)
	checkStatus(status)
}

func (this *Matrix) RotateAt(angle float32, point PointF, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderPrepend
	}
	var status gdip.Status
	if order == gdip.MatrixOrderPrepend {
		status = gdip.TranslateMatrix(this.p, point.X, point.Y, order)
		status |= gdip.RotateMatrix(this.p, angle, order)
		status |= gdip.TranslateMatrix(this.p, -point.X, -point.Y, order)
	} else {
		status |= gdip.TranslateMatrix(this.p, -point.X, -point.Y, order)
		status |= gdip.RotateMatrix(this.p, angle, order)
		status = gdip.TranslateMatrix(this.p, point.X, point.Y, order)
	}
	checkStatus(status)
}

func (this *Matrix) Shear(shearX, shearY float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderPrepend
	}
	status := gdip.ShearMatrix(this.p, shearX, shearY, order)
	checkStatus(status)
}

func (this *Matrix) Invert() {
	status := gdip.InvertMatrix(this.p)
	checkStatus(status)
}

func (this *Matrix) TransformPoints(pts []Point) {
	status := gdip.TransformMatrixPointsI(this.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Matrix) TransformPointsF(pts []PointF) {
	status := gdip.TransformMatrixPoints(this.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Matrix) TransformVectors(pts []Point) {
	status := gdip.VectorTransformMatrixPointsI(this.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Matrix) TransformVectorsF(pts []PointF) {
	status := gdip.VectorTransformMatrixPoints(this.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Matrix) IsInvertible() bool {
	var b win32.BOOL
	status := gdip.IsMatrixInvertible(this.p, &b)
	checkStatus(status)
	return b != 0
}

func (this *Matrix) IsIdentity() bool {
	var b win32.BOOL
	status := gdip.IsMatrixIdentity(this.p, &b)
	checkStatus(status)
	return b != 0
}

func (this *Matrix) Equals(other *Matrix) bool {
	var b win32.BOOL
	status := gdip.IsMatrixEqual(this.p, other.p, &b)
	checkStatus(status)
	return b != 0
}
