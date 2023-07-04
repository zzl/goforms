package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"runtime"
	"unsafe"
)

type Brush struct {
	p *gdip.Brush
}

func newBrush(s *Scope, p *gdip.Brush) *Brush {
	brush := &Brush{p}
	if s != nil {
		s.Add(brush)
	}
	runtime.SetFinalizer(brush, (*Brush).Dispose)
	return brush
}

func (this *Brush) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DeleteBrush(this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Brush) Clone(s *Scope) *Brush {
	var p2 *gdip.Brush
	gdip.CloneBrush(this.p, &p2)
	return newBrush(s, p2)
}

type SolidBrush struct {
	Brush
	color *Color
}

func newSolidBrush(s *Scope, p *gdip.SolidFill) *SolidBrush {
	brush := &SolidBrush{Brush: Brush{&p.Brush}}
	if s != nil {
		s.Add(brush)
	}
	return brush
}

func NewSolidBrush(s *Scope, color Color) *SolidBrush {
	var p *gdip.SolidFill
	status := gdip.CreateSolidFill(color.Argb(), &p)
	checkStatus(status)
	return newSolidBrush(s, p)
}

func (this *SolidBrush) Clone(s *Scope) *SolidBrush {
	var p2 *gdip.Brush
	gdip.CloneBrush(this.p, &p2)
	return newSolidBrush(s, (*gdip.SolidFill)(unsafe.Pointer(p2)))
}

func (this *SolidBrush) P() *gdip.SolidFill {
	return (*gdip.SolidFill)(unsafe.Pointer(this.p))
}

func (this *SolidBrush) SetColor(color Color) {
	if this.color != nil && *this.color == color {
		return
	}
	this.color = &color
	gdip.SetSolidFillColor(this.P(), color.Argb())
}

func (this *SolidBrush) GetColor() Color {
	if this.color != nil {
		return *this.color
	}
	var argb gdip.ARGB
	gdip.GetSolidFillColor(this.P(), &argb)
	color := ColorOf(argb)
	this.color = &color
	return color
}

func (this *SolidBrush) AsBrush() *Brush {
	return &this.Brush
}

type HatchBrush struct {
	Brush
}

func newHatchBrush(s *Scope, p *gdip.Hatch) *HatchBrush {
	brush := &HatchBrush{Brush: Brush{&p.Brush}}
	if s != nil {
		s.Add(brush)
	}
	return brush
}

func NewHatchBrush(s *Scope, style gdip.HatchStyle, foreColor, backColor Color) *HatchBrush {
	var p *gdip.Hatch
	status := gdip.CreateHatchBrush(style, foreColor.Argb(), backColor.Argb(), &p)
	checkStatus(status)
	return newHatchBrush(s, p)
}

func (this *HatchBrush) Clone(s *Scope) *HatchBrush {
	var p2 *gdip.Brush
	gdip.CloneBrush(this.p, &p2)
	return newHatchBrush(s, (*gdip.Hatch)(unsafe.Pointer(p2)))
}

func (this *HatchBrush) P() *gdip.Hatch {
	return (*gdip.Hatch)(unsafe.Pointer(this.p))
}

func (this *HatchBrush) GetHatchStyle() gdip.HatchStyle {
	var style gdip.HatchStyle
	gdip.GetHatchStyle(this.P(), &style)
	return style
}

func (this *HatchBrush) GetForeColor() Color {
	var argb gdip.ARGB
	gdip.GetHatchForegroundColor(this.P(), &argb)
	return ColorOf(argb)
}

func (this *HatchBrush) GetBackColor() Color {
	var argb gdip.ARGB
	gdip.GetHatchBackgroundColor(this.P(), &argb)
	return ColorOf(argb)
}

type LinearGradientBrush struct {
	Brush
}

func newLinearGradientBrush(s *Scope, p *gdip.LineGradient) *LinearGradientBrush {
	brush := &LinearGradientBrush{Brush: Brush{&p.Brush}}
	if s != nil {
		s.Add(brush)
	}
	return brush
}

