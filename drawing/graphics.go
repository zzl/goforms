package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"log"
	"runtime"
	"syscall"
)

type Graphics struct {
	p     *gdip.Graphics
	hdc   win32.HDC
	image *Image
}

func newGraphics(s *Scope, p *gdip.Graphics) *Graphics {
	g := &Graphics{p: p}
	if s != nil {
		s.Add(g)
	}
	runtime.SetFinalizer(g, (*Graphics).Dispose)
	return g
}

func NewGraphicsFromHdc(s *Scope, hdc win32.HDC) (*Graphics, error) {
	var pGraphics *gdip.Graphics
	status := gdip.CreateFromHDC(hdc, &pGraphics)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newGraphics(s, pGraphics), nil
}

func NewGraphicsFromHwnd(s *Scope, hWnd win32.HWND) (*Graphics, error) {
	var pGraphics *gdip.Graphics
	status := gdip.CreateFromHWND(hWnd, &pGraphics)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newGraphics(s, pGraphics), nil
}

func NewGraphicsFromBitmap(s *Scope, bitmap *Bitmap) (*Graphics, error) {
	return NewGraphicsFromImage(s, bitmap.AsImage())
}

func NewGraphicsFromImage(s *Scope, image *Image) (*Graphics, error) {
	var pGraphics *gdip.Graphics
	status := gdip.GetImageGraphicsContext(image.p, &pGraphics)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	g := newGraphics(s, pGraphics)
	g.image = image
	return g, nil
}

func (this *Graphics) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DeleteGraphics(this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Graphics) GetHdc() win32.HDC {
	var hdc win32.HDC
	status := gdip.GetDC(this.p, &hdc)
	checkStatus(status)
	this.hdc = hdc
	return this.hdc
}

func (this *Graphics) ReleaseHdc() {
	status := gdip.ReleaseDC(this.p, this.hdc)
	checkStatus(status)
}

func (this *Graphics) Flush() {
	status := gdip.Flush(this.p, gdip.FlushIntentionFlush)
	checkStatus(status)
}

func (this *Graphics) GetCompositingMode() int {
	var mode gdip.CompositingMode
	status := gdip.GetCompositingMode(this.p, &mode)
	checkStatus(status)
	return int(mode)
}

func (this *Graphics) SetCompositingMode(mode gdip.CompositingMode) {
	status := gdip.SetCompositingMode(this.p, mode)
	checkStatus(status)
}

func (this *Graphics) GetRenderingOrigin() Point {
	var x, y int32
	status := gdip.GetRenderingOrigin(this.p, &x, &y)
	checkStatus(status)
	return Point{x, y}
}

func (this *Graphics) SetRenderingOrigin(pt Point) {
	status := gdip.SetRenderingOrigin(this.p, pt.X, pt.Y)
	checkStatus(status)
}

func (this *Graphics) GetCompositingQuality() int {
	var quality gdip.CompositingQuality
	status := gdip.GetCompositingQuality(this.p, &quality)
	checkStatus(status)
	return int(quality)
}

func (this *Graphics) SetCompositingQuality(quality gdip.CompositingQuality) {
	status := gdip.SetCompositingQuality(this.p, quality)
	checkStatus(status)
}

func (this *Graphics) GetTextRenderingHint() int {
	var hint gdip.TextRenderingHint
	status := gdip.GetTextRenderingHint(this.p, &hint)
	checkStatus(status)
	return int(hint)
}

func (this *Graphics) SetTextRenderingHint(hint gdip.TextRenderingHint) {
	status := gdip.SetTextRenderingHint(this.p, hint)
	checkStatus(status)
}

func (this *Graphics) GetTextContrast() uint {
	var contrast uint32
	status := gdip.GetTextContrast(this.p, &contrast)
	checkStatus(status)
	return uint(contrast)
}

func (this *Graphics) SetTextContrast(contrast uint) {
	status := gdip.SetTextContrast(this.p, uint32(contrast))
	checkStatus(status)
}

func (this *Graphics) GetSmoothingMode() int {
	var mode gdip.SmoothingMode
	status := gdip.GetSmoothingMode(this.p, &mode)
	checkStatus(status)
	return int(mode)
}

func (this *Graphics) SetSmoothingMode(mode gdip.SmoothingMode) {
	status := gdip.SetSmoothingMode(this.p, mode)
	checkStatus(status)
}

