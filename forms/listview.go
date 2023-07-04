package forms

import (
	"fmt"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"math"
	"syscall"
	"unsafe"

	"github.com/zzl/goforms/framework/types"
	"github.com/zzl/goforms/layouts/aligns"

	"github.com/zzl/go-win32api/v2/win32"
)

type ListView interface {
	Control

	GetRowCount() int
	ListViewObj() *ListViewObject
}

type ListViewObject struct {
	ControlObject
	super *ControlObject

	SingleRowSelect bool
	SortCompareFunc func(row1Id int, row2Id, colIndex int) int

	OnSelectionChange Event[*ListViewItemEventInfo]
	OnItemDblClick    Event[*ListViewItemEventInfo]

	Columns []*ListViewColumn

	iml *ImageList

	hWndHeader       HWND
	lastSortColIndex int

	uidGen  *UidGen
	dataMap map[int]any

	sortingCol   int
	sortingDesc  bool
	pCompareFunc uintptr

	lastSelectedIndex int
}

type NewListView struct {
	Parent  Container
	Name    string
	Pos     Point
	Size    Size
	Columns []*ListViewColumn
	Rows    []*ListRowInfo
}

func (me NewListView) Create(extraOpts ...*WindowOptions) ListView {
	lv := NewListViewObject()
	lv.name = me.Name
	lv.Columns = me.Columns

	opts := utils.OptionalArg(extraOpts)
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := lv.Create(*opts)
	assertNoErr(err)
	configControlSize(lv, me.Size)

	if len(me.Rows) != 0 {
		lv.AddRowInfos(me.Rows)
	}

	return lv
}

func NewListViewObject() *ListViewObject {
	return virtual.New[ListViewObject]()
}

type ListViewItemEventInfo struct {
	SimpleEventInfo
	Index int
}

type ListViewColumn struct {
	Title string
	Align int //textalign..

	Width      int
	MinWidth   int
	FixedWidth bool

	Sortable        bool
	DefaultSortDesc bool
	SortCompareFunc func(row1Id int, row2Id int) int

	Formatter types.Formatter
}

type ListRowInfo struct {
	Id         int
	Image      int
	Data       any
	CellValues []any
}

func (this *ListViewObject) ListViewObj() *ListViewObject {
	return this
}

var _pCompareFunc uintptr

func (this *ListViewObject) Init() {
	this.super.Init()
	this.lastSelectedIndex = -1

	if _pCompareFunc == 0 {
		_pCompareFunc = syscall.NewCallback(listViewCompareFunc)
	}

	this.uidGen = NewUidGen(-1, -1)
	this.dataMap = make(map[int]any)
	this.pCompareFunc = _pCompareFunc
}

func (this *ListViewObject) Dispose() {
	if this.iml != nil {
		this.iml.Dispose()
	}
	this.super.Dispose()
}

func (this *ListViewObject) SetImageList(iml *ImageList) {
	if this.iml != nil {
		this.iml.Dispose()
	}
	this.iml = iml
	if this.Handle == 0 {
		log.Fatal("please create control first")
	}
	var hIml win32.HIMAGELIST
	if iml != nil {
		hIml = iml.GetHandle()
	}
	_, errno := SendMessage(this.Handle, win32.LVM_SETIMAGELIST,
		win32.LVSIL_SMALL, hIml)
	if errno != win32.NO_ERROR {
		log.Fatal(errno)
	}
}

func (this *ListViewObject) GetWindowClass() string {
	return "SysListView32"
}

func (this *ListViewObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	var style WINDOW_STYLE
	style = WINDOW_STYLE(win32.LVS_REPORT | win32.LVS_SHOWSELALWAYS) //|win32.WS_BORDER
	if this.SingleRowSelect {
		style |= WINDOW_STYLE(win32.LVS_SINGLESEL)
	}
	return style, 0
}

func (this *ListViewObject) Create(options WindowOptions) error {
	//options.Style |= win32.WS_TABSTOP
	options.ExStyle |= win32.WS_EX_CLIENTEDGE

	err := this.super.Create(options)

	ret, _ := SendMessage(this.Handle, win32.LVM_GETHEADER, 0, 0)
	this.hWndHeader = HWND(ret)

	SendMessage(this.Handle, win32.LVM_SETEXTENDEDLISTVIEWSTYLE,
		win32.LVS_EX_FULLROWSELECT, win32.LVS_EX_FULLROWSELECT)

	if len(this.Columns) > 0 {
		for _, column := range this.Columns {
			this._addColumn(column)
		}
	}

	return err
}

