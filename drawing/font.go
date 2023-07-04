package drawing

import (
	"errors"
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/scope"
	"runtime"
	"syscall"
	"unsafe"
)

type FontCollection struct {
	p *gdip.FontCollection
}

func (this *FontCollection) GetFontFamilies(s *Scope) []*FontFamily {
	var count int32
	status := gdip.GetFontCollectionFamilyCount(this.p, &count)
	checkStatus(status)
	pFamilies := make([]*gdip.FontFamily, count)
	var foundCount int32
	status = gdip.GetFontCollectionFamilyList(this.p, count, pFamilies, &foundCount)
	checkStatus(status)

	nFoundCount := int(foundCount)
	families := make([]*FontFamily, nFoundCount)
	for n := 0; n < nFoundCount; n++ {
		families[n] = newFontFamily(s, pFamilies[n])
	}
	return families
}

type InstalledFontCollection struct {
	FontCollection
}

func NewInstallFontCollection() *InstalledFontCollection {
	fc := &InstalledFontCollection{}
	status := gdip.NewInstalledFontCollection(&fc.p)
	checkStatus(status)
	return fc
}

type PrivateFontColleciton struct {
	FontCollection
}

func newPrivateFontCollection(s *Scope, p *gdip.FontCollection) *PrivateFontColleciton {
	col := &PrivateFontColleciton{FontCollection{p}}
	if s != nil {
		s.Add(col)
	}
	runtime.SetFinalizer(col, (*PrivateFontColleciton).Dispose)
	return col
}

func NewPrivateFontCollection(s *Scope) *PrivateFontColleciton {
	var pCol *gdip.FontCollection
	status := gdip.NewPrivateFontCollection(&pCol)
	checkStatus(status)
	return newPrivateFontCollection(s, pCol)
}

func (this *PrivateFontColleciton) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DeletePrivateFontCollection(&this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *PrivateFontColleciton) AddFontFile(filename string) {
	pwsz, _ := syscall.UTF16PtrFromString(filename)
	status := gdip.PrivateAddFontFile(this.p, pwsz)
	checkStatus(status)

	//AddFontResourceEx?
}

func (this *PrivateFontColleciton) AddMemoryFont(memory []byte) {
	status := gdip.PrivateAddMemoryFont(this.p, &memory[0], int32(len(memory)))
	checkStatus(status)
}

type FontFamily struct {
	p *gdip.FontFamily
}

type GenericFontFamilies int

const (
	GenericFontSerif GenericFontFamilies = iota
	GenericFontSansSerif
	GenericFontMonospace
)

func newFontFamily(s *Scope, p *gdip.FontFamily) *FontFamily {
	fontFamily := &FontFamily{p}
	if s != nil {
		s.Add(fontFamily)
	}
	runtime.SetFinalizer(fontFamily, (*FontFamily).Dispose)
	return fontFamily
}

func NewFontFamily(s *Scope, name string) (*FontFamily, error) {
	return NewFontFamilyWithCol(s, name, nil)
}

func NewFontFamilyOrDefault(s *Scope, name string) *FontFamily {
	family, err := NewFontFamily(s, name)
	if err != nil {
		family = GetGenericSansSerif(s)
	}
	return family
}

func NewFontFamilyGeneric(s *Scope, generic GenericFontFamilies) *FontFamily {
	switch generic {
	case GenericFontSerif:
		return GetGenericSerif(s)
	case GenericFontSansSerif:
		return GetGenericSansSerif(s)
	case GenericFontMonospace:
		return GetGenericMonospace(s)
	}
	return nil
}

func GetGenericSansSerif(s *Scope) *FontFamily {
	var pFamily *gdip.FontFamily
	status := gdip.GetGenericFontFamilySansSerif(&pFamily)
	ensureOk(status)
	return newFontFamily(s, pFamily)
}

func GetGenericSerif(s *Scope) *FontFamily {
	var pFamily *gdip.FontFamily
	status := gdip.GetGenericFontFamilySerif(&pFamily)
	ensureOk(status)
	return newFontFamily(s, pFamily)
}

func GetGenericMonospace(s *Scope) *FontFamily {
	var pFamily *gdip.FontFamily
	status := gdip.GetGenericFontFamilyMonospace(&pFamily)
	ensureOk(status)
	return newFontFamily(s, pFamily)
}

func NewFontFamilyWithCol(s *Scope, name string,
	col *FontCollection) (*FontFamily, error) {

	var pFamily *gdip.FontFamily
	pwsz, _ := syscall.UTF16PtrFromString(name)

	var pCol *gdip.FontCollection
	if col != nil {
		pCol = col.p
	}

	status := gdip.CreateFontFamilyFromName(pwsz, pCol, &pFamily)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newFontFamily(s, pFamily), nil
}

func (this *FontFamily) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DeleteFontFamily(this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *FontFamily) GetName() string {
	langId := win32.GetThreadUILanguage()
	return this.GetNameOfLang(langId)
}

func (this *FontFamily) GetNameOfLang(langId win32.LANGID) string {
	wszName := make([]uint16, 33)
	gdip.GetFamilyName(this.p, wszName[:], langId)
	return syscall.UTF16ToString(wszName)
}

func (this *FontFamily) GetEmHeight(style gdip.FontStyle) uint16 {
	var result win32.UINT16
	status := gdip.GetEmHeight(this.p, int32(style), &result)
	checkStatus(status)
	return result
}

func (this *FontFamily) GetLineSpacing(style gdip.FontStyle) uint16 {
	var result win32.UINT16
	status := gdip.GetLineSpacing(this.p, int32(style), &result)
	checkStatus(status)
	return uint16(result)
}

type Font struct {
	p *gdip.Font
}

