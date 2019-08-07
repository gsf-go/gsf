package component

type IComponent interface {
	GetObjectId() string

	UpdateField(name string, value interface{}) bool
}
