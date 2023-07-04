package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"runtime"
)

type Region struct {
	p *gdip.Region
}

func newRegion(s *Scope, p *gdip.Region) *Region {
	region := &Region{p}
	if s != nil {
		s.Add(region)
	}
	runtime.SetFinalizer(region, (*Region).Dispose)
	return region
}

func NewRegion(s *Scope) *Region {
	var pRegion *gdip.Region
	status := gdip.CreateRegion(&pRegion)
	checkStatus(status)
	return newRegion(s, pRegion)
}

func NewRegionOfRect(s *Scope, rect Rect) *Region {
	var pRegion *gdip.Region
	status := gdip.CreateRegionRectI((*gdip.Rect)(&rect), &pRegion)
	checkStatus(status)
	return newRegion(s, pRegion)
}

func NewRegionOfRectF(s *Scope, rect RectF) *Region {
	var pRegion *gdip.Region
	status := gdip.CreateRegionRect((*gdip.RectF)(&rect), &pRegion)
	checkStatus(status)
	return newRegion(s, pRegion)
}

func NewRegionOfPath(s *Scope, path *Path) *Region {
	var pRegion *gdip.Region
	status := gdip.CreateRegionPath(path.p, &pRegion)
	checkStatus(status)
	return newRegion(s, pRegion)
}

func NewRegionOfData(s *Scope, data []byte) *Region {
	var pRegion *gdip.Region
	status := gdip.CreateRegionRgnData(
		(*win32.BYTE)(&data[0]), int32(len(data)), &pRegion)
	checkStatus(status)
	return newRegion(s, pRegion)
}

func NewRegionFromHrgn(s *Scope, hRgn win32.HRGN) *Region {
	var pRegion *gdip.Region
	status := gdip.CreateRegionHrgn(hRgn, &pRegion)
	checkStatus(status)
	return newRegion(s, pRegion)
}

func (this *Region) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DeleteRegion(this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Region) Clone(s *Scope) *Region {
	var pRegion2 *gdip.Region
	status := gdip.CloneRegion(this.p, &pRegion2)
	checkStatus(status)
	return newRegion(s, pRegion2)
}

func (this *Region) MakeInfinite() {
	status := gdip.SetInfinite(this.p)
	checkStatus(status)
}

func (this *Region) MakeEmpty() {
	status := gdip.SetEmpty(this.p)
	checkStatus(status)
}

func (this *Region) InterceptRect(rect Rect) {
	status := gdip.CombineRegionRectI(this.p,
		(*gdip.Rect)(&rect), gdip.CombineModeIntersect)
	checkStatus(status)
}

func (this *Region) InterceptRectF(rect RectF) {
	status := gdip.CombineRegionRect(this.p,
		(*gdip.RectF)(&rect), gdip.CombineModeIntersect)
	checkStatus(status)
}

func (this *Region) InterceptPath(path *Path) {
	status := gdip.CombineRegionPath(this.p, path.p, gdip.CombineModeIntersect)
	checkStatus(status)
}

func (this *Region) InterceptRegion(region *Region) {
	status := gdip.CombineRegionRegion(this.p, region.p, gdip.CombineModeIntersect)
	checkStatus(status)
}

//func (this *Region) ReleaseHrgn() {
//}

func (this *Region) UnionRect(rect Rect) {
	status := gdip.CombineRegionRectI(this.p,
		(*gdip.Rect)(&rect), gdip.CombineModeUnion)
	checkStatus(status)
}

func (this *Region) UnionRectF(rect RectF) {
	status := gdip.CombineRegionRect(this.p,
		(*gdip.RectF)(&rect), gdip.CombineModeUnion)
	checkStatus(status)
}

func (this *Region) UnionPath(path *Path) {
	status := gdip.CombineRegionPath(this.p, path.p, gdip.CombineModeUnion)
	checkStatus(status)
}

