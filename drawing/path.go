package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"runtime"
	"syscall"
	"unsafe"
)

type Path struct {
	p *gdip.Path
}

type PathData struct {
	Count  int32
	Points []PointF
	Types  []byte
}

func newPath(s *Scope, p *gdip.Path) *Path {
	path := &Path{p}
	if s != nil {
		s.Add(path)
	}
	runtime.SetFinalizer(path, (*Path).Dispose)
	return path
}

func NewPath(s *Scope) *Path {
	return NewPathWithMode(s, gdip.FillModeAlternate)
}

func NewPathWithMode(s *Scope, fillMode gdip.FillMode) *Path {
	var pPath *gdip.Path
	status := gdip.CreatePath(fillMode, &pPath)
	checkStatus(status)
	return newPath(s, pPath)
}

func NewPathWithPoints(s *Scope, pts []Point, types []byte) *Path {
	return NewPathWithPointsMode(s, pts, types, gdip.FillModeAlternate)
}

func NewPathWithPointsMode(s *Scope, pts []Point, types []byte, fillMode gdip.FillMode) *Path {
	var pPath *gdip.Path
	status := gdip.CreatePath2I(&pts[0], &types[0],
		int32(len(pts)), fillMode, &pPath)
	checkStatus(status)
	return newPath(s, pPath)
}

func NewPathWithPointsF(s *Scope, pts []PointF, types []byte) *Path {
	return NewPathWithPointsModeF(s, pts, types, gdip.FillModeAlternate)
}

func NewPathWithPointsModeF(s *Scope, pts []PointF, types []byte, fillMode gdip.FillMode) *Path {
	var pPath *gdip.Path
	status := gdip.CreatePath2(&pts[0], &types[0],
		int32(len(pts)), fillMode, &pPath)
	checkStatus(status)
	return newPath(s, pPath)
}

func (this *Path) Dispose() {
	if this.p == nil {
		return
	}
	status := gdip.DeletePath(this.p)
	checkStatus(status)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Path) Clone(s *Scope) *Path {
	var pPath2 *gdip.Path
	gdip.ClonePath(this.p, &pPath2)
	return newPath(s, pPath2)
}

func (this *Path) Reset() {
	status := gdip.ResetPath(this.p)
	checkStatus(status)
}

func (this *Path) GetFillMode() gdip.FillMode {
	var mode gdip.FillMode
	status := gdip.GetPathFillMode(this.p, &mode)
	checkStatus(status)
	return mode
}

func (this *Path) SetFillMode(mode gdip.FillMode) {
	status := gdip.SetPathFillMode(this.p, mode)
	checkStatus(status)
}

func (this *Path) GetPathData() PathData {
	count := this.GetPointCount()
	var data PathData
	data.Points = make([]PointF, count)
	data.Types = make([]byte, count)

	memPathData := struct {
		count  int32
		points *PointF
		types  *byte
	}{
		count:  count,
		points: &data.Points[0],
		types:  &data.Types[0],
	}
	status := gdip.GetPathData(this.p, (*gdip.PathData)(unsafe.Pointer(&memPathData)))
	checkStatus(status)
	return data
}

func (this *Path) StartFigure() {
	status := gdip.StartPathFigure(this.p)
	checkStatus(status)
}

func (this *Path) CloseFigure() {
	status := gdip.ClosePathFigure(this.p)
	checkStatus(status)
}

func (this *Path) CloseAllFigures() {
	status := gdip.ClosePathFigures(this.p)
	checkStatus(status)
}

func (this *Path) SetMarkers() {
	status := gdip.SetPathMarker(this.p)
	checkStatus(status)
}

func (this *Path) ClearMarkers() {
	status := gdip.ClearPathMarkers(this.p)
	checkStatus(status)
}

func (this *Path) Reverse() {
	status := gdip.ReversePath(this.p)
	checkStatus(status)
}

func (this *Path) GetLastPoint() PointF {
	var point PointF
	status := gdip.GetPathLastPoint(this.p, &point)
	checkStatus(status)
	return point
}

