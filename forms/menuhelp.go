package forms

type MenuHelpUI interface {
	EnterMenuHelp()
	ShowMenuHelp(text string)
	ExitMenuHelp()
}

type StatusBarMenuHelpUI struct {
	statusBar  StatusBar
	simpleMode bool

	oriSimple bool
	oriText   string
}

func NewStatusBarMenuHelpUI(statusBar StatusBar, simpleMode bool) *StatusBarMenuHelpUI {
	return &StatusBarMenuHelpUI{
		statusBar:  statusBar,
		simpleMode: simpleMode,
	}
}

func (this *StatusBarMenuHelpUI) EnterMenuHelp() {
	if this.simpleMode {
		this.oriSimple = this.statusBar.IsSimpleMode()
		if !this.oriSimple {
			this.statusBar.SetSimpleMode(true)
		}
	}
	this.oriText = this.statusBar.GetText()
}

func (this *StatusBarMenuHelpUI) ShowMenuHelp(text string) {
	this.statusBar.SetText(text)
}

func (this *StatusBarMenuHelpUI) ExitMenuHelp() {
	if this.simpleMode {
		if !this.oriSimple {
			this.statusBar.SetSimpleMode(false)
		}
	}
	this.statusBar.SetText(this.oriText)
}
