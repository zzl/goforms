package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
)

// ??
func LoadIcon(filePath string, big bool) win32.HICON {
	//hIcon, errno := win32.LoadIcon(0, win32.StrToPwstr(filePath))
	//if errno != win32.NO_ERROR {
	//	//?
	//	return 0
	//}
	//return hIcon

	//img, _ := drawing.NewBitmapFromFile(nil, filePath)
	//icon := img.GetHIcon()
	//return icon
	size := int32(16)
	if big {
		size = 32
		//size, _ = win32.GetSystemMetrics(win32.SM_CXICON)
	}
	hIcon, errno := win32.LoadImage(0, win32.StrToPwstr(filePath),
		win32.IMAGE_ICON, size, size,
		win32.LR_LOADFROMFILE|win32.LR_SHARED)
	if errno != win32.NO_ERROR {
		println("??")
	}
	return hIcon
}
