package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing/utils"
)

type NamedColorId int32

const (
	ColorIdNull NamedColorId = iota
	ColorIdCustom

	ColorIdActiveBorder
	ColorIdActiveCaption
	ColorIdActiveCaptionText
	ColorIdAppWorkspace
	ColorIdButtonFace
	ColorIdButtonHighlight
	ColorIdButtonShadow
	ColorIdControl
	ColorIdControlDark
	ColorIdControlDarkDark
	ColorIdControlLight
	ColorIdControlLightLight
	ColorIdControlText
	ColorIdDesktop
	ColorIdGradientActiveCaption
	ColorIdGradientInactiveCaption
	ColorIdGrayText
	ColorIdHighlight
	ColorIdHighlightText
	ColorIdHotTrack
	ColorIdInactiveBorder
	ColorIdInactiveCaption
	ColorIdInactiveCaptionText
	ColorIdInfo
	ColorIdInfoText
	ColorIdMenu
	ColorIdMenuBar
	ColorIdMenuHighlight
	ColorIdMenuText
	ColorIdScrollBar
	ColorIdWindow
	ColorIdWindowFrame
	ColorIdWindowText

	//
	ColorIdTransparent
	ColorIdBlack
	ColorIdDimGray
	ColorIdGray
	ColorIdDarkGray
	ColorIdSilver
	ColorIdLightGray
	ColorIdGainsboro
	ColorIdWhiteSmoke
	ColorIdWhite
	ColorIdRosyBrown
	ColorIdIndianRed
	ColorIdBrown
	ColorIdFirebrick
	ColorIdLightCoral
	ColorIdMaroon
	ColorIdDarkRed
	ColorIdRed
	ColorIdSnow
	ColorIdMistyRose
	ColorIdSalmon
	ColorIdTomato
	ColorIdDarkSalmon
	ColorIdCoral
	ColorIdOrangeRed
	ColorIdLightSalmon
	ColorIdSienna
	ColorIdSeaShell
	ColorIdChocolate
	ColorIdSaddleBrown
	ColorIdSandyBrown
	ColorIdPeachPuff
	ColorIdPeru
	ColorIdLinen
	ColorIdBisque
	ColorIdDarkOrange
	ColorIdBurlyWood
	ColorIdTan
	ColorIdAntiqueWhite
	ColorIdNavajoWhite
	ColorIdBlanchedAlmond
	ColorIdPapayaWhip
	ColorIdMoccasin
	ColorIdOrange
	ColorIdWheat
	ColorIdOldLace
	ColorIdFloralWhite
	ColorIdDarkGoldenrod
	ColorIdGoldenrod
	ColorIdCornsilk
	ColorIdGold
	ColorIdKhaki
	ColorIdLemonChiffon
	ColorIdPaleGoldenrod
	ColorIdDarkKhaki
	ColorIdBeige
	ColorIdLightGoldenrodYellow
	ColorIdOlive
	ColorIdYellow
	ColorIdLightYellow
	ColorIdIvory
	ColorIdOliveDrab
	ColorIdYellowGreen
	ColorIdDarkOliveGreen
	ColorIdGreenYellow
	ColorIdChartreuse
	ColorIdLawnGreen
	ColorIdDarkSeaGreen
	ColorIdForestGreen
	ColorIdLimeGreen
	ColorIdLightGreen
	ColorIdPaleGreen
	ColorIdDarkGreen
	ColorIdGreen
	ColorIdLime
	ColorIdHoneydew
	ColorIdSeaGreen
	ColorIdMediumSeaGreen
	ColorIdSpringGreen
	ColorIdMintCream
	ColorIdMediumSpringGreen
	ColorIdMediumAquamarine
	ColorIdAquamarine
	ColorIdTurquoise
	ColorIdLightSeaGreen
	ColorIdMediumTurquoise
	ColorIdDarkSlateGray
	ColorIdPaleTurquoise
	ColorIdTeal
	ColorIdDarkCyan
	ColorIdCyan
	ColorIdAqua
	ColorIdLightCyan
	ColorIdAzure
	ColorIdDarkTurquoise
	ColorIdCadetBlue
	ColorIdPowderBlue
	ColorIdLightBlue
	ColorIdDeepSkyBlue
	ColorIdSkyBlue
	ColorIdLightSkyBlue
	ColorIdSteelBlue
	ColorIdAliceBlue
	ColorIdDodgerBlue
	ColorIdSlateGray
	ColorIdLightSlateGray
	ColorIdLightSteelBlue
	ColorIdCornflowerBlue
	ColorIdRoyalBlue
	ColorIdMidnightBlue
	ColorIdLavender
	ColorIdNavy
	ColorIdDarkBlue
	ColorIdMediumBlue
	ColorIdBlue
	ColorIdGhostWhite
	ColorIdSlateBlue
	ColorIdDarkSlateBlue
	ColorIdMediumSlateBlue
	ColorIdMediumPurple
	ColorIdBlueViolet
	ColorIdIndigo
	ColorIdDarkOrchid
	ColorIdDarkViolet
	ColorIdMediumOrchid
	ColorIdThistle
	ColorIdPlum
	ColorIdViolet
	ColorIdPurple
	ColorIdDarkMagenta
	ColorIdFuchsia
	ColorIdMagenta
	ColorIdOrchid
	ColorIdMediumVioletRed
	ColorIdDeepPink
	ColorIdHotPink
	ColorIdLavenderBlush
	ColorIdPaleVioletRed
	ColorIdCrimson
	ColorIdPink
	ColorIdLightPink
)