func (this *ListViewObject) OnReflectNotify(msg *NotifyMessage) {
	this.super.OnReflectNotify(msg)
	pNmhdr := msg.GetNMHDR()
	//println("@@", pNmhdr.Code)
	switch pNmhdr.Code {
	case win32.LVN_DELETEITEM:
		pNmlv := (*win32.NMLISTVIEW)(unsafe.Pointer(pNmhdr))
		id := int(pNmlv.LParam)
		if id != 0 {
			delete(this.dataMap, id)
			if id < 0 {
				this.uidGen.Recycle(id)
			}
		}
	case win32.LVN_COLUMNCLICK:
		pNmlv := (*win32.NMLISTVIEW)(unsafe.Pointer(pNmhdr))
		this.onColumnClick(int(pNmlv.ISubItem))
	//case win32.LVN_ITEMACTIVATE:
	//	pNmia := (*win32.NMITEMACTIVATE)(unsafe.Pointer(pNmhdr))
	//	println(pNmia.IItem, this.GetSelectedIndex())
	case win32.LVN_ITEMCHANGED:
		pNmlv := (*win32.NMLISTVIEW)(unsafe.Pointer(pNmhdr))
		selected := pNmlv.UNewState&(uint32)(win32.LVIS_SELECTED) != 0
		deselected := pNmlv.UOldState&(uint32)(win32.LVIS_SELECTED) != 0
		if selected || deselected {
			if selected {
				this.setLastSelectedIndex(int(pNmlv.IItem))
			}
			Dispatcher.Invoke(func() {
				index := this.GetSelectedIndex()
				this.setLastSelectedIndex(index)
			})
		}
	case win32.NM_DBLCLK:
		this.OnItemDblClick.Fire(this, &ListViewItemEventInfo{
			Index: this.GetSelectedIndex()})
	}
}

func (this *ListViewObject) setLastSelectedIndex(index int) {
	if index != this.lastSelectedIndex {
		this.lastSelectedIndex = index
		this.OnSelectionChange.Fire(this, &ListViewItemEventInfo{Index: index})
	}
}

func listViewCompareFunc(lParam1 LPARAM, lParam2 LPARAM,
	lParamSort LPARAM) win32.LRESULT {

	this := (*ListViewObject)(unsafe.Pointer(lParamSort))
	colIndex := this.sortingCol
	colCompareFunc := this.Columns[colIndex].SortCompareFunc
	var result = 0
	if colCompareFunc != nil {
		result = colCompareFunc(int(lParam1), int(lParam2))
	} else if this.SortCompareFunc != nil {
		result = this.SortCompareFunc(int(lParam1), int(lParam2), colIndex)
	}
	if this.sortingDesc {
		result = -result
	}
	return win32.LRESULT(result)
}

func (this *ListViewObject) onColumnClick(colIndex int) {
	//println("COL CLICK..")
	column := this.Columns[colIndex]
	if !column.Sortable {
		return
	}
	//
	var hdi win32.HDITEM

	hdi.Mask = win32.HDI_FORMAT
	if this.lastSortColIndex != colIndex {
		_, _ = SendMessage(this.hWndHeader, win32.HDM_SETITEM,
			this.lastSortColIndex, unsafe.Pointer(&hdi))
	}

	ret, errno := SendMessage(this.hWndHeader, win32.HDM_GETITEM,
		colIndex, unsafe.Pointer(&hdi))
	if ret == 0 {
		log.Fatal(errno)
	}
	var desc bool
	fmt := hdi.Fmt
	if fmt&win32.HDF_SORTUP == win32.HDF_SORTUP {
		fmt &^= win32.HDF_SORTUP
		fmt |= win32.HDF_SORTDOWN
		desc = true
	} else if fmt&win32.HDF_SORTDOWN == win32.HDF_SORTDOWN {
		fmt &^= win32.HDF_SORTDOWN
		fmt |= win32.HDF_SORTUP
	} else if column.DefaultSortDesc {
		fmt |= win32.HDF_SORTDOWN
		desc = true
	} else {
		fmt |= win32.HDF_SORTUP
	}
	hdi.Fmt = fmt

	this.sortingCol = colIndex
	this.sortingDesc = desc

	ret, errno = SendMessage(this.Handle, win32.LVM_SORTITEMS,
		unsafe.Pointer(this), this.pCompareFunc)

	ret, errno = SendMessage(this.hWndHeader, win32.HDM_SETITEM,
		colIndex, unsafe.Pointer(&hdi))
	if ret == 0 {
		log.Fatal(errno)
	}
	this.lastSortColIndex = colIndex

}

