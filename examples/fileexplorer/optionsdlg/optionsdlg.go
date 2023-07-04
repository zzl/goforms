package optionsdlg

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing/gdi"
	"github.com/zzl/goforms/forms"
	"log"
)

var hInstShell32 win32.HINSTANCE

func Show() {

	hInstShell32, _ = win32.LoadLibrary(win32.StrToPwstr("shell32.dll"))

	var rc win32.RECT
	win32.GetWindowRect(forms.HWndActive, &rc)

	form := forms.NewForm{
		Title:         "File Explorer Options",
		ClientSize:    forms.DpiSize(380, 443),
		Pos:           forms.Point{int(rc.Left + 30), int(rc.Top + 64)},
		NoIcon:        true,
		NoMaximizeBox: true,
		NoMinimizeBox: true,
		NoResize:      true,
	}.Create()
	font := gdi.NewFontPt("Microsoft Sans Serif", 8)
	form.SetFont(font)
	forms.ContextContainer = form

	tabCtrl := forms.NewTabControlObject()
	tabCtrl.CreateIn(form)
	tabCtrl.SetDpiPos(6, 7)
	cx, cy := form.GetDpiClientSize()
	tabCtrl.SetDpiSize(cx-12, cy-44)

	sheetGeneral := forms.NewTabSheet{Title: "General"}.Create()
	populateSheetGeneral(sheetGeneral)

	sheetView := forms.NewTabSheet{Title: "View"}.Create()
	populateSheetView(sheetView)

	sheetSearch := forms.NewTabSheet{Title: "Search"}.Create()

	tabCtrl.SetItems([]forms.TabItem{
		{Sheet: sheetGeneral},
		{Sheet: sheetView},
		{Sheet: sheetSearch},
	})

	forms.NewButton{
		Id:   uint16(win32.IDOK),
		Text: "OK",
		Pos:  forms.DpiPoint(137, 413),
		Size: forms.DpiSize(75, 23),
		Action: func() {
			forms.Alert("OK?")
		},
	}.Create()

	forms.NewButton{
		Id:   uint16(win32.IDCANCEL),
		Text: "Cancel",
		Pos:  forms.DpiPoint(218, 413),
		Size: forms.DpiSize(75, 23),
		Action: func() {
			form.Close()
		},
	}.Create()

	forms.NewButton{
		Text:     "&Apply",
		Pos:      forms.DpiPoint(299, 413),
		Size:     forms.DpiSize(75, 23),
		Disabled: true,
	}.Create()

	//
	forms.SendMessage(form.GetHandle(), win32.DM_SETDEFID,
		win32.WPARAM(win32.IDOK), 0)

	//
	form.ShowModal()
}

