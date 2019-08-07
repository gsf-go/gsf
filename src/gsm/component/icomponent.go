package component

type IComponent interface {
	GetObjectId() string

	Getter(version string) []interface{}
	Setter(args ...interface{}) bool
}
