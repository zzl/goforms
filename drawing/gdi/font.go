package gdi

import (
	"errors"
	"github.com/zzl/go-win32api/v2/win32"
	"log"
	"syscall"
)

type Font struct {
	Handle    win32.HFONT
	Name      string
	Size      int //pixel
	PointSize int
	Bold      bool
	Italic    bool
	Underline bool
	StrikeOut bool

	owned bool
}

func (this *Font) Init() {
	//
}

func (this *Font) Dispose() {
	if this.Handle == 0 {
		return
	}
	if this.owned {
		delete(dbusCache, this.Handle)
		win32.DeleteObject(win32.HGDIOBJ(this.Handle))
	}
	this.Handle = 0
}

func NewFont(name string, size int) *Font {
	return &Font{
		Name:  name,
		Size:  size,
		owned: true,
	}
}

func NewFontPt(name string, pointSize int) *Font {
	return &Font{
		Name:      name,
		PointSize: pointSize,
		owned:     true,
	}
}

func NewFontFromHandle(hFont win32.HFONT, owned bool) *Font {
	f := &Font{
		Handle: hFont,
		owned:  owned,
	}
	//fill fields from handle?
	return f
}

func (this *Font) Create() error {
	lf := this.createLogFont()
	hFont := win32.CreateFontIndirect(lf)
	if hFont == 0 {
		return errors.New("?")
	}
	this.Handle = hFont
	return nil
}

func (this *Font) CopyUnowned() *Font {
	this.EnsureCreated()
	font := *this
	font.owned = false
	return &font
}

var logPixelSy int32

func (this *Font) createLogFont() *win32.LOGFONT {
	lf := new(win32.LOGFONT)

	height := -int32(this.Size)
	if height == 0 && this.PointSize != 0 {
		if logPixelSy == 0 {
			hdc := win32.GetDC(0)
			logPixelSy = win32.GetDeviceCaps(hdc, win32.LOGPIXELSY)
			win32.ReleaseDC(0, hdc)
		}
		height = win32.MulDiv(int32(this.PointSize), logPixelSy, 72)
		height = -height
	}
	lf.LfHeight = height
	if this.Bold {
		lf.LfWeight = int32(win32.FW_BOLD)
	}
	if this.Italic {
		lf.LfItalic = 1
	}
	if this.Underline {
		lf.LfUnderline = 1
	}
	if this.StrikeOut {
		lf.LfStrikeOut = 1
	}
	lf.LfCharSet = win32.DEFAULT_CHARSET
	wsz, _ := syscall.UTF16FromString(this.Name)
	copy(lf.LfFaceName[:], wsz)

	//lf.LfQuality = win32.CLEARTYPE_QUALITY
	//lf.LfQuality = win32.ANTIALIASED_QUALITY

	return lf
}

func (this *Font) EnsureCreated() {
	if this.Handle != 0 {
		return
	}
	err := this.Create()
	if err != nil {
		log.Fatal(err)
	}
}

func (this *Font) Derive(bold bool, italic bool) *Font {
	font2 := &Font{}
	*font2 = *this
	font2.owned = false
	font2.Handle = 0
	font2.Bold = bold
	font2.Italic = italic
	return font2
}
