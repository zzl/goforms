package forms

import (
	"log"

	"github.com/zzl/go-win32api/v2/win32"
)

type Accelerator struct {
	CmdId uint16

	Ctrl  bool
	Alt   bool
	Shift bool
	Key   byte

	Command *Command
	Action  Action
}

type AcceleratorTable struct {
	Handle win32.HACCEL
	Items  []*Accelerator

	idGen          UidGen
	OnHandleChange SimpleEvent
}

func NewAcceleratorTable() *AcceleratorTable {
	a := &AcceleratorTable{
		idGen: UidGen{nextId: AccelGenIdStart, step: -1},
	}
	return a
}

func (this *AcceleratorTable) Dispose() {
	if this.Handle == 0 {
		return
	}
	_ = win32.DestroyAcceleratorTable(this.Handle)
	this.Handle = 0
}

func (this *AcceleratorTable) ReCreate() error {
	oriHandle := this.Handle
	err := this.Create()
	if err == nil {
		_ = win32.DestroyAcceleratorTable(oriHandle)
		this.OnHandleChange.Fire(this, &SimpleEventInfo{})
	} else {
		this.Handle = oriHandle
	}
	return err
}

func (this *AcceleratorTable) Create() error {
	var accels []win32.ACCEL

	for _, item := range this.Items {
		if item.CmdId == 0 {
			item.CmdId = uint16(this.idGen.Gen())
		}
		accel := win32.ACCEL{}

		ctrl := item.Ctrl
		alt := item.Alt
		shift := item.Shift
		key := item.Key

		var ks *KeyStroke
		if item.Command != nil && len(item.Command.ShortcutKeys) > 0 {
			ks = &item.Command.ShortcutKeys[0]
		}
		if key == 0 && ks != nil {
			ctrl = ks.Ctrl
			alt = ks.Alt
			shift = ks.Shift
			key = ks.Key
		}

		if ctrl {
			accel.FVirt |= win32.FCONTROL
		}
		if alt {
			accel.FVirt |= win32.FALT
		}
		if shift {
			accel.FVirt |= win32.FSHIFT
		}
		if key != 0 {
			accel.FVirt |= win32.FVIRTKEY
			accel.Key = uint16(key)
		} else {
			//accel.Key = uint16(item.Char)
			log.Fatal("key undefined")
		}
		accel.Cmd = item.CmdId
		accels = append(accels, accel)
	}

	hAccel, errno := win32.CreateAcceleratorTable(&accels[0], int32(len(accels)))
	if hAccel == 0 {
		log.Fatal(errno)
	}
	this.Handle = hAccel
	return nil
}

func (this *AcceleratorTable) FindById(id uint16) *Accelerator {
	for _, it := range this.Items {
		if it.CmdId == id {
			return it
		}
	}
	return nil
}

func (this *AcceleratorTable) FindByCommand(command *Command) *Accelerator {
	for _, it := range this.Items {
		if it.Command == command {
			return it
		}
	}
	return nil
}

func (this *AcceleratorTable) AddTabPageKeys(tabCtrl TabControl) {
	items := []*Accelerator{
		{Ctrl: true, Key: byte(win32.VK_TAB), Action: tabCtrl.SelectNext},
		{Ctrl: true, Key: byte(win32.VK_NEXT), Action: tabCtrl.SelectNext},
		{Ctrl: true, Shift: true, Key: byte(win32.VK_TAB), Action: tabCtrl.SelectPrev},
		{Ctrl: true, Key: byte(win32.VK_PRIOR), Action: tabCtrl.SelectPrev},
	}
	this.Items = append(this.Items, items...)
}

func (this *AcceleratorTable) onReflectCommand(msg *CommandMessage) {
	accel := this.FindById(msg.GetCmdId())
	if accel == nil {
		return
	}

	command := accel.Command
	if command != nil {
		command.NotifyExecute()
	}
	if accel.Action != nil {
		accel.Action()
	}
}

func (this *AcceleratorTable) ClearItems() {
	this.Items = nil
}

func (this *AcceleratorTable) AddItemsFromCommand(commands []*Command) {
	for _, cmd := range commands {
		for _, sk := range cmd.ShortcutKeys {
			item := &Accelerator{
				CmdId:   uint16(cmd.Id),
				Ctrl:    sk.Ctrl,
				Alt:     sk.Alt,
				Shift:   sk.Shift,
				Key:     sk.Key,
				Command: cmd,
			}
			this.Items = append(this.Items, item)
		}
	}
}