func populateSheetView(sheet forms.Container) {
	defer forms.SetContextContainer(sheet).Restore()

	forms.NewGroupBox{
		Text: "Folder views",
		Pos:  forms.DpiPoint(15, 8),
		Size: forms.DpiSize(330, 89),
	}.Create()

	hIcon, _ := win32.LoadIcon(hInstShell32, win32.MAKEINTRESOURCE(20))
	forms.NewImageBox{
		Pos:  forms.DpiPoint(30, 33),
		Size: forms.DpiSize(32, 32),
		Icon: hIcon,
	}.Create()

	forms.NewLabel{
		Text: "You can apply this view (such as Details or Icons) " +
			"to all folders of this type.",
		Pos:  forms.DpiPoint(84, 23),
		Size: forms.DpiSize(254, 33),
	}.Create()

	forms.NewButton{
		Text:     "Apply to Fo&lders",
		Pos:      forms.DpiPoint(84, 57),
		Size:     forms.DpiSize(116, 24),
		Disabled: true,
	}.Create()

	forms.NewButton{
		Text: "&Reset Folders",
		Pos:  forms.DpiPoint(215, 57),
		Size: forms.DpiSize(116, 24),
	}.Create()

	forms.NewLabel{
		Text: "Advanced settings:",
		Pos:  forms.DpiPoint(15, 114),
		Size: forms.DpiSize(162, 13),
	}.Create()

	//
	tree := forms.NewTreeViewObject()
	tree.NoLines = true
	tree.NoLinesAtRoot = true
	tree.NoButtons = true

	tree.CreateIn(sheet)
	iml := createTreeIml(tree.Handle)
	tree.SetImageList(iml)
	tree.SetDpiPos(15, 130)
	tree.SetDpiSize(330, 195)
	tree.AddNodes(
		forms.TreeNode{
			Text:  "Files and Folders",
			Image: 0,
			Children: []forms.TreeNode{
				{Text: "Always show icons, never thumbnails", Image: 1},
				{Text: "Always show menus", Image: 1},
				{Text: "Display file icon on thumbnails", Image: 2},
				{Text: "Display file size information in folder tips", Image: 2},
				{Text: "Display the full path in the title bar", Image: 1},
				{Text: "Hidden files and folders", Image: 0, Children: []forms.TreeNode{
					{Text: "Don't show hidden files, folder, or drives", Image: 4},
					{Text: "Show hidden files, folder, and drives", Image: 3},
				}},
				{Text: "Hide empty drives", Image: 2},
				{Text: "Hide extensions for known file types", Image: 2},
				{Text: "Hide folder merge conflicts", Image: 2},
				{Text: "Hide protected operation systemfiles (Recommended)", Image: 2},
				{Text: "Launch folder windows in a separate process", Image: 1},
				{Text: "Restore previous folder windows at logon", Image: 1},
				{Text: "Show drive letters", Image: 2},
				{Text: "Show encrypted or compressed NTFS files in color", Image: 1},
				{Text: "Show pop-up description for folder and desktop items", Image: 2},
				{Text: "Show preview handlers in preview pane", Image: 2},
				{Text: "Show status bar", Image: 2},
				{Text: "Show sync provider notifications", Image: 2},
				{Text: "Use check boxes to select items", Image: 1},
				{Text: "Use Sharing Wizard (Recommended)", Image: 2},
				{Text: "When typing into list view", Image: 0, Children: []forms.TreeNode{
					{Text: "Automatically type into the Search Box", Image: 3},
					{Text: "Select the typed item in the view", Image: 4},
				}},
			},
		},
		forms.TreeNode{Text: "Navigation pane", Image: 0, Children: []forms.TreeNode{
			{Text: "Always show availability status", Image: 1},
			{Text: "Expand to open folder", Image: 1},
			{Text: "Show all folders", Image: 1},
			{Text: "Show libraries", Image: 1},
		}})
	tree.ExpandAll()

	forms.NewButton{
		Text: "Restore &Defaults",
		Pos:  forms.DpiPoint(240, 341),
		Size: forms.DpiSize(105, 23),
	}.Create()
}

func populateSheetGeneral(sheet forms.Container) {
	defer forms.SetContextContainer(sheet).Restore()

	forms.NewLabel{
		Pos:  forms.DpiPoint(10, 14),
		Text: "Open File Explorer to:",
	}.Create()
	forms.NewComboBoxObject()

	forms.NewComboBox{
		Text: "This PC",
		Pos:  forms.DpiPoint(120, 10),
		Size: forms.DpiSize(219, 0),
		Items: []string{
			"Quick access",
			"This PC",
		},
	}.Create()

	forms.NewGroupBox{
		Text: "Browse folders",
		Pos:  forms.DpiPoint(11, 36),
		Size: forms.DpiSize(329, 62),
	}.Create()

	hIcon, _ := win32.LoadIcon(hInstShell32, win32.MAKEINTRESOURCE(184))
	forms.NewImageBox{
		Pos:  forms.DpiPoint(21, 55),
		Size: forms.DpiSize(32, 32),
		Icon: hIcon,
	}.Create()

	forms.NewRadioButton{
		Text:    "Open each folder in the sa&me window",
		Pos:     forms.DpiPoint(62, 52),
		Checked: true,
	}.Create()

	forms.NewRadioButton{
		Text: "Open each folder in its own &window",
		Pos:  forms.DpiPoint(62, 73),
	}.Create()

	//
	forms.NewGroupBox{
		Text: "Click items as follows",
		Pos:  forms.DpiPoint(11, 109),
		Size: forms.DpiSize(329, 98),
	}.Create()

	hIcon, _ = win32.LoadIcon(hInstShell32, win32.MAKEINTRESOURCE(186))
	//id = 327
	forms.NewImageBox{
		Pos:  forms.DpiPoint(21, 129),
		Size: forms.DpiSize(32, 32),
		Icon: hIcon,
	}.Create()

	forms.NewRadioButton{
		Text: "&Single-click to open an item (point to select)",
		Pos:  forms.DpiPoint(62, 128),
	}.Create()

	forms.NewRadioButton{
		Text:    "&Double-click to open an item (single-click to select)",
		Pos:     forms.DpiPoint(62, 182),
		Checked: true,
	}.Create()

	forms.NewRadioButton{
		Text:       "Underline icon titles consistent with my &browser",
		Pos:        forms.DpiPoint(81, 146),
		Disabled:   true,
		BeginGroup: true,
	}.Create()

	forms.NewRadioButton{
		Text:     "Underline icon titles only when i &point at them",
		Pos:      forms.DpiPoint(81, 164),
		Checked:  true,
		Disabled: true,
	}.Create()

	//
	forms.NewGroupBox{
		Text: "Privacy",
		Pos:  forms.DpiPoint(11, 218),
		Size: forms.DpiSize(329, 94),
	}.Create()

	hIcon, _ = win32.LoadIcon(hInstShell32, win32.MAKEINTRESOURCE(327))
	forms.NewImageBox{
		Pos:  forms.DpiPoint(21, 235),
		Size: forms.DpiSize(32, 32),
		Icon: hIcon,
	}.Create()

	forms.NewCheckBox{
		Pos:  forms.DpiPoint(62, 234),
		Text: "Show recently used files in Quick access",
	}.Create()

	forms.NewCheckBox{
		Pos:  forms.DpiPoint(62, 255),
		Text: "Show frequently used folders in Quick access",
	}.Create()

	forms.NewLabel{
		Text: "Clear File Explorer history",
		Pos:  forms.DpiPoint(62, 283),
	}.Create()

	forms.NewButton{
		Text: "&Clear",
		Pos:  forms.DpiPoint(261, 278),
		Size: forms.DpiSize(63, 23),
	}.Create()

	forms.NewButton{
		Text: "&Restore Defaults",
		Pos:  forms.DpiPoint(231, 317),
		Size: forms.DpiSize(108, 23),
	}.Create()

}