func (this *Region) UnionRegion(region *Region) {
	status := gdip.CombineRegionRegion(this.p, region.p, gdip.CombineModeUnion)
	checkStatus(status)
}

func (this *Region) XorRect(rect Rect) {
	status := gdip.CombineRegionRectI(this.p,
		(*gdip.Rect)(&rect), gdip.CombineModeXor)
	checkStatus(status)
}

func (this *Region) XorRectF(rect RectF) {
	status := gdip.CombineRegionRect(this.p,
		(*gdip.RectF)(&rect), gdip.CombineModeXor)
	checkStatus(status)
}

func (this *Region) XorPath(path *Path) {
	status := gdip.CombineRegionPath(this.p, path.p, gdip.CombineModeXor)
	checkStatus(status)
}

func (this *Region) XorRegion(region *Region) {
	status := gdip.CombineRegionRegion(this.p, region.p, gdip.CombineModeXor)
	checkStatus(status)
}

func (this *Region) ExcludeRect(rect Rect) {
	status := gdip.CombineRegionRectI(this.p,
		(*gdip.Rect)(&rect), gdip.CombineModeExclude)
	checkStatus(status)
}

func (this *Region) ExcludeRectF(rect RectF) {
	status := gdip.CombineRegionRect(this.p,
		(*gdip.RectF)(&rect), gdip.CombineModeExclude)
	checkStatus(status)
}

func (this *Region) ExcludePath(path *Path) {
	status := gdip.CombineRegionPath(this.p, path.p, gdip.CombineModeExclude)
	checkStatus(status)
}

func (this *Region) ExcludeRegion(region *Region) {
	status := gdip.CombineRegionRegion(this.p, region.p, gdip.CombineModeExclude)
	checkStatus(status)
}

func (this *Region) ComplementRect(rect Rect) {
	status := gdip.CombineRegionRectI(this.p,
		(*gdip.Rect)(&rect), gdip.CombineModeComplement)
	checkStatus(status)
}

func (this *Region) ComplementRectF(rect RectF) {
	status := gdip.CombineRegionRect(this.p,
		(*gdip.RectF)(&rect), gdip.CombineModeComplement)
	checkStatus(status)
}

func (this *Region) ComplementPath(path *Path) {
	status := gdip.CombineRegionPath(this.p, path.p, gdip.CombineModeComplement)
	checkStatus(status)
}

func (this *Region) ComplementRegion(region *Region) {
	status := gdip.CombineRegionRegion(this.p, region.p, gdip.CombineModeComplement)
	checkStatus(status)
}

func (this *Region) CombineRect(rect Rect, mode gdip.CombineMode) {
	status := gdip.CombineRegionRectI(this.p, (*gdip.Rect)(&rect), mode)
	checkStatus(status)
}

func (this *Region) CombineRectF(rect RectF, mode gdip.CombineMode) {
	status := gdip.CombineRegionRect(this.p, (*gdip.RectF)(&rect), mode)
	checkStatus(status)
}

func (this *Region) CombinePath(path *Path, mode gdip.CombineMode) {
	status := gdip.CombineRegionPath(this.p, path.p, mode)
	checkStatus(status)
}

func (this *Region) CombineRegion(region *Region, mode gdip.CombineMode) {
	status := gdip.CombineRegionRegion(this.p, region.p, mode)
	checkStatus(status)
}

func (this *Region) Translate(dx, dy int32) {
	status := gdip.TranslateRegionI(this.p, dx, dy)
	checkStatus(status)
}

func (this *Region) TranslateF(dx, dy float32) {
	status := gdip.TranslateRegion(this.p, dx, dy)
	checkStatus(status)
}

func (this *Region) Transform(matrix *Matrix) {
	status := gdip.TransformRegion(this.p, matrix.p)
	checkStatus(status)
}

func (this *Region) GetBounds(g *Graphics) Rect {
	var rect Rect
	status := gdip.GetRegionBoundsI(this.p, g.p, (*gdip.Rect)(&rect))
	checkStatus(status)
	return rect
}