var namedColorNames = []string{
	"(Unknown)",
	"",
	"ActiveBorder",
	"ActiveCaption",
	"ActiveCaptionText",
	"AppWorkspace",
	"ButtonFace",
	"ButtonHighlight",
	"ButtonShadow",
	"Control",
	"ControlDark",
	"ControlDarkDark",
	"ControlLight",
	"ControlLightLight",
	"ControlText",
	"Desktop",
	"GradientActiveCaption",
	"GradientInactiveCaption",
	"GrayText",
	"Highlight",
	"HighlightText",
	"HotTrack",
	"InactiveBorder",
	"InactiveCaption",
	"InactiveCaptionText",
	"Info",
	"InfoText",
	"Menu",
	"MenuBar",
	"MenuHighlight",
	"MenuText",
	"ScrollBar",
	"Window",
	"WindowFrame",
	"WindowText",
	"Transparent",
	"Black",
	"DimGray",
	"Gray",
	"DarkGray",
	"Silver",
	"LightGray",
	"Gainsboro",
	"WhiteSmoke",
	"White",
	"RosyBrown",
	"IndianRed",
	"Brown",
	"Firebrick",
	"LightCoral",
	"Maroon",
	"DarkRed",
	"Red",
	"Snow",
	"MistyRose",
	"Salmon",
	"Tomato",
	"DarkSalmon",
	"Coral",
	"OrangeRed",
	"LightSalmon",
	"Sienna",
	"SeaShell",
	"Chocolate",
	"SaddleBrown",
	"SandyBrown",
	"PeachPuff",
	"Peru",
	"Linen",
	"Bisque",
	"DarkOrange",
	"BurlyWood",
	"Tan",
	"AntiqueWhite",
	"NavajoWhite",
	"BlanchedAlmond",
	"PapayaWhip",
	"Moccasin",
	"Orange",
	"Wheat",
	"OldLace",
	"FloralWhite",
	"DarkGoldenrod",
	"Goldenrod",
	"Cornsilk",
	"Gold",
	"Khaki",
	"LemonChiffon",
	"PaleGoldenrod",
	"DarkKhaki",
	"Beige",
	"LightGoldenrodYellow",
	"Olive",
	"Yellow",
	"LightYellow",
	"Ivory",
	"OliveDrab",
	"YellowGreen",
	"DarkOliveGreen",
	"GreenYellow",
	"Chartreuse",
	"LawnGreen",
	"DarkSeaGreen",
	"ForestGreen",
	"LimeGreen",
	"LightGreen",
	"PaleGreen",
	"DarkGreen",
	"Green",
	"Lime",
	"Honeydew",
	"SeaGreen",
	"MediumSeaGreen",
	"SpringGreen",
	"MintCream",
	"MediumSpringGreen",
	"MediumAquamarine",
	"Aquamarine",
	"Turquoise",
	"LightSeaGreen",
	"MediumTurquoise",
	"DarkSlateGray",
	"PaleTurquoise",
	"Teal",
	"DarkCyan",
	"Cyan",
	"Aqua",
	"LightCyan",
	"Azure",
	"DarkTurquoise",
	"CadetBlue",
	"PowderBlue",
	"LightBlue",
	"DeepSkyBlue",
	"SkyBlue",
	"LightSkyBlue",
	"SteelBlue",
	"AliceBlue",
	"DodgerBlue",
	"SlateGray",
	"LightSlateGray",
	"LightSteelBlue",
	"CornflowerBlue",
	"RoyalBlue",
	"MidnightBlue",
	"Lavender",
	"Navy",
	"DarkBlue",
	"MediumBlue",
	"Blue",
	"GhostWhite",
	"SlateBlue",
	"DarkSlateBlue",
	"MediumSlateBlue",
	"MediumPurple",
	"BlueViolet",
	"Indigo",
	"DarkOrchid",
	"DarkViolet",
	"MediumOrchid",
	"Thistle",
	"Plum",
	"Violet",
	"Purple",
	"DarkMagenta",
	"Fuchsia",
	"Magenta",
	"Orchid",
	"MediumVioletRed",
	"DeepPink",
	"HotPink",
	"LavenderBlush",
	"PaleVioletRed",
	"Crimson",
	"Pink",
	"LightPink",
}

