package forms

import (
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"math"
	"slices"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type TreeView interface {
	Control

	UseDefaultNodeImages()
	AddIdTextNodes(nodes []IdTextNode)
	GetNoScrollSize() (int, int)

	GetSelectedItem() win32.HTREEITEM
	GetSelectedId() int
	GetSelectedText() string
	SetSelectedId(id int)
	GetItemIdAtCursor() int

	IsLeaf(hItem win32.HTREEITEM) bool

	GetOnNodeDblClick() *SimpleEvent
	GetOnSelectionChange() *ExtraEvent

	ExpandAll()
	TreeViewObj() *TreeViewObject

	AddNode(node TreeNode) int
	AddNodes(nodes ...TreeNode) []int
	AddChildNode(parentId int, node TreeNode) int

	AddItem(item TreeItem) int
	AddItems(items []TreeItem)
	ClearItems()
}

type TreeViewSpi interface {
	ControlSpi

	LoadLazyChildItems(parentId int) []TreeItem
}

type TreeViewObject struct {
	ControlObject
	super *ControlObject

	Checkable                bool
	DefaultNodeImage         int
	DefaultNodeSelectedImage int

	NoLines       bool
	NoLinesAtRoot bool
	NoButtons     bool

	LazyChildrenLoadCallback func(parentId int) []TreeItem

	OnNodeDblClick    SimpleEvent
	OnRightClick      SimpleEvent
	OnSelectionChange ExtraEvent

	nodeIdHandleMap map[int]win32.HTREEITEM
	uidGen          *UidGen

	iml *ImageList

	//LoadingLazyChildren bool
	LazyExpandingItems map[win32.HTREEITEM]bool
}

type NewTreeView struct {
	Parent Container
	Name   string
	Pos    Point
	Size   Size

	Checkable     bool
	NoLines       bool
	NoLinesAtRoot bool
	NoButtons     bool

	ImageList *ImageList
	Items     []TreeItem
	Nodes     []TreeNode
}

func (me NewTreeView) Create(extraOpts ...*WindowOptions) TreeView {
	tree := NewTreeViewObject()
	tree.name = me.Name
	tree.Checkable = me.Checkable
	tree.NoLines = me.NoLines
	tree.NoLinesAtRoot = me.NoLinesAtRoot
	tree.NoButtons = me.NoButtons

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := tree.Create(*opts)
	assertNoErr(err)
	configControlSize(tree, me.Size)

	if me.ImageList != nil {
		tree.SetImageList(me.ImageList)
	}

	if len(me.Items) != 0 {
		tree.AddItems(me.Items)
	}

	if len(me.Nodes) != 0 {
		tree.AddNodes(me.Nodes...)
	}

	return tree
}

func NewTreeViewObject() *TreeViewObject {
	return virtual.New[TreeViewObject]()
}

type TreeItem struct {
	Id            int
	ParentId      int
	Image         int
	SelectedImage int
	Text          string
	HasChild      bool
	Expanded      bool
}

type IdTextNode struct {
	Id       int
	Text     string
	Children []IdTextNode
}

const LazyNodeId int = math.MinInt32 + 'L' + 'A' + 'Z' + 'Y'

var LazyIdTextNode = IdTextNode{Id: LazyNodeId}
var LazyIdTextNodes = []IdTextNode{LazyIdTextNode}

type TextNode struct {
	Text     string
	Children []TextNode
}

type TreeNode struct {
	Id            int
	Text          string
	Image         int
	SelectedImage int
	Children      []TreeNode
	Expanded      bool
}

func (this *TreeViewObject) TreeViewObj() *TreeViewObject {
	return this
}

func (this *TreeViewObject) Init() {
	this.super.Init()
	this.nodeIdHandleMap = make(map[int]win32.HTREEITEM)
	this.uidGen = NewUidGen(-1, -1)
	this.LazyExpandingItems = make(map[win32.HTREEITEM]bool)
}

func (this *TreeViewObject) GetOnNodeDblClick() *SimpleEvent {
	return &this.OnNodeDblClick
}

func (this *TreeViewObject) GetOnSelectionChange() *ExtraEvent {
	return &this.OnSelectionChange
}

func (this *TreeViewObject) GetWindowClass() string {
	return "SysTreeView32"
}

//todo: handle id recycling

func (this *TreeViewObject) Dispose() {
	if this.iml != nil {
		this.iml.Dispose()
	}
	this.super.Dispose()
}

func (this *TreeViewObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	var style WINDOW_STYLE
	if !this.NoButtons {
		style |= WINDOW_STYLE(win32.TVS_HASBUTTONS)
	}
	if !this.NoLines {
		style |= WINDOW_STYLE(win32.TVS_HASLINES)
	}
	style |= WINDOW_STYLE(win32.TVS_SHOWSELALWAYS)
	if !this.NoLinesAtRoot {
		style |= WINDOW_STYLE(win32.TVS_LINESATROOT)
	}
	return style, 0
}

func (this *TreeViewObject) GetDefaultExStyle() WINDOW_EX_STYLE {
	return win32.WS_EX_CLIENTEDGE
}

func (this *TreeViewObject) Create(options WindowOptions) error {

	err := this.super.Create(options)

	if this.Checkable {
		hImlState := this.createCheckBoxes(4)
		SendMessage(this.Handle, win32.TVM_SETIMAGELIST,
			win32.TVSIL_STATE, hImlState)
	}

	if this.iml != nil {
		SendMessage(this.Handle, win32.TVM_SETIMAGELIST,
			win32.TVSIL_NORMAL, this.iml.GetHandle())
	}

	return err
}

func (this *TreeViewObject) OnReflectNotify(msg *NotifyMessage) {
	pNmhdr := msg.GetNMHDR()
	code := pNmhdr.Code
	switch code {
	case win32.TVN_DELETEITEM:
		pNmtv := (*win32.NMTREEVIEW)(unsafe.Pointer(pNmhdr))
		id := int(pNmtv.ItemOld.LParam)
		delete(this.nodeIdHandleMap, id)
		if id < 0 {
			this.uidGen.Recycle(id)
		}
	case win32.NM_CLICK:
		dwPos := win32.GetMessagePos()
		var ht win32.TVHITTESTINFO
		ht.Pt.X, ht.Pt.Y = win32.GET_X_LPARAM(dwPos), win32.GET_Y_LPARAM(dwPos)
		win32.MapWindowPoints(win32.HWND_DESKTOP, this.Handle, &ht.Pt, 1)
		SendMessage(this.Handle, win32.TVM_HITTEST,
			0, unsafe.Pointer(&ht))
		if ht.Flags&win32.TVHT_ONITEMSTATEICON == win32.TVHT_ONITEMSTATEICON {
			this.onCheckBoxClick(ht.HItem)
		}
	case win32.TVN_KEYDOWN:
		pNm := (*win32.NMTVKEYDOWN)(unsafe.Pointer(pNmhdr))
		if this.Checkable && pNm.WVKey == uint16(win32.VK_SPACE) {
			hItem := this.GetSelectedItem()
			this.onCheckBoxClick(hItem)
		}
	case win32.TVN_ITEMEXPANDING:
		pNmtv := (*win32.NMTREEVIEW)(unsafe.Pointer(pNmhdr))
		hItem := pNmtv.ItemNew.HItem
		if this.isLazyItem(hItem) { //avoid load twice?
			//if this.LazyChildrenLoadCallback == nil {
			//	this.SetItemNoChildren(pNmtv.ItemNew.HItem)
			//} else {
			msg.Result = 1
			msg.Handled = true
			//this.LoadingLazyChildren = true
			this.LazyExpandingItems[hItem] = true
			go this.expandLazyItem(hItem, int(pNmtv.ItemNew.LParam))
			//}
		}
	case win32.NM_SETCURSOR:
		if len(this.LazyExpandingItems) != 0 {
			hCursor, _ := win32.LoadCursor(0, win32.IDC_WAIT)
			win32.SetCursor(hCursor)
			msg.Result = 1
			msg.Handled = true
		}
	case win32.NM_DBLCLK:
		this.OnNodeDblClick.Fire(this, &SimpleEventInfo{})
	case win32.NM_RCLICK:
		this.OnRightClick.Fire(this, &SimpleEventInfo{})
	case win32.TVN_SELCHANGED:
		pNmtv := (*win32.NMTREEVIEW)(unsafe.Pointer(pNmhdr))
		this.OnSelectionChange.Fire(this, &ExtraEventInfo{
			Extra: map[string]interface{}{
				"itemOld": pNmtv.ItemOld,
				"itemNew": pNmtv.ItemNew,
			},
		})
	}

}
func (this *TreeViewObject) ExpandItem(hItem win32.HTREEITEM) {
	SendMessage(this.Handle, win32.TVM_EXPAND, uint32(win32.TVE_EXPAND), hItem)
}

func (this *TreeViewObject) LoadLazyChildItems(parentId int) []TreeItem {
	if this.LazyChildrenLoadCallback != nil {
		return this.LazyChildrenLoadCallback(parentId)
	} else {
		return nil
	}
}

func (this *TreeViewObject) expandLazyItem(hItem win32.HTREEITEM, id int) {
	childItems := this.RealObject.(TreeViewSpi).LoadLazyChildItems(id)
	Dispatcher.Invoke(func() {
		for n := range childItems {
			childItem := &childItems[n]
			if childItem.Id == 0 {
				childItem.Id = this.GenId()
			}
			childItem.ParentId = id
		}
		if len(childItems) == 0 {
			this.SetItemNoChildren(hItem)
		} else {
			this.AddItems(childItems)
			this.ExpandItem(hItem)
		}
		delete(this.LazyExpandingItems, hItem)
		win32.PostMessage(this.Handle, win32.WM_SETCURSOR, 0, 0)
	})
}

func (this *TreeViewObject) GenId() int {
	return this.uidGen.Gen()
}

func (this *TreeViewObject) SetItemNoChildren(hItem win32.HTREEITEM) {
	var tvi win32.TVITEM
	tvi.Mask = win32.TVIF_HANDLE | win32.TVIF_CHILDREN
	tvi.HItem = hItem
	SendMessage(this.Handle, win32.TVM_SETITEM, 0, unsafe.Pointer(&tvi))
}

func (this *TreeViewObject) isLazyItem(hItem win32.HTREEITEM) bool {
	ret, _ := SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
		win32.TVGN_CHILD, hItem)
	return ret == 0
}

