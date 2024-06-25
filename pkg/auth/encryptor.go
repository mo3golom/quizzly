package auth

import "crypto/aes"

type defaultEncryptor struct {
	secretKey string
}

func (e *defaultEncryptor) Encrypt(data string) (string, error) {
	cipher, err := aes.NewCipher([]byte(e.secretKey))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(data))
	cipher.Encrypt(ciphertext, []byte(data))

	return string(ciphertext), nil
}

func (e *defaultEncryptor) Decrypt(data string) (string, error) {
	cipher, err := aes.NewCipher([]byte(e.secretKey))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(data))
	cipher.Decrypt(ciphertext, []byte(data))

	return string(ciphertext), nil
}