func (this *Path) IsVisible(x, y int32, g *Graphics) bool {
	return this.IsVisibleInGraph(x, y, nil)
}

func (this *Path) IsVisibleF(x, y float32, g *Graphics) bool {
	return this.IsVisibleInGraphF(x, y, nil)
}

func (this *Path) IsVisibleInGraph(x, y int32, g *Graphics) bool {
	var pGraphics *gdip.Graphics
	if g != nil {
		pGraphics = g.p
	}
	var result win32.BOOL
	status := gdip.IsVisiblePathPointI(this.p, x, y, pGraphics, &result)
	checkStatus(status)
	return result != 0
}

func (this *Path) IsVisibleInGraphF(x, y float32, g *Graphics) bool {
	var pGraphics *gdip.Graphics
	if g != nil {
		pGraphics = g.p
	}
	var result win32.BOOL
	status := gdip.IsVisiblePathPoint(this.p, x, y, pGraphics, &result)
	checkStatus(status)
	return result != 0
}

func (this *Path) IsOutlineVisible(x, y int32, pen *Pen) bool {
	return this.IsOutlineVisibleInGraph(x, y, pen, nil)
}

func (this *Path) IsOutlineVisibleF(x, y float32, pen *Pen) bool {
	return this.IsOutlineVisibleInGraphF(x, y, pen, nil)
}

func (this *Path) IsOutlineVisibleInGraph(x, y int32, pen *Pen, g *Graphics) bool {
	var result win32.BOOL
	var pGraphics *gdip.Graphics
	if g != nil {
		pGraphics = g.p
	}
	status := gdip.IsOutlineVisiblePathPointI(this.p, x, y, pen.p, pGraphics, &result)
	checkStatus(status)
	return result != 0
}

func (this *Path) IsOutlineVisibleInGraphF(x, y float32, pen *Pen, g *Graphics) bool {
	var result win32.BOOL
	var pGraphics *gdip.Graphics
	if g != nil {
		pGraphics = g.p
	}
	status := gdip.IsOutlineVisiblePathPoint(this.p, x, y, pen.p, pGraphics, &result)
	checkStatus(status)
	return result != 0
}

func (this *Path) AddLine(pt1, pt2 Point) {
	status := gdip.AddPathLineI(this.p, pt1.X, pt1.Y, pt2.X, pt2.Y)
	checkStatus(status)
}

func (this *Path) AddLineF(pt1, pt2 PointF) {
	status := gdip.AddPathLine(this.p, pt1.X, pt1.Y, pt2.X, pt2.Y)
	checkStatus(status)
}