func INDEXTOSTATEIMAGEMASK(index int) win32.TREE_VIEW_ITEM_STATE_FLAGS {
	n := (index) << 12
	return win32.TREE_VIEW_ITEM_STATE_FLAGS(n)
}

func STATEIMAGEMASKTOINDEX(state uint32) int {
	n := int(state >> 12)
	return n
}

func (this *TreeViewObject) onCheckBoxClick(hItem win32.HTREEITEM) {
	id := this.GetItemId(hItem)
	println("Check@", id)

	state := this.GetItemCheckState(hItem)
	if state == 1 {
		state = 2
	} else {
		state = 1
	}
	this.SetItemCheckState(hItem, state)
	descendantItems := this.getDescendantItems(hItem)
	for _, item := range descendantItems {
		this.SetItemCheckState(item, state)
	}
	ancestorItems := this.getAncestorItems(hItem)
	for _, item := range ancestorItems {
		cItems := this.GetChildItems(item)
		state := this.GetItemsCheckState(cItems)
		this.SetItemCheckState(item, state)
	}
}

func (this *TreeViewObject) GetItemsCheckState(hItems []win32.HTREEITEM) int {
	checkedCount, uncheckedCount := 0, 0
	for _, hItem := range hItems {
		state := this.GetItemCheckState(hItem)
		if state == 1 {
			uncheckedCount += 1
		} else if state == 2 {
			checkedCount += 1
		}
	}
	totalCount := len(hItems)
	if uncheckedCount == totalCount {
		return 1
	}
	if checkedCount == totalCount {
		return 2
	}
	return 3
}

