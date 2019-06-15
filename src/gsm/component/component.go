package component

type DefaultComponent struct {
}

func (defaultComponent *DefaultComponent) GetName() string {
	return "Default"
}
