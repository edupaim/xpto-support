package domain

type CryptDocument string

func (c *CryptDocument) Encrypt() error {
	crypt, err := encrypt(string(*c))
	if err != nil {
		return err
	}
	*c = CryptDocument(crypt)
	return nil
}

func (c *CryptDocument) Decrypt() error {
	crypt, err := decrypt(string(*c))
	if err != nil {
		return err
	}
	*c = CryptDocument(crypt)
	return nil
}
