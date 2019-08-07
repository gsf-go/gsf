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
	userComponent.SetValue("Account", "111111")
	userComponent.SetValue("Password", "222222")
	return userComponent
}

func (component *UserComponent) GetObjectId() string {
	return "User"
}
