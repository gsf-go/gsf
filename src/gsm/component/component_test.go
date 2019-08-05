package component

import (
	"github.com/sf-go/gsf/src/gsc/bytestream"
	"github.com/sf-go/gsf/src/gsc/property"
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsm/peer"
	"reflect"
	"testing"
)

type UserComponent struct {
	*property.Property

	Account  string
	Password string
}

func NewUserComponent() *UserComponent {
	userComponent := &UserComponent{
		Property: property.NewProperty(),
	}

	userComponent.Register(userComponent)
	return userComponent
}

func (userComponent *UserComponent) GetObjectId() string {
	return reflect.TypeOf(userComponent).Elem().Name()
}

func (userComponent *UserComponent) Setter(cpt IComponent) {

}
func TestComponent(t *testing.T) {

	serialization.PacketManagerInstance.AddPacket("SerializablePacket",
		func(name string, args ...interface{}) serialization.ISerializablePacket {
			p := args[1].(peer.IPeer)
			return p.GetComponent(name).(serialization.ISerializablePacket)
		})

	peer := peer.NewPeer()
	peer.AddComponent(NewUserComponent())

	sut := NewUserComponent()
	sut.SetValue("Account", "account")
	sut.SetValue("Password", "123456")
	writer := serialization.NewSerializable()
	bytes := sut.ToBinaryWriter(writer)

	sut.Account = ""
	sut.Password = ""

	reader := serialization.NewDeserializable(bytestream.NewByteReader2(bytes))
	sut2 := NewUserComponent()
	sut2.FromBinaryReader(reader)
}
