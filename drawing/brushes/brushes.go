package brushes

import (
	"github.com/zzl/goforms/drawing"
	"github.com/zzl/goforms/drawing/colors"
)

var (
	_black       *drawing.Brush
	_white       *drawing.Brush
	_transparent *drawing.Brush
	_blue        *drawing.Brush
	_skyBlue     *drawing.Brush
	_cyan        *drawing.Brush
	_fuchsia     *drawing.Brush
	_gray        *drawing.Brush
	_green       *drawing.Brush
	_lime        *drawing.Brush
	_magenta     *drawing.Brush
	_orange      *drawing.Brush
	_purple      *drawing.Brush
	_violet      *drawing.Brush
	_red         *drawing.Brush
	_silver      *drawing.Brush
	_yellow      *drawing.Brush

	_window     *drawing.Brush
	_windowText *drawing.Brush
	_grayText   *drawing.Brush
)

func Black() *drawing.Brush {
	if _black == nil {
		_black = drawing.NewSolidBrush(nil, colors.Black).AsBrush()
	}
	return _black
}

func White() *drawing.Brush {
	if _white == nil {
		_white = drawing.NewSolidBrush(nil, colors.White).AsBrush()
	}
	return _white
}

func Transparent() *drawing.Brush {
	if _transparent == nil {
		_transparent = drawing.NewSolidBrush(nil, colors.Transparent).AsBrush()
	}
	return _transparent
}

func Blue() *drawing.Brush {
	if _blue == nil {
		_blue = drawing.NewSolidBrush(nil, colors.Blue).AsBrush()
	}
	return _blue
}

func SkyBlue() *drawing.Brush {
	if _skyBlue == nil {
		_skyBlue = drawing.NewSolidBrush(nil, colors.SkyBlue).AsBrush()
	}
	return _skyBlue
}

func Cyan() *drawing.Brush {
	if _cyan == nil {
		_cyan = drawing.NewSolidBrush(nil, colors.Cyan).AsBrush()
	}
	return _cyan
}

func Fuchsia() *drawing.Brush {
	if _fuchsia == nil {
		_fuchsia = drawing.NewSolidBrush(nil, colors.Fuchsia).AsBrush()
	}
	return _fuchsia
}

func Gray() *drawing.Brush {
	if _gray == nil {
		_gray = drawing.NewSolidBrush(nil, colors.Gray).AsBrush()
	}
	return _gray
}

func Green() *drawing.Brush {
	if _green == nil {
		_green = drawing.NewSolidBrush(nil, colors.Green).AsBrush()
	}
	return _green
}

func Lime() *drawing.Brush {
	if _lime == nil {
		_lime = drawing.NewSolidBrush(nil, colors.Lime).AsBrush()
	}
	return _lime
}

func Magenta() *drawing.Brush {
	if _magenta == nil {
		_magenta = drawing.NewSolidBrush(nil, colors.Magenta).AsBrush()
	}
	return _magenta
}

func Orange() *drawing.Brush {
	if _orange == nil {
		_orange = drawing.NewSolidBrush(nil, colors.Orange).AsBrush()
	}
	return _orange
}

func Purple() *drawing.Brush {
	if _purple == nil {
		_purple = drawing.NewSolidBrush(nil, colors.Purple).AsBrush()
	}
	return _purple
}

func Violet() *drawing.Brush {
	if _violet == nil {
		_violet = drawing.NewSolidBrush(nil, colors.Violet).AsBrush()
	}
	return _violet
}

func Red() *drawing.Brush {
	if _red == nil {
		_red = drawing.NewSolidBrush(nil, colors.Red).AsBrush()
	}
	return _red
}

func Silver() *drawing.Brush {
	if _silver == nil {
		_silver = drawing.NewSolidBrush(nil, colors.Silver).AsBrush()
	}
	return _silver
}

func Yellow() *drawing.Brush {
	if _yellow == nil {
		_yellow = drawing.NewSolidBrush(nil, colors.Yellow).AsBrush()
	}
	return _yellow
}

func Window() *drawing.Brush {
	if _window == nil {
		_window = drawing.NewSolidBrush(nil, colors.Window).AsBrush()
	}
	return _window
}

func WindowText() *drawing.Brush {
	if _windowText == nil {
		_windowText = drawing.NewSolidBrush(nil, colors.WindowText).AsBrush()
	}
	return _windowText
}

func GrayText() *drawing.Brush {
	if _grayText == nil {
		_grayText = drawing.NewSolidBrush(nil, colors.GrayText).AsBrush()
	}
	return _grayText
}
