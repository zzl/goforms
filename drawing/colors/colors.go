package colors

import "github.com/zzl/goforms/drawing"

var (
	//Unknown = drawing.ColorOfId(drawing.ColorIdNull)
	Null = drawing.ColorOfId(drawing.ColorIdNull)

	Black = drawing.ColorOfId(drawing.ColorIdBlack)

	White       = drawing.ColorOfId(drawing.ColorIdWhite)
	Transparent = drawing.ColorOfId(drawing.ColorIdTransparent)
	Blue        = drawing.ColorOfId(drawing.ColorIdBlue)
	SkyBlue     = drawing.ColorOfId(drawing.ColorIdSkyBlue)
	Cyan        = drawing.ColorOfId(drawing.ColorIdCyan)
	DarkBlue    = drawing.ColorOfId(drawing.ColorIdDarkBlue)
	DarkCyan    = drawing.ColorOfId(drawing.ColorIdDarkCyan)
	DarkGray    = drawing.ColorOfId(drawing.ColorIdDarkGray)
	DarkGreen   = drawing.ColorOfId(drawing.ColorIdDarkGreen)
	DarkMagenta = drawing.ColorOfId(drawing.ColorIdDarkMagenta)
	DarkOrange  = drawing.ColorOfId(drawing.ColorIdDarkOrange)
	DarkRed     = drawing.ColorOfId(drawing.ColorIdDarkRed)
	Fuchsia     = drawing.ColorOfId(drawing.ColorIdFuchsia)
	Gainsboro   = drawing.ColorOfId(drawing.ColorIdGainsboro)
	Gray        = drawing.ColorOfId(drawing.ColorIdGray)
	Green       = drawing.ColorOfId(drawing.ColorIdGreen)
	Lime        = drawing.ColorOfId(drawing.ColorIdLime)
	Magenta     = drawing.ColorOfId(drawing.ColorIdMagenta)
	Navy        = drawing.ColorOfId(drawing.ColorIdNavy)
	Orange      = drawing.ColorOfId(drawing.ColorIdOrange)
	Purple      = drawing.ColorOfId(drawing.ColorIdPurple)
	Violet      = drawing.ColorOfId(drawing.ColorIdViolet)
	Red         = drawing.ColorOfId(drawing.ColorIdRed)
	Silver      = drawing.ColorOfId(drawing.ColorIdSilver)
	WhiteSmoke  = drawing.ColorOfId(drawing.ColorIdWhiteSmoke)
	Yellow      = drawing.ColorOfId(drawing.ColorIdYellow)

	//
	Window          = drawing.ColorOfId(drawing.ColorIdWindow)
	WindowText      = drawing.ColorOfId(drawing.ColorIdWindowText)
	Highlight       = drawing.ColorOfId(drawing.ColorIdHighlight)
	HighlightText   = drawing.ColorOfId(drawing.ColorIdHighlightText)
	GrayText        = drawing.ColorOfId(drawing.ColorIdGrayText)
	Control         = drawing.ColorOfId(drawing.ColorIdControl)
	ControlDark     = drawing.ColorOfId(drawing.ColorIdControlDark)
	ControlDarkDark = drawing.ColorOfId(drawing.ColorIdControlDarkDark)
	InactiveBorder  = drawing.ColorOfId(drawing.ColorIdInactiveBorder)
)

func init() {
	drawing.ResolveNamedColorValue(drawing.ColorIdWindow)
	drawing.ResolveNamedColorValue(drawing.ColorIdWindowText)
}