func (this *Path) AddLines(points []Point) {
	status := gdip.AddPathLine2I(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddLinesF(points []PointF) {
	status := gdip.AddPathLine2(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddArc(rect Rect, startAngle, sweepAngle float32) {
	status := gdip.AddPathArcI(this.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Path) AddArcF(rect RectF, startAngle, sweepAngle float32) {
	status := gdip.AddPathArc(this.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Path) AddBezier(p1, p2, p3, p4 Point) {
	status := gdip.AddPathBezierI(this.p, p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, p4.X, p4.Y)
	checkStatus(status)
}

func (this *Path) AddBezierF(p1, p2, p3, p4 PointF) {
	status := gdip.AddPathBezier(this.p, p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y, p4.X, p4.Y)
	checkStatus(status)
}

func (this *Path) AddBeziers(points []Point) {
	status := gdip.AddPathBeziersI(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddBeziersF(points []PointF) {
	status := gdip.AddPathBeziers(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddCurve(points []Point) {
	status := gdip.AddPathCurveI(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddCurveF(points []PointF) {
	status := gdip.AddPathCurve(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddCurveTension(points []Point, tension float32) {
	status := gdip.AddPathCurve2I(this.p, &points[0], int32(len(points)), tension)
	checkStatus(status)
}

func (this *Path) AddCurveTensionF(points []PointF, tension float32) {
	status := gdip.AddPathCurve2(this.p, &points[0], int32(len(points)), tension)
	checkStatus(status)
}

func (this *Path) AddClosedCurve(points []Point) {
	status := gdip.AddPathClosedCurveI(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddClosedCurveF(points []PointF) {
	status := gdip.AddPathClosedCurve(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddClosedCurveTension(points []Point, tension float32) {
	status := gdip.AddPathClosedCurve2I(this.p, &points[0], int32(len(points)), tension)
	checkStatus(status)
}

func (this *Path) AddClosedCurveTensionF(points []PointF, tension float32) {
	status := gdip.AddPathClosedCurve2(this.p, &points[0], int32(len(points)), tension)
	checkStatus(status)
}

func (this *Path) AddRectangle(rect Rect) {
	status := gdip.AddPathRectangleI(this.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Path) AddRectangleF(rect RectF) {
	status := gdip.AddPathRectangle(this.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Path) AddRectangles(rects []Rect) {
	status := gdip.AddPathRectanglesI(this.p, (*gdip.Rect)(&rects[0]), int32(len(rects)))
	checkStatus(status)
}

func (this *Path) AddRectanglesF(rects []RectF) {
	status := gdip.AddPathRectangles(this.p, (*gdip.RectF)(&rects[0]), int32(len(rects)))
	checkStatus(status)
}

func (this *Path) AddEllipse(rect Rect) {
	status := gdip.AddPathEllipseI(this.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Path) AddEllipseF(rect RectF) {
	status := gdip.AddPathEllipse(this.p, rect.X, rect.Y, rect.Width, rect.Height)
	checkStatus(status)
}

func (this *Path) AddPie(rect Rect, startAngle, sweepAngle float32) {
	status := gdip.AddPathPieI(this.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Path) AddPieF(rect RectF, startAngle, sweepAngle float32) {
	status := gdip.AddPathPie(this.p, rect.X, rect.Y,
		rect.Width, rect.Height, startAngle, sweepAngle)
	checkStatus(status)
}

func (this *Path) AddPolygon(points []Point) {
	status := gdip.AddPathPolygonI(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddPolygonF(points []PointF) {
	status := gdip.AddPathPolygon(this.p, &points[0], int32(len(points)))
	checkStatus(status)
}

func (this *Path) AddPath(addingPath *Path, connect bool) {
	status := gdip.AddPathPath(this.p, addingPath.p, win32.BoolToBOOL(connect))
	checkStatus(status)
}

func (this *Path) AddString(s string, family *FontFamily,
	style gdip.FontStyle, emSize float32, origin Point, format *StringFormat) {
	wsz, _ := syscall.UTF16FromString(s)
	rect := Rect{origin.X, origin.Y, 0, 0}
	status := gdip.AddPathStringI(this.p, &wsz[0], int32(len(wsz)-1), family.p,
		style, emSize, (*gdip.Rect)(&rect), format.p)
	checkStatus(status)
}

func (this *Path) AddStringF(s string, family *FontFamily,
	style gdip.FontStyle, emSize float32, origin PointF, format *StringFormat) {
	wsz, _ := syscall.UTF16FromString(s)
	rect := RectF{origin.X, origin.Y, 0, 0}
	status := gdip.AddPathString(this.p, &wsz[0], int32(len(wsz)-1), family.p,
		style, emSize, (*gdip.RectF)(&rect), format.p)
	checkStatus(status)
}

func (this *Path) AddStringRect(s string, family *FontFamily,
	style gdip.FontStyle, emSize float32, layoutRect Rect, format *StringFormat) {
	wsz, _ := syscall.UTF16FromString(s)
	status := gdip.AddPathStringI(this.p, &wsz[0], int32(len(wsz)-1), family.p,
		style, emSize, (*gdip.Rect)(&layoutRect), format.p)
	checkStatus(status)
}

func (this *Path) AddStringRectF(s string, family *FontFamily,
	style gdip.FontStyle, emSize float32, layoutRect RectF, format *StringFormat) {
	wsz, _ := syscall.UTF16FromString(s)
	status := gdip.AddPathString(this.p, &wsz[0], int32(len(wsz)-1), family.p,
		style, emSize, (*gdip.RectF)(&layoutRect), format.p)
	checkStatus(status)
}

func (this *Path) Transform(matrix *Matrix) {
	status := gdip.TransformPath(this.p, matrix.p)
	checkStatus(status)
}

func (this *Path) GetBounds(matrix *Matrix, pen *Pen) Rect {
	var bounds Rect
	var pPen *gdip.Pen
	if pen != nil {
		pPen = pen.p
	}
	var pMatrix *gdip.Matrix
	if matrix != nil {
		pMatrix = matrix.p
	}
	gdip.GetPathWorldBoundsI(this.p, (*gdip.Rect)(&bounds), pMatrix, pPen)
	return bounds
}

func (this *Path) GetBoundsF(matrix *Matrix, pen *Pen) RectF {
	var bounds RectF
	var pPen *gdip.Pen
	if pen != nil {
		pPen = pen.p
	}
	var pMatrix *gdip.Matrix
	if matrix != nil {
		pMatrix = matrix.p
	}
	gdip.GetPathWorldBounds(this.p, (*gdip.RectF)(&bounds), pMatrix, pPen)
	return bounds
}

func (this *Path) Flatten(matrix *Matrix, flatness float32) {
	var pMatrix *gdip.Matrix
	if matrix != nil {
		pMatrix = matrix.p
	}
	status := gdip.FlattenPath(this.p, pMatrix, flatness)
	checkStatus(status)
}

func (this *Path) Widen(pen *Pen, matrix *Matrix, flatness float32) {
	var pMatrix *gdip.Matrix
	if matrix != nil {
		pMatrix = matrix.p
	}
	status := gdip.WidenPath(this.p, pen.p, pMatrix, flatness)
	checkStatus(status)
}

func (this *Path) Warp(dstPoints []PointF, srcRect RectF) {
	this.WarpWithMatrixModeFlat(dstPoints, srcRect, nil, gdip.WarpModePerspective, 0.25)
}

func (this *Path) WarpWithMatrix(dstPoints []PointF, srcRect RectF, matrix *Matrix) {
	this.WarpWithMatrixModeFlat(dstPoints, srcRect, matrix, gdip.WarpModePerspective, 0.25)
}

func (this *Path) WarpWithMatrixMode(dstPoints []PointF, srcRect RectF,
	matrix *Matrix, warpMode gdip.WarpMode) {
	this.WarpWithMatrixModeFlat(dstPoints, srcRect, matrix, warpMode, 0.25)
}

func (this *Path) WarpWithMatrixModeFlat(dstPoints []PointF, srcRect RectF,
	matrix *Matrix, warpMode gdip.WarpMode, flatness float32) {
	var pMatrix *gdip.Matrix
	if matrix != nil {
		pMatrix = matrix.p
	}
	status := gdip.WarpPath(this.p, pMatrix, &dstPoints[0], int32(len(dstPoints)),
		srcRect.X, srcRect.Y, srcRect.Width, srcRect.Height, warpMode, flatness)
	checkStatus(status)
}

func (this *Path) GetPointCount() int32 {
	var count int32
	status := gdip.GetPointCount(this.p, &count)
	checkStatus(status)
	return count
}

func (this *Path) GetPathTypes() []byte {
	count := this.GetPointCount()
	types := make([]byte, count)
	status := gdip.GetPathTypes(this.p, (*win32.BYTE)(&types[0]), count)
	checkStatus(status)
	return types
}

func (this *Path) GetPathPoints() []Point {
	count := this.GetPointCount()
	points := make([]Point, count)
	status := gdip.GetPathPointsI(this.p, &points[0], count)
	checkStatus(status)
	return points
}

func (this *Path) GetPathPointsF() []PointF {
	count := this.GetPointCount()
	points := make([]PointF, count)
	status := gdip.GetPathPoints(this.p, &points[0], count)
	checkStatus(status)
	return points
}