func (this *TreeViewObject) GetItemCheckState(hItem win32.HTREEITEM) int {
	ret, _ := SendMessage(this.Handle, win32.TVM_GETITEMSTATE,
		hItem, uint32(win32.TVIS_STATEIMAGEMASK))
	state := STATEIMAGEMASKTOINDEX(uint32(ret))
	return state
}

func (this *TreeViewObject) getAncestorItems(
	hItem win32.HTREEITEM) []win32.HTREEITEM {
	var aItems []win32.HTREEITEM
	hParentItem := hItem
	for {
		ret, _ := SendMessage(this.Handle,
			win32.TVM_GETNEXTITEM, win32.TVGN_PARENT, hParentItem)
		if ret == 0 {
			break
		}
		hParentItem = win32.HTREEITEM(ret)
		aItems = append(aItems, hParentItem)
	}
	return aItems
}

func (this *TreeViewObject) getDescendantItems(
	hRootItem win32.HTREEITEM) []win32.HTREEITEM {
	var items []win32.HTREEITEM
	items = append(items, hRootItem)
	lastItemCount := 0
	var curItemCount int
	for {
		curItemCount = len(items)
		for n := lastItemCount; n < curItemCount; n++ {
			cItems := this.GetChildItems(items[n])
			if len(cItems) > 0 {
				items = append(items, cItems...)
			}
		}
		lastItemCount = curItemCount
		curItemCount = len(items)
		if lastItemCount == curItemCount {
			break
		}
	}
	return items
}