func newFont(s *Scope, p *gdip.Font) *Font {
	font := &Font{p}
	if s != nil {
		s.Add(font)
	}
	runtime.SetFinalizer(font, (*Font).Dispose)
	return font
}

func NewFontDerived(s *Scope, prototype *Font, newStyle gdip.FontStyle) *Font {
	family := prototype.GetFontFamily(nil)
	defer family.Dispose()
	return NewFontWithUnitStyle(s, family,
		prototype.GetSize(), prototype.GetUnit(), newStyle)
}

func NewFontFromHfont(s *Scope, hFont win32.HFONT) *Font {
	var lf win32.LOGFONT
	ret := win32.GetObject(hFont,
		int32(unsafe.Sizeof(lf)), unsafe.Pointer(&lf))
	if ret == 0 {
		println("get logfont from hfont error ?")
		return nil
	}
	font, err := NewFontFromLogFont(s, &lf)
	if err != nil {
		println("new font from logfont error", err.Error())
	}
	return font
}

func NewFontFromLogFont(s *Scope, logFont *win32.LOGFONT) (*Font, error) {
	hdc := win32.GetDC(0)
	defer win32.ReleaseDC(0, hdc)
	return NewFontFromLogFontHdc(s, logFont, hdc)
}

func NewFontFromLogFontHdc(s *Scope, logFont *win32.LOGFONT, hdc win32.HDC) (*Font, error) {
	var pFont *gdip.Font
	status := gdip.CreateFontFromLogfontW(hdc, logFont, &pFont)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	font := newFont(s, pFont)
	return font, nil
}

func NewFontFromHdc(s *Scope, hdc win32.HDC) (*Font, error) {
	var pFont *gdip.Font
	status := gdip.CreateFontFromDC(hdc, &pFont)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newFont(s, pFont), nil
}

func NewFont(s *Scope, familyName string, emSize float32) (*Font, error) {
	family, err := NewFontFamily(nil, familyName)
	if err != nil {
		return nil, err
	}
	return NewFontWithFamilySize(s, family, emSize), nil
}

func NewFontWithFamilySize(s *Scope, family *FontFamily, emSize float32) *Font {
	return NewFontWithUnitStyle(s, family, emSize, gdip.UnitPoint, gdip.FontStyleRegular)
}

func NewFontWithUnit(s *Scope, family *FontFamily,
	emSize float32, unit gdip.Unit) *Font {
	return NewFontWithUnitStyle(s, family, emSize, unit, gdip.FontStyleRegular)
}

func NewFontWithStyle(s *Scope, family *FontFamily,
	emSize float32, style gdip.FontStyle) *Font {
	return NewFontWithUnitStyle(s, family, emSize, gdip.UnitPoint, style)
}

func NewFontWithUnitStyle(s *Scope, family *FontFamily, emSize float32,
	unit gdip.Unit, style gdip.FontStyle) *Font {
	var pFont *gdip.Font
	status := gdip.CreateFont(family.p, emSize, style, unit, &pFont)
	checkStatus(status)
	return newFont(s, pFont)
}

func (this *Font) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DeleteFont(this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Font) Clone(s *Scope) *Font {
	var pFont2 *gdip.Font
	status := gdip.CloneFont(this.p, &pFont2)
	checkStatus(status)
	return newFont(s, pFont2)
}

func (this *Font) GetFontFamily(s *Scope) *FontFamily {
	var pFamily *gdip.FontFamily
	status := gdip.GetFamily(this.p, &pFamily)
	checkStatus(status)
	return newFontFamily(s, pFamily)
}

func (this *Font) GetName() string {
	fam := this.GetFontFamily(nil)
	name := fam.GetName()
	fam.Dispose()
	return name
}

func (this *Font) ToHfont() (win32.HFONT, error) {
	lf := this.ToLogFont()
	hFont := win32.CreateFontIndirect(lf)
	if hFont == 0 {
		return 0, errors.New("")
	}
	return hFont, nil
}

func (this *Font) ToLogFont() *win32.LOGFONT {
	var lf win32.LOGFONT
	hdc := win32.GetDC(0)
	g, _ := NewGraphicsFromHdc(nil, hdc)
	gdip.GetLogFontW(this.p, g.p, &lf)
	g.Dispose()
	return &lf
}

func (this *Font) GetHeight(g *Graphics) float32 {
	var height float32
	status := gdip.GetFontHeight(this.p, g.p, &height)
	checkStatus(status)
	return height
}

func (this *Font) GetStyle() gdip.FontStyle {
	var style int32
	status := gdip.GetFontStyle(this.p, &style)
	checkStatus(status)
	return gdip.FontStyle(style)
}

func (this *Font) GetSize() float32 {
	var size float32
	status := gdip.GetFontSize(this.p, &size)
	checkStatus(status)
	return size
}

func (this *Font) GetPointSize() float32 {
	if this.GetUnit() == gdip.UnitPoint {
		return this.GetSize()
	}
	s := scope.NewScope()
	defer s.Leave()
	fam := this.GetFontFamily(s)
	var emHeightInPoints float32
	hdc := win32.GetDC(0)
	g, _ := NewGraphicsFromHdc(s, hdc)
	pixelsPerPoint := g.GetDpiY() / 72
	lineSpacingInPixels := this.GetHeight(g)
	style := this.GetStyle()
	emHeightInPixels := lineSpacingInPixels *
		float32(fam.GetEmHeight(style)) / float32(fam.GetLineSpacing(style))
	emHeightInPoints = emHeightInPixels / pixelsPerPoint
	return emHeightInPoints
}

func (this *Font) GetUnit() gdip.Unit {
	var unit gdip.Unit
	status := gdip.GetFontUnit(this.p, &unit)
	checkStatus(status)
	return unit
}
