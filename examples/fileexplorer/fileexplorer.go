package main

import (
	"fmt"
	"github.com/zzl/go-com/com"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/examples/fileexplorer/optionsdlg"
	"github.com/zzl/goforms/forms"
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/events"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"github.com/zzl/goforms/layouts"
	"github.com/zzl/goforms/layouts/aligns"
)

func main() {

	com.Initialize()

	forms.SetDpiAware() //?

	shell := &Shell{}
	shell.Init()

	//
	form := forms.NewForm{
		Title:      "File Explorer",
		ClientSize: forms.Size{800, 600},
	}.Create()

	//
	var tree *ShellTree

	menuBar := forms.NewMenuBar(true)
	menuBar.AddItems([]forms.MenuItem{
		{
			Id:   1,
			Text: "&File",
			SubItems: []*forms.MenuItem{
				{
					Text: "&Browse for folder..",
					Action: func() {
						si := shell.BrowserForFolder()
						if si != nil {
							tree.SetSelectedShellItem(si)
						}
					},
				},
				{Text: "-"},
				{
					Text:   "&Close",
					Action: form.Close,
				},
			},
		},
		{
			Text: "&Edit",
			SubItems: []*forms.MenuItem{
				{Text: "&Copy"},
			},
		},
		{
			Text: "&View",
			SubItems: []*forms.MenuItem{
				{
					Name:    "toolbar",
					Text:    "&Toolbar",
					Checked: true,
					Action: func() {
						layoutItem := form.GetLayout().GetItem("toolbar")
						collapsed := !layoutItem.IsCollapsed()
						layoutItem.SetCollapsed(collapsed)
						form.UpdateLayout()
						menuBar.CheckItemByName("toolbar", !collapsed)
					},
				},
				{
					Name:    "statusbar",
					Text:    "&Status Bar",
					Checked: true,
					Action: func() {
						layoutItem := form.GetLayout().GetItem("statusbar")
						collapsed := !layoutItem.IsCollapsed()
						layoutItem.SetCollapsed(collapsed)

						form.GetLayout().GetItem("treeandlist").(*layouts.LinearItem).
							PaddingBottom = utils.If(collapsed, 6)

						form.UpdateLayout()
						menuBar.CheckItemByName("statusbar", !collapsed)
					},
				},
			},
		},
		{
			Text: "&Tools",
			SubItems: []*forms.MenuItem{
				{
					Text:   "&Folder options...",
					Action: optionsdlg.Show,
				},
			},
		},
		{
			Text: "&Help",
			SubItems: []*forms.MenuItem{
				{Text: "&Help"},
				{Text: "-"},
				{
					Text: "&About",
					Action: func() {
						forms.Info("File Explorer v0.1")
					},
				},
			},
		},
	})
	form.SetMenuBar(menuBar)

	//
	tree = virtual.New[ShellTree]()
	tree.shell = shell
	tree.CreateIn(form)

	splitter := forms.NewSplitterObject()
	splitter.Width = 6
	splitter.CreateIn(form)

	list := virtual.New[ShellList]()
	list.shell = shell
	list.CreateIn(form)

	treeAndListLayout := &layouts.LinearLayout{
		ItemDefaults: &layouts.LinearItem{
			Align: aligns.Stretch,
		},
		Items: []*layouts.LinearItem{
			{
				Control: tree,
				Width:   250,
			},
			{
				Control: splitter,
			},
			{
				Control: list,
				Weight:  1,
			},
		},
	}
	tb := forms.NewToolBar{
		Parent:     form,
		LabelStyle: forms.TbLabelInvisible,
	}.Create()

	tbIml2 := forms.NewImageList(16, true)
	tb.SetImageList(tbIml2, nil)
	tb.LoadSystemImages(win32.IDB_VIEW_SMALL_COLOR)
	hIcon := tbIml2.GetHIcon(int(win32.VIEW_PARENTFOLDER))

	tbIml := forms.NewImageList(16, true)
	tb.SetImageList(tbIml, nil)
	tb.LoadSystemImages(win32.IDB_HIST_SMALL_COLOR)
	upImg, _ := tbIml.AddIcon(hIcon)

	navHist := &NavHistory{}
	updateTbState := func() {
		tb.EnableItem(1, navHist.CanBack())
		tb.EnableItem(2, navHist.CanForward())
		tb.EnableItem(3, tree.GetSelectedShellItem().parent != nil)
	}
	tb.AddItems([]*forms.ToolBarItem{
		{
			Id:    1,
			Text:  "Back",
			Image: int(win32.HIST_BACK),
			Action: func() {
				si := navHist.Back()
				tree.SetSelectedShellItem(si)
			},
		},
		{
			Id:    2,
			Text:  "Forward",
			Image: int(win32.HIST_FORWARD),
			Action: func() {
				si := navHist.Forward()
				tree.SetSelectedShellItem(si)
			},
		},
		{
			Id:    3,
			Text:  "Up",
			Image: upImg,
			Action: func() {
				si := tree.GetSelectedShellItem()
				if si.parent != nil {
					tree.SetSelectedShellItem(si.parent)
					for list.Loading {
						forms.DoEvents()
					}
					list.SetSelectedShellItem(si)
				}
			},
		},
	})

	txtAddr := forms.NewEdit{Parent: form}.Create()

	toolbarLayout := &layouts.LinearLayout{
		DebugName: "TB",
		Items: []*layouts.LinearItem{
			{
				Control:      tb,
				PaddingRight: 3,
			},
			{
				Control: txtAddr, Weight: 1,
			},
		},
	}

	sb := forms.NewStatusBar{Parent: form, Simple: true}.Create()

	form.SetLayout(&layouts.LinearLayout{
		Vertical: true,
		Items: []*layouts.LinearItem{
			{
				ItemName:      "toolbar",
				Layout:        toolbarLayout,
				Padding:       6,
				PaddingBottom: consts.Zero,
			},
			{
				ItemName:      "treeandlist",
				Layout:        treeAndListLayout,
				Weight:        1,
				Padding:       6,
				PaddingBottom: consts.Zero,
			},
			{
				ItemName: "statusbar",
				Control:  sb,
			},
		},
	})

	var showingTreeMenu bool
	tree.OnSelectionChange.AddListener(func(ei *events.ExtraEventInfo) {
		if showingTreeMenu {
			return
		}
		si := tree.GetSelectedShellItem()
		list.Load(si)
		navHist.Push(si)
		updateTbState()

		addr := si.GetPath()
		txtAddr.SetText(addr)
	})

	tree.OnRightClick.AddListener(func(ei *events.SimpleEventInfo) {
		showingTreeMenu = true
		defer func() {
			showingTreeMenu = false
		}()
		oriSelId := tree.GetSelectedId()
		id := tree.GetItemIdAtCursor()
		if id != oriSelId {
			tree.SetSelectedId(id)
		}
		si := tree.GetSelectedShellItem()
		if si != nil {
			si.ShowContextMenu()
			ei.SetHandled(true)
		}
		if id != oriSelId {
			tree.SetSelectedId(oriSelId)
		}
	})

	list.OnItemDblClick.AddListener(func(ei *forms.ListViewItemEventInfo) {
		id := list.IndexToId(ei.Index)
		si := list.GetRowData(id).(*ShellItem)
		if si.IsFolder {
			tree.SetSelectedShellItem(si)
		} else {
			si.Open()
		}
	})

	list.OnLoaded.AddListener(func(ei *events.SimpleEventInfo) {
		status := fmt.Sprintf("%d Items", list.GetRowCount())
		sb.SetSimpleText(status)
	})

	list.GetEvent(win32.WM_CONTEXTMENU).AddListener(func(ei *forms.Message) {
		si := list.GetSelectedShellItem()
		if si != nil {
			si.ShowContextMenu()
			ei.SetHandled(true)
		}
	})

	//form.Update()
	form.Show()

	forms.MessageLoop()
}