func (this *ListViewObject) GetColumnCount() int {
	ret, _ := SendMessage(this.hWndHeader, win32.HDM_GETITEMCOUNT, 0, 0)
	return int(ret)
}

func (this *ListViewObject) DeleteColumn(index int) {
	ret, errno := SendMessage(this.Handle,
		win32.LVM_DELETECOLUMN, index, 0)
	if ret == 0 {
		log.Fatal(errno)
	}
}

func (this *ListViewObject) ClearColumns() {
	colCount := this.GetColumnCount()
	for n := 0; n < colCount; n++ {
		this.DeleteColumn(n)
	}
	this.Columns = nil
}

func (this *ListViewObject) SetColumns(columns []ListViewColumn) {
	this.ClearColumns()
	for _, column := range columns {
		this.AddColumn(column)
	}
}

func (this *ListViewObject) ClearRows() {
	_, _ = SendMessage(this.Handle, win32.LVM_DELETEALLITEMS, 0, 0)
}

func (this *ListViewObject) DeleteRow(index int) {
	ret, errno := SendMessage(this.Handle,
		win32.LVM_DELETEITEM, index, 0)
	if ret == 0 {
		log.Fatal(errno)
	}
}

// return index,rowid
func (this *ListViewObject) AddRow(info *ListRowInfo) (int, int) {

	var lvi win32.LVITEM
	lvi.Mask = win32.LVIF_TEXT | win32.LVIF_PARAM
	lvi.IItem = math.MaxInt32

	if info.Image >= 0 {
		lvi.Mask |= win32.LVIF_IMAGE
		lvi.IImage = int32(info.Image)
	}

	if info.Id == 0 {
		info.Id = this.uidGen.Gen()
	}

	lvi.LParam = uintptr(info.Id)
	if info.Data != nil {
		this.dataMap[info.Id] = info.Data
	}

	iItem, errno := SendMessage(this.Handle, win32.LVM_INSERTITEM,
		0, unsafe.Pointer(&lvi))
	if iItem == NegativeOne {
		log.Fatal(errno)
	}
	lvi.LParam = 0

	cellValues := info.CellValues
	nCount := min(len(cellValues), len(this.Columns))
	for n := 0; n < nCount; n++ {
		lvi.ISubItem = int32(n)
		text := this.formatCellText(n, cellValues[n])
		pwsz, _ := syscall.UTF16PtrFromString(text)
		lvi.PszText = pwsz
		ret, errno := SendMessage(this.Handle, win32.LVM_SETITEMTEXT,
			iItem, unsafe.Pointer(&lvi))
		if ret == 0 {
			log.Fatal(errno)
		}
	}
	return int(iItem), info.Id
}

// rename?
func (this *ListViewObject) AddRowInfos(rows []*ListRowInfo) {
	for _, row := range rows {
		this.AddRow(row)
	}
}

func (this *ListViewObject) AddRows(rows [][]any, rowAsData bool) {
	for _, row := range rows {
		var data any
		if rowAsData {
			data = row
		}
		this.AddRow(&ListRowInfo{0, -1, data, row})
	}
}

func (this *ListViewObject) AddRowsWithId(rows [][]any, rowAsData bool) {
	for _, row := range rows {
		rowId := row[0].(int)
		var data any
		if rowAsData {
			data = row
		}
		this.AddRow(&ListRowInfo{rowId, -1, data, row[1:]})
	}
}

func (this *ListViewObject) AddRowsWithData(dataRows [][]any) {
	for _, row := range dataRows {
		data := row[0]
		this.AddRow(&ListRowInfo{0, -1, data, row[1:]})
	}
}

func (this *ListViewObject) AddRowsWithIdAndData(dataRows [][]any) {
	for _, row := range dataRows {
		rowId := row[0].(int)
		data := row[1]
		this.AddRow(&ListRowInfo{rowId, -1, data, row[2:]})
	}
}

// ?
func (this *ListViewObject) AddRowsWithInfo(dataRows [][]any) {
	for _, row := range dataRows {
		info := row[0].(*ListRowInfo)
		info.CellValues = row[1:]
		this.AddRow(info)
	}
}

func (this *ListViewObject) formatCellText(colIndex int, value any) string {
	column := this.Columns[colIndex]
	if column.Formatter != nil {
		return column.Formatter.Format(value)
	}
	return fmt.Sprintf("%v", value)
}