func (this *TreeViewObject) GetChildItems(
	hParentItem win32.HTREEITEM) []win32.HTREEITEM {
	var cItems []win32.HTREEITEM
	hItem, _ := SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
		win32.TVGN_CHILD, hParentItem)
	for hItem != 0 {
		cItems = append(cItems, hItem)
		hItem, _ = SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
			win32.TVGN_NEXT, hItem)
	}
	return cItems
}

func (this *TreeViewObject) SetItemCheckState(hItem win32.HTREEITEM, state int) {
	var tvi win32.TVITEM
	tvi.Mask = win32.TVIF_HANDLE | win32.TVIF_STATE
	tvi.StateMask = win32.TVIS_STATEIMAGEMASK
	tvi.HItem = hItem
	tvi.State = INDEXTOSTATEIMAGEMASK(state)
	ret, errno := SendMessage(this.Handle,
		win32.TVM_SETITEM, 0, unsafe.Pointer(&tvi))
	if ret == 0 {
		log.Println(errno.Error())
	}
}

func (this *TreeViewObject) ClearItems() {
	SendMessage(this.Handle, win32.TVM_DELETEITEM, 0, 0)
	this.nodeIdHandleMap = make(map[int]win32.HTREEITEM)
}

// return id
func (this *TreeViewObject) AddItem(item TreeItem) int {
	hParent := win32.TVI_ROOT
	if item.ParentId != 0 {
		hParent, _ = this.nodeIdHandleMap[item.ParentId]
	}
	var id = item.Id
	if id == 0 {
		id = this.uidGen.Gen()
	}
	var tvis win32.TVINSERTSTRUCT
	var tvi = tvis.Item()                         //Item
	tvi.Mask = win32.TVIF_TEXT | win32.TVIF_PARAM // |
	//win32.TVIF_IMAGE | win32.TVIF_SELECTEDIMAGE

	this.fillImageFields(tvi, item)

	//
	if this.Checkable {
		tvi.Mask |= win32.TVIF_STATE
		tvi.StateMask = win32.TVIS_STATEIMAGEMASK
		tvi.State = INDEXTOSTATEIMAGEMASK(1)
	}

	if item.HasChild {
		tvi.Mask |= win32.TVIF_CHILDREN
		tvi.CChildren = 1
	}

	if item.Expanded {
		tvi.Mask |= win32.TVIF_STATE
		tvi.StateMask |= win32.TVIS_EXPANDED
		tvi.State |= win32.TVIS_EXPANDED
	}

	//if item.Image
	pwsz, _ := syscall.UTF16PtrFromString(item.Text)
	tvi.PszText = pwsz
	tvi.LParam = uintptr(id)

	tvis.HParent = hParent

	ret, errno := SendMessage(this.Handle, win32.TVM_INSERTITEM,
		0, unsafe.Pointer(&tvis))
	if ret == 0 {
		log.Panic(errno)
	}
	hItem := ret
	this.nodeIdHandleMap[id] = hItem
	return id
}

func (this *TreeViewObject) fillImageFields(tvi *win32.TVITEM, item TreeItem) {
	image := item.Image
	if image == 0 {
		image = this.DefaultNodeImage
	}
	if image == 0 {
		return
	}
	if image == consts.Null {
		//?
	} else if image == consts.Zero {
		image = 0
	}
	tvi.Mask |= win32.TVIF_IMAGE
	tvi.IImage = int32(image)

	selectedImage := item.SelectedImage
	if selectedImage == 0 {
		selectedImage = this.DefaultNodeSelectedImage
	}
	if selectedImage != 0 { //?
		if selectedImage == consts.Zero {
			selectedImage = 0
		}
		tvi.Mask |= win32.TVIF_SELECTEDIMAGE
		tvi.ISelectedImage = int32(selectedImage)
	} else {
		tvi.Mask |= win32.TVIF_SELECTEDIMAGE
		tvi.ISelectedImage = tvi.IImage
	}
}

