package genservice

type (
	encoder interface {
		Encode(src any) ([]byte, error)
	}

	encrypter interface {
		SetKey(string)
		Encrypt([]byte) ([]byte, error)
	}
)
