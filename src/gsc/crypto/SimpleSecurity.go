package crypto

import "crypto/des"

type ISimpleSecurity interface {
	Encrypt(key string, data []byte) []byte
	Decrypt(key string, data []byte) []byte
}

type SimpleSecurity struct {
}

func (simpleSecurity *SimpleSecurity) Encrypt(key string, data []byte) []byte {
	ede2Key := []byte("example key 1234")

	var tripleDESKey []byte
	tripleDESKey = append(tripleDESKey, ede2Key[:16]...)
	tripleDESKey = append(tripleDESKey, ede2Key[:8]...)

	_, err := des.NewTripleDESCipher(tripleDESKey)
	if err != nil {
		panic(err)
	}
	return nil
}

func (simpleSecurity *SimpleSecurity) Decrypt(key string, data []byte) []byte {
	return nil
}
