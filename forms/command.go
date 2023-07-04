package forms

type Command struct {
	Id       int
	Name     string
	Text     string //u
	Image    int    //u //default 0 means null. use zero to specify 0.
	Tooltip  string
	Desc     string
	Action   Action
	Category string

	Disabled   bool
	RadioGroup string
	Checked    bool //u

	ShortcutKeys []KeyStroke

	OnChange  SimpleEvent
	OnExecute SimpleEvent
}

func NewCommand(text string, action Action) *Command {
	return &Command{Text: text, Action: action}
}

type CommandAware interface {
	SetCommand(command *Command)
	GetCommand() *Command
}

type CommandManager struct {
	Items []*Command

	idItemMap map[int]*Command
}

func NewCommandManager() *CommandManager {
	obj := &CommandManager{}
	obj.Init()
	return obj
}

func (this *Command) GetNoPrefixText() string {
	bts := []byte(this.Text)
	cb := len(bts)
	removedCount := 0
	for n := 0; n < cb; n++ {
		b := bts[n]
		if b == '&' {
			copy(bts[n:], bts[n+1:])
			n += 1
			removedCount += 1
		}
	}
	if removedCount > 0 {
		bts = bts[:cb-removedCount]
	}
	return string(bts)
}

func (this *Command) SetState(disabled bool, checked bool) {
	changed := false
	if disabled != this.Disabled {
		this.Disabled = disabled
		changed = true
	}
	if checked != this.Checked {
		this.Checked = checked
		changed = true
	}
	if changed {
		this.NotifyChange()
	}
}

func (this *Command) NotifyChange() {
	this.OnChange.Fire(this, &SimpleEventInfo{})
}

func (this *Command) NotifyExecute() {
	if this.Action != nil {
		this.Action()
	}
	this.OnExecute.Fire(this, &SimpleEventInfo{})
}

//
func (this *CommandManager) Init() {
	this.idItemMap = make(map[int]*Command)
}

func (this *CommandManager) AddItems(items []*Command) {
	this.Items = append(this.Items, items...)
	for _, item := range items {
		if item.Id != 0 {
			this.idItemMap[item.Id] = item
		}
	}
}

func (this *CommandManager) Item(id int) *Command {
	return this.idItemMap[id]
}

func shortcutKeysChanged(keys1 []KeyStroke, keys2 []KeyStroke) bool {
	count := len(keys1)
	if count != len(keys2) {
		return true
	}
	for n := 0; n < count; n++ {
		if keys1[n] != keys2[n] {
			return true
		}
	}
	return false
}

func (this *CommandManager) SetShortcuts(shortcuts map[int][]KeyStroke) {
	for _, item := range this.Items {
		keyStrokes := shortcuts[item.Id]
		if shortcutKeysChanged(item.ShortcutKeys, keyStrokes) {
			item.ShortcutKeys = keyStrokes
			item.NotifyChange()
		}
	}
}