func NewLinearGradientBrush(s *Scope, rect Rect, color1, color2 Color, angle float32,
	isAngleScalable bool, wrapMode gdip.WrapMode) *LinearGradientBrush {
	var p *gdip.LineGradient
	status := gdip.CreateLineBrushFromRectWithAngleI(
		(*gdip.Rect)(&rect), color1.Argb(), color2.Argb(), angle,
		win32.BoolToBOOL(isAngleScalable), wrapMode, &p)
	checkStatus(status)
	return newLinearGradientBrush(s, p)
}

func NewLinearGradientBrushF(s *Scope, rect RectF, color1, color2 Color, angle float32,
	isAngleScalable bool, wrapMode gdip.WrapMode) *LinearGradientBrush {
	var p *gdip.LineGradient
	status := gdip.CreateLineBrushFromRectWithAngle(
		(*gdip.RectF)(&rect), color1.Argb(), color2.Argb(), angle,
		win32.BoolToBOOL(isAngleScalable), wrapMode, &p)
	checkStatus(status)
	return newLinearGradientBrush(s, p)
}

func (this *LinearGradientBrush) P() *gdip.LineGradient {
	return (*gdip.LineGradient)(unsafe.Pointer(this.p))
}

func (this *LinearGradientBrush) GetColors() (color1, color2 Color) {
	argbs := make([]gdip.ARGB, 2)
	status := gdip.GetLineColors(this.P(), &argbs[0])
	checkStatus(status)
	return ColorOf(argbs[0]), ColorOf(argbs[1])
}

func (this *LinearGradientBrush) SetColors(color1, color2 Color) {
	status := gdip.SetLineColors(this.P(), color1.Argb(), color2.Argb())
	checkStatus(status)
}

func (this *LinearGradientBrush) GetRect() Rect {
	var rect Rect
	status := gdip.GetLineRectI(this.P(), (*gdip.Rect)(&rect))
	checkStatus(status)
	return rect
}

func (this *LinearGradientBrush) GetRectF() RectF {
	var rect RectF
	status := gdip.GetLineRect(this.P(), (*gdip.RectF)(&rect))
	checkStatus(status)
	return rect
}

func (this *LinearGradientBrush) GetGammaCorrection() bool {
	var bGamma win32.BOOL
	status := gdip.GetLineGammaCorrection(this.P(), &bGamma)
	checkStatus(status)
	return bGamma != 0
}

func (this *LinearGradientBrush) SetGammaCorrection(useGammaCorrection bool) {
	bGamma := win32.BoolToBOOL(useGammaCorrection)
	status := gdip.SetLineGammaCorrection(this.P(), bGamma)
	checkStatus(status)
}

func (this *LinearGradientBrush) GetBlendCount() int {
	var count int32
	status := gdip.GetLineBlendCount(this.P(), &count)
	checkStatus(status)
	return int(count)
}

func (this *LinearGradientBrush) GetBlend(count int) (blends []float32, positions []float32) {
	blends = make([]float32, count)
	positions = make([]float32, count)
	status := gdip.GetLineBlend(this.P(), &blends[0], &positions[0], int32(count))
	checkStatus(status)
	return
}

func (this *LinearGradientBrush) SetBlend(blends []float32, positions []float32) {
	count := len(blends)
	status := gdip.SetLineBlend(this.P(), &blends[0], &positions[0], int32(count))
	checkStatus(status)
}

func (this *LinearGradientBrush) GetInterpolationColorCount() int {
	var count int32
	status := gdip.GetLinePresetBlendCount(this.P(), &count)
	checkStatus(status)
	return int(count)
}

func (this *LinearGradientBrush) GetInterpolationColors(count int) (
	colors []Color, positions []float32) {

	argbs := make([]gdip.ARGB, count)
	positions = make([]float32, count)
	status := gdip.GetLinePresetBlend(this.P(), &argbs[0], &positions[0], int32(count))
	checkStatus(status)
	colors = make([]Color, count)
	for n := 0; n < count; n++ {
		colors[n] = ColorOf(argbs[n])
	}
	return
}

