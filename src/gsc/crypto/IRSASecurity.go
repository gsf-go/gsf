package crypto

type IRSASecurity interface {
	Encrypt(key string, publicKey string) []byte
	Decrypt(data []byte, publicKey string, privateKey string) []byte
}

type RSASecurity struct {
}

func (rsa *RSASecurity) Encrypt(key string, publicKey string) []byte {
	return nil
}

func (rsa *RSASecurity) Decrypt(data []byte, publicKey string, privateKey string) []byte {
	return nil
}
