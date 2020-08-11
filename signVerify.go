package main

import (
	// "crypto"
	//"crypto/rand"
	//"os"

	// "crypto/rsa"
	"crypto/sha256"
	// "time"

	//"encoding/json"
	//"encoding/base64"
	"fmt"
	// "os"
	//"time"
	"crypto/rand"
	"log"
	rsapk "rsapk"

	cryptogrpghy "./model/cryptogrpghy"
	//"./model/transaction"
)

// EncryptWithPublicKey encrypts data with Public key
func EncryptWithPublicKey(ciphertext []byte, pub *rsapk.PublicKey) ([]byte, error) {
	hash := sha256.New()
	plaintext, err := rsapk.EncryptOAEP(hash, rand.Reader, pub, ciphertext, nil)
	if err != nil {
		//log.Error(err)
		log.Fatalln(err)
	}
	return plaintext, err
}

func main() {
	senderPubKey := "1Ap9xCRC8pQkZqte1pW2G16F84XnFoHS6u"
	receiverPubKey := "1BT7bFdxyyMTvzsTdCgWhWRcsCUCtbNjru"
	amount := fmt.Sprintf("%f", 2.0)
	//timeNow := time.Now().Add(time.Hour * 4).Format("2006-01-02T15:04:05Z07:00") //time.Now().UTC().Format("2006-01-02T03:04:05+00:00")
	secretMessage := senderPubKey + receiverPubKey + amount // + timeNow
	private := cryptogrpghy.ParsePEMtoRSAprivateKey("-----BEGIN RSA PRIVATE KEY-----\nMIICXgIBAAKBgQDIkF7/oFOYFqu6FulTUCRrPrNoBK0ewbMsARA6hORVdPzOmzok\nURduzStquRrnwZXj6/EECPS510zLxGKNSDBGkBSW9/1ctiYMhmWLLfQsXzd2Abky\nJ0kPUlex0xo683cW2+2G6IMczdaE51SgGTh3lcjglMt4sJQtVgXma1s7jQIDAQAB\nAoGBALoHz1Xj7CXBwX9WCQ3R5DXlbpso2zsQB5TlV5wv72qknGk26fMNlGKdw4u2\nLhKRKOrDykYn2HcYEI9glNjfAIaM6oe4TYydJuo33Yo0ZnrzTyq+/Q/atjoT0hc1\nxVp+oJZKocJeMPzhP8DbwZSlq4teR1+BKj2LaAS4k6CHAXIxAkEA9wNqRe3LQ2rI\n/oycaN68r83Zx9yEZvNEyooEOVFBArHXV6B0TOxttLa9IY2oUxKYZecip47Q/Hpc\nUUnh4KRIdwJBAM/cWDqjYgikFk5JnrIDBuYrCaUVgNNT9+/oUJkB/zrk7jULCM3O\nu43pLOoi3/Ut+FivzfYXHLAnBdOPqLkh4RsCQQDvTzK1nwTvQtSJsKaT/z8kv7U/\nKUhpCURbSU2ATlVCjBOKBJzILcK3ctdXW4t5OCnHiB+N4BJemRk5c+/PGLpPAkEA\nqBdiahkSADbhqvGyGfaEr8GCDTQ0d7FhwWq3MuUAh5n2YILJ3dUeqwYzwivtvJIu\nUVnqTuYl1vXXqlx0bzJMnQJAbucZhjDN0MUasAU1zoV54v92BehDovK4jdzV37+x\nnd4LFmLP0YAeWROe7Tll+si80rbZOIBrnq+oyEE/+2LOFA==\n-----END RSA PRIVATE KEY-----\n")
	signature := cryptogrpghy.SignPKCS1v15(secretMessage, *private)
	//alicePublicKey := cryptogrpghy.ParsePEMtoRSApublicKey("-----BEGIN RSA PUBLIC KEY-----\nMIGJAoGBAMcGY0PmySh5009OgIFZ7DxU729b835OnCpdh9ooC80kxDnQMTMG3egX\n34H19JuBt/ud3f0wQIvj/JGZhQn9lc2KLU81mLke1MPEy5JxNi/Skz//93YTcQKB\nHSeGoq0smpixMD6hisLxpCsomFhROvC/zh7dxNKq8pcA4oBjMw27AgMBAAE=\n-----END RSA PUBLIC KEY-----\n")
	fmt.Println("Signature:", signature)
}
