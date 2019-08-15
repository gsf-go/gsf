package component

type IComponent interface {
	GetObjectId() string

	GetterCallback(version string) []interface{}
	SetterCallback(args ...interface{}) bool
	Synchronize() (string, []interface{})
}
