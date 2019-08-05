package component

type IComponent interface {
	GetObjectId() string

	Update(name string, value interface{}) bool
}