func createTreeIml(hWnd win32.HWND) *forms.ImageList {
	hTmeme := win32.OpenThemeData(hWnd, win32.StrToPwstr("button"))
	if hTmeme == 0 {
		log.Panic("?")
	}
	defer win32.CloseThemeData(hTmeme)

	hDc := win32.GetDC(hWnd)
	defer win32.ReleaseDC(hWnd, hDc)

	var size win32.SIZE
	win32.GetThemePartSize(hTmeme, hDc, int32(win32.BP_CHECKBOX),
		int32(win32.CBS_UNCHECKEDNORMAL), nil, win32.TS_DRAW, &size)

	nodeIconSize := forms.DpiScale(int(size.Cx))
	iml := forms.NewImageList(nodeIconSize, true)

	//
	hBmp := win32.CreateCompatibleBitmap(hDc, int32(nodeIconSize)*4, int32(nodeIconSize))
	hDcMem := win32.CreateCompatibleDC(hDc)
	hOriBmp := win32.SelectObject(hDcMem, hBmp)

	//
	var rc win32.RECT
	rc.Right = int32(nodeIconSize)
	rc.Bottom = int32(nodeIconSize)

	win32.DrawThemeBackground(hTmeme, hDcMem, int32(win32.BP_CHECKBOX),
		int32(win32.CBS_UNCHECKEDNORMAL), &rc, nil)

	win32.OffsetRect(&rc, int32(nodeIconSize), 0)
	win32.DrawThemeBackground(hTmeme, hDcMem, int32(win32.BP_CHECKBOX),
		int32(win32.CBS_CHECKEDNORMAL), &rc, nil)

	win32.OffsetRect(&rc, int32(nodeIconSize), 0)
	win32.DrawThemeBackground(hTmeme, hDcMem, int32(win32.BP_RADIOBUTTON),
		int32(win32.RBS_UNCHECKEDNORMAL), &rc, nil)

	win32.OffsetRect(&rc, int32(nodeIconSize), 0)
	win32.DrawThemeBackgroundEx(hTmeme, hDcMem, int32(win32.BP_RADIOBUTTON),
		int32(win32.RBS_CHECKEDNORMAL), &rc, nil)

	win32.SelectObject(hDcMem, hOriBmp)
	win32.DeleteDC(hDcMem)

	var hIcon win32.HICON
	win32.LoadIconWithScaleDown(hInstShell32, win32.MAKEINTRESOURCE(4),
		int32(nodeIconSize), int32(nodeIconSize), &hIcon)
	iml.AddIcon(hIcon)
	iml.Add(hBmp)
	win32.DeleteObject(hBmp)
	return iml
}