func (this *TreeViewObject) AddItems(items []TreeItem) {
	for _, item := range items {
		this.AddItem(item)
	}
}

func (this *TreeViewObject) AddChildNode(parentId int, node TreeNode) int {
	if node.Id == 0 {
		node.Id = this.uidGen.Gen() //?
	}
	item := TreeItem{
		Id:            node.Id,
		ParentId:      parentId,
		Image:         node.Image,
		SelectedImage: node.SelectedImage,
		Text:          node.Text,
		Expanded:      node.Expanded,
	}
	lazyNode := isLazyNode(node)
	if lazyNode {
		item.HasChild = true
	}
	this.AddItem(item)
	if !lazyNode {
		for _, childNode := range node.Children {
			this.AddChildNode(node.Id, childNode)
		}
	}
	return node.Id
}

func (this *TreeViewObject) AddNode(node TreeNode) int {
	return this.AddChildNode(0, node)
}

func (this *TreeViewObject) AddNodes(nodes ...TreeNode) []int {
	var ids []int
	for _, node := range nodes {
		id := this.AddNode(node)
		ids = append(ids, id)
	}
	return ids
}

func isLazyIdTextNode(node IdTextNode) bool {
	return len(node.Children) == 1 && node.Children[0].Id == LazyNodeId
}

func isLazyNode(node TreeNode) bool {
	return len(node.Children) == 1 && node.Children[0].Id == LazyNodeId
}

func (this *TreeViewObject) AddChildIdTextNode(parentId int, node IdTextNode) {
	item := TreeItem{Id: node.Id, ParentId: parentId, Text: node.Text}
	lazyNode := isLazyIdTextNode(node)
	item.HasChild = lazyNode
	this.AddItem(item)
	if !lazyNode {
		for _, childNode := range node.Children {
			this.AddChildIdTextNode(node.Id, childNode)
		}
	}
}

func (this *TreeViewObject) AddIdTextNode(node IdTextNode) {
	this.AddChildIdTextNode(0, node)
}

func (this *TreeViewObject) AddIdTextNodes(nodes []IdTextNode) {
	for _, node := range nodes {
		this.AddIdTextNode(node)
	}
}

func (this *TreeViewObject) AddChildTextNode(parentId int, node TextNode) {
	id := this.AddItem(TreeItem{ParentId: parentId, Text: node.Text})
	for _, childNode := range node.Children {
		this.AddChildTextNode(id, childNode)
	}
}

func (this *TreeViewObject) AddTextNode(node TextNode) {
	this.AddChildTextNode(0, node)
}

func (this *TreeViewObject) AddTextNodes(nodes []TextNode) {
	for _, node := range nodes {
		this.AddTextNode(node)
	}
}

func (this *TreeViewObject) SetImageList(iml *ImageList) {
	if this.Handle != 0 {
		SendMessage(this.Handle, win32.TVM_SETIMAGELIST, 0, iml.GetHandle())
	}
	if this.iml != nil {
		this.iml.Dispose()
	}
	this.iml = iml
}

func (this *TreeViewObject) UseDefaultNodeImages() {

	var ssii win32.SHSTOCKICONINFO
	ssii.CbSize = uint32(unsafe.Sizeof(ssii))
	flags := win32.SHGSI_ICON | win32.SHGSI_SMALLICON
	hr := win32.SHGetStockIconInfo(win32.SIID_FOLDER, flags, &ssii)
	if !win32.SUCCEEDED(hr) {
		log.Fatal("?")
	}
	hIconFolder := ssii.HIcon
	_ = win32.SHGetStockIconInfo(win32.SIID_FOLDEROPEN, flags, &ssii)
	hIconFolderOpen := ssii.HIcon

	iml, _ := CreateImageListFromIcons(16, hIconFolder, hIconFolderOpen)

	win32.DestroyIcon(hIconFolder)
	win32.DestroyIcon(hIconFolderOpen)

	this.SetImageList(iml)

	this.DefaultNodeImage = consts.Zero
	this.DefaultNodeSelectedImage = 1
}

