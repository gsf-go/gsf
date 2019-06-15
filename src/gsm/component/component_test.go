package component

import (
	"github.com/sf-go/gsf/src/gsc/property"
	"github.com/sf-go/gsf/src/gsc/serialization"
	"github.com/sf-go/gsf/src/gsf/peer"
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

func (userComponent *UserComponent) GetName() string {
	return reflect.TypeOf(userComponent).Elem().Name()
}

func TestComponent(t *testing.T) {

	serialization.GetPacketManagerInstance().AddPacket("SerializablePacket",
		func(args ...interface{}) serialization.ISerializablePacket {
			name := args[0].(string)
			p := args[1].(peer.IPeer)
			return p.GetComponent(name).(serialization.ISerializablePacket)
		})

	peer := peer.NewPeer()
	peer.AddComponent(NewUserComponent())

	sut := NewUserComponent()
	sut.SetValue("Account", "account")
	sut.SetValue("Password", "123456")
	writer := serialization.NewEndianBinaryWriter()
	sut.ToBinaryWriter(writer)
	bytes := writer.ToBytes()
	sut.Account = ""
	sut.Password = ""
	reader := serialization.NewEndianBinaryReader(bytes, peer)
	sut2 := NewUserComponent()
	sut2.FromBinaryReader(reader)
}