func (this *LinearGradientBrush) SetInterpolationColors(colors []Color, positions []float32) {
	count := len(colors)
	argbs := make([]gdip.ARGB, count)
	for n := 0; n < count; n++ {
		argbs[n] = colors[n].Argb()
	}
	status := gdip.SetLinePresetBlend(this.P(), &argbs[0], &positions[0], int32(count))
	checkStatus(status)
}

func (this *LinearGradientBrush) SetBlendBellShape(focus float32, scale float32) {
	status := gdip.SetLineSigmaBlend(this.P(), focus, scale)
	checkStatus(status)
}

func (this *LinearGradientBrush) SetBlendTriangularShape(focus float32, scale float32) {
	status := gdip.SetLineLinearBlend(this.P(), focus, scale)
	checkStatus(status)
}

func (this *LinearGradientBrush) GetWrapMode() gdip.WrapMode {
	var mode gdip.WrapMode
	status := gdip.GetLineWrapMode(this.P(), &mode)
	checkStatus(status)
	return gdip.WrapMode(mode)
}

func (this *LinearGradientBrush) SetWrapMode(mode gdip.WrapMode) {
	status := gdip.SetLineWrapMode(this.P(), mode)
	checkStatus(status)
}

func (this *LinearGradientBrush) GetTransform(s *Scope) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.GetLineTransform(this.P(), pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func (this *LinearGradientBrush) SetTransform(matrix *Matrix) {
	status := gdip.SetLineTransform(this.P(), matrix.p)
	checkStatus(status)
}

func (this *LinearGradientBrush) ResetTransform() {
	status := gdip.ResetLineTransform(this.P())
	checkStatus(status)
}

func (this *LinearGradientBrush) MultiplyTransform(matrix *Matrix, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.MultiplyLineTransform(this.P(), matrix.p, order)
	checkStatus(status)
}

func (this *LinearGradientBrush) TranslateTransform(dx, dy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.TranslateLineTransform(this.P(), dx, dy, order)
	checkStatus(status)
}

func (this *LinearGradientBrush) ScaleTransform(sx, sy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.ScaleLineTransform(this.P(), sx, sy, order)
	checkStatus(status)
}

func (this *LinearGradientBrush) RotateTransform(angle float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.RotateLineTransform(this.P(), angle, order)
	checkStatus(status)
}

type PathGradientBrush struct {
	Brush
}

func newPathGradientBrush(s *Scope, p *gdip.PathGradient) *PathGradientBrush {
	brush := &PathGradientBrush{Brush: Brush{&p.Brush}}
	if s != nil {
		s.Add(brush)
	}
	return brush
}

func NewPathGradientBrush(s *Scope, path *Path) *PathGradientBrush {
	var p *gdip.PathGradient
	status := gdip.CreatePathGradientFromPath(path.p, &p)
	checkStatus(status)
	return newPathGradientBrush(s, p)
}

func NewPathGradientBrushFromPoints(s *Scope, points []Point,
	wrapMode gdip.WrapMode) *PathGradientBrush {
	var p *gdip.PathGradient
	status := gdip.CreatePathGradientI(&points[0], int32(len(points)), wrapMode, &p)
	checkStatus(status)
	return newPathGradientBrush(s, p)
}

func NewPathGradientBrushFromPointsF(s *Scope, points []PointF,
	wrapMode gdip.WrapMode) *PathGradientBrush {
	var p *gdip.PathGradient
	status := gdip.CreatePathGradient(&points[0], int32(len(points)), wrapMode, &p)
	checkStatus(status)
	return newPathGradientBrush(s, p)
}

func (this *PathGradientBrush) P() *gdip.PathGradient {
	return (*gdip.PathGradient)(unsafe.Pointer(this.p))
}

func (this *PathGradientBrush) GetCenterColor() Color {
	var argb gdip.ARGB
	status := gdip.GetPathGradientCenterColor(this.P(), &argb)
	checkStatus(status)
	return ColorOf(argb)
}

func (this *PathGradientBrush) SetCenterColor(color Color) {
	status := gdip.SetPathGradientCenterColor(this.P(), color.Argb())
	checkStatus(status)
}

func (this *PathGradientBrush) GetSurroundColors() []Color {
	var count int32
	gdip.GetPathGradientSurroundColorCount(this.P(), &count)
	argbs := make([]gdip.ARGB, count)
	gdip.GetPathGradientSurroundColorsWithCount(this.P(), &argbs[0], &count)
	colors := make([]Color, count)
	for n := 0; n < int(count); n++ {
		colors[n] = ColorOf(argbs[n])
	}
	return colors
}

func (this *PathGradientBrush) SetSurroundColors(colors []Color) {
	count := len(colors)
	argbs := make([]gdip.ARGB, count)
	for n, color := range colors {
		argbs[n] = color.Argb()
	}
	count32 := int32(count)
	status := gdip.SetPathGradientSurroundColorsWithCount(
		this.P(), &argbs[0], &count32)
	checkStatus(status)
}

func (this *PathGradientBrush) GetPath(s *Scope) *Path {
	var pPath *gdip.Path
	status := gdip.GetPathGradientPath(this.P(), pPath)
	checkStatus(status)
	return newPath(s, pPath)
}

func (this *PathGradientBrush) SetPath(path *Path) {
	status := gdip.SetPathGradientPath(this.P(), path.p)
	checkStatus(status)
}

func (this *PathGradientBrush) GetCenterPoint() Point {
	var point Point
	status := gdip.GetPathGradientCenterPointI(this.P(), &point)
	checkStatus(status)
	return point
}

func (this *PathGradientBrush) SetCenterPoint(point Point) {
	status := gdip.SetPathGradientCenterPointI(this.P(), &point)
	checkStatus(status)
}

func (this *PathGradientBrush) GetCenterPointF() PointF {
	var point PointF
	status := gdip.GetPathGradientCenterPoint(this.P(), &point)
	checkStatus(status)
	return point
}

func (this *PathGradientBrush) SetCenterPointF(point PointF) {
	status := gdip.SetPathGradientCenterPoint(this.P(), &point)
	checkStatus(status)
}

func (this *PathGradientBrush) GetRectangle() Rect {
	var rect Rect
	status := gdip.GetPathGradientRectI(this.P(), (*gdip.Rect)(&rect))
	checkStatus(status)
	return rect
}

func (this *PathGradientBrush) GetRectangleF() RectF {
	var rect RectF
	status := gdip.GetPathGradientRect(this.P(), (*gdip.RectF)(&rect))
	checkStatus(status)
	return rect
}

func (this *PathGradientBrush) GetGammaCorrection() bool {
	var bGamma win32.BOOL
	status := gdip.GetPathGradientGammaCorrection(this.P(), &bGamma)
	checkStatus(status)
	return bGamma != 0
}

func (this *PathGradientBrush) SetGammaCorrection(useGammaCorrection bool) {
	bGamma := win32.BoolToBOOL(useGammaCorrection)
	status := gdip.SetPathGradientGammaCorrection(this.P(), bGamma)
	checkStatus(status)
}

func (this *PathGradientBrush) GetBlendCount() int {
	var count int32
	status := gdip.GetPathGradientBlendCount(this.P(), &count)
	checkStatus(status)
	return int(count)
}

func (this *PathGradientBrush) GetBlend(count int) (blends []float32, positions []float32) {
	blends = make([]float32, count)
	positions = make([]float32, count)
	status := gdip.GetPathGradientBlend(this.P(), &blends[0], &positions[0], int32(count))
	checkStatus(status)
	return
}

func (this *PathGradientBrush) SetBlend(blends []float32, positions []float32) {
	count := len(blends)
	status := gdip.SetPathGradientBlend(this.P(), &blends[0], &positions[0], int32(count))
	checkStatus(status)
}

func (this *PathGradientBrush) GetInterpolationColorCount() int {
	var count int32
	status := gdip.GetPathGradientPresetBlendCount(this.P(), &count)
	checkStatus(status)
	return int(count)
}

func (this *PathGradientBrush) GetInterpolationColors(count int) (
	colors []Color, positions []float32) {

	argbs := make([]gdip.ARGB, count)
	positions = make([]float32, count)
	status := gdip.GetPathGradientPresetBlend(this.P(), &argbs[0], &positions[0], int32(count))
	checkStatus(status)
	colors = make([]Color, count)
	for n := 0; n < count; n++ {
		colors[n] = ColorOf(argbs[n])
	}
	return
}

func (this *PathGradientBrush) SetInterpolationColors(colors []Color, positions []float32) {
	count := len(colors)
	argbs := make([]gdip.ARGB, count)
	for n := 0; n < count; n++ {
		argbs[n] = colors[n].Argb()
	}
	status := gdip.SetPathGradientPresetBlend(this.P(), &argbs[0], &positions[0], int32(count))
	checkStatus(status)
}

func (this *PathGradientBrush) SetBlendBellShape(focus float32, scale float32) {
	status := gdip.SetPathGradientSigmaBlend(this.P(), focus, scale)
	checkStatus(status)
}

func (this *PathGradientBrush) SetBlendTriangularShape(focus float32, scale float32) {
	status := gdip.SetPathGradientLinearBlend(this.P(), focus, scale)
	checkStatus(status)
}

func (this *PathGradientBrush) GetWrapMode() gdip.WrapMode {
	var mode gdip.WrapMode
	status := gdip.GetPathGradientWrapMode(this.P(), &mode)
	checkStatus(status)
	return gdip.WrapMode(mode)
}

func (this *PathGradientBrush) SetWrapMode(mode gdip.WrapMode) {
	status := gdip.SetPathGradientWrapMode(this.P(), mode)
	checkStatus(status)
}

func (this *PathGradientBrush) GetTransform(s *Scope) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.GetPathGradientTransform(this.P(), pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func (this *PathGradientBrush) SetTransform(matrix *Matrix) {
	status := gdip.SetPathGradientTransform(this.P(), matrix.p)
	checkStatus(status)
}

func (this *PathGradientBrush) ResetTransform() {
	status := gdip.ResetPathGradientTransform(this.P())
	checkStatus(status)
}

func (this *PathGradientBrush) MultiplyTransform(matrix *Matrix, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.MultiplyPathGradientTransform(this.P(), matrix.p, order)
	checkStatus(status)
}

func (this *PathGradientBrush) TranslateTransform(dx, dy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.TranslatePathGradientTransform(this.P(), dx, dy, order)
	checkStatus(status)
}

func (this *PathGradientBrush) ScaleTransform(sx, sy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.ScalePathGradientTransform(this.P(), sx, sy, order)
	checkStatus(status)
}

func (this *PathGradientBrush) RotateTransform(angle float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.RotatePathGradientTransform(this.P(), angle, order)
	checkStatus(status)
}

func (this *PathGradientBrush) GetFocusScales() (xScale, yScale float32) {
	status := gdip.GetPathGradientFocusScales(this.P(), &xScale, &yScale)
	checkStatus(status)
	return
}

func (this *PathGradientBrush) SetFocusScales(xScale, yScale float32) {
	status := gdip.SetPathGradientFocusScales(this.P(), xScale, yScale)
	checkStatus(status)
}

func (this *PathGradientBrush) AsBrush() *Brush {
	return &this.Brush
}

type TextureBrush struct {
	Brush
}

func newTextureBrush(s *Scope, p *gdip.Texture) *TextureBrush {
	brush := &TextureBrush{Brush: Brush{&p.Brush}}
	if s != nil {
		s.Add(brush)
	}
	return brush
}

func NewTextureBrush(s *Scope, image *Image, wrapMode gdip.WrapMode) *TextureBrush {
	var p *gdip.Texture
	status := gdip.CreateTexture(image.p, wrapMode, &p)
	checkStatus(status)
	return newTextureBrush(s, p)
}

func NewTextureBrushRect(s *Scope, image *Image, wrapMode gdip.WrapMode, rect Rect) *TextureBrush {
	var p *gdip.Texture
	status := gdip.CreateTexture2I(image.p, wrapMode,
		rect.X, rect.Y, rect.Width, rect.Height, &p)
	checkStatus(status)
	return newTextureBrush(s, p)
}

func NewTextureBrushRectAttr(s *Scope, image *Image, wrapMode gdip.WrapMode,
	rect Rect, attrs *ImageAttributes) *TextureBrush {
	var p *gdip.Texture
	status := gdip.CreateTextureIAI(image.p, attrs.p,
		rect.X, rect.Y, rect.Width, rect.Height, &p)
	checkStatus(status)
	return newTextureBrush(s, p)
}

func NewTextureBrushRectF(s *Scope, image *Image,
	wrapMode gdip.WrapMode, rect RectF) *TextureBrush {
	var p *gdip.Texture
	status := gdip.CreateTexture2(image.p, wrapMode,
		rect.X, rect.Y, rect.Width, rect.Height, &p)
	checkStatus(status)
	return newTextureBrush(s, p)
}

func NewTextureBrushRectAttrF(s *Scope, image *Image, wrapMode gdip.WrapMode,
	rect RectF, attrs *ImageAttributes) *TextureBrush {
	var p *gdip.Texture
	status := gdip.CreateTextureIA(image.p, attrs.p,
		rect.X, rect.Y, rect.Width, rect.Height, &p)
	checkStatus(status)
	return newTextureBrush(s, p)
}

func (this *TextureBrush) Clone(s *Scope) *TextureBrush {
	var p2 *gdip.Brush
	status := gdip.CloneBrush(this.p, &p2)
	checkStatus(status)
	return newTextureBrush(s, (*gdip.Texture)(unsafe.Pointer(p2)))
}

func (this *TextureBrush) P() *gdip.Texture {
	return (*gdip.Texture)(unsafe.Pointer(this.p))
}

func (this *TextureBrush) GetTransform(s *Scope) *Matrix {
	var pMatrix *gdip.Matrix
	status := gdip.GetTextureTransform(this.P(), pMatrix)
	checkStatus(status)
	return newMatrix(s, pMatrix)
}

func (this *TextureBrush) SetTransform(matrix *Matrix) {
	status := gdip.SetTextureTransform(this.P(), matrix.p)
	checkStatus(status)
}

func (this *TextureBrush) GetWrapMode() gdip.WrapMode {
	var mode gdip.WrapMode
	status := gdip.GetTextureWrapMode(this.P(), &mode)
	checkStatus(status)
	return gdip.WrapMode(mode)
}

func (this *TextureBrush) SetWrapMode(mode gdip.WrapMode) {
	status := gdip.SetTextureWrapMode(this.P(), mode)
	checkStatus(status)
}

func (this *TextureBrush) GetImage(s *Scope) *Image {
	var pImage *gdip.Image
	gdip.GetTextureImage(this.P(), &pImage)
	return newImage(s, pImage)
}

func (this *TextureBrush) ResetTransform() {
	status := gdip.ResetTextureTransform(this.P())
	checkStatus(status)
}

func (this *TextureBrush) MultiplyTransform(matrix *Matrix, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.MultiplyTextureTransform(this.P(), matrix.p, order)
	checkStatus(status)
}

func (this *TextureBrush) TranslateTransform(dx, dy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.TranslateTextureTransform(this.P(), dx, dy, order)
	checkStatus(status)
}

func (this *TextureBrush) ScaleTransform(sx, sy float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.ScaleTextureTransform(this.P(), sx, sy, order)
	checkStatus(status)
}

func (this *TextureBrush) RotateTransform(angle float32, append bool) {
	order := gdip.MatrixOrderPrepend
	if append {
		order = gdip.MatrixOrderAppend
	}
	status := gdip.RotateTextureTransform(this.P(), angle, order)
	checkStatus(status)
}