func (this *TreeViewObject) GetItemId(hItem win32.HTREEITEM) int {
	var tvi win32.TVITEM
	tvi.HItem = hItem
	tvi.Mask = win32.TVIF_PARAM
	ret, errno := SendMessage(this.Handle,
		win32.TVM_GETITEM, 0, unsafe.Pointer(&tvi))
	if ret == 0 {
		log.Println(errno.Error())
	}
	id := int(tvi.LParam)
	return id
}

func (this *TreeViewObject) GetSelectedItem() win32.HTREEITEM {
	ret, errno := win32.SendMessage(this.Handle,
		win32.TVM_GETNEXTITEM, WPARAM(win32.TVGN_CARET), 0)
	_ = errno
	return win32.HTREEITEM(ret)
}

func (this *TreeViewObject) GetSelectedId() int {
	hItem := this.GetSelectedItem()
	if hItem == 0 {
		return consts.Null
	}
	return this.GetItemId(hItem)
}

func (this *TreeViewObject) GetSelectedText() string {
	hItem := this.GetSelectedItem()
	if hItem == 0 {
		return ""
	}
	text := this.GetItemText(hItem)
	return text
}

func (this *TreeViewObject) SetSelectedItem(hItem win32.HTREEITEM) {
	win32.SendMessage(this.Handle, win32.TVM_SELECTITEM,
		WPARAM(win32.TVGN_CARET), hItem)
}

func (this *TreeViewObject) SetSelectedId(id int) {
	var hItem win32.HTREEITEM
	if id != consts.Null {
		var ok bool
		hItem, ok = this.nodeIdHandleMap[id]
		if !ok {
			return //?
		}
	}
	SendMessage(this.Handle, win32.TVM_SELECTITEM,
		win32.TVGN_CARET, hItem)
}

func (this *TreeViewObject) GetNextItem(hItem win32.HTREEITEM) win32.HTREEITEM {
	ret, errno := SendMessage(this.Handle,
		win32.TVM_GETNEXTITEM, win32.TVGN_NEXT, hItem)
	if ret == 0 {
		_ = errno
		return 0
	}
	hItem = win32.HTREEITEM(ret)
	return hItem
}

func (this *TreeViewObject) GetPrevItem(hItem win32.HTREEITEM) win32.HTREEITEM {
	ret, errno := SendMessage(this.Handle,
		win32.TVM_GETNEXTITEM, win32.TVGN_PREVIOUS, hItem)
	if ret == 0 {
		_ = errno
		return 0
	}
	hItem = win32.HTREEITEM(ret)
	return hItem
}

// ?
func (this *TreeViewObject) GetParentItem(id int) win32.HTREEITEM {
	hItem, ok := this.nodeIdHandleMap[id]
	if !ok {
		return 0
	}
	ret, errno := SendMessage(this.Handle,
		win32.TVM_GETNEXTITEM, win32.TVGN_PARENT, hItem)
	if ret == 0 {
		_ = errno
		//root? log.Println(errno.Error())
		return 0
	}
	hItem = win32.HTREEITEM(ret)
	return hItem
}

func (this *TreeViewObject) GetParentId(id int) int {
	hItem := this.GetParentItem(id)
	if hItem == win32.TVI_ROOT || hItem == 0 {
		return 0
	}
	var tvi win32.TVITEM
	tvi.HItem = hItem
	tvi.Mask = win32.TVIF_PARAM
	ret, errno := SendMessage(this.Handle,
		win32.TVM_GETITEM, 0, unsafe.Pointer(&tvi))
	if ret == 0 {
		log.Fatal(errno)
	}
	parentId := int(tvi.LParam)
	return parentId
}

func (this *TreeViewObject) GetPathIds(id int) []int {
	var pathIds []int
	for id != 0 {
		pathIds = append(pathIds, id)
		id = this.GetParentId(id)
	}
	slices.Reverse(pathIds)
	return pathIds
}

func (this *TreeViewObject) GetItem(id int) win32.HTREEITEM {
	hItem, _ := this.nodeIdHandleMap[id]
	return hItem
}

func (this *TreeViewObject) GetText(id int) string {
	hItem, _ := this.nodeIdHandleMap[id]
	return this.GetItemText(hItem)
}