func (this *ListViewObject) AddColumn(column ListViewColumn) {
	this.Columns = append(this.Columns, &column)
	if this.Handle == 0 {
		return
	}
	this._addColumn(&column)
}

func (this *ListViewObject) _addColumn(column *ListViewColumn) {
	colCount := this.GetColumnCount()
	var lvc win32.LVCOLUMN
	lvc.Mask = win32.LVCF_WIDTH | win32.LVCF_TEXT | win32.LVCF_FMT
	dummyColumnAdded := false
	if colCount == 0 {
		lvc.Cx = 0
		ret, errno := SendMessage(this.Handle, win32.LVM_INSERTCOLUMN,
			0, unsafe.Pointer(&lvc))
		if ret == NegativeOne {
			log.Fatal(errno)
		}
		dummyColumnAdded = true
	}

	if column.FixedWidth {
		lvc.Fmt |= win32.LVCFMT_FIXED_WIDTH
	}
	lvc.Fmt &^= win32.LVCFMT_JUSTIFYMASK
	if column.Align == aligns.Right {
		lvc.Fmt |= win32.LVCFMT_RIGHT
	} else if column.Align == aligns.Center {
		lvc.Fmt |= win32.LVCFMT_CENTER
	}
	pwsz, _ := syscall.UTF16PtrFromString(column.Title)
	lvc.PszText = pwsz
	if column.Width != 0 {
		lvc.Cx = int32(column.Width)
	} else {
		cx, _ := MeasureText(this.hWndHeader, column.Title)
		lvc.Cx = int32(cx) + 16
	}
	if column.MinWidth != 0 {
		lvc.Mask |= win32.LVCF_MINWIDTH
		lvc.CxMin = int32(column.MinWidth)
	}
	index := colCount
	if dummyColumnAdded {
		index += 1
	}
	ret, errno := SendMessage(this.Handle, win32.LVM_INSERTCOLUMN,
		index, unsafe.Pointer(&lvc))
	if ret == NegativeOne {
		log.Fatal(errno)
	}
	if dummyColumnAdded {
		_, _ = SendMessage(this.Handle, win32.LVM_DELETECOLUMN, 0, 0)
	}
}

func (this *ListViewObject) GetSelectedIndex() int {
	ret, errno := SendMessage(this.Handle,
		win32.LVM_GETNEXTITEM, NegativeOne, win32.LVNI_SELECTED)
	_ = errno
	return int(ret)
}

func (this *ListViewObject) GetSelectedIndexes() []int {
	index := -1
	var indexes []int
	for {
		ret, errno := SendMessage(this.Handle,
			win32.LVM_GETNEXTITEM, index, win32.LVNI_SELECTED)
		_ = errno
		if ret == NegativeOne {
			break
		}
		index = int(ret)
		indexes = append(indexes, index)
	}
	return indexes
}

func (this *ListViewObject) SetSelectedIndex(index int) {
	var lvi win32.LVITEM
	lvi.StateMask = win32.LVIS_SELECTED
	if !this.SingleRowSelect {
		SendMessage(this.Handle, win32.LVM_SETITEMSTATE,
			NegativeOne, unsafe.Pointer(&lvi))
	}
	lvi.State = win32.LVIS_SELECTED
	SendMessage(this.Handle, win32.LVM_SETITEMSTATE,
		index, unsafe.Pointer(&lvi))
}

func (this *ListViewObject) GetSelectedData() any {
	rowId := this.GetSelectedId()
	if rowId == 0 {
		return nil
	}
	if data, ok := this.dataMap[rowId]; ok {
		return data
	}
	return nil
}

func (this *ListViewObject) GetSelectedId() int {
	index := this.GetSelectedIndex()
	return this.IndexToId(index)
}

func (this *ListViewObject) IndexToId(index int) int {
	if index == -1 {
		return 0
	}
	var lvi win32.LVITEM
	lvi.IItem = int32(index)
	lvi.Mask = win32.LVIF_PARAM
	ret, errno := SendMessage(this.Handle, win32.LVM_GETITEM,
		0, unsafe.Pointer(&lvi))
	if ret == 0 {
		log.Fatal(errno)
	}
	rowId := int(lvi.LParam)
	return rowId
}

func (this *ListViewObject) GetRowData(id int) any {
	data, _ := this.dataMap[id]
	return data
}

func (this *ListViewObject) SetRowData(id int, data any) {
	this.dataMap[id] = data
}

func (this *ListViewObject) GetRowCount() int {
	ret, _ := win32.SendMessage(this.Handle, win32.LVM_GETITEMCOUNT, 0, 0)
	return int(ret)
}
