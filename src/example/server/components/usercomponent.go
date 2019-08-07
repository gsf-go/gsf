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
	return userComponent
}

func (component *UserComponent) GetObjectId() string {
	return "User"
}

//func (component *UserComponent) Setter(cpt component.IComponent) {
//	userComponent := cpt.(*UserComponent)
//	component.Account = userComponent.Account
//	component.Password = userComponent.Password
//}