var sysColorIds = []NamedColorId{
	ColorIdScrollBar, ColorIdDesktop, ColorIdActiveCaption, ColorIdInactiveCaption,
	ColorIdMenu, ColorIdWindow, ColorIdWindowFrame, ColorIdMenuText,
	ColorIdWindowText, ColorIdActiveCaptionText, ColorIdActiveBorder, ColorIdInactiveBorder,
	ColorIdAppWorkspace, ColorIdHighlight, ColorIdHighlightText, ColorIdControl,
	ColorIdControlDark, ColorIdGrayText, ColorIdControlText, ColorIdInactiveCaptionText,
	ColorIdControlLightLight, ColorIdControlDarkDark, ColorIdControlLight, ColorIdInfoText,
	ColorIdInfo, 0, ColorIdHotTrack, ColorIdGradientActiveCaption,
	ColorIdGradientInactiveCaption, ColorIdGradientInactiveCaption,
	ColorIdMenuHighlight, ColorIdMenuBar,
}

var namedColorValues = []gdip.ARGB{
	0, //Null
	0, //Custom

	0, //ActiveBorder
	0, //ActiveCaption
	0, //ActiveCaptionText
	0, //AppWorkspace
	0, //ButtonFace
	0, //ButtonHighlight
	0, //ButtonShadow
	0, //Control
	0, //ControlDark
	0, //ControlDarkDark
	0, //ControlLight
	0, //ControlLightLight
	0, //ControlText
	0, //Desktop
	0, //GradientActiveCaption
	0, //GradientInactiveCaption
	0, //GrayText
	0, //Highlight
	0, //HighlightText
	0, //HotTrack
	0, //InactiveBorder
	0, //InactiveCaption
	0, //InactiveCaptionText
	0, //Info
	0, //InfoText
	0, //Menu
	0, //MenuBar
	0, //MenuHighlight
	0, //MenuText
	0, //ScrollBar
	0, //Window
	0, //WindowFrame
	0, //WindowText

	0,          //Transparent
	0xFF000000, //Black
	0xFF696969, //DimGray
	0xFF808080, //Gray
	0xFFA9A9A9, //DarkGray
	0xFFC0C0C0, //Silver
	0xFFD3D3D3, //LightGray
	0xFFDCDCDC, //Gainsboro
	0xFFF5F5F5, //WhiteSmoke
	0xFFFFFFFF, //White
	0xFFBC8F8F, //RosyBrown
	0xFFCD5C5C, //IndianRed
	0xFFA52A2A, //Brown
	0xFFB22222, //Firebrick
	0xFFF08080, //LightCoral
	0xFF800000, //Maroon
	0xFF8B0000, //DarkRed
	0xFFFF0000, //Red
	0xFFFFFAFA, //Snow
	0xFFFFE4E1, //MistyRose
	0xFFFA8072, //Salmon
	0xFFFF6347, //Tomato
	0xFFE9967A, //DarkSalmon
	0xFFFF7F50, //Coral
	0xFFFF4500, //OrangeRed
	0xFFFFA07A, //LightSalmon
	0xFFA0522D, //Sienna
	0xFFFFF5EE, //SeaShell
	0xFFD2691E, //Chocolate
	0xFF8B4513, //SaddleBrown
	0xFFF4A460, //SandyBrown
	0xFFFFDAB9, //PeachPuff
	0xFFCD853F, //Peru
	0xFFFAF0E6, //Linen
	0xFFFFE4C4, //Bisque
	0xFFFF8C00, //DarkOrange
	0xFFDEB887, //BurlyWood
	0xFFD2B48C, //Tan
	0xFFFAEBD7, //AntiqueWhite
	0xFFFFDEAD, //NavajoWhite
	0xFFFFEBCD, //BlanchedAlmond
	0xFFFFEFD5, //PapayaWhip
	0xFFFFE4B5, //Moccasin
	0xFFFFA500, //Orange
	0xFFF5DEB3, //Wheat
	0xFFFDF5E6, //OldLace
	0xFFFFFAF0, //FloralWhite
	0xFFB8860B, //DarkGoldenrod
	0xFFDAA520, //Goldenrod
	0xFFFFF8DC, //Cornsilk
	0xFFFFD700, //Gold
	0xFFF0E68C, //Khaki
	0xFFFFFACD, //LemonChiffon
	0xFFEEE8AA, //PaleGoldenrod
	0xFFBDB76B, //DarkKhaki
	0xFFF5F5DC, //Beige
	0xFFFAFAD2, //LightGoldenrodYellow
	0xFF808000, //Olive
	0xFFFFFF00, //Yellow
	0xFFFFFFE0, //LightYellow
	0xFFFFFFF0, //Ivory
	0xFF6B8E23, //OliveDrab
	0xFF9ACD32, //YellowGreen
	0xFF556B2F, //DarkOliveGreen
	0xFFADFF2F, //GreenYellow
	0xFF7FFF00, //Chartreuse
	0xFF7CFC00, //LawnGreen
	0xFF8FBC8B, //DarkSeaGreen
	0xFF228B22, //ForestGreen
	0xFF32CD32, //LimeGreen
	0xFF90EE90, //LightGreen
	0xFF98FB98, //PaleGreen
	0xFF006400, //DarkGreen
	0xFF008000, //Green
	0xFF00FF00, //Lime
	0xFFF0FFF0, //Honeydew
	0xFF2E8B57, //SeaGreen
	0xFF3CB371, //MediumSeaGreen
	0xFF00FF7F, //SpringGreen
	0xFFF5FFFA, //MintCream
	0xFF00FA9A, //MediumSpringGreen
	0xFF66CDAA, //MediumAquamarine
	0xFF7FFFD4, //Aquamarine
	0xFF40E0D0, //Turquoise
	0xFF20B2AA, //LightSeaGreen
	0xFF48D1CC, //MediumTurquoise
	0xFF2F4F4F, //DarkSlateGray
	0xFFAFEEEE, //PaleTurquoise
	0xFF008080, //Teal
	0xFF008B8B, //DarkCyan
	0xFF00FFFF, //Cyan
	0xFF00FFFF, //Aqua
	0xFFE0FFFF, //LightCyan
	0xFFF0FFFF, //Azure
	0xFF00CED1, //DarkTurquoise
	0xFF5F9EA0, //CadetBlue
	0xFFB0E0E6, //PowderBlue
	0xFFADD8E6, //LightBlue
	0xFF00BFFF, //DeepSkyBlue
	0xFF87CEEB, //SkyBlue
	0xFF87CEFA, //LightSkyBlue
	0xFF4682B4, //SteelBlue
	0xFFF0F8FF, //AliceBlue
	0xFF1E90FF, //DodgerBlue
	0xFF708090, //SlateGray
	0xFF778899, //LightSlateGray
	0xFFB0C4DE, //LightSteelBlue
	0xFF6495ED, //CornflowerBlue
	0xFF4169E1, //RoyalBlue
	0xFF191970, //MidnightBlue
	0xFFE6E6FA, //Lavender
	0xFF000080, //Navy
	0xFF00008B, //DarkBlue
	0xFF0000CD, //MediumBlue
	0xFF0000FF, //Blue
	0xFFF8F8FF, //GhostWhite
	0xFF6A5ACD, //SlateBlue
	0xFF483D8B, //DarkSlateBlue
	0xFF7B68EE, //MediumSlateBlue
	0xFF9370DB, //MediumPurple
	0xFF8A2BE2, //BlueViolet
	0xFF4B0082, //Indigo
	0xFF9932CC, //DarkOrchid
	0xFF9400D3, //DarkViolet
	0xFFBA55D3, //MediumOrchid
	0xFFD8BFD8, //Thistle
	0xFFDDA0DD, //Plum
	0xFFEE82EE, //Violet
	0xFF800080, //Purple
	0xFF8B008B, //DarkMagenta
	0xFFFF00FF, //Fuchsia
	0xFFFF00FF, //Magenta
	0xFFDA70D6, //Orchid
	0xFFC71585, //MediumVioletRed
	0xFFFF1493, //DeepPink
	0xFFFF69B4, //HotPink
	0xFFFFF0F5, //LavenderBlush
	0xFFDB7093, //PaleVioletRed
	0xFFDC143C, //Crimson
	0xFFFFC0CB, //Pink
	0xFFFFB6C1, //LightPink
}

