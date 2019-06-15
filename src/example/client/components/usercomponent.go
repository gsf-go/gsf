package components

import (
	"github.com/gsf/gsf/src/gsc/property"
	"reflect"
)

type UserComponent struct {
	*property.Property

	Account  string
	Password string
}

func NewUserComponent() *UserComponent {
	userComponent := &UserComponent{
		Property: property.NewProperty(),
	}

	userComponent.Register(userComponent)
	return userComponent
}

func (userComponent *UserComponent) GetName() string {
	return reflect.TypeOf(userComponent).Elem().Name()
}
