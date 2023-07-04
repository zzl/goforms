package virtual

import (
	"reflect"
	"unsafe"

	"github.com/zzl/goforms/framework/types"
)

// Virtual is an interface that exposes a realObject property
type Virtual interface {
	SetRealObject(object any)
	GetRealObject() any
}

// VirtualObject implements the Virtual interface.
// The Init method of the realObject will be automatically invoked upon setting.
type VirtualObject[T any] struct {
	RealObject T
}

// SetRealObject sets the real object of the virtual object.
// It also initializes the real object if it implements the Initable interface.
func (this *VirtualObject[T]) SetRealObject(realObject any) {
	this.RealObject = realObject.(T)
	if i, ok := realObject.(types.Initable); ok {
		i.Init()
	}
}

// GetRealObject returns the real object associated with the virtual object.
func (this *VirtualObject[T]) GetRealObject() any {
	return this.RealObject
}

// New creates a new virtual object.
// If a pre-constructed object is provided, it is used;
// otherwise, a new object is created.
// If the object implements the Virtual interface,
// Realize is called to set up the virtual object.
func New[T any](preConstructed ...*T) *T {
	var pObj *T
	if len(preConstructed) == 1 {
		pObj = preConstructed[0]
	} else {
		var obj T
		pObj = &obj
	}
	if _, ok := any(pObj).(Virtual); ok {
		Realize(pObj)
	}
	return pObj
}

// Realize sets up a virtual object by automatically assigning the "super" field
// across the whole object hierarchy and then calling
// SetRealObject on it with itself as the realObject.
func Realize(virtualObject any) {
	virtual := virtualObject.(Virtual)
	setupSuper(virtual)
	virtual.SetRealObject(virtualObject)
}

// setupSuper traverses the object hierarchy to find and set the "super" field.
func setupSuper(object Virtual) {
	objectType := reflect.TypeOf(object).Elem()
	objectValue := reflect.ValueOf(object).Elem()

	count := objectType.NumField()
	var superField *reflect.StructField
	var superObject Virtual
	for n := 0; n < count; n++ {
		field := objectType.Field(n)
		if field.Anonymous && superObject == nil {
			addr := objectValue.FieldByIndex(field.Index).Addr()
			if obj, ok := addr.Interface().(Virtual); ok {
				superObject = obj
				if superField != nil {
					break
				}
			}
		} else if field.Name == "super" {
			superField = &field
			if superObject != nil {
				break
			}
		}
	}
	if superObject != nil {
		if superField != nil {
			fieldValue := objectValue.FieldByIndex(superField.Index)
			fieldValue = reflect.NewAt(superField.Type,
				unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
			fieldValue.Set(reflect.ValueOf(superObject).Elem().Addr())
		}
		setupSuper(superObject)
	}
}
