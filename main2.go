package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
)

type pkcs1PrivateKey struct {
	Version int
	N       *big.Int
	E       int
	D       *big.Int
	P       *big.Int
	Q       *big.Int
	// We ignore these values, if present, because rsa will calculate them.
	Dp   *big.Int `asn1:"optional"`
	Dq   *big.Int `asn1:"optional"`
	Qinv *big.Int `asn1:"optional"`

	AdditionalPrimes []pkcs1AdditionalRSAPrime `asn1:"optional,omitempty"`
}

type pkcs1AdditionalRSAPrime struct {
	Prime *big.Int

	// We ignore these values because rsa will calculate them.
	Exp   *big.Int
	Coeff *big.Int
}
type pkcsType int64

const (
	rsaAlgorithmSign = crypto.SHA256

	PKCS1 pkcsType = iota
	PKCS8
)

type XRsa struct {
	keyLen         int
	privateKeyType pkcsType
	publicKey      *rsa.PublicKey
	privateKey     *rsa.PrivateKey
}

func main() {
	dat, _ := ioutil.ReadFile("validator/public.pem")
	publicKeyData := string(dat)
	encryptedkeydata := "d5fce9709e9cae77db9ecf540758e7c3d6787acd"
	hashedkeyenc, _ := PublicEncrypt(publicKeyData, encryptedkeydata)
	fmt.Println("RSA :", hashedkeyenc, " the length is ", len(hashedkeyenc))
}
func PublicEncrypt(publicKey, data string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)

	partLen := pub.N.BitLen()/8 - 11
	p := string(partLen)
	chunks := strings.Split(data, p)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bts, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(chunk))
		if err != nil {
			return "", err
		}
		buffer.Write(bts)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}