func (this *Region) GetBoundsF(g *Graphics) RectF {
	var rect RectF
	status := gdip.GetRegionBounds(this.p, g.p, (*gdip.RectF)(&rect))
	checkStatus(status)
	return rect
}

func (this *Region) GetHrgn(g *Graphics) win32.HRGN {
	var hrgn win32.HRGN
	status := gdip.GetRegionHRgn(this.p, g.p, &hrgn)
	checkStatus(status)
	return hrgn
}

func (this *Region) IsEmpty(g *Graphics) bool {
	var result win32.BOOL
	status := gdip.IsEmptyRegion(this.p, g.p, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) IsInfinite(g *Graphics) bool {
	var result win32.BOOL
	status := gdip.IsInfiniteRegion(this.p, g.p, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) Equals(region *Region, g *Graphics) bool {
	var result win32.BOOL
	status := gdip.IsEqualRegion(this.p, region.p, g.p, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) GetRegionData() []byte {
	var cb uint32
	status := gdip.GetRegionDataSize(this.p, &cb)
	checkStatus(status)
	data := make([]byte, cb)
	var cbFilled uint32
	status = gdip.GetRegionData(this.p, (*win32.BYTE)(&data[0]), cb, &cbFilled)
	checkStatus(status)
	return data[:cbFilled]
}

func (this *Region) IsVisiblePoint(point Point) bool {
	var result win32.BOOL
	status := gdip.IsVisibleRegionPointI(this.p, point.X, point.Y, nil, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) IsVisiblePointF(point PointF) bool {
	var result win32.BOOL
	status := gdip.IsVisibleRegionPoint(this.p, point.X, point.Y, nil, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) IsVisiblePointG(point Point, g *Graphics) bool {
	var result win32.BOOL
	status := gdip.IsVisibleRegionPointI(this.p, point.X, point.Y, g.p, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) IsVisiblePointGF(point PointF, g *Graphics) bool {
	var result win32.BOOL
	status := gdip.IsVisibleRegionPoint(this.p, point.X, point.Y, g.p, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) IsVisibleRect(rect Rect) bool {
	return this.IsVisibleRectG(rect, nil)
}

func (this *Region) IsVisibleRectF(rect RectF) bool {
	return this.IsVisibleRectGF(rect, nil)
}

func (this *Region) IsVisibleRectG(rect Rect, g *Graphics) bool {
	var pGraphics *gdip.Graphics
	if g != nil {
		pGraphics = g.p
	}
	var result win32.BOOL
	status := gdip.IsVisibleRegionRectI(this.p,
		rect.X, rect.Y, rect.Width, rect.Height, pGraphics, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) IsVisibleRectGF(rect RectF, g *Graphics) bool {
	var pGraphics *gdip.Graphics
	if g != nil {
		pGraphics = g.p
	}
	var result win32.BOOL
	status := gdip.IsVisibleRegionRect(this.p,
		rect.X, rect.Y, rect.Width, rect.Height, pGraphics, &result)
	checkStatus(status)
	return result != 0
}

func (this *Region) GetRegionScans(matrix *Matrix) []Rect {
	var count uint32
	status := gdip.GetRegionScansCount(this.p, &count, matrix.p)
	checkStatus(status)
	rects := make([]Rect, count)
	nCount := int32(count)
	status = gdip.GetRegionScansI(this.p, (*gdip.Rect)(&rects[0]), &nCount, matrix.p)
	checkStatus(status)
	return rects[:nCount]
}

func (this *Region) GetRegionScansF(matrix *Matrix) []RectF {
	var count uint32
	status := gdip.GetRegionScansCount(this.p, &count, matrix.p)
	checkStatus(status)
	rects := make([]RectF, count)
	nCount := int32(count)
	status = gdip.GetRegionScans(this.p, (*gdip.RectF)(&rects[0]), &nCount, matrix.p)
	checkStatus(status)
	return rects[:nCount]
}