func (this *TreeViewObject) createCheckBoxes(rightSpace int) win32.HIMAGELIST {
	hdcWin := win32.GetDC(this.Handle)
	pwsz, _ := syscall.UTF16PtrFromString("button")
	hTheme := win32.OpenThemeData(this.Handle, pwsz)
	var frames = 4
	cx, cy := 14, 14
	if hTheme != 0 {
		var size win32.SIZE
		win32.GetThemePartSize(hTheme, hdcWin,
			int32(win32.BP_CHECKBOX), int32(win32.CBS_UNCHECKEDNORMAL), nil,
			win32.TS_DRAW, &size)
		cx, cy = int(size.Cx), int(size.Cy)
	}
	cxSpace := cx + rightSpace
	var bi = win32.BITMAPINFO{
		BmiHeader: win32.BITMAPINFOHEADER{
			BiWidth:    int32(cxSpace * frames),
			BiHeight:   int32(cy),
			BiPlanes:   1,
			BiBitCount: 32,
		},
	}
	bi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bi.BmiHeader))
	var p unsafe.Pointer
	hbmCheckboxes, _ := win32.CreateDIBSection(hdcWin, &bi,
		win32.DIB_RGB_COLORS, unsafe.Pointer(&p), 0, 0)
	hdcMem := win32.CreateCompatibleDC(hdcWin)
	hbmOld := win32.SelectObject(hdcMem, win32.HGDIOBJ(hbmCheckboxes))

	rcBmp := win32.RECT{Left: 0, Top: 0,
		Right: bi.BmiHeader.BiWidth, Bottom: bi.BmiHeader.BiHeight}
	hbr := win32.GetSysColorBrush(win32.COLOR_WINDOW)
	win32.FillRect(hdcMem, &rcBmp, hbr)

	rc := win32.RECT{Left: 0, Top: 0, Right: int32(cx), Bottom: int32(cy)}
	win32.OffsetRect(&rc, int32(cxSpace), 0)
	if hTheme != 0 {
		//F1
		win32.DrawThemeBackground(hTheme, hdcMem, int32(win32.BP_CHECKBOX),
			int32(win32.CBS_UNCHECKEDNORMAL), &rc, nil)
		win32.OffsetRect(&rc, int32(cxSpace), 0)

		//F2
		win32.DrawThemeBackground(hTheme, hdcMem, int32(win32.BP_CHECKBOX),
			int32(win32.CBS_CHECKEDNORMAL), &rc, nil)
		win32.OffsetRect(&rc, int32(cxSpace), 0)

		//F3
		win32.DrawThemeBackground(hTheme, hdcMem, int32(win32.BP_CHECKBOX),
			int32(win32.CBS_MIXEDNORMAL), &rc, nil)
		win32.OffsetRect(&rc, int32(cxSpace), 0)

		win32.CloseThemeData(hTheme)
	} else {
		var baseFlags win32.DFCS_STATE = win32.DFCS_FLAT | win32.DFCS_BUTTONCHECK

		//F1
		win32.DrawFrameControl(hdcMem, &rc, win32.DFC_BUTTON, baseFlags)
		win32.OffsetRect(&rc, int32(cxSpace), 0)

		//F2
		win32.DrawFrameControl(hdcMem, &rc, win32.DFC_BUTTON,
			baseFlags|win32.DFCS_CHECKED)
		win32.OffsetRect(&rc, int32(cxSpace), 0)

		//F3
		win32.DrawFrameControl(hdcMem, &rc, win32.DFC_BUTTON,
			baseFlags|win32.DFCS_CHECKED|win32.DFCS_BUTTON3STATE)
		win32.OffsetRect(&rc, int32(cxSpace), 0)
	}
	win32.SelectObject(hdcMem, hbmOld)
	win32.DeleteDC(hdcMem)
	win32.ReleaseDC(this.Handle, hdcWin)

	//
	hIml := win32.ImageList_Create(
		int32(cxSpace), int32(cy), win32.ILC_COLOR, int32(frames), int32(frames))
	if hIml == 0 {
		log.Fatal("?")
	}
	win32.ImageList_Add(hIml, hbmCheckboxes, 0)
	return hIml
}