func (this *Graphics) GetPixelOffsetMode() int {
	var mode gdip.PixelOffsetMode
	status := gdip.GetPixelOffsetMode(this.p, &mode)
	checkStatus(status)
	return int(mode)
}

func (this *Graphics) SetPixelOffsetMode(mode gdip.PixelOffsetMode) {
	status := gdip.SetPixelOffsetMode(this.p, mode)
	checkStatus(status)
}

func (this *Graphics) GetInterpolationMode() int {
	var mode gdip.InterpolationMode
	status := gdip.GetInterpolationMode(this.p, &mode)
	checkStatus(status)
	return int(mode)
}

func (this *Graphics) SetInterpolationMode(mode gdip.InterpolationMode) {
	status := gdip.SetInterpolationMode(this.p, mode)
	checkStatus(status)
}

func (this *Graphics) GetTransform(s *Scope) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.GetWorldTransform(this.p, pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func (this *Graphics) SetTransform(matrix *Matrix) {
	status := gdip.SetWorldTransform(this.p, matrix.p)
	checkStatus(status)
}

func (this *Graphics) GetPageUnit() gdip.Unit {
	var unit gdip.Unit
	status := gdip.GetPageUnit(this.p, &unit)
	checkStatus(status)
	return unit
}

func (this *Graphics) SetPageUnit(unit gdip.Unit) {
	status := gdip.SetPageUnit(this.p, unit)
	checkStatus(status)
}

func (this *Graphics) GetPageScale() float32 {
	var scale float32
	status := gdip.GetPageScale(this.p, &scale)
	checkStatus(status)
	return scale
}

func (this *Graphics) SetPageScale(scale float32) {
	status := gdip.SetPageScale(this.p, scale)
	checkStatus(status)
}

func (this *Graphics) GetDpiX() float32 {
	var result float32
	status := gdip.GetDpiX(this.p, &result)
	checkStatus(status)
	return result
}

func (this *Graphics) GetDpiY() float32 {
	var result float32
	status := gdip.GetDpiY(this.p, &result)
	checkStatus(status)
	return result
}

func (this *Graphics) CopyFromScreen(srcUpperLeft Point, dstUpperLeft Point, size Size, rop int) error {
	hdcScreen := win32.GetDC(0)
	hdc := this.GetHdc()
	defer func() {
		win32.ReleaseDC(0, hdcScreen)
		this.ReleaseHdc()
	}()
	ok, errno := win32.BitBlt(hdc, dstUpperLeft.X, dstUpperLeft.Y,
		size.Width, size.Height, hdcScreen,
		srcUpperLeft.X, srcUpperLeft.Y, win32.ROP_CODE(rop))
	if ok == win32.FALSE {
		return errno
	}
	return nil
}

func (this *Graphics) ResetTransform() {
	status := gdip.ResetWorldTransform(this.p)
	checkStatus(status)
}

func (this *Graphics) MultiplyTransform(matrix *Matrix, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.MultiplyWorldTransform(this.p, matrix.p, order)
	checkStatus(status)
}

func (this *Graphics) TranslateTransform(dx float32, dy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.TranslateWorldTransform(this.p, dx, dy, order)
	checkStatus(status)
}

func (this *Graphics) ScaleTransform(sx float32, sy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.ScaleWorldTransform(this.p, sx, sy, order)
	checkStatus(status)
}

func (this *Graphics) RotateTransform(angle float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.RotateWorldTransform(this.p, angle, order)
	checkStatus(status)
}

func (this *Graphics) TransformPointsF(dstSpace gdip.CoordinateSpace,
	srcSpace gdip.CoordinateSpace, pts []PointF) {
	status := gdip.TransformPoints(this.p, dstSpace,
		srcSpace, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) TransformPoints(dstSpace gdip.CoordinateSpace,
	srcSpace gdip.CoordinateSpace, pts []Point) {
	status := gdip.TransformPointsI(this.p, dstSpace,
		srcSpace, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) DrawLine(pen *Pen, x1, y1, x2, y2 int32) {
	status := gdip.DrawLineI(this.p, pen.p, x1, y1, x2, y2)
	checkStatus(status)
}

func (this *Graphics) DrawLinePt(pen *Pen, pt1, pt2 Point) {
	this.DrawLine(pen, pt1.X, pt1.Y, pt2.X, pt2.Y)
}

func (this *Graphics) DrawLineF(pen *Pen, x1, y1, x2, y2 float32) {
	status := gdip.DrawLine(this.p, pen.p, x1, y1, x2, y2)
	checkStatus(status)
}

func (this *Graphics) DrawLinePtF(pen *Pen, pt1, pt2 PointF) {
	this.DrawLineF(pen, pt1.X, pt1.Y, pt2.X, pt2.Y)
}

func (this *Graphics) DrawLines(pen *Pen, points []Point) {
	status := gdip.DrawLinesI(this.p, pen.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Graphics) DrawLinesF(pen *Pen, points []PointF) {
	status := gdip.DrawLines(this.p, pen.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Graphics) DrawArc(pen *Pen, x, y, w, h int32, startAngle, sweepAngle float32) {
	status := gdip.DrawArcI(this.p, pen.p, x, y, w, h, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawArcF(pen *Pen, x, y, w, h, startAngle, sweepAngle float32) {
	status := gdip.DrawArc(this.p, pen.p, x, y, w, h, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawArcRect(pen *Pen, rect Rect, startAngle, sweepAngle float32) {
	status := gdip.DrawArcI(this.p, pen.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawArcRectF(pen *Pen, rect RectF, startAngle, sweepAngle float32) {
	status := gdip.DrawArc(this.p, pen.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawBezier(pen *Pen, pt1, pt2, pt3, pt4 Point) {
	status := gdip.DrawBezierI(this.p, pen.p, pt1.X, pt1.Y,
		pt2.X, pt2.Y, pt3.X, pt3.Y, pt4.X, pt4.Y)
	checkStatus(status)
}

func (this *Graphics) DrawBezierF(pen *Pen, pt1, pt2, pt3, pt4 PointF) {
	status := gdip.DrawBezier(this.p, pen.p, pt1.X, pt1.Y,
		pt2.X, pt2.Y, pt3.X, pt3.Y, pt4.X, pt4.Y)
	checkStatus(status)
}

func (this *Graphics) DrawBeziers(pen *Pen, points []Point) {
	status := gdip.DrawBeziersI(this.p, pen.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Graphics) DrawBeziersF(pen *Pen, points []PointF) {
	status := gdip.DrawBeziers(this.p, pen.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Graphics) DrawRectangle(pen *Pen, x, y, width, height int32) {
	status := gdip.DrawRectangleI(this.p, pen.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) DrawRectangleF(pen *Pen, x, y, width, height float32) {
	status := gdip.DrawRectangle(this.p, pen.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) DrawRectangleRect(pen *Pen, rect Rect) {
	status := gdip.DrawRectangleI(this.p, pen.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) DrawRectangleRectF(pen *Pen, rect RectF) {
	status := gdip.DrawRectangle(this.p, pen.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) DrawRectangles(pen *Pen, rects []Rect) {
	status := gdip.DrawRectanglesI(this.p, pen.p,
		(*gdip.Rect)(&rects[0]), int32(len(rects)))
	checkStatus(status)
}

func (this *Graphics) DrawRectanglesF(pen *Pen, rects []RectF) {
	status := gdip.DrawRectangles(this.p, pen.p,
		(*gdip.RectF)(&rects[0]), int32(len(rects)))
	checkStatus(status)
}

func (this *Graphics) DrawEllipse(pen *Pen, x, y, width, height int32) {
	status := gdip.DrawEllipseI(this.p, pen.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) DrawEllipseF(pen *Pen, x, y, width, height float32) {
	status := gdip.DrawEllipse(this.p, pen.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) DrawEllipseRect(pen *Pen, rect Rect) {
	status := gdip.DrawEllipseI(this.p, pen.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) DrawEllipseRectF(pen *Pen, rect RectF) {
	status := gdip.DrawEllipse(this.p, pen.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) DrawPie(pen *Pen, x, y, width, height int32, startAngle, sweepAngle float32) {
	status := gdip.DrawPieI(this.p, pen.p, x, y, width, height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawPieF(pen *Pen, x, y, width, height, startAngle, sweepAngle float32) {
	status := gdip.DrawPie(this.p, pen.p, x, y, width, height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawPieRect(pen *Pen, rect Rect, startAngle, sweepAngle float32) {
	status := gdip.DrawPieI(this.p, pen.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawPieRectF(pen *Pen, rect RectF, startAngle, sweepAngle float32) {
	status := gdip.DrawPie(this.p, pen.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) DrawPolygon(pen *Pen, pts []Point) {
	status := gdip.DrawPolygonI(this.p, pen.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) DrawPolygonF(pen *Pen, pts []PointF) {
	status := gdip.DrawPolygon(this.p, pen.p, &pts[0], int32(len(pts)))
	if status != gdip.Ok {
		log.Fatal(status)
	}
}

func (this *Graphics) DrawPath(pen *Pen, path *Path) {
	status := gdip.DrawPath(this.p, pen.p, path.p)
	checkStatus(status)
}

func (this *Graphics) DrawCurve(pen *Pen, pts []Point) {
	status := gdip.DrawCurveI(this.p, pen.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) DrawCurveF(pen *Pen, pts []PointF) {
	status := gdip.DrawCurve(this.p, pen.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) DrawCurve2(pen *Pen, pts []Point, tension float32) {
	status := gdip.DrawCurve2I(this.p, pen.p, &pts[0], int32(len(pts)), tension)
	checkStatus(status)
}

func (this *Graphics) DrawCurve2F(pen *Pen, pts []PointF, tension float32) {
	status := gdip.DrawCurve2(this.p, pen.p, &pts[0], int32(len(pts)), tension)
	checkStatus(status)
}

func (this *Graphics) DrawClosedCurve(pen *Pen, pts []Point) {
	status := gdip.DrawClosedCurveI(this.p, pen.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) DrawClosedCurveF(pen *Pen, pts []PointF) {
	status := gdip.DrawClosedCurve(this.p, pen.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) DrawClosedCurve2(pen *Pen, pts []Point, tension float32) {
	status := gdip.DrawClosedCurve2I(this.p, pen.p, &pts[0], int32(len(pts)), tension)
	checkStatus(status)
}

func (this *Graphics) DrawClosedCurve2F(pen *Pen, pts []PointF, tension float32) {
	status := gdip.DrawClosedCurve2(this.p, pen.p, &pts[0], int32(len(pts)), tension)
	checkStatus(status)
}

func (this *Graphics) Clear(color Color) {
	status := gdip.GraphicsClear(this.p, color.Argb())
	checkStatus(status)
}

func (this *Graphics) FillRectangle(brush *Brush, x, y, width, height int32) {
	status := gdip.FillRectangleI(this.p, brush.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) FillRectangleF(brush *Brush, x, y, width, height float32) {
	status := gdip.FillRectangle(this.p, brush.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) FillRectangleRect(brush *Brush, rect Rect) {
	status := gdip.FillRectangleI(this.p, brush.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) FillRectangleRectF(brush *Brush, rect RectF) {
	status := gdip.FillRectangle(this.p, brush.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) FillRectangles(brush *Brush, rects []Rect) {
	status := gdip.FillRectanglesI(this.p, brush.p,
		(*gdip.Rect)(&rects[0]), int32(len(rects)))
	checkStatus(status)
}

func (this *Graphics) FillRectanglesF(brush *Brush, rects []RectF) {
	status := gdip.FillRectangles(this.p, brush.p,
		(*gdip.RectF)(&rects[0]), int32(len(rects)))
	checkStatus(status)
}

func (this *Graphics) FillPolygon(brush *Brush, pts []Point, fillWinding bool) {
	fillMode := gdip.FillModeAlternate
	if fillWinding {
		fillMode = gdip.FillModeWinding
	}
	status := gdip.FillPolygonI(this.p, brush.p, &pts[0], int32(len(pts)), fillMode)
	checkStatus(status)
}

func (this *Graphics) FillPolygonF(brush *Brush, pts []PointF, fillWinding bool) {
	fillMode := gdip.FillModeAlternate
	if fillWinding {
		fillMode = gdip.FillModeWinding
	}
	status := gdip.FillPolygon(this.p, brush.p, &pts[0], int32(len(pts)), fillMode)
	checkStatus(status)
}

func (this *Graphics) FillEllipse(brush *Brush, x, y, width, height int32) {
	status := gdip.FillEllipseI(this.p, brush.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) FillEllipseF(brush *Brush, x, y, width, height float32) {
	status := gdip.FillEllipse(this.p, brush.p, x, y, width, height)
	checkStatus(status)
}

func (this *Graphics) FillEllipseRect(brush *Brush, rect Rect) {
	status := gdip.FillEllipseI(this.p, brush.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) FillEllipseRectF(brush *Brush, rect RectF) {
	status := gdip.FillEllipse(this.p, brush.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Graphics) FillPie(brush *Brush, x, y, w, h int32, startAngle, sweepAngle float32) {
	status := gdip.FillPieI(this.p, brush.p, x, y, w, h, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) FillPieF(brush *Brush, x, y, w, h, startAngle, sweepAngle float32) {
	status := gdip.FillPie(this.p, brush.p, x, y, w, h, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) FillPieRect(brush *Brush, rect Rect, startAngle, sweepAngle float32) {
	status := gdip.FillPieI(this.p, brush.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) FillPieRectF(brush *Brush, rect RectF, startAngle, sweepAngle float32) {
	status := gdip.FillPie(this.p, brush.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Graphics) FillPath(brush *Brush, path *Path) {
	status := gdip.FillPath(this.p, brush.p, path.p)
	checkStatus(status)
}

func (this *Graphics) FillClosedCurve(brush *Brush, pts []Point) {
	status := gdip.FillClosedCurveI(this.p, brush.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) FillClosedCurveF(brush *Brush, pts []PointF) {
	status := gdip.FillClosedCurve(this.p, brush.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) FillClosedCurve2(brush *Brush, pts []Point, tension float32, fillWinding bool) {
	fillMode := gdip.FillModeAlternate
	if fillWinding {
		fillMode = gdip.FillModeWinding
	}
	status := gdip.FillClosedCurve2I(this.p, brush.p, &pts[0], int32(len(pts)), tension, fillMode)
	checkStatus(status)
}

func (this *Graphics) FillClosedCurve2F(brush *Brush, pts []PointF, tension float32, fillWinding bool) {
	fillMode := gdip.FillModeAlternate
	if fillWinding {
		fillMode = gdip.FillModeWinding
	}
	status := gdip.FillClosedCurve2(this.p, brush.p, &pts[0], int32(len(pts)), tension, fillMode)
	checkStatus(status)
}

func (this *Graphics) FillRegion(brush *Brush, region *Region) {
	status := gdip.FillRegion(this.p, brush.p, region.p)
	checkStatus(status)
}

func (this *Graphics) DrawString(str string, font *Font, brush *Brush,
	x int, y int, format *StringFormat) {
	this.DrawStringRectF(str, font, brush, RectF{X: float32(x), Y: float32(y)}, format)
}

func (this *Graphics) DrawStringF(str string, font *Font, brush *Brush,
	x float32, y float32, format *StringFormat) {
	this.DrawStringRectF(str, font, brush, RectF{X: x, Y: y}, format)
}

func (this *Graphics) DrawStringRect(str string, font *Font, brush *Brush,
	layoutRect Rect, format *StringFormat) {
	this.DrawStringRectF(str, font, brush, toRectF(layoutRect), format)
}

func (this *Graphics) DrawStringRectF(str string, font *Font, brush *Brush,
	layoutRect RectF, format *StringFormat) {
	wsz, _ := syscall.UTF16FromString(str)
	var pFormat *gdip.StringFormat
	if format != nil {
		pFormat = format.p
	}
	status := gdip.DrawString(this.p, &wsz[0], int32(len(wsz)-1),
		font.p, (*gdip.RectF)(&layoutRect), pFormat, brush.p)
	checkStatus(status)
}

func (this *Graphics) MeasureString(str string, font *Font, layoutArea SizeF, format *StringFormat) (
	size SizeF, charsFitted, linesFilled int32) {
	wsz, _ := syscall.UTF16FromString(str)
	layoutRect := RectF{0, 0, layoutArea.Width, layoutArea.Height}
	var boundBox RectF
	var pFormat *gdip.StringFormat
	if format != nil {
		pFormat = format.p
	}
	status := gdip.MeasureString(this.p, &wsz[0], int32(len(wsz)-1),
		font.p, (*gdip.RectF)(&layoutRect), pFormat,
		(*gdip.RectF)(&boundBox), &charsFitted, &linesFilled)
	checkStatus(status)
	size = SizeF{boundBox.Width, boundBox.Height}
	return
}

func (this *Graphics) MeasureCharacterRanges(str string, font *Font,
	layoutRect RectF, format *StringFormat) []*Region {
	wsz, _ := syscall.UTF16FromString(str)
	var pFormat *gdip.StringFormat
	if format != nil {
		pFormat = format.p
	}

	var count int32
	gdip.GetStringFormatMeasurableCharacterRangeCount(pFormat, &count)

	pRegions := make([]*gdip.Region, count)
	status := gdip.MeasureCharacterRanges(this.p, &wsz[0], int32(len(wsz)-1),
		font.p, (*gdip.RectF)(&layoutRect), pFormat, count, &pRegions[0])
	checkStatus(status)

	regions := make([]*Region, count)
	for n := 0; n < int(count); n++ {
		regions[n] = &Region{pRegions[n]}
	}
	return regions
}

func (this *Graphics) DrawIcon(icon *Icon, x, y int32) {
	if this.image != nil {
		this.DrawBitmap(icon.ToBitmap(), x, y)
	} else {
		icon.Draw(this, x, y)
	}
}

func (this *Graphics) DrawIconRect(icon *Icon, targetRect Rect) {
	if this.image != nil {
		this.DrawBitmapRect(icon.ToBitmap(),
			targetRect.X, targetRect.Y, targetRect.Width, targetRect.Height)
	} else {
		icon.DrawRect(this, targetRect)
	}
}

func (this *Graphics) DrawImage(image *Image, x int32, y int32) {
	status := gdip.DrawImageI(this.p, image.p, x, y)
	checkStatus(status)
}

func (this *Graphics) DrawImageF(image *Image, x float32, y float32) {
	status := gdip.DrawImage(this.p, image.p, x, y)
	checkStatus(status)
}

func (this *Graphics) DrawImageRect(image *Image, x, y, w, h int32) {
	gdip.DrawImageRectI(this.p, image.p, x, y, w, h)
}

func (this *Graphics) DrawImageRectF(image *Image, x, y, w, h float32) {
	gdip.DrawImageRect(this.p, image.p, x, y, w, h)
}

func (this *Graphics) DrawImageRectRect(image *Image,
	dstRect Rect, srcRect Rect, srcUnit gdip.Unit) {
	this.DrawImageRectRectAttr(image, dstRect, srcRect, srcUnit, nil)
}

func (this *Graphics) DrawImageRectRectAttr(image *Image,
	dstRect Rect, srcRect Rect, srcUnit gdip.Unit, attr *ImageAttributes) {
	var pAttr *gdip.ImageAttributes
	if attr != nil {
		pAttr = attr.p
	}
	status := gdip.DrawImageRectRectI(this.p, image.p,
		dstRect.X, dstRect.Y,
		dstRect.Width, dstRect.Height,
		srcRect.X, srcRect.Y,
		srcRect.Width, srcRect.Height,
		srcUnit, pAttr, 0, nil)
	checkStatus(status)
}

func (this *Graphics) DrawImageRectRectF(image *Image,
	dstRect RectF, srcRect RectF, srcUnit gdip.Unit) {
	this.DrawImageRectRectAttrF(image, dstRect, srcRect, srcUnit, nil)
}

func (this *Graphics) DrawImageRectRectAttrF(image *Image,
	dstRect RectF, srcRect RectF, srcUnit gdip.Unit, attr *ImageAttributes) {
	var pAttr *gdip.ImageAttributes
	if attr != nil {
		pAttr = attr.p
	}
	status := gdip.DrawImageRectRect(this.p, image.p,
		dstRect.X, dstRect.Y, dstRect.Width, dstRect.Height,
		srcRect.X, srcRect.Y, srcRect.Width, srcRect.Height,
		srcUnit, pAttr, 0, nil)
	checkStatus(status)
}

func (this *Graphics) DrawImagePointRect(image *Image,
	dstPt Point, srcRect Rect, srcUnit gdip.Unit) {
	status := gdip.DrawImagePointRectI(this.p, image.p,
		dstPt.X, dstPt.Y, srcRect.X, srcRect.Y,
		srcRect.Width, srcRect.Height, srcUnit)
	checkStatus(status)
}

func (this *Graphics) DrawImagePointRectF(image *Image,
	dstPt PointF, srcRect RectF, srcUnit gdip.Unit) {
	status := gdip.DrawImagePointRect(this.p, image.p, dstPt.X, dstPt.Y,
		srcRect.X, srcRect.Y, srcRect.Width, srcRect.Height, srcUnit)
	checkStatus(status)
}

func (this *Graphics) DrawImagePointsRect(image *Image,
	dstPts []Point, srcRect Rect, srcUnit gdip.Unit) {
	status := gdip.DrawImagePointsRectI(this.p, image.p, &dstPts[0],
		int32(len(dstPts)), srcRect.X, srcRect.Y, srcRect.Width, srcRect.Height,
		srcUnit, nil, 0, nil)
	checkStatus(status)
}

func (this *Graphics) DrawImagePointsRectF(image *Image,
	dstPts []PointF, srcRect RectF, srcUnit gdip.Unit) {
	status := gdip.DrawImagePointsRect(this.p, image.p, &dstPts[0],
		int32(len(dstPts)), srcRect.X, srcRect.Y, srcRect.Width, srcRect.Height,
		srcUnit, nil, 0, nil)
	checkStatus(status)
}

func (this *Graphics) SetClip(rect Rect) {
	this.SetClipRect(rect, gdip.CombineModeReplace)
}

func (this *Graphics) SetClipF(rect RectF) {
	this.SetClipRectF(rect, gdip.CombineModeReplace)
}

func (this *Graphics) SetClipRect(rect Rect, combineMode gdip.CombineMode) {
	status := gdip.SetClipRectI(this.p, rect.X, rect.Y,
		rect.Width, rect.Height, combineMode)
	checkStatus(status)
}

func (this *Graphics) SetClipRectF(rect RectF, combineMode gdip.CombineMode) {
	status := gdip.SetClipRect(this.p, rect.X, rect.Y,
		rect.Width, rect.Height, combineMode)
	checkStatus(status)
}

func (this *Graphics) SetClipGraphics(g *Graphics, combineMode gdip.CombineMode) {
	status := gdip.SetClipGraphics(this.p, g.p, combineMode)
	checkStatus(status)
}

func (this *Graphics) SetClipPath(path *Path, combineMode gdip.CombineMode) {
	status := gdip.SetClipPath(this.p, path.p, combineMode)
	checkStatus(status)
}

func (this *Graphics) SetClipRegion(region *Region, combineMode gdip.CombineMode) {
	status := gdip.SetClipRegion(this.p, region.p, combineMode)
	checkStatus(status)
}

func (this *Graphics) IntersectClip(rect Rect) {
	this.SetClipRect(rect, gdip.CombineModeIntersect)
}

func (this *Graphics) IntersectClipF(rect RectF) {
	this.SetClipRectF(rect, gdip.CombineModeIntersect)
}

func (this *Graphics) IntersectClipRegion(region *Region) {
	this.SetClipRegion(region, gdip.CombineModeIntersect)
}

func (this *Graphics) ExcludeClip(rect Rect) {
	this.SetClipRect(rect, gdip.CombineModeExclude)
}

func (this *Graphics) ExcludeClipRegion(region *Region) {
	this.SetClipRegion(region, gdip.CombineModeExclude)
}

func (this *Graphics) ResetClip() {
	status := gdip.ResetClip(this.p)
	checkStatus(status)
}

func (this *Graphics) TranslateClip(dx, dy int32) {
	status := gdip.TranslateClipI(this.p, dx, dy)
	checkStatus(status)
}

func (this *Graphics) TranslateClipF(dx, dy float32) {
	status := gdip.TranslateClip(this.p, dx, dy)
	checkStatus(status)
}

func (this *Graphics) GetClipBounds() Rect {
	var rect Rect
	status := gdip.GetClipBoundsI(this.p, (*gdip.Rect)(&rect))
	checkStatus(status)
	return rect
}

func (this *Graphics) GetClipBoundsF() RectF {
	var rect RectF
	status := gdip.GetClipBounds(this.p, (*gdip.RectF)(&rect))
	checkStatus(status)
	return rect
}

func (this *Graphics) GetClip() *Region {
	region := &Region{}
	status := gdip.GetClip(this.p, region.p)
	checkStatus(status)
	return region
}

func (this *Graphics) IsClipEmpty() bool {
	var result win32.BOOL
	status := gdip.IsClipEmpty(this.p, &result)
	checkStatus(status)
	return result == win32.TRUE
}

func (this *Graphics) GetVisibleClipBounds() Rect {
	var rect Rect
	status := gdip.GetVisibleClipBoundsI(this.p, (*gdip.Rect)(&rect))
	checkStatus(status)
	return rect
}

func (this *Graphics) GetVisibleClipBoundsF() RectF {
	var rect RectF
	status := gdip.GetVisibleClipBounds(this.p, (*gdip.RectF)(&rect))
	checkStatus(status)
	return rect
}

func (this *Graphics) IsVisibleClipEmpty() bool {
	var result win32.BOOL
	status := gdip.IsVisibleClipEmpty(this.p, &result)
	checkStatus(status)
	return result == win32.TRUE
}

func (this *Graphics) IsVisiblePoint(pt Point) bool {
	var result win32.BOOL
	status := gdip.IsVisiblePointI(this.p, pt.X, pt.Y, &result)
	checkStatus(status)
	return result == win32.TRUE
}

func (this *Graphics) IsVisiblePointF(pt PointF) bool {
	var result win32.BOOL
	status := gdip.IsVisiblePoint(this.p, pt.X, pt.Y, &result)
	checkStatus(status)
	return result == win32.TRUE
}

func (this *Graphics) IsVisibleRect(rect Rect) bool {
	var result win32.BOOL
	status := gdip.IsVisibleRectI(this.p, rect.X, rect.Y, rect.Width, rect.Height, &result)
	checkStatus(status)
	return result == win32.TRUE
}

func (this *Graphics) IsVisibleRectF(rect RectF) bool {
	var result win32.BOOL
	status := gdip.IsVisibleRect(this.p, rect.X, rect.Y, rect.Width, rect.Height, &result)
	checkStatus(status)
	return result == win32.TRUE
}

func (this *Graphics) Save() gdip.GraphicsState {
	var state gdip.GraphicsState
	status := gdip.SaveGraphics(this.p, &state)
	checkStatus(status)
	return state
}

func (this *Graphics) Restore(state gdip.GraphicsState) {
	status := gdip.RestoreGraphics(this.p, state)
	checkStatus(status)
}

func (this *Graphics) BeginContainer() gdip.GraphicsContainer {
	var state gdip.GraphicsContainer
	status := gdip.BeginContainer2(this.p, &state)
	checkStatus(status)
	return state
}

func (this *Graphics) BeginContainerRect(dstRect Rect,
	srcRect Rect, unit gdip.Unit) gdip.GraphicsContainer {
	var state gdip.GraphicsContainer
	status := gdip.BeginContainerI(this.p, (*gdip.Rect)(&dstRect), (*gdip.Rect)(&srcRect), unit, &state)
	checkStatus(status)
	return state
}

func (this *Graphics) BeginContainerRectF(dstRect RectF,
	srcRect RectF, unit gdip.Unit) gdip.GraphicsContainer {
	var state gdip.GraphicsContainer
	status := gdip.BeginContainer(this.p, (*gdip.RectF)(&dstRect),
		(*gdip.RectF)(&srcRect), unit, &state)
	checkStatus(status)
	return state
}

func (this *Graphics) EndContainer(container gdip.GraphicsContainer) {
	status := gdip.EndContainer(this.p, container)
	checkStatus(status)
}

func (this *Graphics) Translate(dx float32, dy float32) {
	status := gdip.TranslateWorldTransform(this.p, dx, dy, gdip.MatrixOrderPrepend)
	checkStatus(status)
}

func (this *Graphics) DrawBitmap(bitmap *Bitmap, x int32, y int32) {
	status := gdip.DrawImageI(this.p, bitmap.Image.p, x, y)
	checkStatus(status)
}

func (this *Graphics) DrawBitmapRect(bitmap *Bitmap, x, y, cx, cy int32) {
	status := gdip.DrawImageRectI(this.p, bitmap.Image.p, x, y, cx, cy)
	checkStatus(status)
}

func (this *Graphics) DrawRect(pen *Pen, x, y, w, h float32) {
	status := gdip.DrawRectangle(this.p, pen.p, x, y, w, h)
	checkStatus(status)
}

func (this *Graphics) DrawCircle(pen *Pen, x, y, w, h float32) {
	status := gdip.DrawEllipse(this.p, pen.p, x, y, w, h)
	checkStatus(status)
}

func (this *Graphics) Scale(ratio float32) {
	status := gdip.ScaleWorldTransform(this.p, ratio, ratio, gdip.MatrixOrderPrepend)
	checkStatus(status)
}

func (this *Graphics) DrawPolyline(pen *Pen, pts []Point) {
	status := gdip.DrawLinesI(this.p, pen.p, &pts[0], int32(len(pts)))
	checkStatus(status)
}

func (this *Graphics) SetSmoothMode(smooth bool) {
	mode := gdip.SmoothingModeDefault
	if smooth {
		mode = gdip.SmoothingModeHighQuality
	}
	status := gdip.SetSmoothingMode(this.p, mode)
	checkStatus(status)
}
