package components

import (
	"github.com/sf-go/gsf/src/gsm/component"
)

type UserComponent struct {
	*component.Component

	Account  string
	Password string
}

func NewUserComponent() *UserComponent {
	userComponent := &UserComponent{
		Component: component.NewComponent(),
	}
	userComponent.Register(userComponent)
	userComponent.SetValue("Account", "123456")
	userComponent.SetValue("Password", "456789")
	return userComponent
}

func (component *UserComponent) GetObjectId() string {
	return "User"
}
