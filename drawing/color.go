package drawing

import (
	"errors"
	"fmt"
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing/utils"
	"image/color"
	"strings"
)

type Color struct {
	value gdip.ARGB
	id    NamedColorId
}

func ColorOf(argb gdip.ARGB) Color {
	return Color{id: ColorIdCustom, value: argb}
}

func ColorOfId(id NamedColorId) Color {
	argb := namedColorValues[id]
	return Color{value: argb, id: id}
}

func Rgb(r byte, g byte, b byte) Color {
	argb := uint32(b) | (uint32(g) << 8) | (uint32(r) << 16) | 255<<24
	return Color{id: ColorIdCustom, value: gdip.ARGB(argb)}
}

func Rgba(r byte, g byte, b byte, a byte) Color {
	argb := uint32(b) | (uint32(g) << 8) | (uint32(r) << 16) | (uint32(a) << 24)
	return Color{id: ColorIdCustom, value: gdip.ARGB(argb)}
}

func ColorFromGo(color color.Color) Color {
	r, g, b, a := color.RGBA()
	argb := b | (g << 8) | (r << 16) | (a << 24)
	return Color{id: ColorIdCustom, value: gdip.ARGB(argb)}
}

func ColorFromWin32(color win32.COLORREF) Color {
	argb := utils.Win32ColorToArgb(color)
	return Color{id: ColorIdCustom, value: argb}
}

func ColorFromSysColor(colorIndex int) Color {
	id := sysColorIds[colorIndex]
	argb := ResolveNamedColorValue(id)
	return Color{value: argb, id: id}
}

func ColorFromName(name string) (Color, error) {
	name = strings.ToLower(name)
	for n, tName := range namedColorNames {
		tName = strings.ToLower(tName)
		if tName == name {
			return ColorOfId(NamedColorId(n)), nil
		}
	}
	return Color{}, errors.New("Unknown color name " + name)
}

//func ColorFromName(id NamedColorId) Color {
//	argb := namedColorValues[id]
//	if argb == 0 {
//		if id == ColorIdTransparent {
//			//
//		} else {
//			argb = ResolveNamedColorValue(id)
//		}
//	}
//	return Color{value: argb, id: id}
//}

func (this *Color) Argb() gdip.ARGB {
	if this.value != 0 || this.id == 0 {
		return this.value
	}
	argb := namedColorValues[this.id]
	if argb != 0 {
		return argb
	}
	argb = ResolveNamedColorValue(this.id)
	this.value = argb
	return argb
}

func (this *Color) Id() NamedColorId {
	return this.id
}

func (this *Color) ResolveId() NamedColorId {
	if this.id != 0 {
		return this.id
	}
	return ResolveColorId(this.value)
}

func (this *Color) A() byte {
	argb := this.Argb()
	return byte(argb >> 24)
}

func (this *Color) R() byte {
	argb := this.Argb()
	return byte(argb >> 16)
}

func (this *Color) G() byte {
	argb := this.Argb()
	return byte(argb >> 8)
}

func (this *Color) B() byte {
	argb := this.Argb()
	return byte(argb)
}

func (this *Color) Rgb() (byte, byte, byte) {
	argb := this.Argb()
	return byte(argb >> 16), byte(argb >> 8), byte(argb)
}

func (this *Color) Win32Color() win32.COLORREF {
	argb := this.Argb()
	r, g, b := byte(argb>>16), byte(argb>>8), byte(argb)
	return win32.RGB(r, g, b)
}

//func (this *Color) IsNull() bool {
//	return this.id == 0 //this.value == 0 &&
//}

//func (this *Color) Valid() bool {
//	return this.id != ColorIdNull && this.id != ColorIdCustom
//}

func (this *Color) IsTransparent() bool {
	return this.id == ColorIdTransparent
}

func (this *Color) String() string {
	id := this.ResolveId()
	if id != 0 {
		return namedColorNames[id]
	}
	r, g, b := this.Rgb()
	return fmt.Sprintf("%v, %v, %v", r, g, b)
}

func (this *Color) GetName() string {
	id := this.ResolveId()
	return namedColorNames[id]
}
