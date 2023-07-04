package scope

import "github.com/zzl/goforms/framework/types"

//

type Scope struct {
	disposables []types.Disposable
}

func NewScope() *Scope {
	return &Scope{}
}

func (this *Scope) Add(disposable types.Disposable) {
	this.disposables = append(this.disposables, disposable)
}

func (this *Scope) Leave() {
	disposables := this.disposables
	this.disposables = nil
	for n := len(disposables) - 1; n >= 0; n -= 1 {
		disposables[n].Dispose()
	}
}

type ScopedFunc func(s *Scope)

func WithScope(scopedFunc ScopedFunc) {
	s := NewScope()
	defer s.Leave()
	scopedFunc(s)
}
