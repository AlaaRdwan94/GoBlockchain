package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	cryptogrpghy "./model/cryptogrpghy"
)

func keyEncrypt(keyStr string, cryptoText string) string {
	keyBytes := sha256.Sum256([]byte(keyStr))
	return encrypt(keyBytes[:], cryptoText)
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

type DigitalwalletTransaction struct {
	Sender    string
	Receiver  string
	TokenID   string
	Amount    float32
	Time      string
	Signature string
}

func main() {
	//------------------------------------------------
	//------------------------------------------------
	senderPubKey := "1EnUqL2QdxJ4U7mgyXPMrGfLynfZsw8aAN"
	receiverPubKey := "1CWYFbdnjVXWHuiKG4CvZrszUKQHMAQKDL"
	amount := fmt.Sprintf("%f", 2.0)
	secretMessage := senderPubKey + receiverPubKey + amount // + timeNow
	private := cryptogrpghy.ParsePEMtoRSAprivateKey("-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQDIkF7/oFOYFqu6FulTUCRrPrNoBK0ewbMsARA6hORVdPzOmzok\nURduzStquRrnwZXj6/EECPS510zLxGKNSDBGkBSW9/1ctiYMhmWLLfQsXzd2Abky\nJ0kPUlex0xo683cW2+2G6IMczdaE51SgGTh3lcjglMt4sJQtVgXma1s7jQIDAQAB\nAoGBALoHz1Xj7CXBwX9WCQ3R5DXlbpso2zsQB5TlV5wv72qknGk26fMNlGKdw4u2\nLhKRKOrDykYn2HcYEI9glNjfAIaM6oe4TYydJuo33Yo0ZnrzTyq+/Q/atjoT0hc1\nxVp+oJZKocJeMPzhP8DbwZSlq4teR1+BKj2LaAS4k6CHAXIxAkEA9wNqRe3LQ2rI\n/oycaN68r83Zx9yEZvNEyooEOVFBArHXV6B0TOxttLa9IY2oUxKYZecip47Q/Hpc\nUUnh4KRIdwJBAM/cWDqjYgikFk5JnrIDBuYrCaUVgNNT9+/oUJkB/zrk7jULCM3O\nu43pLOoi3/Ut+FivzfYXHLAnBdOPqLkh4RsCQQDvTzK1nwTvQtSJsKaT/z8kv7U/\nKUhpCURbSU2ATlVCjBOKBJzILcK3ctdXW4t5OCnHiB+N4BJemRk5c+/PGLpPAkEA\nqBdiahkSADbhqvGyGfaEr8GCDTQ0d7FhwWq3MuUAh5n2YILJ3dUeqwYzwivtvJIu\nUVnqTuYl1vXXqlx0bzJMnQJAbucZhjDN0MUasAU1zoV54v92BehDovK4jdzV37+x\nnd4LFmLP0YAeWROe7Tll+si80rbZOIBrnq+oyEE/+2LOFA==\n-----END RSA PRIVATE KEY-----\n")
	signature := cryptogrpghy.SignPKCS1v15(secretMessage, *private)
	// fmt.Println("Signature:", signature)
	//-----------------------------------------------
	//-----------------------------------------------
	var Digitalwallet = DigitalwalletTransaction{}
	Digitalwallet.Sender = senderPubKey     //"1EnUqL2QdxJ4U7mgyXPMrGfLynfZsw8aAN"
	Digitalwallet.Receiver = receiverPubKey //"1BT7bFdxyyMTvzsTdCgWhWRcsCUCtbNjru"
	Digitalwallet.TokenID = "0000000000000"
	Digitalwallet.Amount = 2.0
	Digitalwallet.Time = time.Now().Add(time.Hour * 4).Format("2006-01-02T15:04:05Z07:00")
	Digitalwallet.Signature = signature //"OaDDY8a7X9T6XFQcnS09ci4Hh4DVQe4AMi5tqIrLeSWDvBLnOSO4T9kNxDyQ595wVIBhbhDl6Ux+zg2WLipAXhXwxqjSO2ybSkaEZd7m2LGSJ9NfnGBfUWo+e1UDFkkZAcfb2oP7hDqKWfX7Hkkp2WW6cctpcCHvPj+giHFCW3c="
	ddddd, _ := json.Marshal(Digitalwallet)
	AESdata := string(ddddd)
	AESenc := keyEncrypt("d5fce9709e9cae77db9ecf540758e7c3d6787acd", AESdata)
	fmt.Println("AES encryption :", AESenc)
}
