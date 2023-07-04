package cursors

import (
	"github.com/zzl/go-win32api/v2/win32"
)

var (
	_arrow    win32.HCURSOR
	_hand     win32.HCURSOR
	_ibeam    win32.HCURSOR
	_wait     win32.HCURSOR
	_sizeAll  win32.HCURSOR
	_sizeNwse win32.HCURSOR
	_sizeNesw win32.HCURSOR
	_sizeWe   win32.HCURSOR
	_sizeNs   win32.HCURSOR
)

func loadCursor(idc win32.PWSTR) win32.HCURSOR {
	ret, _ := win32.LoadCursor(0, idc)
	return ret
}

func Arrow() win32.HCURSOR {
	if _arrow == 0 {
		_arrow = loadCursor(win32.IDC_ARROW)
	}
	return _arrow
}

func Hand() win32.HCURSOR {
	if _hand == 0 {
		_hand = loadCursor(win32.IDC_HAND)
	}
	return _hand
}

func IBeam() win32.HCURSOR {
	if _ibeam == 0 {
		_ibeam = loadCursor(win32.IDC_IBEAM)
	}
	return _ibeam
}

func Wait() win32.HCURSOR {
	if _wait == 0 {
		_wait = loadCursor(win32.IDC_WAIT)
	}
	return _wait
}

func SizeAll() win32.HCURSOR {
	if _sizeAll == 0 {
		_sizeAll = loadCursor(win32.IDC_SIZEALL)
	}
	return _sizeAll
}

func SizeNwse() win32.HCURSOR {
	if _sizeNwse == 0 {
		_sizeNwse = loadCursor(win32.IDC_SIZENWSE)
	}
	return _sizeNwse
}

func SizeNesw() win32.HCURSOR {
	if _sizeNesw == 0 {
		_sizeNesw = loadCursor(win32.IDC_SIZENESW)
	}
	return _sizeNesw
}

func SizeWe() win32.HCURSOR {
	if _sizeWe == 0 {
		_sizeWe = loadCursor(win32.IDC_SIZEWE)
	}
	return _sizeWe
}

func SizeNs() win32.HCURSOR {
	if _sizeNs == 0 {
		_sizeNs = loadCursor(win32.IDC_SIZENS)
	}
	return _sizeNs
}
