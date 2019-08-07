package component

type IComponent interface {
	GetObjectId() string

	Getter(version string) []interface{}
	Setter(name string, value interface{}) bool

	Update() bool
}
