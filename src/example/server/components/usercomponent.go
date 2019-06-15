package components

type UserComponent struct {
}

func NewUserComponent() *UserComponent {
	return &UserComponent{}
}

func (userComponent *UserComponent) GetName() string {
	return "User"
}
