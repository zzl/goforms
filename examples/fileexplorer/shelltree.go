package main

import (
	"bytes"
	"github.com/zzl/go-com/com"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/forms"
	"github.com/zzl/goforms/framework/consts"
)

type ShellTree struct {
	forms.TreeViewObject
	super *forms.TreeViewObject

	shell              *Shell
	nodeIdShellItemMap map[int]*ShellItem
}

func (this *ShellTree) Init() {
	this.super.Init()
	this.nodeIdShellItemMap = make(map[int]*ShellItem)
}

func (this *ShellTree) PostCreate(opts *forms.WindowOptions) {
	this.super.PostCreate(opts)

	hIml := this.shell.GetImageList()
	this.SetImageList(forms.NewImageListFromHandle(hIml, false))

	siDesktop := this.shell.GetDesktop()

	id := this.AddItem(forms.TreeItem{
		Text:     siDesktop.Name,
		Image:    siDesktop.Icon,
		HasChild: true,
	})
	this.nodeIdShellItemMap[id] = siDesktop
	this.ExpandItem(this.GetItem(id))
	//this.SetSelectedId(id)
}

func (this *ShellTree) Test() {
	win32.SendMessage(this.Handle, win32.TVM_SELECTITEM,
		win32.WPARAM(win32.TVGN_CARET), 0)
}

func (this *ShellTree) LoadLazyChildItems(parentId int) []forms.TreeItem {
	com.Initialize()

	var items []forms.TreeItem

	siParent := this.nodeIdShellItemMap[parentId]
	siChildren := siParent.GetChildren(false)
	for _, si := range siChildren {
		var item forms.TreeItem
		id := this.GenId()
		item.Id = id
		item.Text = si.Name
		item.Image = si.Icon
		item.HasChild = si.HasChild
		items = append(items, item)

		this.nodeIdShellItemMap[id] = si
	}
	return items
}

func (this *ShellTree) GetSelectedShellItem() *ShellItem {
	id := this.GetSelectedId()
	if id == consts.Null {
		return nil
	}
	return this.nodeIdShellItemMap[id]
}

func (this *ShellTree) SetSelectedShellItem(si *ShellItem) {
	//si.getAbsPidl() //???
	pathItems := si.GetPathItems()
	hParent := win32.TVI_ROOT
	for len(pathItems) > 0 {
		pathSi := pathItems[0]
		pathItems = pathItems[1:]
		hChildren := this.GetChildItems(hParent)
		var hPathItem win32.HTREEITEM
		for _, hChild := range hChildren {
			id := this.GetItemId(hChild)
			siChild := this.nodeIdShellItemMap[id]
			if bytes.Equal(siChild.id, pathSi.id) {
				hPathItem = hChild
				break
			}
		}
		if hPathItem == 0 {
			//?
			break
		} else {

		}
		this.ExpandItem(hPathItem)
		for len(this.LazyExpandingItems) != 0 {
			forms.DoEvents()
		}
		hParent = hPathItem
	}
	this.SetSelectedItem(hParent)
	//for _,
	println("?")
}

func (this *ShellTree) HighlightItem(id int) {
	//
}
