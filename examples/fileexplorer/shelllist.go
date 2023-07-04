package main

import (
	"bytes"
	"github.com/zzl/go-com/com"
	"github.com/zzl/goforms/forms"
	"math"
	"strconv"
)

type ShellList struct {
	forms.ListViewObject
	super *forms.ListViewObject

	shell *Shell

	OnLoaded forms.SimpleEvent
	Loading  bool
}

func (this *ShellList) Init() {
	this.super.Init()
	this.Columns = []*forms.ListViewColumn{
		{Title: "Name", Width: 200},
		{Title: "Date modified", Width: 100},
		{Title: "Type", Width: 100},
		{Title: "Size", Width: 80},
	}
}

func (this *ShellList) PostCreate(opts *forms.WindowOptions) {
	hIml := this.shell.GetImageList()
	this.SetImageList(forms.NewImageListFromHandle(hIml, false))
}

func formatNum(num int) string {
	sNum := strconv.Itoa(num)
	cb := len(sNum)
	var buf bytes.Buffer

	for n := 0; n < cb; n++ {
		if (cb-n)%3 == 0 && n != 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte(sNum[n])
	}
	return buf.String()
}

func (this *ShellList) loadRows(folderItem *ShellItem) []*forms.ListRowInfo {
	var rows []*forms.ListRowInfo

	items := folderItem.GetChildren(true)
	for _, item := range items {
		var sModDate string
		if !item.ModTime.IsZero() {
			sModDate = item.ModTime.Format("2006-01-02 15:04")
		}
		var sSize string
		if item.Size != -1 {
			kb := (int)(math.Ceil(float64(item.Size) / 1024))
			sSize = formatNum(kb) + " KB"
		}
		rows = append(rows, &forms.ListRowInfo{
			Data:  item,
			Image: item.Icon,
			CellValues: []any{
				item.Name, sModDate, item.Type, sSize,
			}})
	}
	return rows
}

// async
func (this *ShellList) _load(folderItem *ShellItem) {
	com.Initialize()
	rows := this.loadRows(folderItem)
	forms.Dispatcher.Invoke(func() {
		this.ClearRows()
		this.AddRowInfos(rows)
		this.OnLoaded.Fire(this, &forms.SimpleEventInfo{})
		this.Loading = false
	}, true)
}

func (this *ShellList) Load(folderItem *ShellItem) {
	this.Loading = true
	wc := forms.NewWaitCursor()
	go func() {
		this._load(folderItem)
		wc.Restore()
	}()
}

func (this *ShellList) SetSelectedShellItem(si *ShellItem) {
	count := this.GetRowCount()
	for n := 0; n < count; n++ {
		id := this.IndexToId(n)
		tsi := this.GetRowData(id).(*ShellItem)
		if bytes.Equal(tsi.id, si.id) {
			this.SetSelectedIndex(n)
			return
		}
	}
}

func (this *ShellList) GetSelectedShellItem() *ShellItem {
	id := this.GetSelectedId()
	if id == 0 {
		return nil
	}
	return this.GetRowData(id).(*ShellItem)
}