func (this *TreeViewObject) GetItemText(hItem win32.HTREEITEM) string {
	var tvi win32.TVITEM
	tvi.HItem = hItem
	tvi.Mask = win32.TVIF_TEXT

	const bufSize = 1024

	var buf [bufSize]uint16
	tvi.PszText = &buf[0]
	tvi.CchTextMax = bufSize

	ret, errno := SendMessage(this.Handle,
		win32.TVM_GETITEM, 0, unsafe.Pointer(&tvi))
	if ret == 0 {
		log.Println(errno.Error())
	}
	return syscall.UTF16ToString(buf[:])
}

func (this *TreeViewObject) collectExpandedItems(hParentItem win32.HTREEITEM) []uintptr {
	if hParentItem != win32.TVI_ROOT {
		ret, _ := SendMessage(this.Handle, win32.TVM_GETITEMSTATE,
			hParentItem, uint32(win32.TVIS_EXPANDED))
		if uint32(ret)&uint32(win32.TVIS_EXPANDED) == 0 {
			return nil
		}
	}
	var hItems []uintptr
	hItem, _ := SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
		win32.TVGN_CHILD, hParentItem)
	for hItem != 0 {
		childItems := this.collectExpandedItems(win32.HTREEITEM(hItem))
		for _, childItem := range childItems {
			hItems = append(hItems, childItem)
		}
		hItems = append(hItems, hItem)
		hItem, _ = SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
			win32.TVGN_NEXT, hItem)
	}
	return hItems
}

func (this *TreeViewObject) GetNoScrollSize() (int, int) {
	hItems := this.collectExpandedItems(win32.TVI_ROOT)
	ret, _ := SendMessage(this.Handle, win32.TVM_GETITEMHEIGHT, 0, 0)
	itemHeight := int(ret)
	cy := itemHeight * len(hItems)
	cxMax := 0
	var rcItem win32.RECT
	var maxRight int32
	for _, hItem := range hItems {
		*(*uintptr)(unsafe.Pointer(&rcItem)) = hItem
		SendMessage(this.Handle, win32.TVM_GETITEMRECT,
			1, unsafe.Pointer(&rcItem))
		if rcItem.Right > maxRight {
			maxRight = rcItem.Right
		}
	}
	cxMax = int(maxRight)
	return cxMax + 8, cy + 8
}

func (this *TreeViewObject) IsLeaf(hItem win32.HTREEITEM) bool {
	var tvi win32.TVITEM
	tvi.Mask = win32.TVIF_CHILDREN | win32.TVIF_HANDLE
	tvi.HItem = hItem
	ret, errno := SendMessage(this.Handle, win32.TVM_GETITEM,
		0, unsafe.Pointer(&tvi))
	if ret == 0 {
		log.Println(errno)
	}
	return tvi.CChildren == 0
}

func (this *TreeViewObject) IsParent(id int) bool {
	hItem, _ := this.nodeIdHandleMap[id]
	return !this.IsLeaf(hItem)
}

func (this *TreeViewObject) ExpandAll() {
	items := []win32.HTREEITEM{win32.TVI_ROOT}
	hItemTop, _ := SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
		win32.TVGN_FIRSTVISIBLE, 0)
	for {
		count := len(items)
		if count == 0 {
			break
		}
		hItem := items[count-1]
		items = items[:count-1]
		hChildItem, _ := SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
			win32.TVGN_CHILD, hItem)
		if hChildItem != 0 {
			this.ExpandItem(hItem)
		}
		for hChildItem != 0 {
			items = append(items, hChildItem)
			hChildItem, _ = SendMessage(this.Handle, win32.TVM_GETNEXTITEM,
				win32.TVGN_NEXT, hChildItem)
		}
	}
	win32.SendMessage(this.Handle, win32.TVM_ENSUREVISIBLE, 0, hItemTop)
	println("?")
}

func (this *TreeViewObject) GetItemIdAtCursor() int {
	var hti win32.TVHITTESTINFO
	win32.GetCursorPos(&hti.Pt)
	win32.ScreenToClient(this.Handle, &hti.Pt)
	hti.Flags = win32.TVHT_ONITEM
	hItem, _ := win32.SendMessage(this.Handle,
		win32.TVM_HITTEST, 0, win32.LPARAM(unsafe.Pointer(&hti)))
	if hItem != 0 {
		return this.GetItemId(hItem)
	}
	return 0
}
