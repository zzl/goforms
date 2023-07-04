package utils

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
)

func Win32ColorToArgb(color win32.COLORREF) gdip.ARGB {
	r, g, b := byte(color&0xFF), byte(color>>8&0xFF), byte(color>>16&0xFF)
	argb := uint32(b) | (uint32(g) << 8) | (uint32(r) << 16) | uint32(0xFF000000)
	return gdip.ARGB(argb)
}