func GetColorName(id NamedColorId) string {
	return namedColorNames[id]
}

func ResolveNamedColorValue(id NamedColorId) gdip.ARGB {
	argb := namedColorValues[id]
	if argb != 0 {
		return argb
	}
	switch id {
	case ColorIdTransparent:
		return 0
	case ColorIdCustom:
		return 0
	case ColorIdButtonFace:
		id = ColorIdControl
	case ColorIdButtonHighlight:
		id = ColorIdControlLightLight
	case ColorIdButtonShadow:
		id = ColorIdControlDark
	}
	sysColorIndex := -1
	for n, tid := range sysColorIds {
		if tid == id {
			sysColorIndex = n
			break
		}
	}
	if sysColorIndex == -1 {
		panic("?")
	}
	ret := win32.GetSysColor(win32.SYS_COLOR_INDEX(sysColorIndex))
	argb = utils.Win32ColorToArgb(ret)
	namedColorValues[id] = argb
	return argb
}

var colorValueIds map[gdip.ARGB]NamedColorId

func ResolveColorId(value gdip.ARGB) NamedColorId {
	id, ok := colorValueIds[value]
	if ok {
		return id
	}
	for id := ColorIdTransparent; id <= ColorIdLightPink; id++ {
		if namedColorValues[id] == value {
			if colorValueIds == nil {
				colorValueIds = make(map[gdip.ARGB]NamedColorId)
			}
			colorValueIds[value] = id
			return id
		}
	}
	return 0
}

func ResolveColorName(value gdip.ARGB) string {
	id := ResolveColorId(value)
	return namedColorNames[id]
}
